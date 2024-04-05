package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ImgSaveType int

const (
	ImgSaveTypeNone ImgSaveType = iota

	ImgSaveTypeFullURL

	ImgSaveTypeCloudfare
)

// Image holds the schema definition for the Image entity.
type Image struct {
	ent.Schema
}

// Fields of the Image.
func (Image) Fields() []ent.Field {
	return []ent.Field{
		field.Int("save_type"),
		field.String("save_path"),
		field.Bytes("sha256_checksum"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Image.
func (Image) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sticker", Sticker.Type).
			Ref("images").
			Unique().
			Required(),
	}
}
