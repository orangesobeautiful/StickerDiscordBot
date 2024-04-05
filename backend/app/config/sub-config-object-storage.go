package config

type ObjectStorage interface {
	GetEndpoint() string
	GetBucketName() string
	GetAccessKeyID() string
	GetAccessKeySecret() string
	GetPublicAccessURL() string
}

type objectStorage struct {
	Endpoint        string
	BucketName      string
	AccessKeyID     string
	AccessKeySecret string
	PublicAccessURL string
}

func (o *objectStorage) GetEndpoint() string {
	return o.Endpoint
}

func (o *objectStorage) GetBucketName() string {
	return o.BucketName
}

func (o *objectStorage) GetAccessKeyID() string {
	return o.AccessKeyID
}

func (o *objectStorage) GetAccessKeySecret() string {
	return o.AccessKeySecret
}

func (o *objectStorage) GetPublicAccessURL() string {
	return o.PublicAccessURL
}
