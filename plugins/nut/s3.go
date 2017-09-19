package nut

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// UploadFile upload file
func (p *Controller) UploadFile(name string) ([]*Attachment, error) {
	files, err := p.GetFiles(name)
	if err != nil {
		return nil, err
	}
	creds := credentials.NewStaticCredentials(
		beego.AppConfig.String("awsaccesskeyid"),
		beego.AppConfig.String("awssecretaccesskey"),
		"",
	)

	if _, err = creds.Get(); err != nil {
		return nil, err
	}

	cfg := aws.NewConfig().
		WithRegion(beego.AppConfig.String("awss3region")).
		WithCredentials(creds)

	svc := s3.New(session.New(), cfg)

	var items []*Attachment
	for _, fh := range files {
		att, err := p.upload(svc, fh)
		if err != nil {
			return nil, err
		}
		items = append(items, att)
	}

	return items, nil
}

func (p *Controller) upload(svc *s3.S3, fh *multipart.FileHeader) (*Attachment, error) {
	bucket := beego.AppConfig.String("awss3bucketname")
	fd, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	fn := uuid.New().String() + filepath.Ext(fh.Filename)

	buffer := make([]byte, fh.Size)
	if _, err = fd.Read(buffer); err != nil {
		return nil, err
	}

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	path := "/upload/" + fn

	params := &s3.PutObjectInput{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String(bucket),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(fh.Size),
		ContentType:   aws.String(fileType),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		return nil, err
	}

	beego.Info(awsutil.StringValue(resp))
	item := Attachment{
		Length:       fh.Size,
		Title:        fh.Filename,
		MediaType:    fileType,
		ResourceID:   DefaultResourceID,
		ResourceType: DefaultResourceType,
		URL:          "https://s3-" + beego.AppConfig.String("awss3region") + ".amazonaws.com/" + bucket + path, // FIXME
		User:         p.CurrentUser(),
	}

	if _, err = orm.NewOrm().Insert(&item); err != nil {
		return nil, err
	}

	return &item, err
}
