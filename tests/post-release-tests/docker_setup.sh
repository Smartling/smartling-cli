#!/bin/bash

# Docker Test Environment Setup Script
# This script helps run the Smartling CLI test suite in a Docker container

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_BINARY="${CLI_BINARY:-./smartling-cli}"
CONFIG_FILE="${CONFIG_FILE:-./smartling.yml}"
TEST_SUITE="${TEST_SUITE:-./test-suite.sh}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $*"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed or not in PATH"
        exit 1
    fi
    
    # Check CLI binary
    if [[ ! -f "$CLI_BINARY" ]]; then
        log_error "CLI binary not found: $CLI_BINARY"
        log_info "Please ensure the binary exists or set CLI_BINARY environment variable"
        exit 1
    fi
    
    # Check config file
    if [[ ! -f "$CONFIG_FILE" ]]; then
        log_error "Config file not found: $CONFIG_FILE"
        log_info "Please ensure smartling.yml exists or set CONFIG_FILE environment variable"
        exit 1
    fi
    
    # Check test suite
    if [[ ! -f "$TEST_SUITE" ]]; then
        log_error "Test suite not found: $TEST_SUITE"
        log_info "Please ensure test-suite.sh exists or set TEST_SUITE environment variable"
        exit 1
    fi
    
    # Make CLI binary executable
    chmod +x "$CLI_BINARY"
    chmod +x "$TEST_SUITE"
    
    log_success "All prerequisites checked"
}

# Run tests in Docker
run_docker_tests() {
    log_info "Starting Docker test environment..."
    
    local abs_cli_binary=$(realpath "$CLI_BINARY")
    local abs_config_file=$(realpath "$CONFIG_FILE")
    local abs_test_suite=$(realpath "$TEST_SUITE")
    local abs_script_dir=$(realpath "$SCRIPT_DIR")
    
    log_info "CLI Binary: $abs_cli_binary"
    log_info "Config File: $abs_config_file"
    log_info "Test Suite: $abs_test_suite"
    
    # Create temporary directory for test results
    local results_dir="${SCRIPT_DIR}/test-results-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$results_dir"
    
    log_info "Test results will be saved to: $results_dir"
    
    # Run Docker container
    log_info "Launching Docker container..."
    
    local docker_exit_code=0
    docker run -it --rm \
        -v "$abs_test_suite:/test/test-suite.sh:ro" \
        -v "$abs_config_file:/test/smartling.yml:ro" \
        -v "$abs_cli_binary:/test/smartling-cli:ro" \
        -v "$results_dir:/test/results" \
        -w /test \
        -e "CLI_BINARY=/test/smartling-cli" \
        -e "CONFIG_FILE=/test/smartling.yml" \
        ubuntu:20.04 \
        bash -c "
            apt-get update -qq && apt-get install -y -qq curl wget ca-certificates > /dev/null 2>&1
            echo 'Starting Smartling CLI Test Suite...'
            chmod +x /test/smartling-cli /test/test-suite.sh
            /test/test-suite.sh
            # Copy results to mounted volume
            if [ -d /tmp/tmp.* ]; then
                find /tmp -name 'tmp.*' -type d -exec cp -r {}/* /test/results/ \; 2>/dev/null || true
            fi
        " || docker_exit_code=$?
    
    # Report results
    if [[ $docker_exit_code -eq 0 ]]; then
        log_success "All tests passed!"
    else
        log_error "Some tests failed (exit code: $docker_exit_code)"
    fi
    
    # Show results location
    if [[ -d "$results_dir" ]] && [[ "$(ls -A "$results_dir" 2>/dev/null)" ]]; then
        log_info "Test results available in: $results_dir"
        
        # Show summary if results.json exists
        local results_json="$results_dir/results.json"
        if [[ -f "$results_json" ]]; then
            log_info "Test Summary:"
            if command -v jq &> /dev/null; then
                jq -r '"  Total: " + (.total_tests | tostring) + ", Passed: " + (.passed | tostring) + ", Failed: " + (.failed | tostring) + ", Success Rate: " + (.success_rate | tostring) + "%"' "$results_json"
            else
                cat "$results_json"
            fi
        fi
        
        # Show log file location
        local log_file="$results_dir/test.log"
        if [[ -f "$log_file" ]]; then
            log_info "Detailed log: $log_file"
        fi
    else
        log_warning "No test results found in $results_dir"
    fi
    
    return $docker_exit_code
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Run Smartling CLI test suite in Docker container"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -c, --cli-binary PATH   Path to CLI binary (default: ./smartling-cli)"
    echo "  -f, --config-file PATH  Path to config file (default: ./smartling.yml)"
    echo "  -t, --test-suite PATH   Path to test suite (default: ./test-suite.sh)"
    echo ""
    echo "Environment variables:"
    echo "  CLI_BINARY    Override CLI binary path"
    echo "  CONFIG_FILE   Override config file path"
    echo "  TEST_SUITE    Override test suite path"
    echo ""
    echo "Example:"
    echo "  $0"
    echo "  $0 -c ./my-cli -f ./my-config.yml"
    echo "  CLI_BINARY=./smartling.linux $0"
    echo ""
    echo "Docker command equivalent:"
    echo "  docker run -it --rm \\"
    echo "    -v ./test-suite.sh:/test/test-suite.sh:ro \\"
    echo "    -v ./smartling.yml:/test/smartling.yml:ro \\"
    echo "    -v ./smartling-cli:/test/smartling-cli:ro \\"
    echo "    -w /test \\"
    echo "    ubuntu:20.04 bash test-suite.sh"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -c|--cli-binary)
                CLI_BINARY="$2"
                shift 2
                ;;
            -f|--config-file)
                CONFIG_FILE="$2"
                shift 2
                ;;
            -t|--test-suite)
                TEST_SUITE="$2"
                shift 2
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Main function
main() {
    echo "========================================"
    echo "  Smartling CLI Docker Test Runner"
    echo "========================================"
    echo ""
    
    parse_args "$@"
    check_prerequisites
    run_docker_tests
    
    local exit_code=$?
    echo ""
    echo "========================================"
    if [[ $exit_code -eq 0 ]]; then
        log_success "Test suite completed successfully!"
    else
        log_error "Test suite completed with failures"
    fi
    echo "========================================"
    
    exit $exit_code
}

# Run main function
main "$@"