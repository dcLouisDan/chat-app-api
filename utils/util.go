package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	errerr := WriteJSON(w, status, map[string]string{"error": err.Error()})
	if errerr != nil {
		log.Fatal("Write JSON error")
	}
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		parts := strings.Split(tokenAuth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ""
		}

		tokenString := parts[1]
		return tokenString
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}
