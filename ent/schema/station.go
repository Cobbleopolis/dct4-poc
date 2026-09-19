package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Station holds the schema definition for the Station entity.
type Station struct {
	ent.Schema
}

// Fields of the Station.
func (Station) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("name"),
		field.Enum("status").
			Values("DEFAULT", "CHECKED_OUT", "NOT_AVAILABLE").
			Default("DEFAULT"),
		field.Time("checkoutTime").
			Optional().
			Nillable(),
		field.String("currentPlayer").
			Optional().
			Nillable(),
		field.Int("orderPriority").
			Default(0).
			Annotations(
				entgql.OrderField("ORDER_PRIORITY"),
			),
	}
}

// Edges of the Station.
func (Station) Edges() []ent.Edge {
	return nil
}

func (Station) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
