param(
  [Parameter(Mandatory=$true)]
  [string]$ModulePath
)

$ErrorActionPreference = "Stop"

if ($ModulePath -notmatch '^github\.com/[^/]+/[^/]+$') {
  throw "ModulePath should look like github.com/yourname/vps-auto-kit"
}

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
  throw "Go is not installed or not available in PATH. Install Go first, reopen PowerShell, then run this script again."
}

$oldPaths = @(
  "github.com/example/vps-auto-kit",
  "github.com/SehAie/VPS-Builder"
)

go mod edit -module $ModulePath

Get-ChildItem -Path . -Recurse -Include *.go | ForEach-Object {
  $content = Get-Content $_.FullName -Raw -Encoding UTF8
  foreach ($old in $oldPaths) {
    $content = $content.Replace($old, $ModulePath)
  }
  Set-Content $_.FullName $content -NoNewline -Encoding UTF8
}

go mod tidy

Write-Host "Module path changed to $ModulePath"
Write-Host "Next: go test ./..."
