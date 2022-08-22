package repo

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

const (
	testEmail = "test@todanni.com"
)

func Test_SaveAndGetUserCredentials(t *testing.T) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = os.Getenv("VAULT_ADDRESS")

	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Fatalf("unable to initialize a Vault client: %v", err)
	}
	vaultClient.SetToken(os.Getenv("VAULT_TOKEN"))
	credentialsRepo := NewCredentialsRepo(vaultClient)

	ctx := context.Background()
	expiryTime, err := time.Parse(time.RFC3339, "2022-08-11T16:01:16.265935+01:00")

	err = credentialsRepo.SaveUserCredentials(ctx, testEmail, oauth2.Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
		Expiry:       expiryTime,
	})
	require.NoError(t, err)

	result, err := credentialsRepo.GetUserCredentials(ctx, testEmail)
	require.NoError(t, err)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token", result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Equal(t, expiryTime, result.Expiry)
}
