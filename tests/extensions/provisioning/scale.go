package provisioning

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/rancher/tfp-automation/config"
	framework "github.com/rancher/tfp-automation/framework/set"
	"github.com/stretchr/testify/require"
)

// Scale is a function that will run terraform apply and scale the provisioned
// cluster, according to user's desired amount.
func Scale(t *testing.T, clusterName, poolName string, terraformOptions *terraform.Options, clusterConfig *config.TerratestConfig) {
	err := framework.ConfigTF(nil, clusterConfig, clusterName, poolName, "")
	require.NoError(t, err)

	terraform.Apply(t, terraformOptions)
}
