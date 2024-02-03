package config

type VectorDatabase interface {
	GetQdrant() Qdrant

	GetToIntializeCollection() bool

	GetCollectionName() string
}

var _ VectorDatabase = (*vectorDatabase)(nil)

type vectorDatabase struct {
	Qdrant *qdrant

	InitializeCollection bool

	CollectionName string
}

func (v *vectorDatabase) GetQdrant() Qdrant {
	return v.Qdrant
}

func (v *vectorDatabase) GetToIntializeCollection() bool {
	return v.InitializeCollection
}

func (v *vectorDatabase) GetCollectionName() string {
	return v.CollectionName
}

type Qdrant interface {
	GetAddr() string

	GetAPICredentials() QdrantAPICredentials
}

var _ Qdrant = (*qdrant)(nil)

type qdrant struct {
	Addr string

	APICredentials *qdrantAPICredentials
}

func (q *qdrant) GetAddr() string {
	return q.Addr
}

func (q *qdrant) GetAPICredentials() QdrantAPICredentials {
	return q.APICredentials
}

type QdrantAPICredentials interface {
	GetAPIKey() string

	GetAPIKeyRequireTransportSecurity() bool
}

var _ QdrantAPICredentials = (*qdrantAPICredentials)(nil)

type qdrantAPICredentials struct {
	APIKey string

	APIKeyRequireTransportSecurity bool
}

func (q *qdrantAPICredentials) GetAPIKey() string {
	return q.APIKey
}

func (q *qdrantAPICredentials) GetAPIKeyRequireTransportSecurity() bool {
	return q.APIKeyRequireTransportSecurity
}
