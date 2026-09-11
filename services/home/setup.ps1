<#
    Registers/updates the Windows Scheduled Task that runs dns-update.ps1 on startup.
    Run this from an elevated (Administrator) PowerShell prompt.
#>

#Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

$ScriptDir    = $PSScriptRoot
$DnsScript    = Join-Path $ScriptDir "dns-update.ps1"
$TaskName     = "UpdateDNS"

if (-not (Test-Path $DnsScript)) {
    Write-Error "Expected to find dns-update.ps1 at: $DnsScript`nMake sure this setup script lives alongside dns-update.ps1."
    exit 1
}

Write-Host "Configuring scheduled task '$TaskName'..."

$action = New-ScheduledTaskAction `
    -Execute "powershell.exe" `
    -Argument "-ExecutionPolicy Bypass -File `"$DnsScript`"" `
    -WorkingDirectory $ScriptDir

$trigger = New-ScheduledTaskTrigger -AtStartup
$trigger.Delay = "PT15S"   # network start

$principal = New-ScheduledTaskPrincipal `
    -UserId "SYSTEM" `
    -LogonType ServiceAccount `
    -RunLevel Highest

$settings = New-ScheduledTaskSettingsSet `
    -StartWhenAvailable `
    -AllowStartIfOnBatteries `
    -DontStopIfGoingOnBatteries

# Remove any existing task with the same name so this script is safely re-runnable
if (Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue) {
    Write-Host "Existing task '$TaskName' found — removing before re-registering..."
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
}

Register-ScheduledTask `
    -TaskName $TaskName `
    -Action $action `
    -Trigger $trigger `
    -Principal $principal `
    -Settings $settings `
    -Description "Updates DNSExit record with the Tailscale IP on startup" | Out-Null

Write-Host ""
Write-Host "Done. '$TaskName' is registered and will run at next startup."
Write-Host "Testing it now..."

Start-ScheduledTask -TaskName $TaskName
Start-Sleep -Seconds 2
Get-ScheduledTaskInfo -TaskName $TaskName | Format-List LastRunTime, LastTaskResult