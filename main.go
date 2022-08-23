package main

import (
	b64 "encoding/base64"
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
	decodedCredentials, err := b64.StdEncoding.DecodeString(googleCredentials)

	config, err := google.ConfigFromJSON([]byte(decodedCredentials), calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	credentialsRepo := repo.NewCredentialsRepo(vaultClient)
	router := mux.NewRouter()
	services.NewCredentialsService(credentialsRepo, router, logger, config)
	services.NewCalendarService(credentialsRepo, router, logger, config)

	http.ListenAndServe(":8083", router)
}
