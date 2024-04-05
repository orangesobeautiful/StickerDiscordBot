package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DiscordUser holds the schema definition for the DiscordUser entity.
type DiscordUser struct {
	ent.Schema
}

// Fields of the DiscordUser.
func (DiscordUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("discord_id").
			NotEmpty(),
		field.String("guild_id").
			NotEmpty(),
		field.String("name"),
		field.String("avatar_url"),
	}
}

func (DiscordUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("discord_id", "guild_id").
			Unique(),
	}
}

// Edges of the DiscordUser.
func (DiscordUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("web_login_session", WebLoginSession.Type),
	}
}
