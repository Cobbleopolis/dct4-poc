//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithConfigPath("gqlgen.yml"),
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("ent.graphql"),
		entgql.WithWhereInputs(true),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
	}
	genConfig := &gen.Config{
		Features: []gen.Feature{
			gen.FeatureSchemaConfig,
		},
	}
	if err := entc.Generate("./ent/schema", genConfig, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
