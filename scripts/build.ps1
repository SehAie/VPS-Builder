$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path dist | Out-Null

go mod tidy

$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o dist/vps-auto-kit-windows-amd64.exe ./cmd/vps-auto-kit

$env:GOARCH="arm64"
go build -o dist/vps-auto-kit-windows-arm64.exe ./cmd/vps-auto-kit

Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Write-Host "Windows build artifacts are in dist/."
Write-Host "Generated:"
Write-Host "  dist/vps-auto-kit-windows-amd64.exe"
Write-Host "  dist/vps-auto-kit-windows-arm64.exe"
