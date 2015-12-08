package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func main() {
	svc := s3.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})

	file, err := os.Open("/Users/G/Desktop/batman.jpg")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	params := &s3.PutObjectInput{
		Bucket:      aws.String("image.agentdesks.com"),
		Key:         aws.String("path/test.jpg"),
		ACL:         aws.String("public-read"),
		Body:        fileBytes,
		ContentType: aws.String("image/jpg"),
	}
	result, err := svc.PutObject(params)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(result)
	}
}
