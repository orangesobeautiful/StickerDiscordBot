package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ChatHistory struct {
	ent.Schema
}

func (ChatHistory) Fields() []ent.Field {
	return []ent.Field{
		field.String("request_message"),
		field.String("reply_message"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

func (ChatHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chatroom", Chatroom.Type).
			Ref("chat_histories").
			Unique().
			Required(),
		edge.To("detail", ChatHistoryDetail.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

type Embed struct {
	ent.Schema
}

func (Embed) Fields() []ent.Field {
	return []ent.Field{
		field.String("input"),
		field.Bytes("content"),
		field.JSON("metadata", map[string]any{}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}
