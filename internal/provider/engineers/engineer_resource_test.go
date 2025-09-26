package engineers_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testProviderConfig + `
resource "dob_engineer" "test" {
  name = "user"
  email = "user@liatrio.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("dob_engineer.test", "name", "user"),
					// Verify first order item
					resource.TestCheckResourceAttr("dob_engineer.test", "email", "user@liatrio.com"),
					// Verify first coffee item has Computed attributes filled.
					resource.TestCheckResourceAttrSet("dob_engineer.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testProviderConfig + `
resource "dob_engineer" "test" {
  name = "user"
  email = "user@liatrio.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("dob_engineer.test", "name", "user"),
					resource.TestCheckResourceAttr("dob_engineer.test", "email", "user@liatrio.com"),
					// Verify first coffee item has Computed attributes updated.
					resource.TestCheckResourceAttrSet("dob_engineer.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
