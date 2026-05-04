$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path dist | Out-Null

go mod tidy

$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o dist/VPS-Builder-windows-amd64.exe ./cmd/VPS-Builder

$env:GOARCH="arm64"
go build -o dist/VPS-Builder-windows-arm64.exe ./cmd/VPS-Builder

Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Write-Host "Windows build artifacts are in dist/."
Write-Host "Generated:"
Write-Host "  dist/VPS-Builder-windows-amd64.exe"
Write-Host "  dist/VPS-Builder-windows-arm64.exe"
