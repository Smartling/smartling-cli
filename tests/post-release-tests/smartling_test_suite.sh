#!/bin/bash

# Smartling CLI Comprehensive Test Suite
# Post-release smoke testing for user-facing functionality

set -euo pipefail

# Configuration
CLI_BINARY="${CLI_BINARY:-./smartling-cli}"
CONFIG_FILE="${CONFIG_FILE:-smartling.yml}"
TEST_FILES_DIR="./test_files"
SNAPSHOT_DIR="./snapshots"
TEST_DIR="$(mktemp -d)"
LOG_FILE="${TEST_DIR}/test.log"
RESULTS_FILE="${TEST_DIR}/results.json"

# Test counters
TESTS_TOTAL=0
TESTS_PASSED=0
TESTS_FAILED=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

log_info() {
    echo -e "\n${BLUE}[INFO]${NC} $*" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $*" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $*" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $*" | tee -a "$LOG_FILE"
}

# Test framework functions
test_start() {
    local test_name="$1"
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    log_info "Starting test: $test_name"
}

test_pass() {
    local test_name="$1"
    TESTS_PASSED=$((TESTS_PASSED + 1))
    log_success "Test passed: $test_name"
}

test_fail() {
    local test_name="$1"
    local reason="$2"
    TESTS_FAILED=$((TESTS_FAILED + 1))
    log_error "Test failed: $test_name - $reason"
}

# CLI wrapper function
run_cli() {
    local cmd="$*"
    log "Executing: $CLI_BINARY $cmd"
    if ! $CLI_BINARY $cmd 2>&1 | tee -a "$LOG_FILE"; then
        return 1
    fi
    return 0
}

# CLI wrapper for expected failures
run_cli_expect_fail() {
    local cmd="$*"
    local expected_pattern="$1"
    shift
    cmd="$*"
    
    log "Executing (expecting failure): $CLI_BINARY $cmd"
    if $CLI_BINARY $cmd 2>&1 | tee -a "$LOG_FILE" | grep -q "$expected_pattern"; then
        return 0  # Expected failure occurred
    else
        return 1  # Unexpected success or wrong error
    fi
}

# Snapshot testing function
snapshot_test() {
    local test_name="$1"
    local cmd="$2"
    local snapshot_file="${SNAPSHOT_DIR}/${test_name}.snapshot"
    
    # Remove timestamps and dynamic content for comparison
    local output
    output=$($CLI_BINARY $cmd 2>&1 | sed 's/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z/<TIMESTAMP>/g')
    
    if [[ -f "$snapshot_file" ]]; then
        if echo "$output" | diff -u "$snapshot_file" - > /dev/null; then
            test_pass "$test_name (snapshot)"
        else
            test_fail "$test_name (snapshot)" "Output differs from snapshot"
            echo "$output" > "${snapshot_file}.new"
            log_warning "New output saved to ${snapshot_file}.new"
        fi
    else
        echo "$output" > "$snapshot_file"
        test_pass "$test_name (snapshot created)"
        log_info "Created new snapshot: $snapshot_file"
    fi
}

# Test file generation
generate_test_files() {
    log_info "Generating test files"
    mkdir -p "$TEST_FILES_DIR"
    
    # Simple text file for main testing
    cat > "${TEST_FILES_DIR}/test.txt" << 'EOF'
Hello world!
EOF

    # Simple properties file for additional testing
    cat > "${TEST_FILES_DIR}/test.properties" << 'EOF'
app.title=Test Application
app.description=This is a test application
app.version=1.0.0
EOF

    # German translation file for import testing
    cat > "${TEST_FILES_DIR}/test_de.properties" << 'EOF'
app.title=Test Application DE
app.description=This is a test application DE
app.version=1.0.0
EOF

    # French translation file for import testing
    cat > "${TEST_FILES_DIR}/test_fr-FR.properties" << 'EOF'
app.title=Test Application fr-FR
app.description=This is a test application fr-FR
app.version=1.0.0
EOF
}

# Test Categories

# 1. Happy Path Tests
test_authentication() {
    test_start "Authentication & Config"
    
    if [[ ! -f "$CONFIG_FILE" ]]; then
        test_fail "Authentication" "Config file $CONFIG_FILE not found"
        return 1
    fi
    
    if run_cli "projects list --short"; then
        test_pass "Authentication"
    else
        test_fail "Authentication" "Failed to list projects"
    fi
}

