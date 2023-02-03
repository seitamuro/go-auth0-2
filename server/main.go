package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	v1 "github.com/seitamuro/go-auth0-2/server/handlers/v1"
	"github.com/seitamuro/go-auth0-2/server/handlers/v1/users/me"
	"github.com/seitamuro/go-auth0-2/server/middlewares/auth0"
)

func main() {
	// .envファイル読み込み
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("API_URL")
	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_CLIENT_ID")

	// 公開鍵を取得する
	jwks, err := auth0.FetchJWKS(domain)
	if err != nil {
		log.Fatal(err)
	}
	// domain, clientID, 公開鍵を元にJWTMiddlewareを作成する
	jwtMiddleware, err := auth0.NewMiddleware(domain, clientID, jwks)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	// /v1へのリクエストの場合のハンドラを登録
	mux.HandleFunc("/v1", v1.HandleIndex)
	// /v1/users/meへのリクエストの場合のハンドラを登録
	// auth0.UseJWTでラップし、ハンドラを呼ぶ前にJWT認証を行う
	mux.Handle("/v1/users/me", auth0.UseJWT(http.HandlerFunc(me.HandleIndex)))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	// リクエスト前にJWTMiddlewareをContextに埋め込むためのMiddlewareを追加
	wrappedMux := auth0.WithJWTMiddleware(jwtMiddleware)(mux)
	wrappedMux = c.Handler(wrappedMux)

	log.Printf("Listening on %s", url)
	if err := http.ListenAndServe(url, wrappedMux); err != nil {
		log.Fatal(err)
	}
}
