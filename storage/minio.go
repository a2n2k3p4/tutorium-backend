package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Client   *minio.Client
	Bucket   string
	Endpoint string
	UseSSL   bool
}

func NewClientFromEnv() (*Client, error) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(".env file not found, using system environment variables")
	}
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	useSSL := false
	if s := os.Getenv("MINIO_USE_SSL"); s != "" {
		b, _ := strconv.ParseBool(s)
		useSSL = b
	}

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, errors.New("minio config missing (MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET)")
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client:   minioClient,
		Bucket:   bucket,
		Endpoint: endpoint,
		UseSSL:   useSSL,
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// UploadBytes uploads bytes to bucket with objectName and returns public URL.
// folder: "users", "classes", "reports"
// filename: "user_42_123456.png"
func (c *Client) UploadBytes(ctx context.Context, folder, filename string, b []byte) (string, error) {
	if len(b) == 0 {
		return "", errors.New("empty file")
	}

	// detect content-type safely (first 512 bytes)
	sz := 512
	if len(b) < sz {
		sz = len(b)
	}
	contentType := http.DetectContentType(b[:sz])

	objectName := fmt.Sprintf("%s/%s", strings.Trim(folder, "/"), filename)
	_, err := c.Client.PutObject(ctx, c.Bucket, objectName, bytes.NewReader(b), int64(len(b)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	scheme := "http"
	if c.UseSSL {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s/%s/%s", scheme, c.Endpoint, c.Bucket, objectName)
	return url, nil
}

// DecodeBase64Image handles data: URI or plain base64
func DecodeBase64Image(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}
	if strings.HasPrefix(s, "data:") {
		comma := strings.IndexByte(s, ',')
		if comma < 0 {
			return nil, errors.New("invalid data URI")
		}
		s = s[comma+1:]
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	return b, nil
}

// GenerateFilename chooses extension based on content type and attaches timestamp
func GenerateFilename(contentType string) string {
	ext := ".bin"
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	}
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}
