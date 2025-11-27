package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	Client     *s3.Client
	BucketName string
	PublicURL  string
}

func NewS3Service() (*S3Service, error) {
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_ACCESS")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	s3URL := os.Getenv("S3_URL")
	publicURL := os.Getenv("S3_PUBLIC_URL")
	region := os.Getenv("S3_REGION")

	if region == "" {
		region = "auto"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s3URL)
	})

	return &S3Service{
		Client:     client,
		BucketName: bucketName,
		PublicURL:  publicURL,
	}, nil
}

func (s *S3Service) UploadFile(file multipart.File, filename string, contentType string) (string, error) {
	log.Println("Uploading file to S3...")
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Construct public URL
	fileURL := fmt.Sprintf("%s/%s", s.PublicURL, filename)
	return fileURL, nil
}
