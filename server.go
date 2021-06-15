package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/nipe0324/gqlgen-todos-example/dataloader"
	"github.com/nipe0324/gqlgen-todos-example/db"
	"github.com/nipe0324/gqlgen-todos-example/graph"
	"github.com/nipe0324/gqlgen-todos-example/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	conn, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}

	resolver := graph.NewResolver(conn)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	r := chi.NewRouter()
	r.Use(dataLoaderMiddleware(conn)) // ミドルウェアを利用

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r)) // chiのrouterを利用
}

// DataLoaderをContextに設定するミドルウェア
func dataLoaderMiddleware(conn *db.SQLHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			loader := dataloader.NewDataLoader(conn)
			ctx := context.WithValue(r.Context(), dataloader.CtxKey, loader)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
