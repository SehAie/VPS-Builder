# Windows usage guide

## 1. Open PowerShell

Go to the project directory:

```powershell
cd "D:\your\path\VPS-Builder-github-ready"
```

## 2. Build the exe

```powershell
go mod tidy
go test ./...
go build -o VPS-Builder.exe ./cmd/VPS-Builder
```

Or build release files:

```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
.\scripts\release-local.ps1
```

## 3. Create a config file

```powershell
.\VPS-Builder.exe init -o config.yaml
```

You can edit the config manually:

```powershell
notepad .\config.yaml
```

## 4. Open required ports in your VPS provider panel

Default ports:

```text
TCP 22      current SSH port
TCP 2222    new SSH port
UDP 8443    Hysteria 2 port
```

If you changed the ports in `config.yaml`, open your custom ports instead.

## 5. Dry run

```powershell
.\VPS-Builder.exe deploy -c config.yaml --dry-run --ask-password
```

## 6. Deploy

```powershell
.\VPS-Builder.exe deploy -c config.yaml --ask-password
```

Type `yes` when asked to confirm.

## 7. Find local output files

```powershell
explorer $env:APPDATA\VPS-Builder
```

## 8. SSH login after deployment

Example:

```powershell
ssh -i "$env:APPDATA\VPS-Builder\keys\1.2.3.4_VPS-Builder-Key" -p 2222 root@1.2.3.4
```

Replace `1.2.3.4` with your real VPS IP.

