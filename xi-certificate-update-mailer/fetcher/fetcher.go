package fetcher

import (
	"io"
	"path/filepath"

	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client interface {
	ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

type Fetcher struct {
	S3           S3Client
	BucketName   string
	BucketPrefix string
}

func NewFetcher(s3 S3Client, bucket, prefix string) *Fetcher {
	return &Fetcher{
		S3:           s3,
		BucketName:   bucket,
		BucketPrefix: prefix,
	}
}

func (f *Fetcher) FetchXML(object *s3.Object) *XmlFile {
	if object == nil {
		logger.Log.Fatal(
			"No file found for today. Has the file been uploaded?",
			logger.String("bucket", f.BucketName),
			logger.String("prefix", f.BucketPrefix),
		)
	}

	objectOutput, err := f.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(f.BucketName),
		Key:    object.Key,
	})

	if err != nil {
		logger.Log.Fatal(
			"Error occurred while getting object.",
			logger.String("object", *object.Key),
			logger.String("bucket", f.BucketName),
		)
	}

	bytes := make([]byte, *objectOutput.ContentLength)

	defer objectOutput.Body.Close()

	_, err = io.ReadFull(objectOutput.Body, bytes)

	if err != nil {
		logger.Log.Fatal(
			"Error occurred while reading file",
			logger.String("object", *object.Key),
			logger.String("bucket", f.BucketName),
		)
	}

	return &XmlFile{
		Key:           aws.StringValue(object.Key),
		LoadedOn:      object.LastModified.Format("2006-01-02"),
		ContentLength: aws.Int64Value(objectOutput.ContentLength),
		Content:       bytes,
	}
}

func (f *Fetcher) FetchFileObject(date string) *s3.Object {
	resp, err := f.S3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(f.BucketName),
		Prefix: aws.String(f.BucketPrefix),
	})

	if err != nil {
		logger.Log.Fatal(
			"Error occurred while listing objects",
			logger.String("bucket", f.BucketName),
			logger.String("prefix", f.BucketPrefix),
		)
	}

	for _, item := range resp.Contents {
		if item.LastModified.Format("2006-01-02") == date {
			logger.Log.Info(
				"Found file on date",
				logger.String("object", *item.Key),
				logger.String("loaded_on", item.LastModified.Format("2006-01-02")),
			)
			return item
		}
	}

	return nil
}

type XmlFile struct {
	Key           string
	LoadedOn      string
	ContentLength int64
	Content       []byte
}

func (x XmlFile) Filename() string {
	if x.Key != "" {
		return filepath.Base(x.Key)
	} else {
		return ""
	}
}
