# Gohan Wiki Documentation

This directory contains the documentation for the Gohan GitHub wiki.

## Files Created

- **Home.md** - Project overview and quick start
- **Installation.md** - Complete installation guide
- **Theme-Management.md** - Comprehensive theme management guide
- **Troubleshooting.md** - Common issues and solutions
- **Development.md** - Contributing and development guide

## Pushing to GitHub Wiki

GitHub wikis are actually Git repositories. Here's how to push these files:

### Method 1: Clone and Push (Recommended)

```bash
# Clone the wiki repository
git clone https://github.com/rebelopsio/gohan.wiki.git

# Navigate to wiki directory
cd gohan.wiki

# Copy wiki files from main repository
cp ../gohan/wiki/*.md .

# Add and commit
git add *.md
git commit -m "Add comprehensive documentation

- Home page with project overview
- Installation guide with multiple methods
- Theme management complete guide
- Troubleshooting common issues
- Development guide for contributors"

# Push to wiki
git push origin master
```

### Method 2: Manual Upload

1. Go to https://github.com/rebelopsio/gohan/wiki
2. Click "New Page" for each file
3. Copy/paste the markdown content
4. Use the filename (without .md) as the page title
5. Save each page

## File Descriptions

### Home.md
- Project introduction
- Quick start guide
- Feature overview
- Architecture summary
- Links to other documentation

**Size**: ~4.5KB

### Installation.md
- Prerequisites
- Multiple installation methods (source, binary, go install)
- Post-installation setup
- Directory structure
- Upgrading and uninstallation
- Troubleshooting installation issues

**Size**: ~5.3KB

### Theme-Management.md
- Complete theme guide
- All 5 built-in themes
- Applying and previewing themes
- Theme history and rollback
- Customization and templates
- Component support details

**Size**: ~11KB

### Troubleshooting.md
- Installation issues
- Theme application problems
- Configuration issues
- Component reload problems
- General troubleshooting
- Getting help resources

**Size**: ~9.6KB

### Development.md
- Development setup
- Architecture overview
- BDD → ATDD → TDD workflow
- Testing guidelines
- Code standards
- Contributing guide
- Release process

**Size**: ~11KB

## Total Documentation

- **5 comprehensive pages**
- **~41KB of documentation**
- **Covers all major aspects of Gohan**

## Maintenance

When updating documentation:

1. Edit the markdown files in `wiki/`
2. Test locally (can use any markdown viewer)
3. Push to wiki repository
4. Verify on GitHub wiki

## Additional Pages to Consider

Future documentation that could be added:

- **Quick-Start.md** - 5-minute getting started guide
- **CLI-Reference.md** - Complete command reference
- **Configuration.md** - Template customization guide
- **Package-Installation.md** - Hyprland installation details
- **API-Reference.md** - Internal API documentation
- **FAQ.md** - Frequently asked questions

## Notes

- GitHub wiki pages are created automatically when you push markdown files
- Page names come from filenames (without .md extension)
- Spaces in filenames become hyphens in URLs
- Links between wiki pages use format: `[Link Text](Page-Name)`
- Images can be uploaded to the wiki and referenced

## Verification

After pushing, verify at:
- https://github.com/rebelopsio/gohan/wiki

Each page should be accessible and formatted correctly.
