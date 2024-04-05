package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

type ChatMessage struct {
	Type string `json:"type"`

	Content string `json:"content"`
}

type ChatHistoryDetail struct {
	ent.Schema
}

func (ChatHistoryDetail) Fields() []ent.Field {
	return []ent.Field{
		field.String("model"),
		field.JSON("full_request_message", []ChatMessage{}),
		field.Any("request"),
		field.Any("response"),
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
