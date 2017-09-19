package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestS3Upload(t *testing.T) {

	token := ""
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsAecretAccessKey, token)

	_, err := creds.Get()
	if err != nil {
		t.Fatal(err)
	}

	cfg := aws.NewConfig().WithRegion(awsS3Region).WithCredentials(creds)

	svc := s3.New(session.New(), cfg)
	file, err := os.Open("s3_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()

	size := fileInfo.Size()

	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	// fileType := http.DetectContentType(buffer)
	path := "/media/" + file.Name()

	params := &s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(awsS3BucketName),
		Key:    aws.String(path),
		Body:   fileBytes,
		// ContentLength: aws.Int64(size),
		// ContentType:   aws.String(fileType),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("response %s", awsutil.StringValue(resp))
}
