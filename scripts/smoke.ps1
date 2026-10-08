$ErrorActionPreference = "Stop"

$BaseUrl = $env:API_GATEWAY_URL
if ([string]::IsNullOrWhiteSpace($BaseUrl)) {
    $BaseUrl = "http://localhost:8088"
}

function Invoke-HealthCheck {
    param(
        [string]$Name,
        [string]$Url
    )

    Write-Host "Checking $Name at $Url"
    $response = Invoke-RestMethod -Method Get -Uri $Url -TimeoutSec 10
    if ($response.status -ne "ok") {
        throw "$Name health check failed"
    }
}

Invoke-HealthCheck -Name "api-gateway" -Url "$BaseUrl/health"

$services = @(
    @{ Name = "user-service"; Url = "http://localhost:8080/health" },
    @{ Name = "auth-service"; Url = "http://localhost:8081/health" },
    @{ Name = "ride-service"; Url = "http://localhost:8082/health" },
    @{ Name = "location-service"; Url = "http://localhost:8083/health" },
    @{ Name = "matching-service"; Url = "http://localhost:8084/health" },
    @{ Name = "notification-service"; Url = "http://localhost:8085/health" },
    @{ Name = "pricing-service"; Url = "http://localhost:8086/health" },
    @{ Name = "payment-service"; Url = "http://localhost:8087/health" }
)

foreach ($service in $services) {
    Invoke-HealthCheck -Name $service.Name -Url $service.Url
}

Write-Host "Smoke checks passed."
