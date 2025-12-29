# 🔄 Auto Restart Server

## Script yang Tersedia

### 1. `restart_server.sh` - Restart Server (Recommended)
```bash
./restart_server.sh
```
Script ini akan:
- Stop server yang sedang berjalan
- Start server baru
- Menampilkan log

### 2. `stop_server.sh` - Stop Server
```bash
./stop_server.sh
```

### 3. `start_server.sh` - Start Server
```bash
./start_server.sh
```

### 4. `check_server.sh` - Check Server Status
```bash
./check_server.sh
```

---

## ⚠️ Catatan Penting

**Setiap kali ada perbaikan coding, server HARUS di-restart untuk menerapkan perubahan!**

### Cara Restart:
```bash
./restart_server.sh
```

### Atau Manual:
```bash
# Stop
./stop_server.sh

# Start
./start_server.sh
```

---

## 📋 Checklist Setelah Perbaikan

- [ ] Code sudah diperbaiki
- [ ] Linter errors sudah diperbaiki
- [ ] Server sudah di-restart
- [ ] Test endpoint untuk memastikan perubahan bekerja

---

## 🐛 Troubleshooting

Jika server tidak restart:
1. Cek apakah port 8080 bebas: `lsof -i:8080`
2. Kill process manual: `lsof -ti:8080 | xargs kill -9`
3. Coba restart lagi: `./restart_server.sh`


