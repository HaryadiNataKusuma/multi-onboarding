# Multi-Onboarding Platform

A modern multi-program onboarding platform built with Vue.js (frontend) and Go (backend). Supports Vehicle, Travel, Property, and Health insurance product onboarding.

## Features

- ✨ Create, read, update, and delete travel products
- 🎨 Beautiful and modern UI with Orange Theme
- 🚀 Fast and responsive
- 📱 Mobile-friendly design
- 🔄 Real-time updates

## Tech Stack

- **Frontend**: Vue.js 3 + Vite
- **Backend**: Go (Golang) with Gorilla Mux
- **Styling**: Modern CSS with Orange Theme

## Project Structure

```
Multi-Onboarding/
├── backend/
│   ├── delivery/
│   │   └── http/
│   │       └── handler.go          # HTTP handlers (delivery layer)
│   ├── repository/
│   │   ├── product_repository.go   # Data access layer
│   │   └── errors.go               # Repository errors
│   ├── usecase/
│   │   ├── product_usecase.go      # Business logic layer
│   │   └── errors.go               # Usecase errors
│   ├── model/
│   │   └── product.go              # Domain models/entities
│   ├── main.go                     # Application entry point
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ProductForm.vue
│   │   │   └── ProductList.vue
│   │   ├── services/
│   │   │   └── api.js
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
└── README.md
```

## Setup Instructions

### Backend Setup

1. Navigate to the backend directory:
```bash
cd Multi-Onboarding
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run main.go
```

The backend will run on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd Multi-Onboarding/frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run the development server:
```bash
npm run dev
```

The frontend will run on `http://localhost:5173`

## API Endpoints

- `GET /api/products` - Get all products
- `GET /api/products/{id}` - Get a specific product
- `POST /api/products` - Create a new product
- `PUT /api/products/{id}` - Update a product
- `DELETE /api/products/{id}` - Delete a product

## Product Model

```json
{
  "id": 1,
  "name": "Bali Paradise Package",
  "description": "Experience the beauty of Bali",
  "destination": "Bali, Indonesia",
  "price": 1299.99,
  "duration": 7,
  "category": "beach",
  "status": "draft",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Categories

- adventure
- beach
- cultural
- luxury
- family
- honeymoon

## Status Options

- draft
- published
- archived

## Development

### Backend Architecture
- **Clean Architecture** with layered structure:
  - **Delivery Layer** (`delivery/http`): HTTP handlers for API endpoints
  - **Usecase Layer** (`usecase`): Business logic and validation
  - **Repository Layer** (`repository`): Data access abstraction
  - **Model Layer** (`model`): Domain entities
- Uses in-memory storage (replace with database for production)
- CORS enabled for local development
- RESTful API design
- Proper error handling and validation

### Frontend
- Vue 3 Composition API
- Axios for HTTP requests
- Responsive design
- Modern UI/UX with Orange Theme

## Production Deployment

For production:
1. Replace in-memory storage with a database (PostgreSQL, MySQL, etc.)
2. Add authentication and authorization
3. Add input validation and sanitization
4. Add error handling and logging
5. Build frontend: `npm run build`
6. Deploy backend and frontend to your hosting service

## License

MIT

