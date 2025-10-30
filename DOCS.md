# Documentation Guide

This guide explains how to work with Gohan's documentation, which is built using [MkDocs](https://www.mkdocs.org/) with the [Material theme](https://squidfunk.github.io/mkdocs-material/).

## 📚 Documentation Structure

```
gohan/
├── docs/                      # Documentation source files
│   ├── index.md              # Homepage
│   ├── installation.md       # Installation guide
│   ├── theme-management.md   # Theme management guide
│   ├── troubleshooting.md    # Troubleshooting guide
│   └── development.md        # Development guide
├── mkdocs.yml                # MkDocs configuration
├── .github/workflows/docs.yml # Auto-deployment workflow
└── site/                     # Generated static site (git-ignored)
```

## 🚀 Quick Start

### Prerequisites

- Python 3.8 or higher
- Make (for using Makefile targets)

### Local Development

```bash
# Build documentation (creates .venv and installs dependencies)
make docs-build

# Serve documentation locally with live reload
make docs-serve

# Open http://127.0.0.1:8000 in your browser
```

The `docs-serve` command will:
- Create a Python virtual environment (`.venv/`)
- Install MkDocs and all required plugins
- Start a local development server
- Watch for changes and rebuild automatically

## 🛠️ Available Commands

### Using Make (Recommended)

```bash
# Create Python virtual environment only
make docs-venv

# Install MkDocs dependencies
make docs-install

# Build documentation with strict mode (fail on warnings)
make docs-build

# Serve documentation locally at http://127.0.0.1:8000
make docs-serve

# Deploy to GitHub Pages (requires push access)
make docs-deploy

# Clean documentation artifacts and venv
make docs-clean
```

### Using MkDocs Directly

If you prefer to use MkDocs directly:

```bash
# Activate virtual environment
source .venv/bin/activate

# Build documentation
mkdocs build --strict

# Serve documentation
mkdocs serve

# Deploy to GitHub Pages
mkdocs gh-deploy --force

# Deactivate virtual environment
deactivate
```

## 📝 Writing Documentation

### Markdown Basics

Documentation files use Markdown with additional features from the Material theme:

```markdown
# Page Title

## Section Heading

### Subsection

Regular paragraph with **bold** and *italic* text.

- Bullet list item 1
- Bullet list item 2

1. Numbered list item 1
2. Numbered list item 2

[Link text](https://example.com)
[Internal link](other-page.md)

`inline code`

\```bash
# Code block with syntax highlighting
make docs-build
\```
```

### Material Theme Features

#### Admonitions (Callout Boxes)

```markdown
!!! note "Optional Title"
    This is a note admonition.

!!! tip
    This is a tip without a custom title.

!!! warning "Important"
    This is a warning.

!!! danger
    This is a danger alert.

!!! success
    This is a success message.

!!! info
    This is an info box.

!!! failure
    This is a failure/error message.
```

#### Collapsible Admonitions

```markdown
??? note "Click to expand"
    This content is collapsed by default.

???+ tip "Expanded by default"
    This content is expanded by default but can be collapsed.
```

#### Tabbed Content

```markdown
=== "Tab 1"

    Content for tab 1

=== "Tab 2"

    Content for tab 2

=== "Tab 3"

    Content for tab 3
```

#### Code Blocks with Titles

```markdown
\```bash title="Install Gohan"
git clone https://github.com/rebelopsio/gohan.git
cd gohan
go build -o gohan ./cmd/gohan
\```
```

#### Task Lists

```markdown
- [x] Completed task
- [ ] Uncompleted task
- [ ] Another task
```

## 🎨 Customization

### Theme Configuration

The theme is configured in `mkdocs.yml`:

```yaml
theme:
  name: material
  palette:
    # Light/dark mode with automatic switching
    - media: "(prefers-color-scheme: light)"
      scheme: default
      primary: deep purple
      accent: purple
    - media: "(prefers-color-scheme: dark)"
      scheme: slate
      primary: deep purple
      accent: purple

  features:
    - navigation.instant      # Instant loading
    - navigation.tracking     # URL tracking
    - navigation.tabs         # Top-level sections as tabs
    - navigation.sections     # Sections in sidebar
    - search.suggest          # Search suggestions
    - content.code.copy      # Copy button for code blocks
```

### Navigation Structure

Navigation is defined in `mkdocs.yml`:

```yaml
nav:
  - Home: index.md
  - Getting Started:
      - Installation: installation.md
  - User Guide:
      - Theme Management: theme-management.md
      - Troubleshooting: troubleshooting.md
  - Development:
      - Contributing: development.md
```

