# build and install a Go CLI app on Windows, with PowerShell completions

param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$App
)

# exit if something fails
$ErrorActionPreference = "Stop"

$destDir = Join-Path $env:LOCALAPPDATA "Programs\bin"
$completionsDir = Join-Path $HOME ".config\powershell\completions"

Push-Location $App
try {
    # build the app
    go mod tidy
    go build -o "$App.exe" main.go
    Write-Host "Compiled $App"

    # shell completions
    New-Item -ItemType Directory -Force -Path $completionsDir | Out-Null
    $completionFile = Join-Path $completionsDir "$App.ps1"
    & ".\$App.exe" completion powershell | Out-File -FilePath $completionFile -Encoding utf8
    Write-Host "Generated PowerShell completion script at $completionFile"

    # make sure the completion script gets loaded on new shells
    if (-not (Test-Path $PROFILE)) {
        New-Item -ItemType File -Path $PROFILE -Force | Out-Null
    }
    $sourceLine = ". `"$completionFile`""
    $alreadyWired = Select-String -Path $PROFILE -Pattern ([regex]::Escape($sourceLine)) -Quiet -ErrorAction SilentlyContinue
    if (-not $alreadyWired) {
        Add-Content -Path $PROFILE -Value $sourceLine
        Write-Host "Added completion sourcing to `$PROFILE ($PROFILE)"
    }

    # install in the user's bin
    # TODO systemwide install (e.g. C:\Program Files)
    New-Item -ItemType Directory -Force -Path $destDir | Out-Null
    Move-Item -Force "$App.exe" (Join-Path $destDir "$App.exe")
    Write-Host "Moved $App.exe to $destDir"
}
finally {
    Pop-Location
}

# warn if destDir isn't on PATH
$pathEntries = $env:Path -split ";"

if ($pathEntries -notcontains $destDir) {
    Write-Warning "$destDir is not in your PATH. Add it with:`n  [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$destDir', 'User')"
} else {
    Write-Host "Installation complete. $destDir is already in your PATH."
}