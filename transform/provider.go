package transform

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"transform_group_by_value": dataSourceGroupByValue(),
			"transform_glob_map":       dataSourceGlobMap(),
		},
	}
}
