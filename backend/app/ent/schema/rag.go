package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type RAGReferencePool struct {
	ent.Schema
}

func (RAGReferencePool) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

func (RAGReferencePool) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", DiscordGuild.Type).
			Ref("rag_reference_pool").
			Unique().
			Required(),
		edge.From("chatroom", Chatroom.Type).
			Ref("rag_reference_pool"),
		edge.To("texts", RAGReferenceText.Type),
	}
}

type RAGReferenceText struct {
	ent.Schema
}

func (RAGReferenceText) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("text"),
		field.Bytes("embed_content"),
		field.JSON("embed_metadata", EmbedMetadata{}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

type EmbedMetadata map[string]any

func (RAGReferenceText) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ref", RAGReferencePool.Type).
			Ref("texts").
			Unique().
			Required(),
	}
}
