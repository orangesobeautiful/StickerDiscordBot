package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Sticker holds the schema definition for the Sticker entity.
type Sticker struct {
	ent.Schema
}

// Fields of the Sticker.
func (Sticker) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

func (Sticker) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

// Edges of the Sticker.
func (Sticker) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", DiscordGuild.Type).
			Ref("sticker").
			Unique().
			Required(),
		edge.To("images", Image.Type),
	}
}
