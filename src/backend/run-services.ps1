run-services.ps1

$Services = @(
    @{
        Name = "wiki-service"
        Dir = "wiki"
        Port = "8001"
        Executable = "wiki-service.exe"
        LogFile = "wiki-service.log"
        PIDFile = "wiki-service.pid"
    },
    @{
        Name = "entry-service"
        Dir = "entry"
        Port = "8002"
        Executable = "entry-service.exe"
        LogFile = "entry-service.log"
        PIDFile = "entry-service.pid"
    },
    @{
        Name = "comment-service"
        Dir = "comment"
        Port = "8003"
        Executable = "comment-service.exe"
        LogFile = "comment-service.log"
        PIDFile = "comment-service.pid"
    },
    @{
        Name = "version-service"
        Dir = "version"
        Port = "8005"
        Executable = "version-service.exe"
        LogFile = "version-service.log"
        PIDFile = "version-service.pid"
    },
    @{
        Name = "media-service"
        Dir = "media"
        Port = "8081"
        Executable = "media-service.exe"
        LogFile = "media-service.log"
        PIDFile = "media-service.pid"
    },
    @{
        Name = "auth-service"
        Dir = "auth"
        Port = "8082"
        Executable = "auth-service.exe"
        LogFile = "auth-service.log"
        PIDFile = "auth-service.pid"
    },
    @{
        Name = "api-gateway"
        Dir = "gateway"
        Port = "8000"
        Executable = "api-gateway.exe"
        LogFile = "api-gateway.log"
        PIDFile = "api-gateway.pid"
    }
)

Function Start-Service {
    param($Service)
    Write-Host "Starting $($Service.Name)..."

    # Change to service directory
    Push-Location $Service.Dir

    # Build the service
    go build -o $Service.Executable .

    # Prepare start info
    $StartInfo = New-Object System.Diagnostics.ProcessStartInfo
    $StartInfo.FileName = ".\$($Service.Executable)"
    $StartInfo.WorkingDirectory = (Get-Location).Path
    $StartInfo.RedirectStandardOutput = $true
    $StartInfo.RedirectStandardError = $true
    $StartInfo.UseShellExecute = $false

    # Start the process
    $Process = New-Object System.Diagnostics.Process
    $Process.StartInfo = $StartInfo
    $Process.Start() | Out-Null

    # Async read output
    $StdOut = $Process.StandardOutput
    $StdErr = $Process.StandardError
    Start-Job -ScriptBlock {
        param($Stream, $LogFile)
        while (!$Stream.EndOfStream) {
            $Line = $Stream.ReadLine()
            $Line | Out-File $LogFile -Append
        }
    } -ArgumentList $StdOut, "../$($Service.LogFile)" | Out-Null
    Start-Job -ScriptBlock {
        param($Stream, $LogFile)
        while (!$Stream.EndOfStream) {
            $Line = $Stream.ReadLine()
            $Line | Out-File $LogFile -Append
        }
    } -ArgumentList $StdErr, "../$($Service.LogFile)" | Out-Null

    # Save PID
    $Process.Id | Out-File "../$($Service.PIDFile)"

    # Return to previous directory
    Pop-Location
}

Function Start-AllServices {
    foreach ($Service in $Services) {
        # Start dependencies first
        if ($Service.ContainsKey("DependsOn")) {
            foreach ($Dependency in $Service.DependsOn) {
                $DependentService = $Services | Where-Object { $_.Name -eq $Dependency }
                if ($null -eq $DependentService) {
                    Write-Warning "Cannot find dependency $Dependency for service $($Service.Name)"
                } elseif (-not (Get-Process -Name $DependentService.Executable -ErrorAction SilentlyContinue)) {
                    Start-Service $DependentService
                }
            }
        }

        # Start the service if it's not already running
        if (-not (Get-Process -Name $Service.Executable -ErrorAction SilentlyContinue)) {
            Start-Service $Service
        } else {
            Write-Host "$($Service.Name) is already running."
        }
    }
}

Function Stop-AllServices {
    foreach ($Service in $Services) {
        $ProcessName = [System.IO.Path]::GetFileNameWithoutExtension($Service.Executable)
        Write-Host "Stopping $($Service.Name)..."
        Get-Process -Name $ProcessName -ErrorAction SilentlyContinue | Stop-Process -Force
        # Remove PID file
        $PIDFilePath = "$($Service.PIDFile)"
        if (Test-Path $PIDFilePath) {
            Remove-Item $PIDFilePath -Force
        }
    }
}

# Main logic
param(
    [string] $Action = "start"
)

switch ($Action) {
    "start" { Start-AllServices }
    "stop"  { Stop-AllServices }
    default { Write-Host "Invalid action. Use 'start' or 'stop'." }
}
