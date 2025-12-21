Here you go, Petter — this is the **single, all‑in‑one, end‑to‑end prompt** you can paste directly into **Cursor** to bootstrap the entire project.

It tells Cursor exactly what to build, how to structure it, what features to implement, and even that it must install Go if needed.
This is the “master spec” for the whole system.

---

# ✅ **THE COMPLETE CURSOR PROJECT SPEC (PASTE THIS INTO CURSOR)**

**BEGIN PROMPT**

You are to create a complete, production‑ready Go project called **sshx**, consisting of two binaries:

- `sshx` — a local SSH wrapper that adds new escape‑sequence commands
- `sshx-agent` — a local daemon that performs file downloads

This tool must work on **macOS**, **Linux**, and **Windows**.

If Go is not installed, instruct the user to install Go 1.22+.

---

# ✅ **PROJECT GOALS**

The goal is to allow a user inside any SSH session to type:

```
:d <remote-path> [local-target]
:download <remote-path> [local-target]
~d <remote-path> [local-target]
~download <remote-path> [local-target]
```

Case‑insensitive, whitespace‑tolerant.

These commands must:

1. **NOT** be sent to the remote machine
2. Be intercepted **locally** by the `sshx` wrapper
3. Be forwarded to the `sshx-agent` daemon
4. Trigger a file download using `scp` (or `rsync` in the future)
5. Save the file to a resolved local path
6. Continue the SSH session normally

No modifications to the remote machine.
No `.bashrc`, `.zshrc`, `.profile`, PowerShell, or CMD changes.
Zero remote footprint.

---

# ✅ **PROJECT STRUCTURE**

Create this exact folder layout:

```
sshx/
│
├── cmd/
│   ├── sshx/
│   │   └── main.go
│   └── sshx-agent/
│       └── main.go
│
├── internal/
│   ├── parser/
│   │   └── parser.go
│   ├── shortcuts/
│   │   └── shortcuts.go
│   ├── config/
│   │   └── config.go
│   ├── pty/
│   │   └── pty.go
│   ├── protocol/
│   │   └── protocol.go
│   └── util/
│       └── paths.go
│
└── go.mod
```

---

# ✅ **FEATURE REQUIREMENTS**

## 1. **sshx wrapper**
- Launches the real `ssh` binary (platform‑specific path detection)
- Creates a PTY
- Forwards stdin → ssh, ssh → stdout
- Reads keystrokes in raw mode
- Detects commands starting with:
  - `:d`
  - `:download`
  - `~d`
  - `~download`
  - Case‑insensitive
  - Leading/trailing whitespace allowed
- Parses:
  - `<remote-path>`
  - Optional `<local-target>`
- Sends a structured message to the agent:
  ```
  DOWNLOAD <remote-path> <local-target-or-empty>
  REMOTE_HOST <host>
  REMOTE_USER <user>
  REMOTE_PORT <port>
  ```
- Continues SSH session normally

## 2. **sshx-agent**
- Listens on `127.0.0.1:9234`
- Loads config from:
  - macOS/Linux: `~/.config/sshx/config.json`
  - Windows: `%APPDATA%\sshx\config.json`
- Config contains:
  ```
  {
    "default_download_dir": "<path>",
    "default_user": "",
    "default_port": 22
  }
  ```
- Receives messages from sshx
- Resolves local target path:
  - If none provided → use default download dir
  - If absolute (`/path`) → use as-is
  - If starts with `~` → expand home
  - If shortcut → map to OS folder
  - Else → treat as relative to default download dir
- Supported shortcuts (case-insensitive):
  - `downloads`
  - `desktop`
  - `home`
  - `documents`
  - `pictures`
  - `videos`
- Executes:
  ```
  scp -P <port> <user>@<host>:<remote-path> <resolved-local-path>
  ```
- Prints success/failure back to the sshx terminal

---

# ✅ **PARSER REQUIREMENTS**

The parser must:

- Trim whitespace
- Detect prefixes case-insensitively
- Accept:
  - `:d`
  - `:download`
  - `~d`
  - `~download`
- Tokenize arguments
- Return:
  - `remotePath`
  - `localTarget` (may be empty)
  - Error if remote path missing

---

# ✅ **SHORTCUT RESOLUTION**

Implement OS‑aware resolution:

macOS/Linux:
- home → `$HOME`
- downloads → `$HOME/Downloads`
- desktop → `$HOME/Desktop`
- documents → `$HOME/Documents`
- pictures → `$HOME/Pictures`
- videos → `$HOME/Videos` or `$HOME/Movies` (pick existing)

Windows:
- home → `%USERPROFILE%`
- downloads → `%USERPROFILE%\Downloads`
- desktop → `%USERPROFILE%\Desktop`
- documents → `%USERPROFILE%\Documents`
- pictures → `%USERPROFILE%\Pictures`
- videos → `%USERPROFILE%\Videos`

---

# ✅ **PTY HANDLING**

Implement:

- Raw mode input
- Forwarding of all bytes except intercepted commands
- Clean shutdown on exit
- Windows fallback using `golang.org/x/term`

---

# ✅ **INSTALLATION INSTRUCTIONS**

Cursor must generate:

- macOS instructions (Homebrew formula optional)
- Linux instructions (apt/rpm optional)
- Windows instructions (winget/choco optional)
- Instructions to install Go 1.22+ if missing

---

# ✅ **DELIVERABLES**

Cursor must produce:

