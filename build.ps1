#!/usr/bin/env pwsh
# Build script for SimpleAI - injects version from wails.json into binary and creates installer

$ErrorActionPreference = "Stop"

Write-Host "Reading version from wails.json..." -ForegroundColor Cyan
$version = (Get-Content wails.json | ConvertFrom-Json).info.productVersion
Write-Host "Building SimpleAI v$version..." -ForegroundColor Green

# Build the application
wails build -ldflags "-X main.Version=$version"

if ($LASTEXITCODE -ne 0) {
    Write-Host "`nBuild failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "`nBuild successful! Executable: build\bin\SimpleAI.exe" -ForegroundColor Green

# Build NSIS installer
Write-Host "`nCreating NSIS installer..." -ForegroundColor Cyan
wails build -nsis -ldflags "-X main.Version=$version"

if ($LASTEXITCODE -eq 0) {
    Write-Host "`nInstaller created successfully!" -ForegroundColor Green
    Write-Host "Installer: build\bin\SimpleAI-amd64-installer.exe" -ForegroundColor Green
}
else {
    Write-Host "`nInstaller creation failed!" -ForegroundColor Yellow
    Write-Host "Application executable is still available in build\bin\" -ForegroundColor Yellow
}
