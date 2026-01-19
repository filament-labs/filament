package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Wallet holds the schema definition for the Wallet entity.
type Wallet struct {
	ent.Schema
}

// Fields of the Wallet.
func (Wallet) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_default").Default(false),
		field.String("actor_id").Nillable(),
		field.String("name").NotEmpty(),
		field.Bytes("encrypted_seed").Sensitive().NotEmpty(),
		field.Bytes("encrypted_key_json").Sensitive().NotEmpty(),
		field.Bytes("salt").Sensitive().NotEmpty(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Nillable(),
	}
}

// Edges of the Wallet.
func (Wallet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("addresses", Address.Type), // One wallet has many addresses
		//edge.To("transactions", Transaction.Type), // One wallet has many transactions
	}
}
