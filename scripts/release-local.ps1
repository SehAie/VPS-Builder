$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path dist | Out-Null

go mod tidy
go test ./...

$env:GOOS="windows"
$env:GOARCH="amd64"
go build -trimpath -ldflags="-s -w" -o dist/VPS-Builder-windows-amd64.exe ./cmd/VPS-Builder

$env:GOARCH="arm64"
go build -trimpath -ldflags="-s -w" -o dist/VPS-Builder-windows-arm64.exe ./cmd/VPS-Builder

Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Copy-Item README.md dist\README.md -Force
Copy-Item LICENSE dist\LICENSE -Force
Copy-Item examples\config.example.yaml dist\config.example.yaml -Force

Compress-Archive -Path dist\VPS-Builder-windows-amd64.exe,dist\README.md,dist\LICENSE,dist\config.example.yaml -DestinationPath dist\VPS-Builder-windows-amd64.zip -Force
Compress-Archive -Path dist\VPS-Builder-windows-arm64.exe,dist\README.md,dist\LICENSE,dist\config.example.yaml -DestinationPath dist\VPS-Builder-windows-arm64.zip -Force

Get-FileHash dist\*.exe,dist\*.zip -Algorithm SHA256 | ForEach-Object { "$($_.Hash)  $([System.IO.Path]::GetFileName($_.Path))" } | Set-Content dist\SHA256SUMS.txt

Write-Host "Release files are ready in dist/:"
Get-ChildItem dist
