package domain

type MeilisearchIndexName struct {
	prefix string
}

func NewMeilisearchIndexName(prefix string) MeilisearchIndexName {
	return MeilisearchIndexName{
		prefix: prefix,
	}
}

func (m MeilisearchIndexName) Sticker() string {
	return m.prefix + "sticker"
}
