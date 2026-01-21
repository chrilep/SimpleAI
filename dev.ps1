#!/usr/bin/env pwsh
# Development script for SimpleAI - runs with hot reload and version injection

$ErrorActionPreference = "Stop"

Write-Host "Reading version from wails.json..." -ForegroundColor Cyan
$version = (Get-Content wails.json | ConvertFrom-Json).info.productVersion
Write-Host "Starting SimpleAI v$version in development mode..." -ForegroundColor Green
Write-Host "Hot reload enabled. Browser access: http://localhost:34115`n" -ForegroundColor Yellow

wails dev -ldflags "-X main.Version=$version"
