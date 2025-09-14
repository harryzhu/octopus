package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	R2Client *s3.Client
)

func initR2() {
	accountID := GetEnv("CFR2ID", "")
	accessKeyID := GetEnv("CFR2KEYID", "")
	accessKeySecret := GetEnv("CFR2KEYSECRET", "")

	if IsAnyEmpty(accountID, accessKeyID, accessKeySecret) {
		DebugWarn("R2init", "cannot get env vars: accountID / accessKeyID / accessKeySecret, R2 service is not available")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		PrintError("initR2", err)
	}
	R2Client = s3.NewFromConfig(cfg)

}

func R2Ping(bucketName string) error {
	_, err := R2Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		PrintError("R2Ping: R2 service is not available:", err)
		return err
	}

	return nil
}

func R2Save(bucketName string, r2Key string, fileName string, fileMIME string) error {
	if IsAnyEmpty(bucketName, r2Key, fileName, fileMIME) {
		return NewError("cannot be empty")
	}

	r2Key = strings.ToLower(r2Key)

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		PrintError("R2Save.1", err)
		return err
	}

	_, err = R2Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(r2Key),
		Body:        file,
		ContentType: aws.String(fileMIME),
	})

	if err != nil {
		PrintError("R2Save.2: Couldn't save file to S3:", err)
		return err
	}

	return nil
}

func R2Delete(bucketName string, r2Key string) error {
	if IsAnyEmpty(bucketName, r2Key) {
		return NewError("cannot be empty")
	}

	r2Key = strings.ToLower(r2Key)

	_, err := R2Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(r2Key),
	})

	if err != nil {
		PrintError("R2Delete: cannot delete file to S3:", err)
		return err
	}

	return nil
}
