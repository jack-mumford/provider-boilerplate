package engineers_test

import (
	"testing"

	providerpkg "terraform-provider-devops/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Local test configuration for acceptance tests in this package.
const testProviderConfig = `
provider "dob" {
  endpoint = "http://localhost:8080"
}
`

var testProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"dob": providerserver.NewProtocol6WithError(providerpkg.New("test")()),
}

func TestAccEngineerDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig + `data "dob_engineer" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Minimal assertion: engineers list exists (count attribute is set)
					resource.TestCheckResourceAttrSet("data.dob_engineer.test", "engineers.#"),

					resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.4.name", "Jack"),
					resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.4.email", "jack@liatrio.com"),
					resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.4.id", "C3UM6"),

					resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.1.name", "Madi"),
					resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.1.email", "madi@liatrio.com"),
					resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.1.id", "EW7D7"),
				),
			},
		},
	})
}
