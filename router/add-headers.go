package router

import (
	"net/http"

	"github.com/knadh/koanf/v2"
)

func withHeaders(h http.HandlerFunc, k *koanf.Koanf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", k.String("cors.allowed-origins")) // Or set to a specific origin
		// TODO: in future only keep HEAD for presigned cookies no need for body in presigned cookies
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		h(w, r)
	}
}
