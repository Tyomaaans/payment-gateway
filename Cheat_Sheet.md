# DOKU Golang SDK — Integration Cheat Sheet

Panduan lengkap dan praktis untuk integrasi **DOKU Golang SDK** (SNAP Standard API) ke dalam aplikasi Anda.

---

## 1. Setup

### Kebutuhan Wajib
Pastikan parameter berikut sudah disiapkan sebelum melakukan inisialisasi:
* `Client ID`
* `Secret Key`
* `Private Key`
* `Public Key`
* `DOKU Public Key`
* `IsProduction` (boolean)

### Instalasi SDK
```bash
go get github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library@latest
```

### Inisialisasi Klien
```go
import "github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library"

snap := doku.Snap{
    PrivateKey:   privateKey,
    ClientId:     clientId,
    IsProduction: isProduction,
    SecretKey:    secretKey,
    Issuer:       issuer,
    PublicKey:    publicKey,
}
```

---

## 2. Virtual Account

### A. Create VA
Digunakan untuk membuat Virtual Account baru.
* **Function:** `snap.CreateVa(request)`
* **Data Utama:**
  * `partnerServiceId`
  * `customerNo`
  * `virtualAccountNo`
  * `virtualAccountName`
  * `trxId`
  * `totalAmount`
  * `additionalInfo.channel`
  * `virtualAccountTrxType`
* **Tipe VA:**
  * `C` → Closed Amount
  * `O` → Open Amount
* **Contoh Channel:** `VIRTUAL_ACCOUNT_BANK_CIMB`

### B. Update VA
Digunakan untuk mengubah data atau status Virtual Account.
* **Function:** `snap.UpdateVa(request)`
* *Catatan: Data yang digunakan mirip dengan Create VA.*

### C. Delete VA
Digunakan untuk menghapus Virtual Account.
* **Function:** `snap.DeletePaymentCode(request)`
* **Data Utama:**
  * `partnerServiceId`
  * `customerNo`
  * `virtualAccountNo`
  * `trxId`
  * `additionalInfo.channel`

### D. Check VA Status
Digunakan untuk mengecek status pembayaran atau VA.
* **Function:** `snap.CheckStatusVa(request)`
* **Data Utama:**
  * `partnerServiceId`
  * `customerNo`
  * `virtualAccountNo`
  * `inquiryRequestId` *(optional)*
  * `paymentRequestId` *(optional)*

---

## 3. Jenis Virtual Account

| Jenis | Deskripsi & Kegunaan |
| :--- | :--- |
| **DGPC** | • VA yang sudah disediakan/generate oleh DOKU.<br>• Cocok untuk transaksi satu kali (one-time). |
| **MGPC** | • VA yang dibuat oleh merchant.<br>• Cocok untuk sistem top up. |
| **DIPC** | • VA yang sudah didaftarkan di sistem merchant.<br>• DOKU melakukan inquiry ke server merchant ketika customer melakukan pembayaran.<br>• *Merchant wajib menyediakan endpoint inquiry.* |

---

## 4. Account Binding
Digunakan sebelum melakukan pembayaran menggunakan akun Direct Debit / E-Wallet.

### A. Account Binding
Digunakan untuk menghubungkan akun customer dengan merchant.
* **Supported Channels:** Allo Bank, CIMB, OVO
* **Function:** `snap.DoAccountBinding(request, deviceId, ipAddress)`
* **Data Penting:** `phoneNo`, `custIdMerchant`, `channel`, `successRegistrationUrl`, `failedRegistrationUrl`, `deviceModel`, `osType`, `channelId`
* *Catatan: Customer kemudian melakukan proses verifikasi/OTP sesuai channel.*

### B. Account Unbinding
Digunakan untuk melepas akun customer yang sudah di-bind.

**Urutan Alur Unbinding:**
```text
Account Binding → authCode → GetTokenB2B2C() → accessToken → DoAccountUnbinding()
```

**Function Terkait:**
```go
// 1. Dapatkan Access Token via Auth Code
token, err := snap.GetTokenB2B2C(authCode)

// 2. Lakukan Unbinding
response, err := snap.DoAccountUnbinding(request, ipAddress)
```

---

