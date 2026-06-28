#!/bin/bash
# Security scanning script for zara-jira-mcp
# Run in CI/CD pipeline or before production deployment

set -e

echo "🔒 zara-jira-mcp Security Scan"
echo "=============================="

# 1. Check for hardcoded secrets in source code files (but not in env files)
echo "1. Scanning for hardcoded secrets in Go source files..."

SECRETS_FOUND=0

# Check Go source files for potential hardcoded secrets (but skip .env files)
if grep -rnE "(API_TOKEN=|PASSWORD=|SECRET=|KEY=|TOKEN=)[A-Za-z0-9_\.]{20,}" --include="*.go" . | grep -v test | grep -v example; then
    echo "   ❌ Found potential hardcoded secrets in Go source code"
    SECRETS_FOUND=$((SECRETS_FOUND + 1))
else
    echo "   ✅ No hardcoded secrets found in Go source code"
fi

# Check for common weak patterns
if grep -rn "md5\|sha1\|des3\|rc4" --include="*.go" . | grep -v test | grep -v example; then
    echo "   ⚠️  Found potential weak cryptographic functions"
else
    echo "   ✅ No weak crypto found"
fi

# 2. Validate Go module security
echo ""
echo "2. Checking Go module security..."

if [ -f "go.sum" ]; then
    echo "   ✅ go.sum file present"
else
    echo "   ❌ go.sum file missing"
fi

# 3. Run vulnerability scanner if available
if command -v govulncheck &> /dev/null; then
    echo ""
    echo "3. Running vulnerability scanner (govulncheck)..."
    govulncheck ./...
else
    echo ""
    echo "3. ℹ️  govulncheck not installed (run: go install golang.org/x/vuln/cmd/govulncheck@latest)"
fi

# 4. Summary
echo ""
echo "📋 Security Scan Summary"
echo "========================"

if [ $SECRETS_FOUND -gt 0 ]; then
    echo "❌ CRITICAL: Found $SECRETS_FOUND secret(s) in source code"
    echo "   Remediation: Move secrets to environment files or secret management"
    exit 1
else
    echo "✅ No critical secret exposure found in source code"
fi

# 5. Provide recommendations
echo ""
echo "📋 Recommendations:"
echo "   1. Store secrets in environment variables (in .env file, gitignored)"
echo "   2. Use secret management (AWS Secrets Manager, HashiCorp Vault) for production"
echo "   3. Run govulncheck to check for CVEs"
echo "   4. Enable Dependabot for automated updates"
echo "   5. Run this scan in CI/CD pipeline"

echo ""
echo "✅ Security scan completed"
