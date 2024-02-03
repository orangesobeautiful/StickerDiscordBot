package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
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
		edge.To("request_message_embed", Embed.Type),
	}
}

type ChatMessage struct {
	Type string `json:"type"`

	Content string `json:"content"`
}

type ChatMessageRequestArgument struct {
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`

	MaxTokens uint `json:"max_tokens,omitempty"`

	N uint `json:"n,omitempty"`

	PresencePenalty float32 `json:"presence_penalty,omitempty"`

	Temperature float32 `json:"temperature,omitempty"`

	TopP float32 `json:"top_p,omitempty"`
}

type ChatHistoryDetail struct {
	ent.Schema
}

func (ChatHistoryDetail) Fields() []ent.Field {
	return []ent.Field{
		field.String("model"),
		field.JSON("request_argument", ChatMessageRequestArgument{}),
		field.JSON("full_request_message", []ChatMessage{}),
		field.Uint("prompt_tokens"),
		field.Uint("completion_tokens"),
		field.Float("prompt_price").
			GoType(decimal.Decimal{}).
			SchemaType(
				map[string]string{
					dialect.MySQL:    "decimal(6,2)",
					dialect.Postgres: "numeric",
				},
			),
		field.Float("completion_price").
			GoType(decimal.Decimal{}).
			SchemaType(
				map[string]string{
					dialect.MySQL:    "decimal(6,2)",
					dialect.Postgres: "numeric",
				},
			),
	}
}

func (ChatHistoryDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ref", ChatHistory.Type).
			Ref("detail").
			Unique().
			Required(),
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
