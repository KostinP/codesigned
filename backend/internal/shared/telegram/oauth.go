package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
)

type AuthData struct {
	ID        int64
	Username  string
	FirstName string
	Hash      string
}

func ParseTelegramAuth(r *http.Request) AuthData {
	return AuthData{
		ID:        parseInt64(r.FormValue("id")),
		Username:  r.FormValue("username"),
		FirstName: r.FormValue("first_name"),
		Hash:      r.FormValue("hash"),
	}
}

func VerifyTelegramAuth(authData AuthData, botToken string) bool {
	secretKey := sha256.Sum256([]byte(botToken))

	dataCheckString := fmt.Sprintf("first_name=%s\nid=%d\nusername=%s", authData.FirstName, authData.ID, authData.Username)

	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(h.Sum(nil))

	return hash == authData.Hash
}

func parseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
