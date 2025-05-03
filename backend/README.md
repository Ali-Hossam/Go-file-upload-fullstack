# GoSheet Backend

A high-performance Go backend service for processing, storing, and retrieving CSV data with real-time updates.

## Features

- Fast CSV file processing with Go
- Real-time upload progress tracking via WebSockets
- Concurrent file processing for better performance
- Clean architecture with separation of concerns
- PostgreSQL database integration using GORM
- RESTful API for data retrieval with filtering, sorting, and pagination

## Tech Stack

- Go 1.24+
- Echo web framework
- GORM ORM for PostgreSQL
- Gorilla WebSockets for real-time updates
- PostgreSQL for data storage

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 14+ database
- Git

## Environment Setup

Create a `.env` file in the root directory with:

```
DB_DSN_LOCAL=postgres://username:password@localhost:5432/database_name?sslmode=disable
SERVER_PORT=8080
```

Adjust the PostgreSQL connection string to match your database setup.

## Installation

```bash
# Clone the repository
git clone git@github.com:username/Go-file-upload-fullstack.git
cd Go-file-upload-fullstack/backend

# Install dependencies
go mod download
```

## Development

```bash
# Run the application
go run cmd/app/main.go
```

This will start the server on port 8080 (or the port specified in your `.env`).

## API Endpoints

### File Upload

- `POST /api/upload` - Upload CSV files
- `GET /api/upload/status/:uploadID` - WebSocket endpoint for tracking upload progress

### Data Retrieval

- `GET /api/students` - Get student records with filtering, sorting, and pagination
  - Query parameters:
    - `page` - Page number (default: 1)
    - `size` - Records per page (default: 100)
    - `sort_by` - Field to sort by (Student_name, Subject, Grade)
    - `sort_order` - Sort direction (asc, desc)
    - `name` - Filter by student name (partial match)
    - `subject` - Filter by subject

## Project Structure

```
backend/
├── cmd/
│   └── app/                # Application entry point
├── config/                 # Application configuration
├── database/               # Database connection and models
│   ├── model/              # GORM models
│   └── repository/         # Data access layer
├── internal/               # Internal application code
│   ├── api/                # API handlers and routes
│   │   └── handler/        # Request handlers
│   └── service/            # Business logic
│       └── csv/            # CSV processing services
└── .env                    # Environment variables
```

## Architecture

The application follows a clean architecture pattern:

1. **Models**: Define data structures for the database
2. **Repositories**: Handle data access and storage
3. **Services**: Implement business logic
4. **Handlers**: Process HTTP requests and responses
5. **Router**: Define API endpoints

## CSV Processing

CSV files are processed concurrently with real-time progress updates:

1. Files are uploaded via HTTP
2. Processing is done in background goroutines
3. Progress is reported via WebSockets
4. Data is inserted into the database in batches for performance

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Performance Considerations

- Files are processed in chunks to minimize memory usage
- Database insertions are batched for better performance
- Background processing prevents blocking the HTTP server
- Multiple files are processed concurrently