test_project_operations() {
    test_start "Project Operations"
    
    # Test project listing
    if run_cli "projects list"; then
        test_pass "Project list"
    else
        test_fail "Project list" "Command failed"
        return 1
    fi
    
    # Test project info
    if run_cli "projects info"; then
        test_pass "Project info"
    else
        test_fail "Project info" "Command failed"
    fi
    
    # Test locales listing
    if run_cli "projects locales"; then
        test_pass "Project locales"
    else
        test_fail "Project locales" "Command failed"
    fi
    
    # Test short format
    if run_cli "projects locales --short"; then
        test_pass "Project locales (short)"
    else
        test_fail "Project locales (short)" "Command failed"
    fi
}

test_file_operations_workflow() {
    test_start "File Operations Workflow"
    
    local test_file="${TEST_FILES_DIR}/test.txt"
    local test_prop_file="${TEST_FILES_DIR}/test.properties"
    local test_prop_translation_file_de="${TEST_FILES_DIR}/test_de.properties"
    local test_prop_translation_file_fr="${TEST_FILES_DIR}/test_fr-FR.properties"
    local file_uri1="/test-workflow/test.txt"
    local file_uri2="/test-workflow/test-copy.txt"
    local file_uri_renamed="/test-workflow/test-renamed.txt"
    local download_dir="${TEST_DIR}/test1"
    
    mkdir -p "$download_dir"
    
    # 1. Upload test file
    log_info "Step 1: Upload test file"
    if run_cli "files push $test_file $file_uri1"; then
        test_pass "File upload (original)"
    else
        test_fail "File upload (original)" "Upload failed"
        return 1
    fi
    
    # 2. Upload same file with different URI
    log_info "Step 2: Upload same file with different URI"
    # Directive is required to avoid issue with namespace when we Rename file (see test below)
    if run_cli "files push $test_file $file_uri2 --directive 'file_uri_as_namespace=false'"; then
        test_pass "File upload (copy)"
    else
        test_fail "File upload (copy)" "Upload failed"
    fi
    
    # 3. Import translations
    log_info "Step 3: Import translations"
    if run_cli "files push $test_prop_file test.properties"; then
        test_pass "File upload (original)"
    else
        test_fail "File upload (original)" "Upload failed"
        return 1
    fi

    if run_cli "files import test.properties $test_prop_translation_file_de de --published"; then
        test_pass "Import translations (de)"
    else
        test_fail "Import translations (de)" "Import failed"
    fi
    
    if run_cli "files import test.properties $test_prop_translation_file_fr fr-FR --published"; then
        test_pass "Import translations (fr-FR)"
    else
        test_fail "Import translations (fr-FR)" "Import failed"
    fi
    
    # 4. List files
    log_info "Step 4: List files"
    # todo : It doesn't work with '**/test*.txt'. Server return that no files in TMS
    if run_cli "files list"; then
        test_pass "File listing"
    else
        test_fail "File listing" "List failed"
    fi
    
    # 5. Download one file to current folder
    log_info "Step 5: Download file to current folder"
    if run_cli "files pull $file_uri1 --source"; then
        test_pass "File download (single)"
    else
        test_fail "File download (single)" "Download failed"
    fi
    
    # 6. Rename file
    log_info "Step 6: Rename file"
    if run_cli "files rename $file_uri2 $file_uri_renamed"; then
        test_pass "File rename"
    else
        test_fail "File rename" "Rename failed"
    fi
    
    # 7. Download all files to subfolder
    log_info "Step 7: Download all files to subfolder"
    if run_cli "files pull '**/test*.txt' --source -d $download_dir"; then
        test_pass "File download (bulk)"
    else
        test_fail "File download (bulk)" "Bulk download failed"
    fi
    
    # 8. Delete uploaded files
    log_info "Step 8: Delete uploaded files"
    if run_cli "files delete $file_uri1" && run_cli "files delete $file_uri_renamed"; then
        test_pass "File deletion"
    else
        test_fail "File deletion" "Deletion failed"
    fi
}

# 2. Error Handling Tests
test_error_handling() {
    test_start "Error Handling"
    
    # Test invalid project ID
    if run_cli_expect_fail "project" "projects list -p invalid-project-id"; then
        test_pass "Invalid project ID error"
    else
        test_fail "Invalid project ID error" "Expected error not found"
    fi
    
    # Test upload for missing file
    if run_cli_expect_fail "no files found by specified patterns" "files push non-existent-file.txt"; then
        test_pass "Upload missing file error"
    else
        test_fail "Upload missing file error" "Expected error not found"
    fi
    
    # Test download for missing file in TMS
    if run_cli_expect_fail "no files found on the remote server matching provided pattern" "files pull /non-existent-file.txt --source"; then
        test_pass "Download missing file from TMS error"
    else
        test_fail "Download missing file from TMS error" "Expected error not found"
    fi
}

