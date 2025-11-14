package schema

import "entgo.io/ent"

// CoreUser holds the schema definition for the CoreUser entity.
type CoreUser struct {
	ent.Schema
}

// Fields of the CoreUser.
func (CoreUser) Fields() []ent.Field {
	return nil
}

// Edges of the CoreUser.
func (CoreUser) Edges() []ent.Edge {
	return nil
}
