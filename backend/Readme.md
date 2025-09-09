# Mess Registration Backend

A production-level REST API for mess registration system built with Go, Gin, GORM, and PostgreSQL.

## Features

- **Student Registration**: Register students to different mess halls (0-3)
- **User Management**: Complete user profile management
- **Swap System**: Request and manage mess swaps between mess A and B
- **Authentication**: JWT-based authentication with OAuth support
- **Role-based Access**: Different access levels for students, mess staff, and hostel office
- **API Key Protection**: Secure endpoints for mess scanning
- **Auto-matching**: Automatic swap matching for public requests

## Tech Stack

- **Backend**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens
- **Containerization**: Docker & Docker Compose
- **Hot Reload**: Air for development

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd mess-registration/backend
   ```

2. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start with Docker (Recommended)**
   ```bash
   make setup
   ```
   This will:
   - Start PostgreSQL database
   - Start the application
   - Set up hot reload for development

4. **Or run locally**
   ```bash
   # Start database only
   docker-compose up -d db

   # Install dependencies
   make deps

   # Run the application
   make run
   ```

### Development

```bash
# Start development server with hot reload
make dev

# View logs
make docker-logs

# Reset everything
make quick-dev
```

## API Endpoints

### Students

- `POST /api/register?mess=0` - Register student
- `GET /api/user` - Get user profile (JWT required)

### Swap System

- `POST /api/swap-request` - Create swap request (JWT required)
- `DELETE /api/swap-request` - Delete swap request (JWT required)
- `GET /api/get-swaps` - Get available swaps (JWT required)

### Mess Staff

- `GET /api/scanning?rollno=xxx` - Scan student (API key required)

### Hostel Office

- `GET /api/admin/user?rollno=xxx` - Get user by roll number
- `POST /api/admin/update?rollno=xxx` - Update user
- `GET /api/admin/users` - Get all users
- `DELETE /api/admin/reset` - Reset all data

## Authentication

### JWT Token
```bash
Authorization: Bearer <jwt_token>
```

### API Key (for mess staff)
```bash
X-API-Key: <api_key>
```

## Database Schema

### Users Table
```sql
- id (UUID, Primary Key)
- user_id (String, Unique)
- name (String)
- email (String, Unique)
- roll_no (String, Unique)
- user_type (String) - student, admin, mess_staff
- veg_type (String) - veg, non-veg
- mess (Integer) - 0, 1, 2, 3
- created_at, updated_at
```

### Swap Requests Table
```sql
- id (UUID, Primary Key)
- requester_id (UUID, Foreign Key)
- name (String)
- email (String)
- type (String) - friend, public
- status (String) - pending, approved, rejected
- created_at, updated_at
```

## API Response Format

All endpoints return responses in this format:

```json
{
  "message": "Success message",
  "error": "Error message (if any)",
  "data": {
    // Response data
  }
}
```

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Development Commands

```bash
# Build the application
make build

# Run tests
make test

# Format code
make fmt

# Clean artifacts
make clean

# View all commands
make help
```

## Docker Commands

```bash
# Start all services
make docker-up

# Stop all services
make docker-down

# View logs
make docker-logs

# Reset database
make db-reset
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `mess_registration` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `JWT_SECRET` | JWT signing secret | (change in production) |
| `GIN_MODE` | Gin mode | `debug` |
| `MESS_API_KEY` | API key for mess staff | (change in production) |
| `PORT` | Server port | `8080` |

## Production Deployment

1. **Set environment variables**
   - Change `JWT_SECRET` to a strong, random string
   - Change `MESS_API_KEY` to a secure API key
   - Set `GIN_MODE=release`

2. **Build and deploy**
   ```bash
   # Build production image
   docker build -t mess-registration .

   # Run with production environment
   docker run -p 8080:8080 --env-file .env.prod mess-registration
   ```

## Features to Add

- [ ] OAuth integration (Google/Microsoft)
- [ ] Special dinner requests
- [ ] Day pass system
- [ ] Advanced swap matching algorithm
- [ ] Email notifications
- [ ] CSV export functionality
- [ ] Analytics dashboard
- [ ] Rate limiting
- [ ] Request logging

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.
