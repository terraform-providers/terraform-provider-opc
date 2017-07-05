package storage

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_ClientTestUser   = "test-user"
	_ClientTestDomain = "test-domain"
)

func newAuthenticatingServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ORACLE_LOG") != "" {
			log.Printf("[DEBUG] Received request: %s, %s\n", r.Method, r.URL)
		}

		if r.URL.Path == "/authenticate/" {
			http.SetCookie(w, &http.Cookie{Name: "testAuthCookie", Value: "cookie value"})
		} else {
			handler(w, r)
		}
	}))
}

func getStorageTestClient(c *opc.Config) (*StorageClient, error) {
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
		apiEndpoint, err := url.Parse(os.Getenv("OPC_STORAGE_ENDPOINT"))
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

	return NewStorageClient(c)
}
