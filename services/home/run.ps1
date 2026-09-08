# Start-Process wsl.exe -ArgumentList "-d Debian -u root -e sleep infinity" -WindowStyle Hidden


# Get WSL's IPv4 address
$wslIp = $null
for ($i = 0; $i -lt 10; $i++) {
    $wslIp = (wsl.exe hostname -I 2>$null).Trim().Split(" ")[0]
    if ($wslIp) { break }
    Start-Sleep -Seconds 3
}

if (-not $wslIp) {
    Write-Error "Could not find WSL IPv4 address."
    exit 1
}

Write-Host "WSL IP: $wslIp"

$ports = @(80, 443, 9096, 5001, 5173, 8000)

foreach ($port in $ports) {
    # Remove existing rule
    netsh interface portproxy delete v4tov4 `
        listenaddress=0.0.0.0 `
        listenport=$port 2>$null

    # Add new rule
    netsh interface portproxy add v4tov4 `
        listenaddress=0.0.0.0 `
        listenport=$port `
        connectaddress=$wslIp `
        connectport=$port
}

Write-Host ""
Write-Host "Port forwarding configured:"
netsh interface portproxy show v4tov4

wsl.exe hostname -I

Read-Host "Press Enter to close"