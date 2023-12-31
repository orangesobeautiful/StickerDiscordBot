package objectstorage

import (
	"context"
	"io"
	"net/url"

	"backend/app/config"
	"backend/app/pkg/hserr"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/xerrors"
)

type BucketObjectDataConverter interface {
	GetObjectDirectURL(key string) (result string)
}

type BucketObjectOperator interface {
	Upload(ctx context.Context, key string, bodyReader io.Reader, opts ...putObjectOptionFunc) (err error)
	DeleteObjects(ctx context.Context, keys ...string) (err error)
}

type BucketObjectHandler interface {
	BucketObjectDataConverter
	BucketObjectOperator
}

var _ BucketObjectHandler = (*bucketHandler)(nil)

type bucketHandler struct {
	s3Client *s3.Client
	uploader *manager.Uploader

	bucketName *string

	publicAccessURL *url.URL
}

func NewBucketHandler(ctx context.Context, cfg *config.CfgInfo) (handler BucketObjectHandler, err error) {
	publicAccessURL, err := url.Parse(cfg.ObjectStorage.PublicAccessURL)
	if err != nil {
		return nil, xerrors.Errorf("parse public access url: %w", err)
	}

	accessKeyID := cfg.ObjectStorage.AccessKeyID
	accessKeySecret := cfg.ObjectStorage.AccessKeySecret

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: cfg.ObjectStorage.Endpoint,
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithEndpointResolverWithOptions(r2Resolver),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
	)
	if err != nil {
		return nil, xerrors.Errorf("failed to load config, %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	const bufferSize = 25 * 1024 * 1024
	s3Uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(bufferSize)
	})

	bh := &bucketHandler{
		s3Client:        s3Client,
		uploader:        s3Uploader,
		bucketName:      aws.String(cfg.ObjectStorage.BucketName),
		publicAccessURL: publicAccessURL,
	}
	return bh, nil
}

type putObjectOptionFunc func(*putObjectOption)

type putObjectOption struct {
	contentLength      int64
	contentDisposition string
	contentType        string
}

func PutObjectWithContentLength(contentLength int64) putObjectOptionFunc {
	return func(o *putObjectOption) {
		o.contentLength = contentLength
	}
}

func PutObjectWithContentDisposition(contentDisposition string) putObjectOptionFunc {
	return func(o *putObjectOption) {
		o.contentDisposition = contentDisposition
	}
}

func PutObjectWithContentType(contentType string) putObjectOptionFunc {
	return func(o *putObjectOption) {
		o.contentType = contentType
	}
}

func (b *bucketHandler) Upload(
	ctx context.Context,
	key string, bodyReader io.Reader, opts ...putObjectOptionFunc,
) (err error) {
	var o putObjectOption
	for _, opt := range opts {
		opt(&o)
	}

	putInput := &s3.PutObjectInput{
		Bucket: b.bucketName,
		Key:    aws.String(key),
		Body:   bodyReader,
	}
	if o.contentLength != 0 {
		putInput.ContentLength = aws.Int64(o.contentLength)
	}
	if o.contentDisposition != "" {
		putInput.ContentDisposition = aws.String(o.contentDisposition)
	}
	if o.contentType != "" {
		putInput.ContentType = aws.String(o.contentType)
	}

	_, err = b.uploader.Upload(ctx, putInput)
	if err != nil {
		return hserr.NewInternalError(err, "s3 put object")
	}
	return nil
}

func (b *bucketHandler) GetObjectDirectURL(key string) (result string) {
	return b.publicAccessURL.JoinPath(key).String()
}

func (b *bucketHandler) DeleteObjects(ctx context.Context, keys ...string) (err error) {
	var ojbectIDs []types.ObjectIdentifier
	for _, key := range keys {
		ojbectIDs = append(ojbectIDs, types.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	_, err = b.s3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: b.bucketName,
		Delete: &types.Delete{Objects: ojbectIDs},
	})
	if err != nil {
		return hserr.NewInternalError(err, "s3 delete objects")
	}

	return nil
}
