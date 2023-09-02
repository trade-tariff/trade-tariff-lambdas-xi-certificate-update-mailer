package fetcher

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func setup() {
	logger.Initialize()
}

type mockS3Client struct {
	mockListOutput   *s3.ListObjectsV2Output
	mockGetObjectOut *s3.GetObjectOutput
}

func (m *mockS3Client) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return m.mockListOutput, nil
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return m.mockGetObjectOut, nil
}

func TestFetchFileObject(t *testing.T) {
	setup()
	mockClient := &mockS3Client{
		mockListOutput: &s3.ListObjectsV2Output{
			Contents: []*s3.Object{
				{
					Key:          aws.String("path/to/file_2023-08-30.xml"),
					LastModified: aws.Time(time.Date(2023, 8, 30, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
	}

	bucket := "mockBucket"
	prefix := "mockPrefix"
	fetcher := NewFetcher(mockClient, bucket, prefix)

	// Date to search for in our mock data.
	date := "2023-08-30"
	obj := fetcher.FetchFileObject(date)

	assert.Equal(t, "path/to/file_2023-08-30.xml", *obj.Key)
	assert.Equal(t, time.Date(2023, 8, 30, 0, 0, 0, 0, time.UTC), *obj.LastModified)
}

func TestFetchXML(t *testing.T) {
	setup()

	xmlContent := []byte("<xml><data>test data</data></xml>")
	mockObject := &s3.Object{
		Key:          aws.String("path/to/file_2023-08-30.xml"),
		LastModified: aws.Time(time.Date(2023, 8, 30, 0, 0, 0, 0, time.UTC)),
	}

	mockClient := &mockS3Client{
		mockGetObjectOut: &s3.GetObjectOutput{
			ContentLength: aws.Int64(int64(len(xmlContent))),
			Body:          io.NopCloser(bytes.NewReader(xmlContent)),
		},
	}

	fetcher := NewFetcher(mockClient, "mockBucket", "mockPrefix")

	xmlFile := fetcher.FetchXML(mockObject)

	assert.Equal(t, "path/to/file_2023-08-30.xml", xmlFile.Key)
  assert.Equal(t, xmlFile.LoadedOn, "2023-08-30")
	assert.Equal(t, int64(len(xmlContent)), xmlFile.ContentLength)
	assert.Equal(t, xmlContent, xmlFile.Content)
}
