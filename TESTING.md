# CodeQuest CLI Testing Guide

This guide walks you through testing the complete development and distribution flow.

## 🧪 Complete Testing Flow

### Phase 1: Local Development Testing

1. **Run local tests:**
   ```bash
   ./test-complete-flow.sh
   ```
   This tests:
   - Go module setup
   - Local build
   - Basic CLI functionality
   - Challenge fetch/test workflow

2. **Test cross-platform builds:**
   ```bash
   ./test-cross-platform-build.sh
   ```
   This tests:
   - Builds for Linux, macOS, Windows
   - Binary size verification
   - Checksum generation

### Phase 2: Repository Setup and First Release

1. **Initialize repository:**
   ```bash
   ./setup.sh
   ```

2. **Push to GitHub:**
   ```bash
   git push -u origin main
   ```

3. **Create first release:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

4. **Monitor GitHub Actions:**
   - Go to: https://github.com/crisecheverria/codequest/actions
   - Should see "Release" workflow running
   - Should complete successfully and create release

### Phase 3: Installation Testing

1. **Test installation methods:**
   ```bash
   ./test-installation.sh
   ```

2. **Test go install method:**
   ```bash
   # In a different directory
   go install github.com/crisecheverria/codequest@latest
   codequest --help
   ```

3. **Test binary download:**
   ```bash
   # Download for your platform
   curl -L https://github.com/crisecheverria/codequest/releases/latest/download/codequest-$(uname -s)-$(uname -m) -o codequest
   chmod +x codequest
   ./codequest --help
   ```

4. **Test complete workflow:**
   ```bash
   codequest list --language go
   codequest fetch go-multiple-return-values
   cd challenge-go-multiple-return-values
   # Edit solution.go to implement the function
   codequest test
   ```

### Phase 4: Update Flow Testing

1. **Test update process:**
   ```bash
   ./test-update-flow.sh
   ```

2. **Quick version bump:**
   ```bash
   ./bump-version.sh v1.0.1
   ```

## 🎯 Expected Results

### ✅ Successful Local Testing
- `./test-complete-flow.sh` shows all green checkmarks
- CLI builds without errors
- Basic commands work (list, fetch, test)

### ✅ Successful GitHub Actions
- Workflow completes without errors
- Release is created with 5 binary files + checksums
- Binaries are downloadable

### ✅ Successful Installation
- `go install` works and CLI is in PATH
- Direct binary download works
- CLI shows correct version with `--version`
- Complete workflow (list → fetch → test) works

### ✅ Successful Updates
- New tag triggers GitHub Actions
- New release is created
- Users can download updated version
- `go install @latest` gets the new version

## 🔍 Troubleshooting

### Build Failures
- Check Go version (needs 1.23+)
- Verify module name in go.mod
- Check import paths in .go files

### GitHub Actions Failures
- Verify .github/workflows/release.yml exists
- Check repository permissions
- Ensure tag format is `v*.*.*`

### Installation Issues
- Check GitHub release has all binaries
- Verify download URLs are correct
- Test with curl -I to check headers

### Runtime Issues
- Ensure Go 1.23+ installed (for Go challenges)
- Ensure Node.js 18+ installed (for JS/TS challenges)
- Check file permissions on binary

## 📊 Testing Checklist

- [ ] Local build works
- [ ] Cross-platform builds work
- [ ] GitHub repository created
- [ ] Initial code pushed
- [ ] First tag created (v1.0.0)
- [ ] GitHub Actions runs successfully
- [ ] Release created with binaries
- [ ] go install works
- [ ] Binary download works
- [ ] CLI shows correct version
- [ ] Complete workflow works
- [ ] Update process works
- [ ] New version released successfully

## 🚀 Ready for Production

Once all tests pass, your CLI is ready for:
- Public distribution
- User downloads
- Community contributions
- Regular updates via GitHub Actions

Share the installation instructions from the README.md with your users!