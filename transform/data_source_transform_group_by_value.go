package transform

import (
	"github.com/hashicorp/terraform/helper/schema"
	"sort"
)

func dataSourceGroupByValue() *schema.Resource {
	return &schema.Resource{
		Read: groupByValue,

		Schema: map[string]*schema.Schema{
			// "Inputs"
			FieldInput: {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			FieldExtract: {
				Type:     schema.TypeString,
				Required: true,
			},
			// "Outputs"
			FieldItems: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func groupByValue(d *schema.ResourceData, m interface{}) error {
	input := d.Get(FieldInput).(map[string]interface{})
	extract := d.Get(FieldExtract).(string)

	var result []string

	for k, v := range input {
		valueStr := v.(string)
		if valueStr == extract {
			result = append(result, k)
		}
	}

	sort.Strings(result)

	d.Set(FieldItems, result)

	hash, err := getSHA256(result)
	if err != nil {
		return err
	}

	d.SetId(hash)

	return nil
}
