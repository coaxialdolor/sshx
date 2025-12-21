# SSHX

```
   ________   ________   __  __
  / ____/ /  / ____/ /  / / / /
 / /   / /  / / __/ /  / /_/ /
/ /___/ /__/ /_/ / /__/ __  /
\____/____/\____/____/_/ /_/
      SSHX — Local-first SSH Superpowers
```

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![Platform](https://img.shields.io/badge/Platforms-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Alpha-orange)
![PRs](https://img.shields.io/badge/PRs-Welcome-blue)
![Security](https://img.shields.io/badge/Security-Zero%20Remote%20Footprint-critical)

> **SSHX** is a local‑first SSH wrapper that gives you superpowers inside any SSH session — without installing anything on the remote machine.
>
> Type `:d file.txt` or `~download /var/log/auth.log` inside SSH, and SSHX downloads the file directly to your machine using your existing SSH keys.
>
> No remote footprint. No shell modifications. No hacks. Just clean, universal, cross‑platform SSH enhancements.

## Features

- **Zero Remote Footprint**: No installation or configuration needed on remote machines
- **Cross-Platform**: Works on macOS, Linux, and Windows
- **Simple Commands**: Use `:d`, `:download`, `~d`, or `~download` to download files
- **Smart Path Resolution**: Supports shortcuts like `desktop`, `downloads`, `documents`, etc.
- **Configurable**: Customize default download directory and SSH settings

## Installation

### Prerequisites

- Go 1.22 or later ([download](https://go.dev/dl/))
- SSH client (usually pre-installed)

### Quick Install (Recommended)

#### macOS & Linux

Run the installer script:

```bash
git clone <repository-url>
cd sshdownload
./install.sh
```

The installer will:
1. Build the binaries
2. Install them to `/usr/local/bin` (or a custom location)
3. Optionally add an alias so `ssh` automatically uses `sshx`

#### Windows

Run the PowerShell installer:

```powershell
git clone <repository-url>
cd sshdownload
powershell -ExecutionPolicy Bypass -File install.ps1
```

The installer will:
1. Build the binaries
2. Install them to `%LOCALAPPDATA%\sshx`
3. Add the directory to your PATH
4. Optionally add a PowerShell alias so `ssh` automatically uses `sshx`

### Manual Build from Source

If you prefer to build manually:

```bash
git clone <repository-url>
cd sshdownload
go build ./cmd/sshx
go build ./cmd/sshx-agent
```

Then copy the binaries to a directory in your PATH:

**macOS/Linux:**
```bash
sudo cp sshx sshx-agent /usr/local/bin/
```

**Windows:**
Copy `sshx.exe` and `sshx-agent.exe` to a directory in your PATH, or add the directory containing them to your PATH.

## Usage

### 1. Start the Agent

First, start the `sshx-agent` daemon in a separate terminal:

```bash
sshx-agent
```

The agent will listen on `127.0.0.1:9234` and handle download requests.

### 2. Use SSHX

Connect to a remote server. If you enabled the alias during installation, you can use:

```bash
ssh user@example.com
```

Otherwise, use:

```bash
sshx user@example.com
```

> **Note:** The alias feature is optional and safe. It doesn't modify your system SSH binary—it only adds a shell alias in your profile. You can remove it anytime with `sshx uninstall-alias`.

### 3. Download Files

Once connected to the remote server, type any of these commands:

```bash
:d /path/to/remote/file.txt
:download /path/to/remote/file.txt
~d /path/to/remote/file.txt
~download /path/to/remote/file.txt
```

You can optionally specify a local destination:

```bash
:d /remote/file.txt ~/Downloads/
:d /remote/file.txt desktop
:d /remote/file.txt /absolute/path/to/save
```

## Managing the SSH Alias

### Removing the Alias

If you enabled the alias during installation and want to remove it:

```bash
sshx uninstall-alias
```

This will detect your shell, find the profile file, and remove only the SSHX alias while leaving the rest of your profile untouched.

### Manually Adding the Alias

If you didn't enable the alias during installation, you can add it manually:

**macOS/Linux (zsh/bash):**
Add to `~/.zshrc` or `~/.bashrc`:
```bash
alias ssh="sshx"
```

**macOS/Linux (fish):**
Add to `~/.config/fish/config.fish`:
```fish
alias ssh="sshx"
```

**Windows (PowerShell):**
Add to your PowerShell profile (run `$PROFILE` to see the path):
```powershell
Set-Alias ssh sshx
```

Then restart your shell or source the profile file.

## Configuration

Configuration file location:
- **macOS/Linux**: `~/.config/sshx/config.json`
- **Windows**: `%APPDATA%\sshx\config.json`

Example configuration:

```json
{
  "default_download_dir": "/Users/username/Downloads",
  "default_user": "",
  "default_port": 22
}
```

### Shortcuts

The following shortcuts are supported (case-insensitive):

- `home` → Home directory
- `downloads` → Downloads folder
- `desktop` → Desktop folder
- `documents` → Documents folder
- `pictures` → Pictures folder
- `videos` → Videos folder (or Movies on macOS if Videos doesn't exist)

## Examples

```bash
# Download to default directory
:d /var/log/app.log

# Download to desktop
:d /var/log/app.log desktop

# Download with absolute path
:d /var/log/app.log ~/Documents/logs/

# Download with custom name
:d /remote/file.txt ~/local/custom-name.txt
```

## How It Works

1. `sshx` wraps the standard `ssh` command
2. It intercepts special commands (`:d`, `:download`, etc.) locally
3. Download commands are sent to `sshx-agent` via TCP (localhost:9234)
4. The agent uses `scp` to download files using your SSH keys
5. Files are saved to the resolved local path
6. Normal SSH session continues uninterrupted

## Project Structure

```
sshx/
│
├── cmd/
│   ├── sshx/          # SSH wrapper binary
│   └── sshx-agent/    # Agent daemon binary
│
├── internal/
│   ├── alias/         # Alias management
│   ├── parser/        # Command parser
│   ├── shortcuts/     # Path shortcuts
│   ├── config/        # Configuration management
│   ├── pty/           # PTY handling
│   ├── protocol/      # Agent communication protocol
│   └── util/          # Utilities
│
├── install.sh         # Installer for macOS/Linux
├── install.ps1        # Installer for Windows
├── go.mod
└── README.md
```

## Development

### Build

```bash
go build ./cmd/sshx
go build ./cmd/sshx-agent
```

### Run Tests

```bash
go test ./...
```

## License

MIT

## Status

This project is in **Alpha** status. Core functionality is implemented, but expect bugs and missing features.

## Roadmap

- [ ] Upload functionality (`:u`, `:upload`)
- [ ] Directory sync with rsync
- [ ] Bookmarks for common remote paths
- [ ] Session tools (`:pwd`, `:ls`, etc.)
- [ ] Progress bars
- [ ] Parallel downloads

## Contributing

PRs welcome! Please ensure your code:
- Follows Go best practices
- Includes tests where appropriate
- Works on macOS, Linux, and Windows

## Security

SSHX operates entirely locally and uses your existing SSH keys. No modifications are made to remote machines. All communication between `sshx` and `sshx-agent` happens over localhost.

