package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/heather92115/verdure-admin/graph"
	"github.com/heather92115/verdure-admin/internal/db"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8090"

func main() {
	fmt.Println("Starting the gql server")

	dsn := db.GetDatabaseURL()

	err := db.CreatePool(dsn)
	if err != nil {
		fmt.Printf("Failed DB connections, %v\n", err)
		return
	}

	port := os.Getenv("GQL_PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/admin/gql", playground.Handler("GraphQL playground", "/admin"))
	http.Handle("/admin", srv)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
