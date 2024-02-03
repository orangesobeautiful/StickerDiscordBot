package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type DiscordGuild struct {
	ent.Schema
}

func (DiscordGuild) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			Unique().
			Immutable(),
	}
}

func (DiscordGuild) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("activate_chatroom", Chatroom.Type).
			Unique(),
		edge.To("chatrooms", Chatroom.Type),
	}
}

type Chatroom struct {
	ent.Schema
}

func (Chatroom) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

func (Chatroom) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", DiscordGuild.Type).
			Ref("chatrooms").
			Unique().
			Required(),
		edge.To("chat_histories", ChatHistory.Type),
	}
}