## 🚀 Deployment

### Automatic Deployment (GitHub Actions)

Documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch:

1. Make changes to documentation in `docs/` directory
2. Commit and push to `main` branch:
   ```bash
   git add docs/ mkdocs.yml
   git commit -m "docs: update documentation"
   git push origin main
   ```
3. GitHub Actions workflow (`.github/workflows/docs.yml`) will:
   - Install Python and MkDocs
   - Build documentation with `--strict` mode
   - Deploy to `gh-pages` branch
   - Make it available at https://rebelopsio.github.io/gohan

### Manual Deployment

If you need to deploy manually:

```bash
# Build and deploy to GitHub Pages
make docs-deploy

# Or using mkdocs directly
source .venv/bin/activate
mkdocs gh-deploy --force
```

### First-Time Setup

To enable GitHub Pages for the first time:

1. Go to your repository on GitHub
2. Navigate to Settings → Pages
3. Under "Source", select:
   - Branch: `gh-pages`
   - Folder: `/ (root)`
4. Click "Save"
5. Your documentation will be available at: https://rebelopsio.github.io/gohan

## 🔍 Testing

### Local Testing

Always test your documentation locally before pushing:

```bash
# Build with strict mode (fails on warnings/errors)
make docs-build

# Serve locally and check in browser
make docs-serve
```

### Common Issues

#### Broken Links

The build will fail in strict mode if there are broken links:

```
WARNING - Doc file 'index.md' contains a link 'non-existent.md', but the target is not found
```

**Fix**: Update the link to point to an existing file or remove it.

#### Wrong Anchor Links

```
WARNING - Doc file contains a link '#wrong-anchor', but there is no such anchor
```

**Fix**: Ensure anchor links match the heading exactly (lowercase with hyphens):
- Heading: `## Theme Management`
- Link: `#theme-management`

#### Case Sensitivity

Use lowercase filenames consistently:
- ✅ `installation.md`
- ❌ `Installation.md`

## 📦 Dependencies

The documentation uses these Python packages (installed automatically):

- `mkdocs` - Static site generator
- `mkdocs-material` - Material theme for MkDocs
- `mkdocs-minify-plugin` - Minify HTML/JS/CSS for faster loading

Dependencies are specified in the Makefile and installed in `.venv/`.

## 🔄 Workflow

Typical documentation workflow:

1. **Create or edit** documentation files in `docs/`
2. **Preview locally** with `make docs-serve`
3. **Test build** with `make docs-build`
4. **Commit changes** to git
5. **Push to main** - GitHub Actions deploys automatically
6. **Verify deployment** at https://rebelopsio.github.io/gohan

## 📚 Resources

- [MkDocs Documentation](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
- [Markdown Guide](https://www.markdownguide.org/)
- [Python Markdown Extensions](https://python-markdown.github.io/extensions/)

## 🤝 Contributing

When contributing documentation:

1. Follow the existing structure and style
2. Use Material theme features (admonitions, tabs, etc.)
3. Test locally before pushing
4. Ensure `make docs-build` passes
5. Keep language clear and concise
6. Include code examples where helpful
7. Add screenshots when illustrating UI

## 💡 Tips

- Use admonitions to highlight important information
- Break long pages into sections with clear headings
- Use tabbed content for alternatives (different OS, methods, etc.)
- Add copy buttons to code blocks with the `title` attribute
- Test all links and ensure they work
- Keep navigation structure logical and shallow
- Use search-friendly headings and keywords

## 🐛 Troubleshooting

### Virtual Environment Issues

If you encounter issues with the virtual environment:

```bash
# Clean everything and rebuild
make docs-clean
make docs-build
```

### Build Fails in CI

Check the GitHub Actions logs:
1. Go to your repository on GitHub
2. Click "Actions" tab
3. Click on the failing workflow run
4. Review the "Build documentation" step

### Documentation Not Updating

If your changes don't appear:
1. Clear browser cache
2. Wait a few minutes for GitHub Pages to update
3. Check GitHub Actions completed successfully
4. Verify changes were pushed to `main` branch

## 📝 Checklist

Before pushing documentation changes:

- [ ] Tested locally with `make docs-serve`
- [ ] Build passes with `make docs-build`
- [ ] All links work correctly
- [ ] Code examples are tested
- [ ] Admonitions used for important notes
- [ ] Navigation updated if new pages added
- [ ] Spelling and grammar checked
- [ ] Consistent formatting and style
- [ ] Screenshots added/updated if needed
- [ ] Changes committed with descriptive message