# 3. Edge Case Tests
test_edge_cases() {
    test_start "Edge Cases"

    # Test empty file upload
    local empty_file="${TEST_FILES_DIR}/empty.txt"
    touch "$empty_file"

    if run_cli_expect_fail "failed to upload 1 files" "files push $empty_file /test-edge/empty.txt"; then
        test_pass "Empty file upload"
    else
        test_fail "Empty file upload" "Upload failed"
    fi

    # Test file with special characters in name
    local special_file="${TEST_FILES_DIR}/test file with spaces.txt"
    echo "test content" > "$special_file"
    
    if run_cli "files push \"$special_file\" \"/test-edge/special file.txt\""; then
        test_pass "Special characters in filename"
        run_cli "files delete '/test-edge/special file.txt'" || true
    else
        test_fail "Special characters in filename" "Upload failed"
    fi
}

# 4. Snapshot Tests
test_snapshots() {
    test_start "Snapshot Tests"
    
    mkdir -p "$SNAPSHOT_DIR"
    
    # Test help output
    snapshot_test "help_main" "--help"
    snapshot_test "help_files" "files --help"
    snapshot_test "help_projects" "projects --help"
    
    # Test list formatting
    snapshot_test "projects_list_short" "projects list --short"
    snapshot_test "locales_source" "projects locales --source"
}

# Cleanup function
cleanup() {
    log_info "Cleaning up test environment"
    
    # Clean up any remaining test files in Smartling
    log_info "Cleaning up remote test files"
#    $CLI_BINARY files list '**/test-*' --short 2>/dev/null | while read -r file_uri; do
#        if [[ -n "$file_uri" ]]; then
#            log "Deleting: $file_uri"
#            $CLI_BINARY files delete "$file_uri" 2>/dev/null || true
#        fi
#    done
    $CLI_BINARY files list '**/test-*' --short | $CLI_BINARY files delete -
    $CLI_BINARY files delete "test.properties" 2>/dev/null || true

    # Clean up local test directory
    if [[ -d "$TEST_DIR" ]]; then
        rm -rf "$TEST_DIR"
    fi
}

# Results reporting
generate_report() {
    log_info "Generating test report"
    
    cat > "$RESULTS_FILE" << EOF
{
    "timestamp": "$(date -Iseconds)",
    "total_tests": $TESTS_TOTAL,
    "passed": $TESTS_PASSED,
    "failed": $TESTS_FAILED,
    "success_rate": $(( TESTS_PASSED * 100 / TESTS_TOTAL ))
}
EOF

    echo ""
    echo "=================================="
    echo "         TEST RESULTS"
    echo "=================================="
    echo "Total Tests:    $TESTS_TOTAL"
    echo "Passed:         $TESTS_PASSED"
    echo "Failed:         $TESTS_FAILED"
    echo "Success Rate:   $(( TESTS_PASSED * 100 / TESTS_TOTAL ))%"
    echo ""
    echo "Log file:       $LOG_FILE"
    echo "Results file:   $RESULTS_FILE"
    echo "Test directory: $TEST_DIR"
    echo "=================================="
    
    if [[ $TESTS_FAILED -gt 0 ]]; then
        exit 1
    fi
}

# Main execution
main() {
    log_info "Starting Smartling CLI Test Suite"
    log_info "CLI Binary: $CLI_BINARY"
    log_info "Config File: $CONFIG_FILE"
    log_info "Test Directory: $TEST_DIR"
    
    # Setup
    generate_test_files
    
    # Set trap for cleanup
    trap cleanup EXIT
    
    # Run test categories
    test_authentication
    test_project_operations
    test_file_operations_workflow
    test_error_handling
    test_edge_cases
    test_snapshots
    
    # Generate report
    generate_report
}

# Check prerequisites
if [[ ! -x "$CLI_BINARY" ]]; then
    echo "Error: CLI binary not found or not executable: $CLI_BINARY"
    echo "Set CLI_BINARY environment variable or place binary in current directory"
    exit 1
fi

if [[ ! -f "$CONFIG_FILE" ]]; then
    echo "Error: Config file not found: $CONFIG_FILE"
    echo "Set CONFIG_FILE environment variable or create smartling.yml in current directory"
    exit 1
fi

# Run main function
main "$@"