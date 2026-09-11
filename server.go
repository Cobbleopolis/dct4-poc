package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cobbleopolis/dragoncontimer/ent"
	"github.com/cobbleopolis/dragoncontimer/ent/station"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello, World!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DB_CONN_STR")

	fmt.Printf("Connection String: %s\n", connStr)

	sc := ent.SchemaConfig{
		Station: "dct",
	}

	client, err := ent.Open("postgres", connStr, ent.AlternateSchema(sc))
	if err != nil {
		log.Fatalf("Error connecting to the postgres: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed to create the db schema: %v", err)
	}

	testStation := client.Station.
		Create().
		SetName("Test Station").
		SaveX(ctx)

	client.Station.
		UpdateOne(testStation).
		SetStatus(station.StatusCHECKED_OUT).
		SetCurrentPlayer("John Doe").
		SetCheckoutTime(time.Now()).
		SaveX(ctx)
}
