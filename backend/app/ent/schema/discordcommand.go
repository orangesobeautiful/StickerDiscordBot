package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type DiscordCommand struct {
	ent.Schema
}

func (DiscordCommand) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("discord_id").
			NotEmpty().
			Unique(),
		field.Bytes("sha256_checksum"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Immutable().
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
