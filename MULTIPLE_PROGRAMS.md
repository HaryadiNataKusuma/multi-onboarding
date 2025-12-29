# Multiple Onboarding Programs Architecture

## Overview

`main.go` sekarang mendukung multiple program onboarding dalam satu server:
- **Vehicle** (Motor Vehicle) - `/api/vehicle/*`
- **Travel** - `/api/travel/*`
- **Property** (Future) - `/api/property/*`
- **Health** (Future) - `/api/health/*`

## Configuration

### Default Program

Default program ditentukan via environment variable:
```bash
export DEFAULT_ONBOARDING_PROGRAM=vehicle  # atau "travel", "property", "health"
```

Jika tidak di-set, default adalah `vehicle`.

### Program Status

Program bisa di-enable/disable di function `getProgramsConfig()`:

```go
Programs: map[string]ProgramConfig{
    "vehicle": {
        Enabled: true,  // Set false untuk disable
        Name:    "Motor Vehicle",
        Prefix:  "vehicle",
    },
    "travel": {
        Enabled: true,
        Name:    "Travel",
        Prefix:  "travel",
    },
    // ...
}
```

## Route Structure

### Vehicle Program
- **Main routes**: `/api/vehicle/*`
- **Backward compatibility** (jika vehicle adalah default): `/api/*` → `/api/vehicle/*`
- **Endpoints**:
  - `/api/vehicle/insurances`
  - `/api/vehicle/products`
  - `/api/vehicle/addons`
  - `/api/vehicle/templates`
  - `/api/vehicle/commissions`
  - `/api/vehicle/histories`
  - dll.

### Travel Program
- **Main routes**: `/api/travel/*`
- **Backward compatibility** (jika travel adalah default): `/api/*` → `/api/travel/*`
- **Endpoints**:
  - `/api/travel/insurances`
  - `/api/travel/products`
  - `/api/travel/addons`
  - `/api/travel/templates`
  - `/api/travel/commissions`
  - `/api/travel/histories`
  - `/api/travel/regions`
  - `/api/travel/countries`
  - dll.

## Adding New Programs

### 1. Create Domain Structure
```
domain/
  └── [program-name]/
      ├── insurance/
      ├── product/
      ├── addon/
      └── ...
```

### 2. Add Program Config
Di `getProgramsConfig()`, tambahkan:
```go
"[program-name]": {
    Enabled: true,
    Name:    "Program Name",
    Prefix:  "program-name",
},
```

### 3. Create Setup Function
Buat function `setup[Program]Program()`:
```go
func setupPropertyProgram(e *echo.Echo, config ProgramsConfig) {
    // Initialize repositories, usecases, handlers
    // Setup routes under /api/property/*
    // Setup backward compatibility if needed
}
```

### 4. Call Setup in main()
```go
if programsConfig.Programs["property"].Enabled {
    log.Println("🏠 Initializing Property onboarding program...")
    setupPropertyProgram(e, programsConfig)
}
```

## Backward Compatibility

Untuk setiap program, jika program tersebut adalah default program, route `/api/*` akan di-map ke `/api/[program]/*`.

Contoh:
- Jika `DEFAULT_ONBOARDING_PROGRAM=vehicle`:
  - `/api/insurances` → `/api/vehicle/insurances`
  - `/api/products` → `/api/vehicle/products`

- Jika `DEFAULT_ONBOARDING_PROGRAM=travel`:
  - `/api/insurances` → `/api/travel/insurances`
  - `/api/products` → `/api/travel/products`

## Legacy Routes

Route legacy yang shared across programs tetap di `/api/*`:
- `/api/product-rules/*`
- `/api/addon-rules/*`
- `/api/addon-rule-details/*`
- `/api/commissions/generate`
- `/api/all-result/download`
- `/api/html/convert`

## Example Usage

### Start with Vehicle as default:
```bash
export DEFAULT_ONBOARDING_PROGRAM=vehicle
go run main.go
```

### Start with Travel as default:
```bash
export DEFAULT_ONBOARDING_PROGRAM=travel
go run main.go
```

### Access endpoints:
```bash
# Vehicle endpoints
curl http://localhost:8080/api/vehicle/insurances
curl http://localhost:8080/api/insurances  # If vehicle is default

# Travel endpoints
curl http://localhost:8080/api/travel/insurances
curl http://localhost:8080/api/insurances  # If travel is default
```

## Benefits

1. **Single Server**: Semua program onboarding dalam satu server
2. **Isolated Routes**: Setiap program punya route prefix sendiri
3. **Easy Scaling**: Tambah program baru tanpa mengganggu yang existing
4. **Backward Compatible**: Frontend existing tetap bisa pakai `/api/*`
5. **Configurable**: Enable/disable program via config
6. **Clean Architecture**: Setiap program punya domain structure sendiri

