#!/bin/bash

# Install Git hooks for pre-commit checks
# This script can be used as an alternative to the pre-commit framework

set -e

HOOKS_DIR=".git/hooks"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Installing Git hooks..."

# Ensure hooks directory exists
mkdir -p "$HOOKS_DIR"

# Create pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash

# Pre-commit hook for Gohan
# Runs format, lint, and tests before allowing commit

set -e

echo "Running pre-commit checks..."

# Format code
echo "1. Formatting code..."
make fmt

# Run linter
echo "2. Running linter..."
if ! make lint; then
    echo "❌ Linting failed. Please fix the issues before committing."
    exit 1
fi

# Run unit tests
echo "3. Running unit tests..."
if ! make test-unit; then
    echo "❌ Unit tests failed. Please fix the issues before committing."
    exit 1
fi

echo "✅ Pre-commit checks passed!"
EOF

# Create pre-push hook
cat > "$HOOKS_DIR/pre-push" << 'EOF'
#!/bin/bash

# Pre-push hook for Gohan
# Runs comprehensive tests before allowing push

set -e

echo "Running pre-push checks..."

# Run race detector
echo "1. Running race detector..."
if ! make test-race; then
    echo "❌ Race detector found issues. Please fix before pushing."
    exit 1
fi

# Run all tests
echo "2. Running all tests..."
if ! make test-all; then
    echo "❌ Tests failed. Please fix before pushing."
    exit 1
fi

echo "✅ Pre-push checks passed!"
EOF

# Make hooks executable
chmod +x "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-push"

echo "✅ Git hooks installed successfully!"
echo ""
echo "The following hooks are now active:"
echo "  - pre-commit: runs fmt, lint, and unit tests"
echo "  - pre-push: runs race detector and all tests"
echo ""
echo "To skip hooks temporarily, use:"
echo "  git commit --no-verify"
echo "  git push --no-verify"
