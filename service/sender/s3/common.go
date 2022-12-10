package s3

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/labstack/gommon/log"
	"github.com/trungnghia250/malo-api/database"
	"mime/multipart"
)

func getImageBuffer(body multipart.File) ([]byte, error) {
	buffer := make([]byte, 512)
	_, err := body.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func UploadImage(fh *multipart.FileHeader, bucket string) (string, error) {
	body, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Errorf(err.Error())
		}
	}()

	buffer, err := getImageBuffer(body)

	if !filetype.IsImage(buffer) {
		return "", errors.New("image not safe")
	}

	uploader := s3manager.NewUploader(database.Sess)
	kind, err := filetype.Match(buffer)
	if err != nil {
		return "", err
	}

	filename := uuid.New().String()
	path := fmt.Sprintf("%s/%s.%s", bucket, filename, kind.Extension)
	resp, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("malo-app"),
		ACL:         aws.String("public-read"),
		Key:         aws.String(path),
		ContentType: aws.String(kind.MIME.Value),
		Body:        body,
	})

	if err != nil {
		return "", err
	}

	return resp.Location, nil
}
