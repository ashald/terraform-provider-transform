package transform

import (
	"github.com/gobwas/glob"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGlobMap() *schema.Resource {
	return &schema.Resource{
		Read: readFilter,

		Schema: map[string]*schema.Schema{
			// "Inputs"
			FieldInput: {
				Type:     schema.TypeMap,
				Required: true,
			},
			FieldPattern: {
				Type:     schema.TypeString,
				Required: true,
			},
			FieldSeparator: {
				Type:     schema.TypeString,
				Optional: true,
			},
			FieldExclude: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// "Outputs"
			FieldOutput: {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func readFilter(d *schema.ResourceData, m interface{}) error {
	input := d.Get(FieldInput).(map[string]interface{})

	pattern := d.Get(FieldPattern).(string)
	separator := d.Get(FieldSeparator).(string)
	exclude := d.Get(FieldExclude).(bool)

	inputConverted := make(map[string]string)

	for k, v := range input {
		inputConverted[k] = v.(string)
	}

	result, err := filterMap(inputConverted, pattern, separator, exclude)
	if err != nil {
		return err
	}

	d.Set(FieldOutput, result)

	hash, err := getSHA256(result)
	if err != nil {
		return err
	}

	d.SetId(hash)

	return nil
}

func filterMap(input map[string]string, pattern string, separator string, exclude bool) (map[string]string, error) {
	patternCompiled, err := glob.Compile(pattern, []rune(separator)...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	for k, v := range input {
		// Go does not have XOR for bools
		// but != works like XOR for bools :cactus:
		if patternCompiled.Match(k) != exclude {
			result[k] = v
		}
	}
	return result, nil
}