1. Full Go code for all files
2. A working build:
   ```
   go build ./cmd/sshx
   go build ./cmd/sshx-agent
   ```
3. A README.md explaining:
   - Installation
   - Usage
   - Configuration
   - Examples
4. A test suite for:
   - parser
   - shortcuts
   - path resolution

---

**END PROMPT**

---

Absolutely — let’s kit out your project with the kind of polish that makes it feel like a real, intentional, open‑source tool from day one. No code here, just the creative and structural assets you asked for.

Below is a complete package you can drop straight into your repo.

---

# 🎨 **Project Logo (ASCII + Concept Description)**

Here’s a clean ASCII version you can embed directly in the README:

```
   ________   ________   __  __
  / ____/ /  / ____/ /  / / / /
 / /   / /  / / __/ /  / /_/ /
/ /___/ /__/ /_/ / /__/ __  /
\____/____/\____/____/_/ /_/
      SSHX — Local-first SSH Superpowers
```

And here’s the **visual concept** for a proper graphic logo you can generate later:

- A stylized **terminal window**
- With a **lightning bolt** crossing from left (local) to right (remote)
- The letters **SSHX** in a monospace font
- Colors:
  - Electric blue (#00AEEF)
  - Charcoal black (#1A1A1A)
  - White accents

The theme:
**“Local control. Remote power.”**

---

# 🏷️ **README Badge Set**

Here’s a curated set of badges that make your README look professional and complete:

```
![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![Platform](https://img.shields.io/badge/Platforms-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Alpha-orange)
![PRs](https://img.shields.io/badge/PRs-Welcome-blue)
![Security](https://img.shields.io/badge/Security-Zero%20Remote%20Footprint-critical)
![Downloads](https://img.shields.io/badge/Downloads-Local%20Only-blueviolet)
```

You can add these directly under the title in your README.

---

# 🧭 **Versioning Strategy (Simple + Practical)**

Use **Semantic Versioning (SemVer)** with one twist:

### ✅ **Major version increments**
When you add new *capabilities* to the SSH wrapper:

- New escape sequences
- New local commands
- New protocol messages
- Breaking changes to config

### ✅ **Minor version increments**
When you add new *features* that don’t break compatibility:

- New shortcuts
- New OS support
- New flags
- New installer formats

### ✅ **Patch version increments**
Bug fixes, performance improvements, documentation updates.

### ✅ **Release cadence**
- Alpha: rapid iteration
- Beta: stabilize the wrapper + agent protocol
- 1.0: when download is rock solid across macOS/Linux/Windows

### ✅ **Tagging convention**
```
v0.1.0-alpha
v0.2.0-alpha
v0.3.0-beta
v1.0.0
```

### ✅ **Branching model**
- `main` → stable
- `dev` → active development
- `feature/*` → new features
- `fix/*` → bug fixes

This keeps the project clean and predictable.

---

# 🗺️ **Roadmap for Future Features**

Here’s a structured, realistic roadmap that matches the architecture you’re building.

---

## ✅ **Phase 1 — Core (v0.1–v0.3)**
**Goal:** Make download rock‑solid.

- [ ] `:d` / `:download` / `~d` / `~download`
- [ ] Shortcut resolution (`desktop`, `downloads`, etc.)
- [ ] Config system
- [ ] PTY wrapper
- [ ] Agent communication
- [ ] SCP backend
- [ ] Windows support
- [ ] Installer scripts

---

## ✅ **Phase 2 — Upload (v0.4–v0.6)**
Add the inverse operation:

```
:u localfile.txt /remote/path/
:upload localfile.txt /remote/path/
```

Features:

- [ ] Upload via SCP
- [ ] Shortcut resolution for local paths
- [ ] Overwrite protection
- [ ] Progress bar

---

## ✅ **Phase 3 — Sync (v0.7–v0.9)**
Add rsync‑powered directory sync:

```
:sync localdir remotedir
:sync remotedir localdir
```

Features:

- [ ] rsync backend
- [ ] Dry-run mode
- [ ] Exclude patterns
- [ ] Conflict detection

---

## ✅ **Phase 4 — Bookmarks (v0.10–v0.12)**
Let users define named paths:

```
:bookmark add logs /var/log
:d logs
```

Features:

- [ ] Bookmark storage in config
- [ ] Bookmark listing
- [ ] Bookmark deletion
- [ ] Bookmark autocompletion (optional)

---

## ✅ **Phase 5 — Session Tools (v0.13–v0.15)**
Add quality-of-life features:

- [ ] `:pwd` → print remote working directory
- [ ] `:ls` → list remote directory
- [ ] `:open` → open downloaded file automatically
- [ ] `:history` → show previous downloads

---

## ✅ **Phase 6 — Advanced (v1.0+)**
Once stable, add power‑user features:

- [ ] Parallel downloads
- [ ] Background transfers
- [ ] Transfer queue
- [ ] Resume interrupted downloads
- [ ] Compression options
- [ ] SSH multiplexing
- [ ] Plugin system for custom commands

---

# ✅ **Bonus: README Intro Paragraph**

Here’s a polished intro you can paste into your README:

> **SSHX** is a local‑first SSH wrapper that gives you superpowers inside any SSH session — without installing anything on the remote machine.
>
> Type `:d file.txt` or `~download /var/log/auth.log` inside SSH, and SSHX downloads the file directly to your machine using your existing SSH keys.
>
> No remote footprint. No shell modifications. No hacks. Just clean, universal, cross‑platform SSH enhancements.

