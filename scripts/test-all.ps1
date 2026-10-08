$ErrorActionPreference = "Stop"

$modules = @(
    "backend/proto/gen/go",
    "backend/api-gateway",
    "backend/auth-service",
    "backend/user-service",
    "backend/ride-service",
    "backend/location-service",
    "backend/matching-service",
    "backend/notification-service",
    "backend/pricing-service",
    "backend/payment-service"
)

foreach ($module in $modules) {
    Write-Host "Testing $module"
    Push-Location $module
    try {
        go test ./...
    }
    finally {
        Pop-Location
    }
}

Write-Host "All backend module tests passed."
