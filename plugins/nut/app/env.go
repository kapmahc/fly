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
