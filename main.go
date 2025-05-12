package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/djchanahcjd/go-rss/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
    DB *db.Queries
}

func main() {
    // 测试rss功能
    feed, error := urlToFeed("https://wagslane.dev/index.xml")
    if error!= nil {
        log.Fatal("Error fetching feed:", error)
    }
    log.Println("Feed Title:", feed.Channel.Title)

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("PORT is not found in the environment")
    }

    dbUrl := os.Getenv("DB_URL")
    if dbUrl == "" {
        log.Fatal("DB_URL is not found in the environment")
    }

    conn, err := sql.Open("postgres", dbUrl)
    if err != nil {
        log.Fatal("Cannot connect to db:", err)
    }

    db := db.New(conn)
    apiCfg := apiConfig{
        DB: db,
    }

    go startScraping(db, 10, time.Minute)

    router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
    }))


    v1Router := chi.NewRouter()
    v1Router.Get("/healthz", handlerReadiness)  //  检查服务是否准备好
    v1Router.Get("/error", errorHandler)  //  测试错误处理

    v1Router.Post("/users", apiCfg.createUserHandler)
    v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.getUserHandler))

    v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.createFeedHandler))
    v1Router.Get("/feeds", apiCfg.getAllFeedsHandler)

    v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.createFeedFollowsHandler))
    v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.getFeedFollowsByUserHandler))
    v1Router.Delete("/feed_follows/{feedID}", apiCfg.middlewareAuth(apiCfg.deleteFeedFollowHandler))

    v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.getPostsForUserHandler))

    router.Mount("/v1", v1Router)


    log.Printf("Server is running on port: %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}