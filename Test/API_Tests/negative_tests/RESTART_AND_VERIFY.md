# 🔄 RESTART SERVER & VERIFIKASI BUG FIX

## ⚠️ **STATUS SAAT INI**

- ✅ **Kode sudah diperbaiki** di `source/app/person_attributes/person_attributes.go`
- ❌ **Server masih menjalankan kode lama** (masih mengembalikan 500)
- ⚠️ **Perlu restart server** untuk menerapkan perubahan

---

## 📋 **LANGKAH RESTART SERVER**

### **1. Hentikan Server yang Berjalan**

Di terminal tempat server berjalan:
- Tekan `Ctrl + C` untuk menghentikan server
- Tunggu sampai server benar-benar berhenti

### **2. Jalankan Server dengan Kode Baru**

```powershell
# Navigate ke direktori aplikasi
cd "c:\RepoGit\person-service - v2\source\app"

# Set environment variables
$env:DATABASE_URL="postgresql://postgres:postgres@localhost:5432/person_service?sslmode=disable"
$env:PERSON_API_KEY_GREEN="person-service-key-82aca3c8-8e5d-42d4-9b00-7bc2f3077a58"
$env:PORT="3000"

# Jalankan server
go run main.go
```

### **3. Verifikasi Server Berjalan**

Tunggu sampai melihat pesan:
```
INFO: Server ready on port 3000
```

---

## ✅ **VERIFIKASI BUG FIX**

### **Test Manual dengan PowerShell:**

```powershell
# Test GET non-existent person
$headers = @{"x-api-key" = "person-service-key-82aca3c8-8e5d-42d4-9b00-7bc2f3077a58"}
try {
    $response = Invoke-WebRequest -Uri "http://localhost:3000/persons/00000000-0000-0000-0000-000000000000/attributes" -Headers $headers -Method GET
    Write-Host "❌ FAILED: Got status $($response.StatusCode) - Expected 404"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 404) {
        Write-Host "✅ SUCCESS: Got 404 as expected"
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response: $responseBody"
    } else {
        Write-Host "❌ FAILED: Got status $statusCode - Expected 404"
    }
}
```

**Expected Output:**
```
✅ SUCCESS: Got 404 as expected
Response: {"message":"Person not found"}
```

### **Jalankan Automated Test:**

```powershell
cd "c:\RepoGit\person-service - v2\Test"

# Test GET endpoint
npm test -- API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js --testNamePattern="Should return empty or 404 for non-existent personId"

# Test POST endpoint
npm test -- API_Tests/negative_tests/POST_persons_personId_attributes_negative.test.js --testNamePattern="Should reject request with non-existent personId"

# Atau jalankan semua negative tests
npm run test:negative
```

---

## 🎯 **EXPECTED RESULTS SETELAH RESTART**

### **Sebelum Fix (Kode Lama):**
```
Status Code: 500
Response: {"message":"Failed to verify person"}
```

### **Setelah Fix (Kode Baru):**
```
Status Code: 404
Response: {"message":"Person not found"}
```

---

## 🔍 **TROUBLESHOOTING**

### **Masih Mendapat 500 Setelah Restart?**

1. **Pastikan server benar-benar restart:**
   - Stop server (Ctrl+C)
   - Tunggu 2-3 detik
   - Start lagi

2. **Pastikan kode sudah benar:**
   ```powershell
   # Cek apakah fix ada di kode
   Select-String -Path "source\app\person_attributes\person_attributes.go" -Pattern "pgx.ErrNoRows"
   ```
   Harus menemukan minimal 2 hasil (untuk GetAllAttributes dan CreateAttribute)

3. **Pastikan tidak ada binary lama:**
   ```powershell
   # Hapus binary lama jika ada
   Remove-Item "source\app\person-service.exe" -ErrorAction SilentlyContinue
   ```

4. **Build ulang:**
   ```powershell
   cd "c:\RepoGit\person-service - v2\source\app"
   go clean
   go run main.go
   ```

---

## 📊 **CHECKLIST VERIFIKASI**

Setelah restart server, pastikan:

- [ ] Server berjalan tanpa error
- [ ] GET `/persons/00000000-0000-0000-0000-000000000000/attributes` mengembalikan **404**
- [ ] POST `/persons/00000000-0000-0000-0000-000000000000/attributes` mengembalikan **404**
- [ ] Response message adalah `"Person not found"` (bukan `"Failed to verify person"`)
- [ ] Tidak ada error 500 untuk non-existent person ID
- [ ] Automated test pass untuk test case utama

---

## ✅ **SETELAH VERIFIKASI BERHASIL**

Jika semua test pass, update status di:
- `BUG_FIX_STATUS.md` - Update status test menjadi ✅ PASS
- `BUG_FIX_REQUEST_ENGINEER.md` - Update status menjadi ✅ FIXED

---

**Catatan:** Pastikan server di-restart setiap kali ada perubahan kode agar perubahan diterapkan.
