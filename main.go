package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-opc/opc"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opc.Provider})
}
