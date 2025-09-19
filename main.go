package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/djchanahcjd/go-rss/handlers"
	"github.com/djchanahcjd/go-rss/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 测试rss功能
	feed, error := urlToFeed("https://wagslane.dev/index.xml")
	if error != nil {
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
	apiCfg := handlers.ApiConfig{
		DB: db,
	}
	r := init_router(apiCfg)

	go startScraping(db, 10, time.Minute)

	log.Printf("Server is running on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func init_router(apiCfg handlers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()
	// 先设置CORS中间件
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// 静态文件服务
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	// 前端入口页面
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "frontend/index.html")
	})

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlers.HealthzHandler) //  检查服务是否准备好

	v1Router.Post("/users", apiCfg.CreateUser)
	v1Router.Get("/users", apiCfg.AuthMiddleware(apiCfg.GetUser))

	v1Router.Post("/feeds", apiCfg.AuthMiddleware(apiCfg.CreateFeed))
	v1Router.Get("/feeds", apiCfg.GetAllFeeds)

	v1Router.Post("/feed_follows", apiCfg.AuthMiddleware(apiCfg.CreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.AuthMiddleware(apiCfg.GetFeedFollowsByUser))
	v1Router.Delete("/feed_follows/{feedID}", apiCfg.AuthMiddleware(apiCfg.DeleteFeedFollow))

	v1Router.Get("/posts", apiCfg.AuthMiddleware(apiCfg.GetPostsForUser))

	r.Mount("/v1", v1Router)

	return r
}
