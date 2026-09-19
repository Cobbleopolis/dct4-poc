package main

import (
	"context"
	"log"
	"net/http"
	"os"

	// "time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cobbleopolis/dragoncontimer/ent"

	// "github.com/cobbleopolis/dragoncontimer/ent/station"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	log.Println("Hello, World!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DB_CONN_STR")

	log.Printf("Connection String: %s\n", connStr)

	sc := ent.SchemaConfig{
		Station: "dct",
	}

	client, err := ent.Open("postgres", connStr, ent.AlternateSchema(sc))
	if err != nil {
		log.Fatalf("Error connecting to the postgres: %v", err)
	}
	defer client.Close()

	ctx := ent.NewContext(context.Background(), client)
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed to create the db schema: %v", err)
	}

	// testStation := client.Station.
	// 	Create().
	// 	SetName("Test Station").
	// 	SaveX(ctx)

	// client.Station.
	// 	UpdateOne(testStation).
	// 	SetStatus(station.StatusCHECKED_OUT).
	// 	SetCurrentPlayer("John Doe").
	// 	SetCheckoutTime(time.Now()).
	// 	SaveX(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// srv := handler.New(NewExecutableSchema(Config{Resolvers: &Resolver{}}))
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: &Resolver{
			client: client,
		},
	}))
	// srv := handler.NewDefaultServer(dct.NewExecutableSchema())

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	srv.Use(entgql.Transactioner{TxOpener: client})
	// srv.Use(entgql.Qu)

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		graphql.GetOperationContext(ctx).DisableIntrospection = false
		return next(ctx)
	})

	http.Handle("/gql/playground", playground.Handler("GraphQL playground", "/gql"))
	http.Handle("/gql", srv)

	log.Printf("connect to http://localhost:%s/gql/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
