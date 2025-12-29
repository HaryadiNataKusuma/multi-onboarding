# Cara Menjalankan Server - Multi-Onboarding Platform

## Prerequisites
- Go 1.24.0 atau lebih baru
- MySQL database running di `127.0.0.1:3306`
- Database: `travel_service_development`
- User: `root`
- Password: `Qoala123**`

## Menjalankan Server

```bash
# Install dependencies
go mod download

# Jalankan server
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Routes yang Tersedia

### Products
- `GET /api/products` - Get all products (backward compatibility, redirects to travel)
- `GET /api/products/travel` - Get travel products
- `GET /api/products/vehicle` - Get vehicle products

### Regions
- `GET /api/regions` - Get all regions

### Insurances
- `GET /api/insurances` - Get all insurances

### Dan lainnya...

## Troubleshooting

Jika mendapatkan error `ERR_CONNECTION_REFUSED`:
1. Pastikan server sudah berjalan dengan `go run main.go`
2. Pastikan port 8080 tidak digunakan aplikasi lain
3. Pastikan database MySQL berjalan dan bisa diakses
4. Cek log di console untuk error messages


