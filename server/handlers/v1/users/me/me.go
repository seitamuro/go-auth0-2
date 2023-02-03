package me

import (
	"encoding/json"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
	"github.com/seitamuro/go-auth0-2/server/middlewares/auth0"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	subToUsers = map[string]User{
		"auth0|61a8178b21127500715968e2": {
			Name: "kourin",
			Age:  15,
		},
	}
)

// subを元にUserを取得する関数
// 実際のAPIではDBなどに照会し、subに紐づくUserを取得するなどをする
func getUser(sub string) *User {
	user, ok := subToUsers[sub]
	if !ok {
		return nil
	}
	return &user
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	token := auth0.GetJWT(r.Context())

	// token.Claimsをjwt.MapClaimsへ変換
	claims := token.Claims.(jwt.MapClaims)
	// claimsの中にペイロードの情報が入っている
	sub := claims["sub"].(string)

	// userを取得する
	user := getUser(sub)
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// レスポンスを返す
	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
