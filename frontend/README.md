# GoSheet Frontend

A modern, responsive React application for uploading and managing CSV data using the GoSheet backend.

## Features

- Drag and drop file uploading with real-time progress indicators
- CSV validation on client-side
- WebSocket connection for live upload status updates
- Data visualization with filtering, sorting, and pagination
- Mobile-responsive design with modern UI

## Tech Stack

- React 19
- TypeScript
- Vite 6
- Tailwind CSS 4
- React Router
- Lucide React for icons

## Prerequisites

- Node.js (v18+ recommended)
- Bun, npm, or Yarn package manager

## Environment Setup

Create a `.env` file in the root directory with:

```
VITE_SERVER_URL=localhost:8080
```

Adjust the server URL to match your backend deployment.

## Installation

```bash
# Install dependencies
bun install
# or npm install
# or yarn install
```

## Development

```bash
# Start development server
bun dev
# or npm run dev
# or yarn dev
```

This will start a local development server at http://localhost:5173.

## Building for Production

```bash
# Build for production
bun run build
# or npm run build
# or yarn build
```

The build output will be in the `dist` directory.

## Project Structure

```
frontend/
├── public/               # Static assets
├── src/                  # Source code
│   ├── assets/           # Images, fonts, etc.
│   ├── components/       # Reusable UI components
│   ├── pages/            # Page components
│   ├── types/            # TypeScript type definitions
│   ├── App.tsx           # Main application component
│   ├── index.css         # Global styles
│   └── main.tsx          # Application entry point
├── .env                  # Environment variables
├── index.html            # HTML template
├── package.json          # Project dependencies
├── tsconfig.json         # TypeScript configuration
└── vite.config.ts        # Vite configuration
```

## Key Components

- **FileUpload**: Handles CSV file uploading with real-time progress indicators via WebSockets
- **StudentsTable**: Displays uploaded student records
- **Options**: Controls for filtering, sorting, and pagination
- **NameSearch**: Search component for filtering records by student name

## Usage

1. Upload CSV files using the drag-and-drop area or file browser
2. View real-time upload progress
3. Browse records on the dashboard
4. Use filters, sorting, and pagination to navigate the data

## Deployment

The application can be deployed to any static hosting service:

1. Build the project with `bun run build`
2. Upload the contents of the `dist` directory to your hosting provider
3. Configure your server to handle client-side routing (if applicable)
