package docker

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProvisioners map[string]terraform.ResourceProvisioner
var testAccProvisioner *schema.Provisioner

func init() {
	testAccProvisioner = Provisioner().(*schema.Provisioner)
	testAccProvisioners = map[string]terraform.ResourceProvisioner{
		"docker": testAccProvisioner,
	}
}

func TestProvisioner(t *testing.T) {
	if err := Provisioner().(*schema.Provisioner).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisioner_impl(t *testing.T) {
	var _ terraform.ResourceProvisioner = Provisioner()
}

//
//func TestAccProvisioner_basic(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers: testAccProviders,
//		CheckDestroy: testCheckDockerConfigDestroy,
//		Steps: []resource.TestStep{
//			resource.TestStep{
//				Config: `
//				resource "docker_config" "foo" {
//					name = "foo-config"
//					data = "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="
//				}
//				`,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("docker_config.foo", "name", "foo-config"),
//					resource.TestCheckResourceAttr("docker_config.foo", "data", "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="),
//				),
//			},
//		},
//	})
//}
