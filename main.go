package main

import (
	"io/ioutil"
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

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	router := mux.NewRouter()
	logger, _ := zap.NewProduction()

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = os.Getenv("VAULT_ADDRESS")

	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Fatalf("unable to initialize a Vault client: %v", err)
	}
	vaultClient.SetToken(os.Getenv("VAULT_TOKEN"))
	credentialsRepo := repo.NewCredentialsRepo(vaultClient)

	services.NewCredentialsService(credentialsRepo, router, logger, config)
	services.NewCalendarService(credentialsRepo, router, logger, config)

	http.ListenAndServe("localhost:8083", router)
}
