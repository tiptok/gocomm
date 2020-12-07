package handler

import "net/http"

func AllowCorsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("options support"))
		}

		next.ServeHTTP(w, r)
	})
}
