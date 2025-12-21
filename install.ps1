# SSHX Installer for Windows
# Run: powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

# Colors
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

Write-Output "SSHX Installer"
Write-Output "=============="
Write-Output ""

# Check if Go is installed
$goCmd = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCmd) {
    Write-ColorOutput Red "Error: Go is not installed or not in PATH"
    Write-Output "Please install Go 1.22+ from https://go.dev/dl/"
    exit 1
}

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$InstallDir = "$env:LOCALAPPDATA\sshx"

# Build binaries
Write-Output "Building binaries..."
Push-Location $ScriptDir

try {
    go build -o sshx.exe ./cmd/sshx
    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput Red "Error: Failed to build sshx"
        exit 1
    }

    go build -o sshx-agent.exe ./cmd/sshx-agent
    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput Red "Error: Failed to build sshx-agent"
        exit 1
    }

    Write-ColorOutput Green "✓ Binaries built successfully"
    Write-Output ""
}
finally {
    Pop-Location
}

# Create install directory
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

# Install binaries
Write-Output "Installing binaries to $InstallDir..."
Copy-Item "$ScriptDir\sshx.exe" "$InstallDir\sshx.exe" -Force
Copy-Item "$ScriptDir\sshx-agent.exe" "$InstallDir\sshx-agent.exe" -Force

Write-ColorOutput Green "✓ Binaries installed"
Write-Output ""

# Add to PATH if not already there
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
if ($env:Path -notlike "*$InstallDir*") {
    $currentUserPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    [System.Environment]::SetEnvironmentVariable("Path", "$currentUserPath;$InstallDir", "User")
    Write-ColorOutput Green "✓ Added $InstallDir to PATH"
    Write-Output "Note: You may need to restart your terminal for PATH changes to take effect."
    Write-Output ""
}

# Ask about alias
Write-Output "SSHX can enhance your SSH experience by making the 'ssh' command automatically use sshx."
Write-Output ""
Write-Output "This does NOT replace or modify your system SSH binary."
Write-Output "It simply adds a safe alias in your PowerShell profile."
Write-Output ""
$choice = Read-Host "Enable this feature? (1) Yes — make 'ssh' launch sshx (2) No — keep ssh unchanged [1/2]"

if ($choice -eq "1") {
    # Get PowerShell profile path
    $profilePath = $PROFILE

    # Create profile directory if it doesn't exist
    $profileDir = Split-Path -Parent $profilePath
    if (-not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir | Out-Null
    }

    # Check if alias already exists
    $aliasExists = $false
    if (Test-Path $profilePath) {
        $content = Get-Content $profilePath -Raw
        if ($content -match 'Set-Alias ssh sshx') {
            $aliasExists = $true
        }
    }

    if ($aliasExists) {
        Write-ColorOutput Yellow "Alias already exists in $profilePath"
    } else {
        # Add alias to profile
        if (Test-Path $profilePath) {
            Add-Content -Path $profilePath -Value "`n# SSHX alias`nSet-Alias ssh sshx"
        } else {
            "# SSHX alias`nSet-Alias ssh sshx" | Out-File -FilePath $profilePath -Encoding UTF8
        }
        Write-ColorOutput Green "✓ Alias added to $profilePath"
        Write-Output "Restart PowerShell or run: . $profilePath"
    }
} else {
    Write-Output "Alias not added. You can add it manually later or run the installer again."
}

Write-Output ""
Write-Output "=========================================="
Write-Output "Installation complete!"
Write-Output ""
Write-Output "To start using SSHX:"
Write-Output "  sshx user@host"
Write-Output ""
if ($choice -eq "1") {
    Write-Output "If you enabled the alias:"
    Write-Output "  ssh user@host"
    Write-Output ""
}
Write-Output "To undo the alias:"
Write-Output "  sshx uninstall-alias"
Write-Output ""
Write-Output "To start the agent:"
Write-Output "  sshx-agent"
Write-Output ""

