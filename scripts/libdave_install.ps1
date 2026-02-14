# libdave-install.ps1
# Usage: .\libdave-install.ps1 -Version "v0.0.1"

[CmdletBinding(PositionalBinding=$false)]
param (
    [Parameter(Mandatory=$true, Position=0)]
    [string]$Version,
    [switch]$ForceBuild,
    [string]$SslFlavour = "boringssl"
)

$ErrorActionPreference = "Stop"

# --- Configuration ---
$RepoOwner   = "discord"
$RepoName    = "libdave"
$LibDaveRepo = "https://github.com/$RepoOwner/$RepoName"

$InstallBase = Join-Path $env:LOCALAPPDATA "libdave"
$BinDir      = Join-Path $InstallBase "bin"
$LibDir      = Join-Path $InstallBase "lib"
$IncDir      = Join-Path $InstallBase "include"
$PcDir       = Join-Path $env:LOCALAPPDATA "pkgconfig"
$PcFile      = Join-Path $PcDir "dave.pc"

function Log-Info ([string]$Msg) {
    Write-Host "-> $Msg" -ForegroundColor Cyan
}

function Check-Dependencies {
    $deps = @("git","make","cmake")
    foreach ($cmd in $deps) {
        if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
            throw "Missing dependency: $cmd"
        }
    }
}

function Get-Environment {
    $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { "X64" }
        "ARM64" { "ARM64" }
        default { $_ }
    }
    return @{ Arch = $arch }
}

function Install-Prebuilt {
    param($Tag, $Env)

    $AssetPattern = "libdave-Windows-$($Env.Arch)-$SslFlavour.zip"
    $DownloadUrl  = "$LibDaveRepo/releases/download/$Tag/$AssetPattern"
    $TempZip      = Join-Path $env:TEMP "libdave_prebuilt.zip"
    $StageDir     = Join-Path $env:TEMP "libdave_stage"

    Log-Info "Checking for prebuilt asset at: $DownloadUrl"

    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZip -UseBasicParsing
    }
    catch {
        Log-Info "No prebuilt asset found. Falling back to build."
        return $false
    }

    Log-Info "Found prebuilt asset. Extracting..."

    if (Test-Path $StageDir) {
        Remove-Item $StageDir -Recurse -Force
    }

    Expand-Archive -Path $TempZip -DestinationPath $StageDir -Force

    if (Test-Path $InstallBase) {
        Remove-Item $InstallBase -Recurse -Force
    }

    New-Item -ItemType Directory -Path $BinDir,$LibDir,$IncDir -Force | Out-Null

    Copy-Item "$StageDir\include\dave\dave.h" -Destination $IncDir -Recurse
    Copy-Item "$StageDir\bin\libdave.dll"     -Destination $BinDir
    Copy-Item "$StageDir\lib\libdave.lib"     -Destination $LibDir

    Remove-Item $TempZip -Force
    Remove-Item $StageDir -Recurse -Force

    return $true
}

function Build-Manual {
    param($Ref)

    Log-Info "Starting manual build process for ref: $Ref ($SslFlavour)"
    Check-Dependencies

    $WorkDir = Join-Path $env:TEMP "libdave_build_$(New-Guid)"
    New-Item -ItemType Directory -Path $WorkDir | Out-Null

    git clone $LibDaveRepo $WorkDir

    $CurrentDir = Get-Location
    Set-Location (Join-Path $WorkDir "cpp")

    git checkout $Ref
    git submodule update --init --recursive

    Log-Info "Bootstrapping vcpkg..."
    .\vcpkg\bootstrap-vcpkg.bat -disableMetrics

    Log-Info "Compiling shared library..."
    make shared "SSL=$SslFlavour" BUILD_TYPE=Release

    Log-Info "Installing..."

    if (Test-Path $InstallBase) {
        Remove-Item $InstallBase -Recurse -Force
    }

    New-Item -ItemType Directory -Path $BinDir,$LibDir,$IncDir -Force | Out-Null

    Copy-Item "includes\dave\dave.h"        -Destination $IncDir
    Copy-Item "build\Release\libdave.dll"   -Destination $BinDir
    Copy-Item "build\Release\libdave.lib"   -Destination $LibDir

    Set-Location $CurrentDir
    Remove-Item $WorkDir -Recurse -Force
}

function Generate-PkgConfig {
    Log-Info "Generating pkg-config metadata..."

    if (-not (Test-Path $PcDir)) {
        New-Item -ItemType Directory -Path $PcDir -Force | Out-Null
    }

    $Prefix = $InstallBase.Replace('\','/')

$PcContent = @"
prefix=$Prefix
exec_prefix=$Prefix/bin
libdir=$Prefix/lib
includedir=$Prefix/include

Name: dave
Description: Discord Audio & Video End-to-End Encryption (DAVE) Protocol
Version: $Version
URL: $LibDaveRepo
Libs: -L`${libdir} -ldave
Cflags: -I`${includedir}
"@

    Out-File -FilePath $PcFile -InputObject $PcContent -Encoding UTF8
    Log-Info "Created $PcFile"
}

# --- Main ---
$CurrentDir = Get-Location

try {
    $EnvInfo = Get-Environment

    $IsSha = $Version -match "^[0-9a-fA-F]{7,40}$"
    $BuildRef = if ($IsSha) {
        $Version
    } else {
        "$($Version.Replace('/cpp',''))/cpp"
    }

    if ($IsSha -or $ForceBuild) {
        Build-Manual -Ref $BuildRef
    }
    else {
        $ok = Install-Prebuilt -Tag $BuildRef -Env $EnvInfo
        if (-not $ok) {
            Build-Manual -Ref $BuildRef
        }
    }

    Generate-PkgConfig

    Log-Info "Installation successful: libdave $Version ($($EnvInfo.Arch))"
}
finally {
    Set-Location $CurrentDir
}
