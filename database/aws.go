package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/trungnghia250/malo-api/config"
)

var Sess *session.Session

func ConnectAws() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Config.Aws.Region),
		Credentials: credentials.NewStaticCredentials(
			config.Config.Aws.AccessKeyID,
			config.Config.Aws.SecretKeyID,
			""),
	})
	if err != nil {
		return err
	}

	Sess = sess

	return nil
}
