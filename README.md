This repository contains a Go-based file-sharing and management system that allows users to securely upload, manage, and share files. The system is designed to handle concurrent file uploads efficiently and stores file metadata in PostgreSQL. It supports user authentication and authorization using JWT, ensuring each user can manage their own files while providing the ability to share files through public URLs.

Features
User Authentication & Authorization

User registration and login with email and password.
JWT-based session management to secure access.
Role-based access control to ensure users can only manage their own files.
File Upload & Management

Upload files to Amazon S3 (or local storage if S3 integration is unavailable).
Save file metadata (file name, size, upload date, storage URL) in PostgreSQL.
Efficiently handle large file uploads using Go's goroutines.
API endpoints for uploading files and retrieving file metadata.
File Retrieval & Sharing

Retrieve uploaded files' metadata via secure API endpoints.
Generate public URLs for sharing files, allowing others to access files securely.
Users can only access and share their own files.
Technologies Used
Go: Main programming language used for the backend.
PostgreSQL: Database for storing user and file metadata.
Amazon S3: (Optional) Cloud storage for uploaded files.
JWT: Used for user authentication and session management.
Goroutines: To handle concurrency and optimize file uploads.
Getting Started
Prerequisites
Go 1.17+
PostgreSQL 12+
AWS S3 account (Optional, for cloud storage)
Docker (optional, for setting up PostgreSQL)
Installation
Clone the repository:

bash
Copy code
git clone https://github.com/yourusername/file-sharing-management-system.git
cd file-sharing-management-system
Set up environment variables:

Create a .env file in the root directory of the project to configure your PostgreSQL, S3, and JWT settings:

plaintext
Copy code
POSTGRES_USER=yourusername
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=yourdbname
POSTGRES_HOST=localhost
S3_BUCKET=your-s3-bucket-name
S3_REGION=your-s3-region
JWT_SECRET=your_jwt_secret
Install dependencies:

bash
Copy code
go mod download
Run PostgreSQL:

If using Docker, start PostgreSQL with:

bash
Copy code
docker-compose up -d
Run the application:

bash
Copy code
go run main.go
API Endpoints
User Authentication

POST /register: Register a new user.
POST /login: Log in and receive a JWT token.
File Management

POST /upload: Upload a file and save metadata.
GET /files: Retrieve metadata for all uploaded files.
GET /share/:file_id: Generate a public URL for sharing a file.
Contributing
We welcome contributions! Feel free to fork the repository and submit a pull request with your improvements.
