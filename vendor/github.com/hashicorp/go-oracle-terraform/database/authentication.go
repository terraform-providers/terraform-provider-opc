package database

import (
	"encoding/base64"
	"fmt"
)

// Get a new auth token for the storage client
func (c *DatabaseClient) getAuthenticationHeader() *string {
	usernamePassword := []byte(fmt.Sprintf("%s:%s", *c.client.UserName, *c.client.Password))
	authToken := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(usernamePassword))
	return &authToken
}
