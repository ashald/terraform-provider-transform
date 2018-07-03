package transform

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"reflect"
)

const inputGlobMap = `
locals {
  input = {
    "aaa/bbb/111" = "val1"
    "aaa/ccc/111" = "val1"
    "aaa/ddd/222" = "val2"
  }
}

data "transform_glob_map" "include"         { input="${local.input}" pattern="aaa/*/111"                }
data "transform_glob_map" "exclude"         { input="${local.input}" pattern="aaa/*/111" exclude = true }
data "transform_glob_map" "include_w_sep"   { input="${local.input}" pattern="aaa/*"     separator="/"  }
data "transform_glob_map" "include_wo_sep"  { input="${local.input}" pattern="aaa/*"                    }

output "glob_include"			{ value="${data.transform_glob_map.include.output}"			}
output "glob_exclude" 			{ value="${data.transform_glob_map.exclude.output}" 		}
output "glob_include_w_sep"		{ value="${data.transform_glob_map.include_w_sep.output}"	}
output "glob_include_wo_sep"	{ value="${data.transform_glob_map.include_wo_sep.output}"	}
`

func TestGlobMapDataSource(t *testing.T) {
	glob_include := map[string]string{
		"aaa/ccc/111": "val1",
		"aaa/bbb/111": "val1",
	}
	glob_exclude := map[string]string{
		"aaa/ddd/222": "val2",
	}
	glob_include_w_sep := map[string]string{}
	glob_include_wo_sep := map[string]string{
		"aaa/bbb/111": "val1",
		"aaa/ccc/111": "val1",
		"aaa/ddd/222": "val2",
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: inputGlobMap,
				Check: resource.ComposeTestCheckFunc(
					testMapOutputEquals("glob_include", glob_include),
					testMapOutputEquals("glob_exclude", glob_exclude),
					testMapOutputEquals("glob_include_w_sep", glob_include_w_sep),
					testMapOutputEquals("glob_include_wo_sep", glob_include_wo_sep),
				),
			},
		},
	})
}

func testMapOutputEquals(name string, expected map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		output := state.RootModule().Outputs[name]

		if output == nil {
			return fmt.Errorf("missing '%s' output", name)
		}

		outputMap := map[string]string{}

		for k, v := range output.Value.(map[string]interface{}) {
			outputMap[k] = v.(string)
		}

		if !reflect.DeepEqual(outputMap, expected) {
			return fmt.Errorf("output '%s' value '%v' does not match expected '%v'", name, outputMap, expected)
		}
		return nil
	}
}