## 5. Card Registration
Digunakan untuk mendaftarkan kartu customer untuk kebutuhan Direct Debit.
* **Supported Bank:** BRI
* **Function:** `snap.DoCardRegistration(request, deviceId)`
* **Data Penting:** `bankCardNo`, `bankCardType`, `email`, `expiryDate`, `custIdMerchant`, `phoneNo`, `channel`, `successRegistrationUrl`, `failedRegistrationUrl`

### Card Unregistration
Digunakan untuk melepas kartu yang sudah terdaftar.

**Urutan Alur Unregistration:**
```text
Auth Code → GetTokenB2B2C() → Access Token → DoCardRegistrationUnbinding()
```

---

## 6. Direct Debit / E-Wallet Payment
Digunakan setelah customer berhasil melakukan binding atau card registration.
* **Function:** `snap.DoPayment(request, ipAddress, authCode)`

### Channel & Konfigurasi Khusus
| Channel | Channel Code DOKU | Data Tambahan / Spesifik |
| :--- | :--- | :--- |
| **Allo Bank** | `DIRECT_DEBIT_ALLO_SNAP` | `lineItems`, `payOptionDetails.payMethod` (`BALANCE`, `POINT`, `PAYLATER`) |
| **BRI** | `DIRECT_DEBIT_BRI_SNAP` | — |
| **CIMB** | `DIRECT_DEBIT_CIMB_SNAP` | — |
| **OVO** | `EMONEY_OVO_SNAP` | `feeType`, `payOptionDetails.payMethod` (`CASH`, `POINTS`), `paymentType` (`SALE`, `RECURRING`) |

**Data Umum Pembayaran:**
* `partnerReferenceNo`
* `amount.value` & `amount.currency`
* `additionalInfo.channel` & `additionalInfo.remarks`
* `successPaymentUrl` & `failedPaymentUrl`

---

## 7. Payment Jump App
Digunakan untuk pembayaran yang membutuhkan redirect atau deep link ke aplikasi mobile customer.
* **Supported Channels:** DANA, ShopeePay
* **Function:** `snap.DoPaymentJumpApp(request, deviceId, ipAddress)`
* **Data Penting:** `partnerReferenceNo`, `amount`, `urlParam`, `additionalInfo.channel`, `validUpto`, `pointOfInitiation`

### Channel Khusus
* **DANA (`EMONEY_DANA_SNAP`):** Tambahan `orderTitle`, `supportDeepLinkCheckoutUrl`.
* **ShopeePay (`EMONEY_SHOPEE_PAY_SNAP`):** Tambahan `metadata`.

---

## 8. Check Transaction Status
Digunakan untuk mengecek status transaksi pembayaran secara berkala.
* **Function:** `snap.DoCheckStatus(request)`

---

## 9. Refund
Digunakan untuk melakukan pengembalian dana transaksi.
* **Function:** 
```go
snap.DoRefund(request, ipAddress, authCode, deviceId)
```

---

## 10. Balance Inquiry
Digunakan untuk mengecek saldo akun atau payment source.
* **Function:** 
```go
snap.DoBalanceInquiry(request, deviceId, ipAddress, authCode)
```

---

## 11. Error Handling

Setiap pemanggilan fungsi pada SDK wajib mengecek nilai `err`.

```go
response, err := snap.CreateVa(request)
if err != nil {
    // Handle error di sini
}
```

### Error Umum
| Response Code | Arti / Keterangan |
| :---: | :--- |
| `4010000` | Unauthorized (Kesalahan kredensial / token invalid) |
| `4012400` | Virtual Account tidak ditemukan |
| `2002400` | Success / Berhasil |

---

## 12. Flow Implementasi

### A. Virtual Account Flow
```text
[Create VA] 
     ↓
[Customer Melakukan Pembayaran] 
     ↓
[Check VA Status]
```

### B. Direct Debit / E-Wallet Flow
```text
[Account Binding / Card Registration]
     ↓
[Auth Code Didapatkan]
     ↓
[DoPayment()]
     ↓
[Check Status Transaksi]
```

### C. Unbinding Flow
```text
[Auth Code] 
     ↓
[GetTokenB2B2C()] 
     ↓
[Access Token] 
     ↓
[DoAccountUnbinding / Card Unregistration]
```

### D. Refund Flow
```text
[Payment Success] 
     ↓
[DoRefund()]
```
