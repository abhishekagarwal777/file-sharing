package config

import (
    "os"
)

type Config struct {
    StorageType string
    S3Bucket    string
    LocalPath   string
}

func LoadConfig() Config {
    return Config{
        StorageType: os.Getenv("STORAGE_TYPE"),   // "local" or "s3"
        S3Bucket:    os.Getenv("S3_BUCKET"),      // For S3 bucket name
        LocalPath:   os.Getenv("LOCAL_PATH"),     // Path to store files locally
    }
}
