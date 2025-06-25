# Smartling CLI Test Suite

Comprehensive smoke testing suite for Smartling CLI post-release validation.

## Features

- **Structured Test Categories**: Happy path, error handling, edge cases, and snapshot testing
- **Complete Workflow Testing**: Upload → Import → List → Download → Rename → Delete
- **Docker Environment**: Isolated testing environment using Ubuntu container
- **Comprehensive Logging**: Colored output with detailed test logs
- **Snapshot Testing**: Validates CLI output consistency across releases
- **JSON Results**: Machine-readable test results for CI integration

## Quick Start

### Prerequisites

1. Smartling CLI binary (`smartling-cli` or `smartling.linux`)
2. Valid `smartling.yml` configuration file
3. Docker installed (for containerized testing)

### Directory Structure

```
project/
├── smartling-cli          # CLI binary
├── smartling.yml          # Configuration file
├── test-suite.sh          # Main test suite
├── run-docker-tests.sh    # Docker test runner
└── test-results-*/        # Generated test results
```

### Running Tests

#### Option 1: Direct Execution (Linux)
```bash
# Make scripts executable
chmod +x test-suite.sh smartling-cli

# Run tests directly
./test-suite.sh
```

#### Option 2: Docker Environment (Recommended)
```bash
# Using the Docker runner script
chmod +x run-docker-tests.sh
./run-docker-tests.sh

# Or directly with Docker
docker run -it --rm \
  -v ./test-suite.sh:/test/test-suite.sh:ro \
  -v ./smartling.yml:/test/smartling.yml:ro \
  -v ./smartling-cli:/test/smartling-cli:ro \
  -w /test \
  ubuntu:20.04 bash test-suite.sh
```

## Test Categories

### 1. Happy Path Tests
- Authentication and configuration validation
- Project operations (list, info, locales)
- Basic file operations workflow

### 2. File Operations Workflow
Complete end-to-end testing:
1. **Generate test file** (`Hello world!`)
2. **Upload test file** to two different URIs
3. **Import translations** (German and French)
4. **List files** using wildcards
5. **Download single file** to current directory
6. **Rename file** with new URI
7. **Download all files** to subdirectory
8. **Delete uploaded files** for cleanup

### 3. Error Handling Tests
- Invalid project ID handling
- Missing file scenarios
- Non-existent file URI downloads
- Expected error message validation

### 4. Edge Case Tests
- Empty file uploads
- Special characters in filenames
- Boundary conditions

### 5. Snapshot Tests
- CLI help output consistency
- Command formatting validation
- Timestamp normalization for comparisons

## Configuration

### Environment Variables
```bash
export CLI_BINARY="./smartling.linux"    # CLI binary path
export CONFIG_FILE="./smartling.yml"     # Config file path
export TEST_SUITE="./test-suite.sh"      # Test suite path
```

### Docker Runner Options
```bash
./run-docker-tests.sh -h                 # Show help
./run-docker-tests.sh -c ./my-cli         # Custom CLI binary
./run-docker-tests.sh -f ./my-config.yml  # Custom config file
```

## Output and Results

### Test Execution
- **Colored output**: Green (pass), Red (fail), Blue (info), Yellow (warning)
- **Progress tracking**: Real-time test execution status
- **Detailed logging**: All CLI commands and outputs logged

### Result Files
```
test-results-YYYYMMDD-HHMMSS/
├── test.log              # Detailed execution log
├── results.json          # Machine-readable results
├── snapshots/            # Snapshot comparison files
│   ├── help_main.snapshot
│   ├── projects_list_short.snapshot
│   └── ...
└── test_files/           # Generated test files
    ├── test.txt
    ├── test_de-DE.txt
    └── test_fr-FR.txt
```

### Results JSON Format
```json
{
    "timestamp": "2024-01-15T10:30:45+00:00",
    "total_tests": 25,
    "passed": 23,
    "failed": 2,
    "success_rate": 92
}
```

## Integration with CI/CD

### Exit Codes
- `0`: All tests passed
- `1`: One or more tests failed

### Example CI Integration
```bash
#!/bin/bash
# CI script example

# Run tests
./run-docker-tests.sh

# Check results
if [ $? -eq 0 ]; then
    echo "✅ All smoke tests passed - Release is ready"
    exit 0
else
    echo "❌ Smoke tests failed - Release blocked"
    exit 1
fi
```

## Customization

### Adding New Tests
1. Add test function to appropriate category
2. Call function from `main()` 
3. Follow naming convention: `test_category_name()`

### Modifying Test Files
Edit `generate_test_files()` function to create custom test content.

### Snapshot Updates
When CLI output format changes legitimately:
1. Delete old snapshot files
2. Re-run tests to generate new snapshots
3. Review and commit new snapshots

## Troubleshooting

### Common Issues

**"CLI binary not found"**
- Ensure binary exists and is executable
- Check `CLI_BINARY` environment variable
- Verify file permissions

**"Config file not found"**
- Ensure `smartling.yml` exists in correct location
- Check `CONFIG_FILE` environment variable
- Validate YAML syntax

**"Docker not found"**
- Install Docker
- Ensure Docker daemon is running
- Check user permissions for Docker

**"Tests failing in Docker but not locally"**
- Check file permissions (executable bits)
- Verify volume mounts
- Review Docker container logs

### Debug Mode
```bash
# Enable verbose logging
set -x
./test-suite.sh

# Check specific test logs
tail -f test-results-*/test.log
```

## Best Practices

1. **Run tests in clean environment** (Docker recommended)
2. **Review snapshot changes** before accepting them
3. **Clean up test files** after each run (automatic)
4. **Monitor test execution time** for performance regressions
5. **Archive test results** for historical analysis

## Contributing

When adding new test scenarios:
1. Follow existing test structure
2. Include both positive and negative test cases
3. Add appropriate cleanup procedures
4. Update documentation
5. Test with Docker environment