package sanity

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/tfp-automation/defaults/configs"
)

// StandaloneConfigCleanup is a function that will run terraform destroy and cleanup Terraform resources.
func StandaloneConfigCleanup(t *testing.T, terraformOptions *terraform.Options) {
	rancherConfig := new(rancher.Config)
	config.LoadConfig(configs.Rancher, rancherConfig)

	if *rancherConfig.Cleanup {
		terraform.Destroy(t, terraformOptions)
		StandaloneConfigCleanupTF()
	}
}