# GitHub open-source release guide

## 1. Create a GitHub repository

Recommended repository name:

```text
VPS-Builder
```

Create an empty public repository. Do not initialize it with README, .gitignore, or license because this project already includes them.

## 2. Make sure the Go module path is correct

For this package, the module path is already set to:

```text
github.com/SehAie/VPS-Builder
```

If you change the repository name later, run:

```powershell
.\scripts\rename-module.ps1 -ModulePath github.com/YourName/YourRepo
```

## 3. Test locally

```powershell
go mod tidy
go test ./...
.\scripts\release-local.ps1
```

## 4. Upload source to GitHub

```powershell
git init
git add .
git commit -m "Initial open-source release"
git branch -M main
git remote add origin https://github.com/SehAie/VPS-Builder.git
git push -u origin main
```

## 5. Publish a release automatically

```powershell
git tag v0.1.0
git push origin v0.1.0
```

Then open the GitHub repository, go to Actions, and wait for the release workflow to complete.

The generated release will include:

- `VPS-Builder-windows-amd64.exe`
- `VPS-Builder-windows-arm64.exe`
- `VPS-Builder-windows-amd64.zip`
- `VPS-Builder-windows-arm64.zip`
- `SHA256SUMS.txt`

