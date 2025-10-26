# Gohan

**Omakase Hyprland Installer for Debian Sid/Trixie**

Gohan (ご飯) transforms Debian Sid into a polished, production-ready Hyprland environment with Omarchy-level integration.

## Features

- 🎯 **Opinionated Setup**: Zero-configuration transformation of Debian Sid
- 🔄 **Repo-First Strategy**: Uses official Debian packages when available
- 💾 **Comprehensive Backup**: Automatic rollback for safe updates
- 🎨 **Beautiful TUI**: Charmbracelet-powered installation wizard
- 🎭 **Theme System**: Pre-configured themes with one-command switching
- 🔍 **Smart Detection**: Automatic GPU and hardware configuration

## Requirements

- Debian Sid (unstable) or Trixie (testing)
- 10GB+ free disk space
- Internet connection
- Comfort with rolling development (Sid users)

## Installation

```bash
# Quick install (coming soon)
curl -fsSL https://gohan.sh | sh

# Or build from source
git clone https://github.com/rebelopsio/gohan.git
cd gohan
go build -o gohan ./cmd/gohan
sudo mv gohan /usr/local/bin/
```

## Usage

```bash
# Start installation wizard
gohan init

# Theme management
gohan theme list
gohan theme set catppuccin

# System maintenance
gohan update
gohan doctor
gohan rollback
```

## Development

See [CLAUDE.md](../CLAUDE.md) for development workflow and [overview.md](../overview.md) for technical architecture.

## Philosophy

Gohan targets experienced Linux users who want cutting-edge Hyprland without maintaining Arch. We embrace Debian Sid's rolling nature while providing comprehensive safety nets.

## License

MIT

## Credits

Built with ❤️ using [Charmbracelet](https://github.com/charmbracelet) and inspired by [Omarchy](https://github.com/prasanthrangan/hyprdots).
