package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var s3Client *s3.Client

func InitS3Client() {
	region := os.Getenv("S3_REGION")
	endpoint := os.Getenv("S3_ENDPOINT")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == s3.ServiceID && region == region {
					return aws.Endpoint{
						URL: endpoint,
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			}),
		),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	fmt.Println("S3 client initialized")
}

func UploadFile(bucketName, fileName string, file *os.File) error {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %v", file.Name(), bucketName, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", file.Name(), bucketName)
	return nil
}

func DownloadFile(bucketName, fileName string) (io.ReadCloser, error) {
	output, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to download item %q, %v", fileName, err)
	}

	fmt.Printf("Successfully downloaded %q\n", fileName)
	return output.Body, nil
}
