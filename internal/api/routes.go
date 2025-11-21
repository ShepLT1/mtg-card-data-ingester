package api

import "net/http"

type MethodHandler map[string]http.HandlerFunc

func (m MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := m[r.Method]; ok {
		handler(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// NewRouter initializes all routes and returns a ServeMux
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// /ingest endpoint
	mux.Handle("/ingest", MethodHandler{
		http.MethodPost: ingestPOSTHandler,
		http.MethodGet:  ingestGETHandler,
	})

	return mux
}
