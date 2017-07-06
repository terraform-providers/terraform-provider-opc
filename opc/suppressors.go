package opc

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// Suppress Diff on mismatched case
func suppressCaseDifferences(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}
