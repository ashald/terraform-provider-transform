package main

import (
	"github.com/ashald/terraform-provider-transform/transform"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return transform.Provider()
		},
	})
}
