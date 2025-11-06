package http

import (
	"log"
	"net/http"
)

func SetupRouter(handler *BalanceHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/balance/credit", withJSONContentType(handler.Credit))
	mux.HandleFunc("/api/balance/transfer", withJSONContentType(handler.Transfer))
	mux.HandleFunc("/api/balance", withJSONContentType(handler.GetBalance))
	mux.HandleFunc("/health", withJSONContentType(healthCheck))

	return loggingMiddleware(mux)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func withJSONContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
