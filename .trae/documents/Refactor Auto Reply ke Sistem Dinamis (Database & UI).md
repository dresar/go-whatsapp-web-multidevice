Berikut adalah rencana implementasi untuk fitur "Dynamic Auto Reply" yang terintegrasi dengan database dan UI.

### **Rencana Implementasi**

#### **1. Backend: Database & Model**
Kita akan memperluas `SQLiteRepository` yang sudah ada untuk menyimpan data Auto Reply.

*   **Model**: Membuat struct `AutoReply` di `domains/chatstorage` dengan field: `ID`, `Keyword`, `Response`, `MatchType` (Exact/Contains).
*   **Migrasi**: Menambahkan query `CREATE TABLE auto_replies` ke dalam fungsi `getMigrations()` di `src/infrastructure/chatstorage/sqlite_repository.go`.
*   **Repository**: Menambahkan method CRUD (`StoreAutoReply`, `GetAutoReplies`, `DeleteAutoReply`) di `SQLiteRepository`.

#### **2. Backend: API Handler**
Membuat endpoint REST API agar frontend bisa mengelola data.

*   **Handler**: Membuat file baru `src/ui/rest/autoreply.go` (atau lokasi serupa) untuk menangani request:
    *   `GET /api/v1/autoreply`
    *   `POST /api/v1/autoreply`
    *   `DELETE /api/v1/autoreply/:id`
*   **Routing**: Mendaftarkan route baru ini ke sistem routing utama (kemungkinan di `src/cmd/rest.go` atau integrasi manual di `root.go` jika diperlukan).

#### **3. Logika Bot: WhatsApp Handler**
Mengubah logika statis menjadi dinamis.

*   **Modifikasi**: `src/infrastructure/whatsapp/auto_reply.go`.
*   **Logika Baru**:
    1.  Hapus pengecekan `config.WhatsappAutoReplyMessage`.
    2.  Ambil semua rule auto reply dari database via `chatStorageRepo`.
    3.  Looping rule:
        *   Jika `MatchType == "Exact"`: Cek kesamaan persis (`==`).
        *   Jika `MatchType == "Contains"`: Cek apakah mengandung kata (`strings.Contains`).
    4.  Jika cocok, kirim balasan.

#### **4. Frontend: User Interface**
Menambahkan menu manajemen di dashboard.

*   **Component**: Membuat file `src/views/components/AutoReplyManager.js`.
    *   Fitur: Tabel daftar keyword, Form tambah baru, Tombol hapus.
*   **Integrasi**: Mengedit `src/views/index.html`:
    *   Import component `AutoReplyManager.js`.
    *   Registrasi component ke Vue instance.
    *   Menambahkan section "Auto Reply" di layout utama (HTML).

Apakah Anda setuju dengan rencana ini? Silakan konfirmasi untuk memulai eksekusi.