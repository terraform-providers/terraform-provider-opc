package database

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

func GetDatabaseTestClient(c *opc.Config) (*DatabaseClient, error) {
	// Build up config with default values if omitted

	if c.IdentityDomain == nil {
		domain := os.Getenv("OPC_IDENTITY_DOMAIN")
		c.IdentityDomain = &domain
	}

	if c.Username == nil {
		username := os.Getenv("OPC_USERNAME")
		c.Username = &username
	}

	if c.Password == nil {
		password := os.Getenv("OPC_PASSWORD")
		c.Password = &password
	}

	if c.APIEndpoint == nil {
		apiEndpoint, err := url.Parse(os.Getenv("OPC_DATABASE_ENDPOINT"))
		if err != nil {
			return nil, err
		}
		c.APIEndpoint = apiEndpoint
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				TLSHandshakeTimeout: 120 * time.Second},
		}
	}

	return NewDatabaseClient(c)
}
