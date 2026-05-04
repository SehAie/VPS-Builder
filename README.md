# VPS Builder / VPS-Builder Windows

A free and open-source Windows CLI tool for bootstrapping a Debian/Ubuntu VPS and deploying Hysteria 2.

This project is designed to run locally on Windows. Your VPS password and generated SSH private key stay on your own computer.

## Features

- Windows-only local CLI
- Interactive YAML config wizard
- SSH password login for the first connection
- Automatic ed25519 SSH key generation
- Public key upload to the VPS
- SSH port hardening
- Optional SSH password login disablement after key login verification
- Debian/Ubuntu package update
- Optional BBR/sysctl network optimization
- Hysteria 2 server installation
- Self-signed TLS certificate mode
- Local Hysteria 2 client config output
- GitHub Actions workflow for Windows exe releases

## Requirements

Local Windows computer:

- Windows 10 or Windows 11
- Go installed, for building from source
- PowerShell
- Git, if you want to publish to GitHub

Remote VPS:

- Debian or Ubuntu
- Root SSH login available for the first deployment
- Provider firewall/security group should allow:
  - TCP 22, or your current SSH port
  - TCP 2222, or your new SSH port
  - UDP 8443, or your Hysteria 2 port

## Build

```powershell
go mod tidy
go test ./...
go build -o VPS-Builder.exe ./cmd/VPS-Builder
```

Or use the provided script:

```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
.\scripts\release-local.ps1
```

Release files will be generated in `dist/`.

## Basic usage

Generate a config file:

```powershell
.\VPS-Builder.exe init -o config.yaml
```

Preview the deployment without changing the VPS:

```powershell
.\VPS-Builder.exe deploy -c config.yaml --dry-run --ask-password
```

Deploy:

```powershell
.\VPS-Builder.exe deploy -c config.yaml --ask-password
```

The tool will ask for your current VPS SSH/root password. The password is not saved to `config.yaml`.

## Local output files

By default, files are saved under:

```powershell
explorer $env:APPDATA\VPS-Builder
```

Typical files:

```text
keys\<IP>_VPS-Builder-Key
keys\<IP>_VPS-Builder-Key.pub
clients\<IP>_hysteria.yaml
last-deploy-summary.txt
```

## GitHub release

Push a version tag to trigger GitHub Actions and publish Windows exe assets:

```powershell
git tag v0.1.0
git push origin v0.1.0
```

The workflow uploads:

- `VPS-Builder-windows-amd64.exe`
- `VPS-Builder-windows-arm64.exe`
- ZIP packages
- `SHA256SUMS.txt`

## Safety notes

Use a fresh test VPS first. Do not test on an important production server.

Do not commit real passwords, private keys, or generated client configs to GitHub.

