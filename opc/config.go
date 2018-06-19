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

type Config struct {
	User             string
	Password         string
	IdentityDomain   string
	Endpoint         string
	MaxRetries       int
	Insecure         bool
	StorageEndpoint  string
	StorageServiceId string
}

type OPCClient struct {
	computeClient *compute.Client
	storageClient *storage.Client
}

func (c *Config) Client() (*OPCClient, error) {

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

	opcClient := &OPCClient{}

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
		if (c.StorageServiceId) != "" {
			config.IdentityDomain = &c.StorageServiceId
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
