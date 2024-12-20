package rbac

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/rancher/shepherd/clients/rancher"
	ranchFrame "github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/rancher/tfp-automation/config"
	"github.com/rancher/tfp-automation/defaults/configs"
	"github.com/rancher/tfp-automation/framework"
	cleanup "github.com/rancher/tfp-automation/framework/cleanup/rancher2"
	qase "github.com/rancher/tfp-automation/pipeline/qase/results"
	"github.com/rancher/tfp-automation/tests/extensions/provisioning"
	rb "github.com/rancher/tfp-automation/tests/extensions/rbac"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RBACTestSuite struct {
	suite.Suite
	client           *rancher.Client
	session          *session.Session
	rancherConfig    *rancher.Config
	terraformConfig  *config.TerraformConfig
	terratestConfig  *config.TerratestConfig
	terraformOptions *terraform.Options
}

func (r *RBACTestSuite) SetupSuite() {
	testSession := session.NewSession()
	r.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client

	rancherConfig := new(rancher.Config)
	ranchFrame.LoadConfig(configs.Rancher, rancherConfig)

	r.rancherConfig = rancherConfig

	terraformConfig := new(config.TerraformConfig)
	ranchFrame.LoadConfig(config.TerraformConfigurationFileKey, terraformConfig)

	r.terraformConfig = terraformConfig

	terratestConfig := new(config.TerratestConfig)
	ranchFrame.LoadConfig(config.TerratestConfigurationFileKey, terratestConfig)

	r.terratestConfig = terratestConfig

	terraformOptions := framework.Rancher2Setup(r.T(), r.rancherConfig, r.terraformConfig, r.terratestConfig)
	r.terraformOptions = terraformOptions

	provisioning.GetK8sVersion(r.T(), r.client, r.terratestConfig, r.terraformConfig, configs.DefaultK8sVersion)
}

func (r *RBACTestSuite) TestTfpRBAC() {
	nodeRolesDedicated := []config.Nodepool{config.EtcdNodePool, config.ControlPlaneNodePool, config.WorkerNodePool}

	tests := []struct {
		name     string
		rbacRole config.Role
	}{
		{"Cluster Owner", config.ClusterOwner},
		{"Project Owner", config.ProjectOwner},
	}

	for _, tt := range tests {
		terratestConfig := *r.terratestConfig
		terratestConfig.Nodepools = nodeRolesDedicated

		tt.name = tt.name + " Module: " + r.terraformConfig.Module

		testUser, testPassword, clusterName, poolName := configs.CreateTestCredentials()

		r.Run((tt.name), func() {
			defer cleanup.ConfigCleanup(r.T(), r.terraformOptions)

			adminClient, err := provisioning.FetchAdminClient(r.T(), r.client)
			require.NoError(r.T(), err)

			provisioning.Provision(r.T(), r.client, r.rancherConfig, r.terraformConfig, &terratestConfig, testUser, testPassword, clusterName, poolName, r.terraformOptions, nil)
			provisioning.VerifyCluster(r.T(), adminClient, clusterName, r.terraformConfig, &terratestConfig)

			rb.RBAC(r.T(), r.client, r.rancherConfig, r.terraformConfig, &terratestConfig, testUser, testPassword, clusterName, poolName, r.terraformOptions, tt.rbacRole)
		})
	}

	if r.terratestConfig.LocalQaseReporting {
		qase.ReportTest()
	}
}

func TestTfpRBACTestSuite(t *testing.T) {
	suite.Run(t, new(RBACTestSuite))
}
