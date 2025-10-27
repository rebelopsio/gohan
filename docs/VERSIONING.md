# Versioning Guide

Gohan follows [Semantic Versioning 2.0.0](https://semver.org/) with a 0.X.X convention during initial development.

## Version Format

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Breaking changes (stays at 0 until stable API)
- **MINOR**: New features, backwards compatible
- **PATCH**: Bug fixes, backwards compatible

## Development Phase (0.X.X)

During the initial development phase, versions start with `0.`:

```
0.1.0 - Initial release
0.2.0 - Added feature X
0.2.1 - Fixed bug in feature X
0.3.0 - Added feature Y
...
0.9.0 - Pre-stable release
1.0.0 - First stable release
```

### Why 0.X.X?

- **Indicates Pre-Production**: Version 0.X.X signals that the API is not yet stable
- **Allows Breaking Changes**: Minor version bumps (0.1 → 0.2) can include breaking changes
- **Community Expectations**: Users understand the software is still evolving

## Creating a Release

### 1. Determine Version Number

```bash
# First release
v0.1.0

# Added new features
v0.2.0

# Bug fix
v0.2.1

# Breaking change (during 0.x phase)
v0.3.0
```

### 2. Create and Push Tag

```bash
# Create annotated tag
git tag -a v0.1.0 -m "Release v0.1.0"

# Push tag to trigger release workflow
git push origin v0.1.0
```

### 3. GitHub Actions Automation

The tag push automatically:
1. Runs CI tests
2. Builds binaries for multiple platforms
3. Creates GitHub release
4. Uploads release artifacts

## Moving to 1.0.0

Version 1.0.0 signifies:

✅ **Stable API**: Public interfaces won't change without major version bump
✅ **Production Ready**: Safe for production use
✅ **Comprehensive Tests**: Full test coverage
✅ **Documentation**: Complete user and API documentation
✅ **Security**: Security audit completed

### Criteria for 1.0.0 Release

- [ ] API is stable and unlikely to change
- [ ] All critical features implemented
- [ ] Comprehensive test coverage (>80%)
- [ ] Complete documentation
- [ ] Production deployments successful
- [ ] No critical bugs
- [ ] Community feedback incorporated

## Version Bumping Guidelines

### Patch Version (0.X.Y)

Increment for:
- Bug fixes
- Documentation updates
- Internal refactoring
- Performance improvements (without API changes)

```bash
git tag -a v0.2.1 -m "Fix installation retry logic"
```

### Minor Version (0.Y.0)

Increment for:
- New features
- Deprecations (with backwards compatibility)
- API additions
- During 0.x: Breaking changes

```bash
git tag -a v0.3.0 -m "Add SQLite persistence support"
```

### Major Version (X.0.0)

Increment for:
- After 1.0.0: Breaking API changes
- Complete rewrites
- Incompatible changes

```bash
# First stable release
git tag -a v1.0.0 -m "First stable release"

# Future breaking change
git tag -a v2.0.0 -m "Complete API redesign"
```

## Release Checklist

Before tagging a release:

- [ ] All tests pass (`make test-all`)
- [ ] Linter passes (`make lint`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version numbers updated (if hardcoded)
- [ ] Migration guide (for breaking changes)

## Changelog Format

Keep a CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/):

```markdown
# Changelog

## [0.2.0] - 2024-01-15

### Added
- SQLite persistence support
- List installations API endpoint
- Cancel installation API endpoint

### Changed
- Improved error messages
- Updated dependencies

### Fixed
- Race condition in session repository
- Memory leak in package manager

## [0.1.0] - 2024-01-01

### Added
- Initial release
- Preflight validation
- Installation management
- HTTP API server
```

## Pre-release Versions

For alpha/beta/rc versions:

```bash
v0.2.0-alpha.1
v0.2.0-beta.1
v0.2.0-rc.1
```

Create with:

```bash
git tag -a v0.2.0-rc.1 -m "Release Candidate 1 for v0.2.0"
```

## Version in Code

The version is injected at build time via ldflags:

```bash
# Manual build with version
go build -ldflags "-X main.version=0.1.0" ./cmd/gohan

# Or use Makefile
VERSION=0.1.0 make build
```

## Checking Current Version

```bash
gohan version

# Output:
# gohan v0.1.0
#   commit: abc123
#   built:  2024-01-15T10:30:00Z
```

## Further Reading

- [Semantic Versioning 2.0.0](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
