package main

import (
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    port := os.Getenv("PORT")

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
    v1Router.Get("/healthz", readinessHandler)  //  检查服务是否准备好
    v1Router.Get("/error", errorHandler)  //  测试错误处理
    // 添加基本路由
    v1Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("welcome"))
    })

    router.Mount("/v1", v1Router)


    log.Println("Server is running on port: %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}