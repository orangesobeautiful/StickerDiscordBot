package config

type FullTextSearchDatabaseType string

const (
	FullTextSearchDatabaseMeilisearch FullTextSearchDatabaseType = "meilisearch"
)

type FullTextSearchDatabase interface {
	GetType() FullTextSearchDatabaseType

	GetMeilisearch() Meilisearch

	GetToMigrate() bool
}

var _ FullTextSearchDatabase = (*fullTextSearchDatabase)(nil)

type fullTextSearchDatabase struct {
	Meilisearch *meilisearch

	DisableMigrate bool
}

func (f *fullTextSearchDatabase) GetType() FullTextSearchDatabaseType {
	return FullTextSearchDatabaseMeilisearch
}

func (f *fullTextSearchDatabase) GetMeilisearch() Meilisearch {
	return f.Meilisearch
}

func (f *fullTextSearchDatabase) GetToMigrate() bool {
	return !f.DisableMigrate
}

type Meilisearch interface {
	GetAddr() string

	GetAPIKey() string

	GetIndexPrefix() string
}

var _ Meilisearch = (*meilisearch)(nil)

type meilisearch struct {
	Addr string

	APIKey string

	IndexPrefix string
}

func (m *meilisearch) GetAddr() string {
	return m.Addr
}

func (m *meilisearch) GetAPIKey() string {
	return m.APIKey
}

func (m *meilisearch) GetIndexPrefix() string {
	return m.IndexPrefix
}
