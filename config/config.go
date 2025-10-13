package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func EnvGetter(name, def string) func() string {
	return sync.OnceValue(func() string {
		// Try current directory first
		_ = godotenv.Load()
		// Fallback to parent directory
		_ = godotenv.Load("../.env")
		v := os.Getenv(name)
		if v == "" {
			return def
		}
		return v
	})
}

var (
	// Database
	DBUser     = EnvGetter("DB_USER", "default_user")
	DBPassword = EnvGetter("DB_PASSWORD", "default_password")
	DBHost     = EnvGetter("DB_HOST", "localhost")
	DBPort     = EnvGetter("DB_PORT", "5432")
	DBName     = EnvGetter("DB_NAME", "mydb")

	// JWT & App
	JWTSecret = EnvGetter("JWT_SECRET", "secret")
	STATUS    = EnvGetter("STATUS", "development")
	KUAPI     = EnvGetter("KU_API", "xxx.xxx.xxx.xxx/route")

	// MinIO
	MINIOEndpoint  = EnvGetter("MINIO_ENDPOINT", "localhost:9000")
	MINIOAccessKey = EnvGetter("MINIO_ACCESS_KEY", "minioadmin")
	MINIOSecretKey = EnvGetter("MINIO_SECRET_KEY", "minioadmin")
	MINIOBucket    = EnvGetter("MINIO_BUCKET", "tutorium")
	MINIOUseSSL    = EnvGetter("MINIO_USE_SSL", "false")

	// Payments (Omise)
	OMISEPublicKey = EnvGetter("OMISE_PUBLIC_KEY", "")
	OMISESecretKey = EnvGetter("OMISE_SECRET_KEY", "")

	// Payments defaults
	PAYMENTDefaultCurrency = EnvGetter("PAYMENT_DEFAULT_CURRENCY", "THB")
	PAYMENTReturnURI       = EnvGetter("PAYMENT_RETURN_URI", "")

	// Gorm config
	GORMLog = EnvGetter("GORM_LOG", "Warn")
)
