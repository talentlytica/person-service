# 📋 Report Fix Issue — Summary for Engineer

**Lokasi:** `Test/Test status/`  
**Tanggal:** 29 Januari 2026  
**Status test:** 198 passed, **21 failed**, 219 total

---

## 1. Ringkasan Eksekutif

Ada **21 test gagal** yang perlu ditangani. Penyebab terbagi dua:

| Kategori | Jumlah | Tindakan |
|----------|--------|----------|
| **Backend** — API mengembalikan status code salah | **11** | Perbaiki handler (return 404 untuk resource not found; tangani Unicode) |
| **Test** — Assertion atau penanganan response salah | **10** | Perbaiki test (accept 404, perbaiki `toContain`, handle error) |

**File di folder ini untuk engineer:**
- **Failed_testcase.md** — Daftar 21 test gagal + failure message
- **Passed_testcase.md** — Daftar 198 test lulus (referensi)
- **REPORT_FIX_ISSUE.md** — Dokumen ini (ringkasan + action items)

---

## 2. Backend — Yang Perlu Diperbaiki di API (11 isu)

### 2.1 Return 404 (bukan 400) untuk resource not found

**Lokasi kode:** `source/app/person_attributes/person_attributes.go` (dan endpoint terkait)

Saat **personId** atau **attributeId** tidak ada di database, API saat ini mengembalikan **400**. Harusnya **404 Not Found**.

| # | Endpoint | Test case | Expected | Actual | File test |
|---|----------|-----------|----------|--------|-----------|
| 1 | PUT `/persons/:personId/attributes/:attributeId` | Non-existent personId | 404 | 400 | PUT_persons_personId_attributes_attributeId_negative.test.js |
| 2 | PUT `/persons/:personId/attributes/:attributeId` | Non-existent attributeId | 404 | 400 | (same) |
| 3 | DELETE `/persons/:personId/attributes/:attributeId` | Non-existent personId | 404 | 400 | DELETE_persons_personId_attributes_attributeId_negative.test.js |
| 4 | DELETE (same) | Non-existent attributeId | 404 | 400 | (same) |
| 5 | DELETE (same) | Already deleted attribute (idempotent) | 404 | 400 | (same) |
| 6 | DELETE (same) | Ignore request body on DELETE | 404 | 400 | (same) |
| 7 | DELETE (same) | Multiple concurrent DELETEs | 404 | 400 | (same) |
| 8 | DELETE (same) | Invalid Content-Type header | 404 | 400 | (same) |
| 9 | GET `/persons/:personId/attributes/:attributeId` | Non-existent personId | 404 | 400 | GET_persons_personId_attributes_attributeId_negative.test.js |
| 10 | GET (same) | Non-existent attributeId | 404 | 400 | (same) |

**Action:** Di handler **UpdateAttribute**, **DeleteAttribute**, dan **GetAttribute**, pastikan:
- Jika person tidak ditemukan → `return 404` (bukan 400).
- Jika attribute tidak ditemukan → `return 404` (bukan 400).
- Gunakan pengecekan yang sama seperti di **CreateAttribute** / **GetAllAttributes** (mis. `pgx.ErrNoRows` → 404).

### 2.2 Unicode / special characters — 500 Internal Server Error

| # | Test case | Expected | Actual | File test |
|---|-----------|----------|--------|-----------|
| 11 | SPEC: Should handle Unicode and special characters | 200/201 atau 400 | **500** | person_attributes_security_spec.test.js |

**Action:** Di handler person attributes (create/update), tangani input Unicode/special characters tanpa panic. Pastikan encoding/decoding dan query ke DB aman; untuk input invalid kembalikan **400**, bukan **500**.

---

## 3. Test — Yang Perlu Diperbaiki di Sisi Test (10 isu)

### 3.1 GET `/persons/:personId/attributes` — Terima 404 sebagai valid (4 isu)

API **sudah benar** mengembalikan **404** untuk person tidak ada. Test tidak menangani 404 (tidak pakai try/catch / tidak treat 404 sebagai success).

| # | Test case | Masalah | File |
|---|-----------|--------|------|
| 12 | Should return empty or 404 for non-existent personId | Test expects 200 or 404; request throws karena 404 | GET_persons_personId_attributes_negative.test.js |
| 13 | Should reject request with body (GET should not have body) | Same — personId non-existent → 404, test throws | (same) |
| 14 | Should handle invalid Accept header | Same | (same) |
| 15 | Should ignore invalid query parameters | Same | (same) |

**Action:** Untuk request ke **non-existent personId** (mis. `550e8400-e29b-41d4-a716-446655440000`), wrap request di try/catch dan treat **404** sebagai **valid response**. Atau gunakan **personId yang ada** untuk test body/Accept/query, sehingga response 200 dan assertion bisa tetap dicek.

