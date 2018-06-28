package transform

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"fmt"
	"reflect"
)

const input = `
output "result" { value="${data.transform_group_by_value.data.items}" }

data "transform_group_by_value" "data" {
  input = {
    key1 = "val1"
	key2 = "val1"
	key3 = "val3"
  }
  extract = "val1"
}
`


func TestGroupByValueDataSource(t *testing.T) {
	expected := []string{"key1", "key2"}

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: input,
				Check: resource.ComposeTestCheckFunc(
					testListOutputEquals("result", expected),
				),
			},
		},
	})
}

func testListOutputEquals(name string, expected []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		output := state.RootModule().Outputs[name]

		if output == nil {
			return fmt.Errorf("missing '%s' output", name)
		}

		var outputList []string

		for _, v := range output.Value.([]interface{}) {
			outputList = append(outputList, v.(string))
		}

		if !reflect.DeepEqual(outputList, expected) {
			return fmt.Errorf("output '%s' value '%v' does not match expected '%v'", name, output.Value, expected)
		}
		return nil
	}
}
