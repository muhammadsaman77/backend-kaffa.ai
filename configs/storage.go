package configs

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitStorage(ctx context.Context) *s3.Client {

	awsEndpoint := viper.GetString("aws_endpoint")

	os.Setenv("AWS_ACCESS_KEY_ID", viper.GetString("aws_access_key_id"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", viper.GetString("aws_secret_access_key"))
	os.Setenv("AWS_REGION", viper.GetString("aws_region"))
	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		Log.Fatal("unable to load SDK config, " + err.Error())
	}
	s3Client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = &awsEndpoint
		o.UsePathStyle = true
		o.HTTPClient = &http.Client{Timeout: 10 * time.Second}

	})
	Log.Info("S3 client initialized successfully")
	return s3Client

}
