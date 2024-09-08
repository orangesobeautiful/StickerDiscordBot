package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type AdditionalMigrate struct {
	ent.Schema
}

func (AdditionalMigrate) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.Int("version").
			Positive(),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Immutable().
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
