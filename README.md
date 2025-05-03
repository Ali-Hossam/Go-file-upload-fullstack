# CSV Student Records Processor - Full Stack Application

## Overview

This project is a full-stack application for uploading, processing, and displaying CSV student records data. The system consists of:

1. A Go backend for handling file uploads, processing CSV data, and providing a RESTful API
2. A React frontend for file upload with real-time progress tracking and data visualization
   
![Screenshot from 2025-05-03 21-54-01](https://github.com/user-attachments/assets/4acf71d0-80cb-4b70-afac-0d5238557877)
![Screenshot from 2025-05-03 21-54-13](https://github.com/user-attachments/assets/b7344fdd-0411-465a-9ba3-0888e0d372c1)

## Features

### Backend Features
- Fast CSV file processing with Go
- Concurrent file processing for better performance
- Real-time upload progress tracking via WebSockets
- PostgreSQL database for data storage using GORM
- RESTful API with filtering, sorting, and pagination
- Clean architecture with separation of concerns

### Frontend Features
- Drag and drop file uploading with real-time progress indicators
- CSV validation on client-side
- Data visualization with filtering, sorting, and pagination
- Mobile-responsive design with modern UI

## Tech Stack

### Backend
- Go 1.24+
- Echo web framework
- GORM ORM with PostgreSQL
- Gorilla WebSockets for real-time updates

### Frontend
- React 19
- TypeScript
- Vite 6
- Tailwind CSS 4
- React Router

## Getting Started

### Prerequisites
- Docker and Docker Compose
- PostgreSQL (if running locally without Docker)
- Node.js 18+ (for frontend development)
- Go 1.24+ (for backend development)

### Local Development

#### Backend
1. Navigate to the backend directory:
```bash
cd Go-file-upload-fullstack/backend
```

2. Create a `.env` file:
```
DB_DSN_LOCAL=postgres://username:password@localhost:5432/studentsdb?sslmode=disable
SERVER_PORT=8080
```

3. Install dependencies and run:
```bash
go mod download
go run cmd/app/main.go
```

#### Frontend
1. Navigate to the frontend directory:
```bash
cd Go-file-upload-fullstack/frontend
```

2. Create a `.env` file:
```
VITE_SERVER_URL=http://localhost:8080/
```

3. Install dependencies and run:
```bash
bun install
# or npm install
bun dev
# or npm run dev
```

## API Endpoints

### File Upload
- `POST /api/upload` - Upload CSV files
- `GET /api/upload/status/:uploadID` - WebSocket endpoint for upload progress

### Data Retrieval
- `GET /api/students` - Get student records with filtering, sorting, pagination
  - Query parameters:
    - `page` - Page number (default: 1)
    - `size` - Records per page (default: 100)
    - `sort_by` - Field to sort by (Student_name, Subject, Grade)
    - `sort_order` - Sort direction (asc, desc)
    - `name` - Filter by student name (partial match)
    - `subject` - Filter by subject
