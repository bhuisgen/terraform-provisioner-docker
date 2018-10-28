package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provisioner-docker/docker"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: docker.Provisioner})
}
