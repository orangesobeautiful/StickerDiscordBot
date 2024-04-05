package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type WebLoginSession struct {
	ent.Schema
}

func (WebLoginSession) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
	}
}

func (WebLoginSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("discord_user", DiscordUser.Type).
			Ref("web_login_session").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
