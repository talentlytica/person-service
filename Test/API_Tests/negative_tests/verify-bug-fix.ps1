# Script untuk memverifikasi bug fix setelah server restart
# Usage: .\verify-bug-fix.ps1

Write-Host "`n=== VERIFIKASI BUG FIX ===" -ForegroundColor Yellow
Write-Host "Memeriksa apakah server sudah menggunakan kode yang diperbaiki...`n" -ForegroundColor Cyan

$BASE_URL = "http://localhost:3000"
$API_KEY = "person-service-key-82aca3c8-8e5d-42d4-9b00-7bc2f3077a58"
$headers = @{"x-api-key" = $API_KEY}

$testCases = @(
    @{
        Name = "GET non-existent person"
        Method = "GET"
        Url = "$BASE_URL/persons/00000000-0000-0000-0000-000000000000/attributes"
        ExpectedStatus = 404
        ExpectedMessage = "Person not found"
    },
    @{
        Name = "POST non-existent person"
        Method = "POST"
        Url = "$BASE_URL/persons/00000000-0000-0000-0000-000000000000/attributes"
        ExpectedStatus = 404
        ExpectedMessage = "Person not found"
        Body = @{
            key = "email"
            value = "test@example.com"
            meta = @{
                caller = "test"
                reason = "test"
                traceId = "123"
            }
        } | ConvertTo-Json
    }
)

$allPassed = $true

foreach ($testCase in $testCases) {
    Write-Host "Testing: $($testCase.Name)..." -ForegroundColor Cyan -NoNewline
    
    try {
        if ($testCase.Method -eq "GET") {
            $response = Invoke-WebRequest -Uri $testCase.Url -Headers $headers -Method GET -ErrorAction Stop
            $statusCode = $response.StatusCode
            $responseBody = $response.Content
        } else {
            $response = Invoke-WebRequest -Uri $testCase.Url -Headers $headers -Method POST -Body $testCase.Body -ContentType "application/json" -ErrorAction Stop
            $statusCode = $response.StatusCode
            $responseBody = $response.Content
        }
        
        Write-Host " ❌ FAILED" -ForegroundColor Red
        Write-Host "   Expected: Status $($testCase.ExpectedStatus), but got: Status $statusCode" -ForegroundColor Red
        Write-Host "   Response: $responseBody" -ForegroundColor Red
        $allPassed = $false
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        
        if ($statusCode -eq $testCase.ExpectedStatus) {
            $responseJson = $responseBody | ConvertFrom-Json
            if ($responseJson.message -like "*$($testCase.ExpectedMessage)*") {
                Write-Host " ✅ PASSED" -ForegroundColor Green
                Write-Host "   Status: $statusCode" -ForegroundColor Gray
                Write-Host "   Message: $($responseJson.message)" -ForegroundColor Gray
            } else {
                Write-Host " ⚠️  PARTIAL" -ForegroundColor Yellow
                Write-Host "   Status: $statusCode (correct)" -ForegroundColor Gray
                Write-Host "   Message: $($responseJson.message) (expected: $($testCase.ExpectedMessage))" -ForegroundColor Yellow
            }
        } else {
            Write-Host " ❌ FAILED" -ForegroundColor Red
            Write-Host "   Expected: Status $($testCase.ExpectedStatus), but got: Status $statusCode" -ForegroundColor Red
            Write-Host "   Response: $responseBody" -ForegroundColor Red
            $allPassed = $false
        }
    }
}

Write-Host "`n=== HASIL VERIFIKASI ===" -ForegroundColor Yellow

if ($allPassed) {
    Write-Host "✅ SEMUA TEST PASSED!" -ForegroundColor Green
    Write-Host "Bug fix sudah berhasil diterapkan di server." -ForegroundColor Green
    Write-Host "`nSelanjutnya, jalankan automated test:" -ForegroundColor Cyan
    Write-Host '  cd "c:\RepoGit\person-service - v2\Test"' -ForegroundColor White
    Write-Host '  npm run test:negative' -ForegroundColor White
} else {
    Write-Host "❌ BEBERAPA TEST GAGAL" -ForegroundColor Red
    Write-Host "`nKemungkinan penyebab:" -ForegroundColor Yellow
    Write-Host "1. Server belum di-restart setelah perubahan kode" -ForegroundColor White
    Write-Host "2. Server masih menjalankan binary/kode lama" -ForegroundColor White
    Write-Host "3. Ada masalah lain dengan server" -ForegroundColor White
    Write-Host "`nPastikan:" -ForegroundColor Yellow
    Write-Host "1. Server sudah dihentikan (Ctrl+C)" -ForegroundColor White
    Write-Host "2. Server sudah dijalankan ulang dengan: go run main.go" -ForegroundColor White
    Write-Host "3. Server menampilkan: 'Server ready on port 3000'" -ForegroundColor White
}

Write-Host "`n"
