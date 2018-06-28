package transform

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testProviders map[string]terraform.ResourceProvider

var transformProvider *schema.Provider

func init() {
	transformProvider = Provider().(*schema.Provider)
	testProviders = map[string]terraform.ResourceProvider{
		"transform": transformProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := transformProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
