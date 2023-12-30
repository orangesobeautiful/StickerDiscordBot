package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
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
		field.String("channel_id").
			NotEmpty(),
		field.String("name"),
		field.String("avatar_url"),
	}
}

// Edges of the DiscordUser.
func (DiscordUser) Edges() []ent.Edge {
	return nil
}