### 3.2 Assertion `toContain` terbalik / salah (3 isu)

Test menulis `expect([400, 413]).toContain(error.response.status)` dengan maksud “status harus salah satu dari 400, 413”. Di Jest, `toContain` pada array artinya “array mengandung value”, jadi yang benar: **expect([400, 413]).toContain(error.response.status)**. Error message “Expected value: 500, Received array: [400, 413]” menunjukkan ada test yang menulis terbalik atau expect 500.

| # | Test case | Masalah | File |
|---|-----------|--------|------|
| 16 | Should handle or reject extremely long key (10KB) | Assertion: expect array to contain value; fix to `expect([400, 413]).toContain(error.response.status)` | POST_api_key-value_negative.test.js |
| 17 | Should handle or reject extremely long value (1MB) | Same | (same) |
| 18 | Should handle very long key in URL | Same — expect `[400, 404, 414].toContain(status)` | GET_api_key-value_key_negative.test.js |

**Action:** Cek baris assertion di file di atas. Pastikan:  
`expect([400, 404, 413, 414]).toContain(error.response.status)` (atau set status yang dianggap valid). Jangan expect **500** jika API memang mengembalikan 400/413/414.

### 3.3 Test bug — akses `error.response` undefined (1 isu)

| # | Test case | Masalah | File |
|---|-----------|--------|------|
| 19 | Should reject key with only whitespace | `Cannot read properties of undefined (reading 'status')` — `error.response` undefined | POST_api_key-value_negative.test.js |

**Action:** Sebelum akses `error.response.status`, pastikan `error.response` ada (mis. `if (error.response) { expect(...).toContain(error.response.status); }` atau expect 4xx bila request gagal tanpa response body).

### 3.4 GET with body / Invalid Accept — expectation vs API (2 isu)

| # | Test case | Masalah | File |
|---|-----------|--------|------|
| 20 | Should ignore request body on GET | API returns **400**; test mungkin expect 200/404. Putuskan: API boleh 400 untuk GET with body, lalu test harus expect 400; atau API diubah ignore body dan return 200/404. | GET_persons_personId_attributes_attributeId_negative.test.js |
| 21 | Should handle invalid Accept header | Test expects `[404, 406].toContain(400)` → fail. Pastikan assertion: either accept 400 atau expect `expect([400, 404, 406]).toContain(error.response.status)`. | (same) |

**Action:** Tentukan perilaku yang diinginkan (400 vs 404/406), lalu sesuaikan assertion di test.

---

## 4. Checklist untuk Engineer

**Backend (source/app):**
- [ ] PUT `/persons/:personId/attributes/:attributeId` — return **404** jika person atau attribute tidak ada (bukan 400).
- [ ] DELETE `/persons/:personId/attributes/:attributeId` — return **404** untuk not found / already deleted (bukan 400).
- [ ] GET `/persons/:personId/attributes/:attributeId` — return **404** untuk not found (bukan 400).
- [ ] Person attributes — handle Unicode/special characters tanpa 500; return 400 untuk invalid input.

**Test (Test/API_Tests/):**
- [ ] GET_persons_personId_attributes_negative.test.js — untuk non-existent personId, terima **404** (wrap in try/catch atau gunakan existing personId untuk skenario body/Accept/query).
- [ ] POST_api_key-value_negative.test.js — perbaiki assertion toContain untuk extremely long key/value; perbaiki test “whitespace-only key” (cek `error.response`).
- [ ] GET_api_key-value_key_negative.test.js — perbaiki assertion untuk very long key (toContain dengan status yang benar).
- [ ] GET_persons_personId_attributes_attributeId_negative.test.js — sesuaikan assertion “ignore request body” dan “invalid Accept” dengan perilaku API yang dipilih.

---

## 5. Verifikasi Setelah Fix

1. Jalankan semua test:
   ```bash
   cd Test
   npm run test:all
   ```
2. Regenerate status:
   ```bash
   node extract-test-results.js
   ```
3. Cek **Test status/**:  
   - **Failed_testcase.md** harus kosong atau berkurang.  
   - **Passed_testcase.md** bertambah sesuai test yang sekarang lulus.

---

## 6. Referensi File di Folder `Test/Test status/`

| File | Isi |
|------|-----|
| **Failed_testcase.md** | Daftar 21 test gagal + failure message per test |
| **Passed_testcase.md** | Daftar 198 test lulus |
| **REPORT_FIX_ISSUE.md** | Ringkasan ini + action items untuk engineer |

**Path test files (relatif dari repo):**  
`Test/API_Tests/negative_tests/` dan `Test/API_Tests/specification_tests/`
