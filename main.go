package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sclng-backend-test-v1/githubclient"
	"sclng-backend-test-v1/service"
	"strconv"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := newConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", pongHandler)
	router.HandleFunc("/repos", githubRepositoriesHandler)
	// Initialize web server and configure the following routes:
	// GET /repos
	// GET /stats

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}

// Fetch repositories from GitHub API and write the JSON response to the client
func githubRepositoriesHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	client := githubclient.NewClient()
	repoService := service.NewRepositoryService(client)

	languageFilter := r.URL.Query().Get("language")
	licenseFilter := r.URL.Query().Get("license")
	limitStr := r.URL.Query().Get("limit")
	limit := 100

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}
	fmt.Println("Fetching repositories...")
	repos, err := repoService.FetchAndProcessRepositories(languageFilter, licenseFilter, limit)
	if err != nil {
		log.WithError(err).Error(w, fmt.Sprintf("Failed to fetch repositories: %v", err), http.StatusInternalServerError)
		return nil
	}

	jsonResponse, err := json.MarshalIndent(map[string]interface{}{
		"repositories": repos,
	}, "", "  ")
	if err != nil {
		log.WithError(err).Error(w, fmt.Sprintf("Failed to prepare response %v", err), http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	return nil
}
