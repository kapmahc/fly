package nut

import (
	"bytes"
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
func (p *Controller) UploadFile(name string) (*Attachment, error) {
	fd, fh, err := p.GetFile(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	buf := make([]byte, fh.Size)
	if _, err := fd.Read(buf); err != nil {
		return nil, err
	}
	return p.writeToS3(fh.Filename, buf, fh.Size)
}

func (p *Controller) writeToS3(name string, body []byte, size int64) (*Attachment, error) {
	reg := beego.AppConfig.String("awss3region")
	aid := beego.AppConfig.String("awsaccesskeyid")
	aky := beego.AppConfig.String("awssecretaccesskey")
	buk := beego.AppConfig.String("awss3bucketname")

	creds := credentials.NewStaticCredentials(aid, aky, "")
	if _, err := creds.Get(); err != nil {
		return nil, err
	}

	svc := s3.New(
		session.New(),
		aws.NewConfig().WithRegion(reg).WithCredentials(creds),
	)

	fn := "/upload/" + uuid.New().String() + filepath.Ext(name)

	fileBytes := bytes.NewReader(body)
	fileType := http.DetectContentType(body)

	params := &s3.PutObjectInput{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String(buk),
		Key:           aws.String(fn),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		return nil, err
	}

	beego.Info(awsutil.StringValue(resp))
	item := Attachment{
		Length:       size,
		Title:        name,
		MediaType:    fileType,
		ResourceID:   DefaultResourceID,
		ResourceType: DefaultResourceType,
		URL:          "https://s3-" + reg + ".amazonaws.com/" + buk + fn, // FIXME
		User:         p.CurrentUser(),
	}

	if _, err = orm.NewOrm().Insert(&item); err != nil {
		return nil, err
	}

	return &item, nil
}
