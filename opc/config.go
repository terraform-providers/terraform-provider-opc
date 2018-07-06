package opc

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/hashicorp/go-oracle-terraform/storage"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
)

// Config represents the provider configuarion attributes
type Config struct {
	User             string
	Password         string
	IdentityDomain   string
	Endpoint         string
	MaxRetries       int
	Insecure         bool
	StorageEndpoint  string
	StorageServiceID string
}

// Client holder for the OPC (OCI Classic) API Clients
type Client struct {
	computeClient *compute.Client
	storageClient *storage.Client
}

// Client gets the OPC (OCI Classic) API Clients
func (c *Config) Client() (*Client, error) {

	userAgentString := fmt.Sprintf("HashiCorp-Terraform-v%s", terraform.VersionString())

	config := opc.Config{
		IdentityDomain: &c.IdentityDomain,
		Username:       &c.User,
		Password:       &c.Password,
		MaxRetries:     &c.MaxRetries,
		UserAgent:      &userAgentString,
	}

	if logging.IsDebugOrHigher() {
		config.LogLevel = opc.LogDebug
		config.Logger = opcLogger{}
	}

	// Setup HTTP Client based on insecure
	httpClient := cleanhttp.DefaultClient()
	if c.Insecure {
		transport := cleanhttp.DefaultTransport()
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		httpClient.Transport = transport
	}

	config.HTTPClient = httpClient

	opcClient := &Client{}

	if c.Endpoint != "" {
		computeEndpoint, err := url.ParseRequestURI(c.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid Compute Endpoint URI: %s", err)
		}
		config.APIEndpoint = computeEndpoint
		computeClient, err := compute.NewComputeClient(&config)
		if err != nil {
			return nil, err
		}
		opcClient.computeClient = computeClient
	}

	if c.StorageEndpoint != "" {
		storageEndpoint, err := url.ParseRequestURI(c.StorageEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid Storage Endpoint URI: %+v", err)
		}
		config.APIEndpoint = storageEndpoint
		if (c.StorageServiceID) != "" {
			config.IdentityDomain = &c.StorageServiceID
		}
		storageClient, err := storage.NewStorageClient(&config)
		if err != nil {
			return nil, err
		}
		opcClient.storageClient = storageClient
	}

	return opcClient, nil
}

type opcLogger struct{}

func (l opcLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.SetFlags(0)
	log.Print(fmt.Sprintf("go-oracle-terraform: %s", strings.Join(tokens, " ")))
}

func (c *Client) getComputeClient() (*compute.Client, error) {
	if c.computeClient == nil {
		return nil, fmt.Errorf("Compute API client has not been initialized. Ensure the `endpoint` for the Compute Classic REST API Endpoint has been declared in the provider configuration.")
	}
	return c.computeClient, nil
}

func (c *Client) getStorageClient() (*storage.Client, error) {
	if c.storageClient == nil {
		return nil, fmt.Errorf("Storage API client has not been initialized. Ensure the `storage_endpoint` for the Object Storage Classic REST API Endpoint has been declared in the provider configuration.")
	}
	return c.storageClient, nil
}
