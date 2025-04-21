# MedSched (Sistem Manajemen Jadwal Dokter dan Janji Kontrol Pasien)
MedSched adalah aplikasi web berbasis REST API untuk mengelola jadwal dokter dan janji kontrol pasien. Aplikasi ini memungkinkan admin untuk mengelola data dokter, jadwal, pengecualian jadwal (jadwal libur / cuti), serta data pengguna (pasien dan admin), sementara pasien bisa membuat janji kontrol dengan dokter yang tersedia.

## Fitur Utama

1. Pendaftaran dan Login Pengguna
    - Pasien dapat mendaftar dan login ke aplikasi untuk membuat janji kontrol.
    - Admin dapat login dan mengelola data pengguna lainnya.

2. Manajemen Dokter
    - Admin dapat menambah, mengubah, dan menghapus data dokter.
    - Admin dapat melihat daftar semua dokter atau mencari dokter berdasarkan ID.

3. Manajemen Jadwal Dokter
    - Admin dapat membuat, mengubah, menghapus, dan melihat jadwal dokter.
    - Setiap jadwal menyertakan informasi seperti hari, jam praktek, dan jumlah pasien yang bisa diterima.

4. Pengecualian Jadwal
    - Admin dapat menambah atau menghapus pengecualian jadwal (Misalkan ada hari libur atau dokter cuti)

5. Manajemen Janji Kontrol
    - Pasien dapat membuat janji kontrol dengan memilih jadwal dokter yang tersedia.
    - Janji kontrol akan diperiksa untuk memastikan apakah slot tersedia dan apakah tanggalnya valid.
    - Pengingat email dikirim kepada pasien sehari sebelum jadwal kontrol.

## Arsitektur Aplikasi

Aplikasi ini menggunakan Go (Golang) sebagai bahasa pemrograman, dengan framework Gin-Gonic untuk routing, dan menggunakan PostgreSQL untuk penyimpanan data. Semua data seperti pengguna, dokter, jadwal, dan janji kontrol disimpan dalam database. Untuk pengingat jadwal dilakukan via pengiriman pesan lewat email menggunakan GoMail

## API Endpoints
Berikut adalah daftar API endpoints yang tersedia dalam sistem:

### Akses Publik

- POST /register: Mendaftar pengguna baru (admin atau pasien).
- POST /login: Login pengguna dan mendapatkan token JWT.
- GET /doctor_schedule: Melihat jadwal dokter yang sudah ada

### Akses User

**Manajemen User oleh User**
- GET /users/me: Melihat profil pengguna yang sedang login.
- PUT /users/me: Memperbarui profil pengguna yang sedang login.

### Akses Admin

**Manajemen User oleh Admin**
- GET /users: Melihat daftar user.
- GET /users/by: Melihat informasi user berdasarkan username / ID.
- PUT /users/by: Memperbarui informasi user berdasarkan username / ID.
- DELETE /users/by: Menghapus informasi user berdasarkan username / ID.

**Dokter (Doctors)**
- GET /doctors: Melihat daftar dokter.
- GET /doctors/by: Mencari dokter berdasarkan ID.
- POST /doctors: Membuat dokter baru.
- PUT /doctors/by: Memperbarui data dokter berdasarkan ID.
- DELETE /doctors/by: Menghapus dokter berdasarkan ID.

**Jadwal (Schedules)**
- GET /schedules: Melihat semua jadwal dokter.
- GET /schedules/by: Mencari jadwal dokter berdasarkan ID.
- POST /schedules: Membuat jadwal dokter baru.
- PUT /schedules/by: Memperbarui jadwal dokter berdasarkan ID.
- DELETE /schedules/by: Menghapus jadwal dokter berdasarkan ID.

**Pengecualian Jadwal (Schedule Exceptions)**
- POST /schedule_exceptions: Menambah pengecualian jadwal (admin).
- GET /schedule_exceptions: Melihat semua pengecualian jadwal.
- DELETE /schedule_exceptions/by: Menghapus pengecualian jadwal berdasarkan ID.

### Akses User dan Admin
**Janji Kontrol (Appointments)**
- POST /appointments: Membuat janji kontrol untuk pasien.
- GET /appointments: Membuat janji kontrol untuk pasien. (Admin)
- GET /appointments/date: Membuat janji kontrol untuk pasien. (Admin)

### **[Swagger]**(https://medsched-production.up.railway.app/swagger/index.html)