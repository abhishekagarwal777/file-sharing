package files

import (
    "bytes"
    "fmt"
    "mime/multipart"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

// SaveFileToS3 uploads the file to an S3 bucket
func SaveFileToS3(file multipart.File, header *multipart.FileHeader, bucketName string) (string, error) {
    defer file.Close()

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
    })
    if err != nil {
        return "", fmt.Errorf("failed to connect to AWS: %v", err)
    }

    s3Client := s3.New(sess)

    fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)

    // Read the file content into a buffer
    buffer := make([]byte, header.Size)
    file.Read(buffer)

    // Upload the file to S3
    _, err = s3Client.PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(bucketName),
        Key:                  aws.String(fileName),
        Body:                 bytes.NewReader(buffer),
        ContentLength:        aws.Int64(header.Size),
        ContentType:          aws.String(header.Header.Get("Content-Type")),
        ContentDisposition:   aws.String("attachment"),
        ServerSideEncryption: aws.String("AES256"), // Optional, can enable server-side encryption
    })

    if err != nil {
        return "", fmt.Errorf("could not upload file: %v", err)
    }

    return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName), nil
}
