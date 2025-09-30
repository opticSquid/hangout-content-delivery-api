package router

import (
	"net/http"

	"github.com/knadh/koanf/v2"
)

func withCORS(h http.HandlerFunc, k *koanf.Koanf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", k.String("cors.allowed-origins")) // Or set to a specific origin
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		h(w, r)
	}
}
