# 🚀 URL Guide - Multi-Onboarding Platform

## 📍 Base URLs

### Backend API
```
http://localhost:8080
```

### Frontend (jika ada)
```
http://localhost:5173
```

---

## 🔧 Cara Menjalankan

### 1. Backend Server

**Jalankan server:**
```bash
# Option 1: Menggunakan script (recommended)
./start_server.sh

# Option 2: Manual
go run main.go
```

**Server akan berjalan di:** `http://localhost:8080`

### 2. Frontend (jika ada)

**Jalankan frontend:**
```bash
cd frontend
npm install  # jika belum install dependencies
npm run dev
```

**Frontend akan berjalan di:** `http://localhost:5173`

---

## 📋 API Endpoints - Travel Domain

### Products
- `GET    /api/products` - Get all products (backward compatibility)
- `GET    /api/products/travel` - Get travel products
- `POST   /api/products` - Create product
- `GET    /api/products/:id` - Get product by ID
- `PUT    /api/products/:id` - Update product
- `DELETE /api/products/:id` - Delete product
- `GET    /api/products/draft` - Get all draft products
- `POST   /api/products/draft/add` - Add draft product
- `POST   /api/products/draft/clear` - Clear all drafts
- `POST   /api/products/draft/confirm` - Confirm drafts
- `DELETE /api/products/draft/:id` - Delete draft

### Regions
- `GET    /api/regions` - Get all regions
- `POST   /api/regions` - Create region
- `GET    /api/regions/:id` - Get region by ID
- `PUT    /api/regions/:id` - Update region
- `DELETE /api/regions/:id` - Delete region

### Insurances
- `GET    /api/insurances` - Get all insurances
- `POST   /api/insurances` - Create insurance
- `GET    /api/insurances/:id` - Get insurance by ID
- `PUT    /api/insurances/:id` - Update insurance
- `DELETE /api/insurances/:id` - Delete insurance
- `GET    /api/insurances/draft` - Get all draft insurances
- `POST   /api/insurances/draft/add` - Add draft insurance
- `GET    /api/insurances/draft/:id` - Get draft by ID
- `PUT    /api/insurances/draft/:id` - Update draft
- `DELETE /api/insurances/draft/:id` - Delete draft
- `POST   /api/insurances/draft/clear` - Clear all drafts
- `POST   /api/insurances/draft/confirm` - Confirm drafts

### Countries
- `GET    /api/countries` - Get all countries
- `GET    /api/countries/:id` - Get country by ID
- `GET    /api/countries/code/:code` - Get country by code

### Commissions
- `GET    /api/commissions` - Get all commissions
- `POST   /api/commissions` - Create commission
- `GET    /api/commissions/:id` - Get commission by ID
- `PUT    /api/commissions/:id` - Update commission
- `DELETE /api/commissions/:id` - Delete commission
- `GET    /api/commissions/draft` - Get all draft commissions
- `POST   /api/commissions/draft/add` - Add draft commission
- `POST   /api/commissions/draft/clear` - Clear all drafts
- `POST   /api/commissions/draft/confirm` - Confirm drafts

### Addons
- `GET    /api/addons` - Get all addons
- `POST   /api/addons` - Create addon
- `GET    /api/addons/:id` - Get addon by ID
- `PUT    /api/addons/:id` - Update addon
- `DELETE /api/addons/:id` - Delete addon

### Templates
- `GET    /api/templates` - Get all templates
- `POST   /api/templates` - Create template
- `GET    /api/templates/:locale/:id` - Get template by locale and ID
- `PUT    /api/templates/:locale/:id` - Update template
- `DELETE /api/templates/:locale/:id` - Delete template
- `GET    /api/templates/draft` - Get all draft templates
- `POST   /api/templates/draft/add` - Add draft template
- `PUT    /api/templates/draft/:id` - Update draft template
- `DELETE /api/templates/draft/:id` - Delete draft template
- `POST   /api/templates/draft/clear` - Clear all drafts
- `POST   /api/templates/draft/confirm` - Confirm drafts

### Histories
- `GET    /api/histories` - Get all histories
- `GET    /api/histories/:id` - Get history by ID
- `GET    /api/histories/table/:tableName` - Get histories by table name
- `GET    /api/histories/record/:recordId` - Get histories by record ID

---

## 📋 API Endpoints - Vehicle Domain

### Products
- `GET    /api/products/vehicle` - Get vehicle products
- `POST   /api/products/vehicle` - Create vehicle product
- `GET    /api/products/vehicle/:id` - Get vehicle product by ID
- `PUT    /api/products/vehicle/:id` - Update vehicle product
- `DELETE /api/products/vehicle/:id` - Delete vehicle product

### Insurances
- `GET    /api/vehicle/insurances` - Get all vehicle insurances
- `POST   /api/vehicle/insurances/draft/add` - Add draft insurance
- `GET    /api/vehicle/insurances/draft` - Get all draft insurances
- `PUT    /api/vehicle/insurances/draft/:id` - Update draft
- `DELETE /api/vehicle/insurances/draft/:id` - Delete draft
- `POST   /api/vehicle/insurances/draft/clear` - Clear drafts
- `POST   /api/vehicle/insurances/draft/confirm` - Confirm drafts
- `GET    /api/vehicle/insurances` - Get confirmed insurances
- `PUT    /api/vehicle/insurances/:id` - Update confirmed insurance

### Dan lainnya...
Lihat `main.go` untuk daftar lengkap endpoint vehicle domain.

---

## 🧪 Testing dengan cURL

### Test Regions
```bash
curl http://localhost:8080/api/regions
```

### Test Insurances
```bash
curl http://localhost:8080/api/insurances
```

### Test Products
```bash
curl http://localhost:8080/api/products
```

### Test Countries
```bash
curl http://localhost:8080/api/countries
```

---

## 🌐 Akses via Browser

Buka browser dan akses:
- **Regions:** http://localhost:8080/api/regions
- **Insurances:** http://localhost:8080/api/insurances
- **Products:** http://localhost:8080/api/products
- **Countries:** http://localhost:8080/api/countries

---

## 📝 Catatan

1. **Pastikan MySQL berjalan** di `127.0.0.1:3306`
2. **Database:** `travel_service_development`
3. **User:** `root`
4. **Password:** `Qoala123**`
5. **Port Backend:** `8080`
6. **Port Frontend:** `5173` (jika ada)

---

## 🛠️ Troubleshooting

### Server tidak bisa start
```bash
# Cek apakah port 8080 sudah digunakan
lsof -i:8080

# Kill process yang menggunakan port 8080
lsof -ti:8080 | xargs kill -9

# Atau gunakan script
./stop_server.sh
./start_server.sh
```

### Cek status server
```bash
./check_server.sh
```

---

## 📞 Support

Jika ada masalah, cek:
1. Log server di console
2. File `server.log` (jika menggunakan script)
3. Database connection di `main.go`


