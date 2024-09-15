package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "time"
    "fmt"
    "path/filepath"

    "golang.org/x/crypto/bcrypt"
    "file-sharing-platform/internal/db"
    "file-sharing-platform/internal/auth"
    "file-sharing-platform/internal/config"
    "github.com/joho/godotenv"
)

// Credentials struct for user input
type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Register handler
func Register(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    _, err = db.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", creds.Email, hashedPassword)
    if err != nil {
        http.Error(w, "Error inserting record", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    var storedPassword string
    err = db.DB.QueryRow("SELECT password FROM users WHERE email=$1", creds.Email).Scan(&storedPassword)
    if err == sql.ErrNoRows {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateJWT(creds.Email)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    token,
        Expires:  time.Now().Add(1 * time.Hour),
    })
}

// UploadFile handler
func UploadFile(w http.ResponseWriter, r *http.Request) {
    config := config.LoadConfig()

    // Parse the form containing the file upload
    err := r.ParseMultipartForm(10 << 20) // Max file size: 10MB
    if err != nil {
        http.Error(w, "File too big", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Invalid file upload", http.StatusBadRequest)
        return
    }
    defer file.Close()

    var fileURL string
    if config.StorageType == "s3" {
        fileURL, err = files.SaveFileToS3(file, header, config.S3Bucket)
    } else {
        fileName, err := files.SaveFileLocally(file, header, config.LocalPath)
        if err == nil {
            fileURL = filepath.Join(config.LocalPath, fileName)
        }
    }

    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
        return
    }

    // Save metadata to the database
    userID := 1 // Replace with the actual user ID from JWT
    _, err = db.DB.Exec("INSERT INTO files (user_id, file_name, s3_url, upload_date, file_size) VALUES ($1, $2, $3, NOW(), $4)",
        userID, header.Filename, fileURL, header.Size)
    if err != nil {
        http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "File uploaded successfully: %s", fileURL)
}

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Initialize the database
    db.InitDB()

    // Set up HTTP handlers
    http.HandleFunc("/register", Register)
    http.HandleFunc("/login", Login)
    http.HandleFunc("/upload", UploadFile)

    // Start the server
    log.Fatal(http.ListenAndServe(":8080", nil))
}
