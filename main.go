package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	vault "github.com/hashicorp/vault/api"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	"github.com/todanni/calendar-service/repo"
	"github.com/todanni/calendar-service/services"
)

const (
	VaultAddress = "https://vault.todanni.com"
)

func main() {
	logger, _ := zap.NewProduction()

	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = VaultAddress

	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Fatalf("unable to initialize a Vault client: %v", err)
	}
	vaultClient.SetToken(os.Getenv("VAULT_TOKEN"))

	googleCredentials := os.Getenv("GOOGLE_CREDENTIALS")
	config, err := google.ConfigFromJSON([]byte(googleCredentials), calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	credentialsRepo := repo.NewCredentialsRepo(vaultClient)
	router := mux.NewRouter()
	services.NewCredentialsService(credentialsRepo, router, logger, config)
	services.NewCalendarService(credentialsRepo, router, logger, config)

	http.ListenAndServe("localhost:8083", router)
}
