package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"SWE1-project-data-ingester/internal/ingestion"
)

func ingestPOSTHandler(w http.ResponseWriter, r *http.Request) {

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		endpoint, err := ingestion.FetchDefaultCardEndpoint()
		if err != nil {
			fmt.Printf("Failed to fetch endpoint: %v", err)
			return
		}

		externalCards, err := ingestion.FetchExternalCards(endpoint)
		if err != nil {
			fmt.Printf("Failed to fetch cards: %v", err)
			return
		}

		failures, err := ingestion.IngestCardData(ctx, externalCards)
		if err != nil {
			fmt.Printf("Ingestion error: %v", err)
		}

		fmt.Printf("Ingestion completed with %d failures\n", len(failures))
		for i := 0; i < len(failures); i++ {
			fmt.Println(failures[i])
		}
	}()

	w.WriteHeader(http.StatusAccepted)
  w.Write([]byte("Ingestion started. Check logs for progress."))
}

func ingestGETHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Returning ingestion status (GET)")
}

