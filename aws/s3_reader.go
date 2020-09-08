package aws

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var sess *session.Session

func Init() {
	var err error
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		exitErrorf(err.Error())
	}
}

func GetObject(bucket string, key string) ([]byte, error) {
	if sess == nil {
		Init()
	}
	svc := s3.New(sess)
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println(err.Error())
		return []byte{}, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(output.Body)
	return buf.Bytes(), nil
}
