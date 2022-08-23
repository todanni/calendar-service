package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	vault "github.com/hashicorp/vault/api"
	"golang.org/x/oauth2"
)

const (
	credentialsPath = "calendar-tokens"
)

// Credentials - The Google OAuth response after code exchange flow
// Each user's credentials should be stored in a secure long-lived storage
// and retrieved from it when the API is called for the event list to be refreshed
type Credentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type CredentialsRepo interface {
	GetUserCredentials(ctx context.Context, email string) (oauth2.Token, error)
	SaveUserCredentials(ctx context.Context, email string, credentials oauth2.Token) error
}

type credsRepo struct {
	// Depending on what we decide to use,
	// this will be the object that lets us write to the persistence storage
	client *vault.Client
}

func NewCredentialsRepo(client *vault.Client) CredentialsRepo {
	return &credsRepo{
		client: client,
	}
}

func (c *credsRepo) GetUserCredentials(ctx context.Context, email string) (oauth2.Token, error) {
	secret, err := c.client.KVv2(credentialsPath).Get(ctx, email)
	if err != nil {
		log.Fatalf(
			"Unable to read the super secret password from the vault: %v",
			err,
		)
		return oauth2.Token{}, err
	}

	var creds oauth2.Token
	expiryTime, err := time.Parse(time.RFC3339Nano, secret.Data["expiry"].(string))
	if err != nil {
		fmt.Println("cannot parse expiry time" + err.Error())
	}

	creds.AccessToken = secret.Data["access_token"].(string)
	creds.RefreshToken = secret.Data["refresh_token"].(string)
	creds.TokenType = secret.Data["token_type"].(string)
	creds.Expiry = expiryTime

	return creds, nil
}

func (c *credsRepo) SaveUserCredentials(ctx context.Context, email string, creds oauth2.Token) error {
	credsMap := map[string]interface{}{
		"access_token":  creds.AccessToken,
		"refresh_token": creds.RefreshToken,
		"token_type":    creds.TokenType,
		"expiry":        creds.Expiry.Format(time.RFC3339Nano),
	}

	fmt.Println(credsMap["expiry"])
	fmt.Println(creds.Expiry)

	_, err := c.client.KVv2(credentialsPath).Put(ctx, email, credsMap)
	if err != nil {
		log.Fatalf("Unable to write secret: %v to the vault", err)
		return err
	}

	return nil
}
