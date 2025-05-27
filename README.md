# i3-bind

A modern CLI/TUI utility for managing i3 window manager keybindings with ease and style

[AUR](https://aur.archlinux.org/packages/i3-bind)

## Features:
 - 🎯 Add/Remove keybindings with simple commands
 - 🔍 Search keybindings by key or action
 - 💬 Comment management for better organization
 - 🎨 Syntax highlighting with colorized output
 - 🖥️ Interactive TUI mode powered by fzf
 - 🔒 Automatic backups before modifications
 - 📋 List all keybindings in a readable format

## Installation

### Arch Linux (AUR)

```bash
# Using yay
yay -S i3-bind

# Using paru
paru -S i3-bind

# Manual Installation
git clone https://aur.archlinux.org/i3-bind.git
cd i3-bind
makepkg -simple
```

### From Source

```bash
git clone https://github.com/Hanashiko/i3-bind.git
cd i3-bind
go build -o i3-bind
sudo mv i3-bind /usr/bin/
```

### Prerequisites

 - Go 1.19+ (for building from source)
 - `fzf` (required for interactive mode)

## Usage

### Basic commands

#### Add a keybinding
```bash
i3-bind add mod4+Enter "exec alacritty"
i3-bind add mod4+d "exec dmenu_run"
i3-bind add '$mod+shift+q' kill
```

#### Remove a keybinding
```bash
i3-bind remove mod4+q
i3-bind remove mod4+Enter
i3-bind remove '$mod+shift+print'
```

#### List all keybindings
```bash
i3-bind list
```

#### search keybindings
```bash
i3-bind find firefox # find by action
i3-bind find mod4+shift # find by key pattern
i3-bind find terminal # find in comments
```

#### Add/Update comments
```bash
i3-bind comment mod4+r "restart i3"
i3-bind comment mod4+shift+e "exit i3"
i3-bind comment '$mod+return' "launch terminal"
```

#### Interactive mode (TUI)
```bash
i3-bind interactive
# or
i3-bind tui
# or
i3-bind menu
```

### Global Options

 - `--config, -c`: Specift custom i3 config file path
 - `--no-color`: Disable colored output
 - `--help, -h`: Show help information
 - `--version`: Show version information

### Example

```bash
# Use custom config file
i3-bind --config ~/.config/i3/config.backup list

# Disable colors for scripting
i3-bind --no-color list

# Add a keybinding with complex action
i3-bind add 'mod4+Print' 'exec --no-startup-id maim -s | xclip -selection clipboard -t image/png'

# Find all exec commands
i3-bind find exec

# Interactive mode for browsing and managing
i3-bind interactive
```

## Interactive Mode

The interactive TUI mode provides a user-friendly interface for managing keybindings:

1. Browse keybindings with fuzzy search
2. Preview key, action and comments
3. Remove keybinding
4. Add/Update comments
5. View detailed information

Navigation:
 - Use arrow keys or type to search
 - Press Enter to select
 - Ctrl+C to exit

## Configuration

i3-bind automatically detects your i3 config file at `~/.config/i3/config`. You can override this with the `--config` flag.

### Backup System

Before any modification, i3-bind creates a backup of your config file:
 - Backup location: `<config-path>.backup` (Example: `~/.config/i3/config.backup`)
 - Automatic restoration on write failures

## Output Format

i3-bind provides colorized output for better readability:
 - Keys: Cyan and bold
 - Actions: Green
 - Comments: Yellow
 - Errors: Red and bold
 - Success message: Green and bold

## Advanced Usage

### Scripting Integration

```bash
#!/bin/bash
# Add multiple keybindings from a script
keybindings=(
    "mod4+1 workspace 1"
    "mod4+2 workspace 2"
    "mod4+shift+1 move container to workspace 1"
    "mod4+shift+2 move container to workspace 2"
)
for binding in "${keybindings[@]}"; do
    key="${binding%% *}"
    action="${binding#"$key" }"
    i3-bind add "$key" "$action"
done
```

### Comment Organization

Use comments to organize your keybindings by category

```bash
i3-bind comment mod4+Return "Terminal applications"
i3-bind comment mod4+d "Application launcher"
i3-bind comment mod4+shift+q "Window management"
```

## Troubleshooting

### Common Issues

**Config file not found**
```bash
i3-bind --config /path/to/your/i3/config list
```

**Permission denied**
```bash
# Ensure you have write permissions to the config file
chmod 644 ~/.config/i3/config
```

**Interactive mode not working**
```bash
# Install fzf
sudo pacman -S fzf # Arch Linux
sudo apt install fzf # Ubuntu/Debian
brew install fzf # macOS
```

**Duplicate keybinding error***
```bash
# Remove existing binding first
i3-bind remove mod4+Return
i3-bind add mod4+Return "exec new-terminal"
```

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add come amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Requst

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

 - Build with [Cobra](https://github.com/spf13/cobra) for CLI functionality
 - Uses [fatih/color](https://github.com/fatih/color) for teminal colors
 - Interactive mode powered by [fzf](https://github.com/junegunn/fzf)
 - Inspired by i3 window manager community

---

**Author**: [Hanashiko](https://github.com/hanashiko)

**Repository**: [github.com/Hanashiko/i3-bind](https://github.com/hanashiko/i3-bind)
