package util

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig 
	Database DatabaseConfig 
	Ton      TonConfig 
	SMTP     SMTPConfig 
	Redis    RedisConfig 
	Auth     AuthConfig 
	AWS 	AWSConfig
}

type AWSConfig struct {
	AWSRegion string
	AWSAccessKeyID string
	AWSSecretAccessKey string
	AWSBucket string
}


type AppConfig struct {
	BaseUrl            string 
	HttpPort           int   
	DomainName	      string	
	SeedPhrase 	   string 
	CookieSecretKey    string
	JwtSecretKey       string 
	NotificationsEmail string 
	AdminCollectionAddress string
	AuthMetadataID     int64
	AlloweGroupChatID    int64
	BasicUsername            string
	BasicPassword            string
}

type DatabaseConfig struct {
	Dsn           string 
	Automigrate   bool  
}

type TonConfig struct {
	PublicConfig       string
	MaxConcurrentTask  int   
	SharedSecret       string
	ProfLifeTimeSec    int   
	NodeAddress        string
	ApiKey             string 
	AdminWallet	 string
}

type SMTPConfig struct {
	Host     string 
	Port     int   
	Username string
	Password string 
	From     string 
}

type RedisConfig struct {
	Addr string
	Password string
	PoolSize int
}

type AuthConfig struct {
	GithubRedirectUrl string
	GithubClientId    string 
	GithubClientSecret string
	TelegramBotToken  string
}

func LoadConfig() (config Config, err error) {

	env := os.Getenv("APP_ENV") // this could be "dev", "prod", or "test"
	if env == "" {
		err = godotenv.Load("app.dev.env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}	

	httpPort, _ := strconv.Atoi(os.Getenv("APP_HTTP_PORT")) // error handling omitted for brevity

	// seedPhrase := os.Getenv("APP_SEED_PHRASE")


	// convert to int64
	authMetadataID := os.Getenv("APP_AUTH_METADATA_ID")

	int64AuthMetadataID, _ := strconv.ParseInt(authMetadataID, 10, 64)

	int64AlloweGroupChatID, _ := strconv.ParseInt(os.Getenv("APP_ALLOWED_GROUP_CHAT_ID"), 10, 64)

	appConfig := AppConfig{
		BaseUrl:            os.Getenv("APP_BASE_URL"),
		DomainName:         os.Getenv("APP_DOMAIN_NAME"),
		HttpPort:           httpPort,
		CookieSecretKey:    os.Getenv("APP_COOKIE_SECRET_KEY"),
		JwtSecretKey:       os.Getenv("APP_JWT_SECRET_KEY"),
		NotificationsEmail: os.Getenv("APP_NOTIFICATIONS_EMAIL"),
		SeedPhrase: os.Getenv("APP_SEED_PHRASE"),
		AdminCollectionAddress: os.Getenv("APP_ADMIN_COLLECTION_ADDRESS"),
		AuthMetadataID: int64AuthMetadataID,
		AlloweGroupChatID: int64AlloweGroupChatID,
		BasicUsername: os.Getenv("APP_BASIC_USERNAME"),
		BasicPassword: os.Getenv("APP_BASIC_PASSWORD"),
	}

	autoMigrate, _ := strconv.ParseBool(os.Getenv("DATABASE_AUTOMIGRATE"))


	databaseConfig := DatabaseConfig{
		Dsn:           os.Getenv("DATABASE_DSN"),
		Automigrate:   autoMigrate,
	}

	tonMaxConcurrentTask, _ := strconv.Atoi(os.Getenv("TON_MAX_CONCURRENT_TASK"))
	tonProfLifeTimeSec, _ := strconv.Atoi(os.Getenv("TON_PROF_LIFE_TIME_SEC"))

	tonConfig := TonConfig{
		PublicConfig:       os.Getenv("TON_PUBLIC_CONFIG"),
		MaxConcurrentTask:  tonMaxConcurrentTask,
		SharedSecret:       os.Getenv("TON_SHARED_SECRET"),
		ProfLifeTimeSec:    tonProfLifeTimeSec,
		NodeAddress:        os.Getenv("TON_NODE_ADDRESS"),
		ApiKey:             os.Getenv("TON_API_KEY"),
		AdminWallet:	 os.Getenv("TON_ADMIN_WALLET"),
	}

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	smtpConfig := SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     smtpPort,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}

	redisPoolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))

	redisConfig := RedisConfig{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		PoolSize: redisPoolSize,
	}

	authConfig := AuthConfig{
		GithubRedirectUrl: os.Getenv("AUTH_GITHUB_REDIRECT_URL"),
		GithubClientId:    os.Getenv("AUTH_GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("AUTH_GITHUB_CLIENT_SECRET"),
		TelegramBotToken:  os.Getenv("AUTH_TELEGRAM_BOT_TOKEN"),
	}

	awsConfig := AWSConfig{
		AWSRegion: os.Getenv("AWS_REGION"),
		AWSAccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),		
		AWSBucket: os.Getenv("AWS_BUCKET"),
	}

	config = Config{
		App:      appConfig,
		Database: databaseConfig,
		Ton:      tonConfig,
		SMTP:     smtpConfig,
		Redis:    redisConfig,
		Auth:     authConfig,
		AWS: awsConfig,
	}

	return
}