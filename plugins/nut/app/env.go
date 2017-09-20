package app

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

// AWS aws
func AWS(f func(*credentials.Credentials) error) error {
	cfg := viper.GetStringMap("aws")
	aid := cfg["access_key_id"].(string)
	aky := cfg["secret_access_key"].(string)

	creds := credentials.NewStaticCredentials(aid, aky, "")
	if _, err := creds.Get(); err != nil {
		return err
	}
	return f(creds)
}

// S3URL url
func S3URL(u string) string {
	s3f := viper.GetStringMapString("aws.s3")
	buk := s3f["bucket"]
	reg := s3f["region"]
	return "https://s3-" + reg + ".amazonaws.com/" + buk + u
}

// S3 s3
func S3(f func(*s3.S3, string, string) error) error {
	return AWS(func(creds *credentials.Credentials) error {
		s3f := viper.GetStringMapString("aws.s3")
		buk := s3f["bucket"]
		reg := s3f["region"]

		svc := s3.New(
			session.New(),
			aws.NewConfig().WithRegion(reg).WithCredentials(creds),
		)
		return f(svc, reg, buk)
	})
}

// Home home url
func Home() string {
	scheme := "http"
	if viper.GetBool("server.ssl") {
		scheme += "s"
	}
	return scheme + "://" + Name()
}

// Name get server.name
func Name() string {
	return viper.GetString("server.name")
}

// IsProduction production mode ?
func IsProduction() bool {
	return viper.GetString("env") == "production"
}

func init() {

	viper.SetDefault("aws", map[string]interface{}{
		"access_key_id":     "change-me",
		"secret_access_key": "change-me",
		"s3": map[string]string{
			"region": "us-west-2",
			"bucket": "www.change-me.com",
		},
	})
	viper.SetDefault("redis", map[string]interface{}{
		"host": "localhost",
		"port": 6379,
		"db":   8,
	})

	viper.SetDefault("rabbitmq", map[string]interface{}{
		"user":     "guest",
		"password": "guest",
		"host":     "localhost",
		"port":     "5672",
		"virtual":  "fly-dev",
	})

	viper.SetDefault("postgresql", map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"user":     "postgres",
		"password": "",
		"dbname":   "fly_dev",
		"sslmode":  "disable",
		"pool": map[string]int{
			"max_open": 180,
			"max_idle": 6,
		},
	})

	viper.SetDefault("server", map[string]interface{}{
		"port": 8080,
		"ssl":  false,
		"name": "www.change-me.com",
	})

	viper.SetDefault("secret", Random(32))

	viper.SetDefault("elasticsearch", map[string]interface{}{
		"host": "localhost",
		"port": 9200,
	})

}
