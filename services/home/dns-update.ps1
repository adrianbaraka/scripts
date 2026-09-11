$logPath = Join-Path $PSScriptRoot "update-dns.log"
function Write-Log {
    param([string]$Message)
    "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') - $Message" | Set-Content $logPath
}

try {
    # Load .env
    Get-Content (Join-Path $PSScriptRoot ".env") | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]*)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()

            # Remove surrounding quotes
            $value = $value -replace '^["'']|["'']$', ''

            Set-Item "Env:$key" $value
        }
    }

    function Test-EnvExists {
        param (
            [string]$Name
        )

        if (-not (Test-Path "Env:$Name")) {
            throw "The env var '$Name' is not set."
        }
    }

    Test-EnvExists DOMAIN
    Test-EnvExists DNS_EXIT_TOKEN
    Test-EnvExists TAILSCALE_IP

    Write-Host "Updating $env:DOMAIN DNS to $env:TAILSCALE_IP"

    $result = curl.exe -fsS "https://api.dnsexit.com/dns/ud/" `
        -d "apikey=$env:DNS_EXIT_TOKEN" `
        -d "host=$env:DOMAIN,*.$env:DOMAIN" `
        -d "ip=$env:TAILSCALE_IP" 2>&1

    if ($LASTEXITCODE -ne 0) {
        throw "curl failed (exit $LASTEXITCODE): $result"
    }

    Write-Log "SUCCESS - Updated $env:DOMAIN to $env:TAILSCALE_IP - Response: $result"
}
catch {
    Write-Log "FAILED - $($_.Exception.Message)"
    exit 1
}