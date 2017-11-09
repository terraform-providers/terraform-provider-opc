package opc

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/hashicorp/go-oracle-terraform/storage"
	"github.com/hashicorp/terraform/helper/logging"
)

type Config struct {
	User             string
	Password         string
	IdentityDomain   string
	Endpoint         string
	MaxRetries       int
	Insecure         bool
	StorageEndpoint  string
	DatabaseEndpoint string
	JavaEndpoint     string
}

type OPCClient struct {
	computeClient  *compute.ComputeClient
	storageClient  *storage.StorageClient
	databaseClient *database.DatabaseClient
	javaClient     *java.JavaClient
}

func (c *Config) Client() (*OPCClient, error) {
	u, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Invalid endpoint URI: %s", err)
	}

	config := opc.Config{
		IdentityDomain: &c.IdentityDomain,
		Username:       &c.User,
		Password:       &c.Password,
		APIEndpoint:    u,
		MaxRetries:     &c.MaxRetries,
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

	computeClient, err := compute.NewComputeClient(&config)
	if err != nil {
		return nil, err
	}

	opcClient := &OPCClient{
		computeClient: computeClient,
	}

	if c.StorageEndpoint != "" {
		storageEndpoint, err := url.ParseRequestURI(c.StorageEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid storage endpoint URI: %+v", err)
		}
		config.APIEndpoint = storageEndpoint
		storageClient, err := storage.NewStorageClient(&config)
		if err != nil {
			return nil, err
		}
		opcClient.storageClient = storageClient
	}

	if c.DatabaseEndpoint != "" {
		databaseEndpoint, err := url.ParseRequestURI(c.DatabaseEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid database endpoint URI: %+v", err)
		}
		config.APIEndpoint = databaseEndpoint
		databaseClient, err := database.NewDatabaseClient(&config)
		if err != nil {
			return nil, err
		}
		opcClient.databaseClient = databaseClient
	}

	if c.JavaEndpoint != "" {
		javaEndpoint, err := url.ParseRequestURI(c.JavaEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid java endpoint URI: %+v", err)
		}
		config.APIEndpoint = javaEndpoint
		javaClient, err := java.NewJavaClient(&config)
		if err != nil {
			return nil, err
		}
		opcClient.javaClient = javaClient
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
