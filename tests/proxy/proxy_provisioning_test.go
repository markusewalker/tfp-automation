package proxy

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/rancher/shepherd/clients/rancher"
	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	"github.com/rancher/shepherd/extensions/token"
	ranchFrame "github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/rancher/tfp-automation/config"
	"github.com/rancher/tfp-automation/defaults/configs"
	"github.com/rancher/tfp-automation/framework"
	"github.com/rancher/tfp-automation/framework/cleanup"
	"github.com/rancher/tfp-automation/framework/set/resources/proxy"
	resources "github.com/rancher/tfp-automation/framework/set/resources/proxy"
	"github.com/rancher/tfp-automation/framework/set/resources/rancher2"
	qase "github.com/rancher/tfp-automation/pipeline/qase/results"
	"github.com/rancher/tfp-automation/tests/extensions/provisioning"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TfpProxyProvisioningTestSuite struct {
	suite.Suite
	client                     *rancher.Client
	session                    *session.Session
	rancherConfig              *rancher.Config
	terraformConfig            *config.TerraformConfig
	terratestConfig            *config.TerratestConfig
	standaloneTerraformOptions *terraform.Options
	terraformOptions           *terraform.Options
	adminUser                  *management.User
	proxyBastion               string
}

func (p *TfpProxyProvisioningTestSuite) TearDownSuite() {
	keyPath := proxy.KeyPath()
	cleanup.Cleanup(p.T(), p.standaloneTerraformOptions, keyPath)
}

func (p *TfpProxyProvisioningTestSuite) SetupSuite() {
	p.terraformConfig = new(config.TerraformConfig)
	ranchFrame.LoadConfig(config.TerraformConfigurationFileKey, p.terraformConfig)

	p.terratestConfig = new(config.TerratestConfig)
	ranchFrame.LoadConfig(config.TerratestConfigurationFileKey, p.terratestConfig)

	keyPath := proxy.KeyPath()
	standaloneTerraformOptions := framework.Setup(p.T(), p.terraformConfig, p.terratestConfig, keyPath)
	p.standaloneTerraformOptions = standaloneTerraformOptions

	proxyBastion, err := resources.CreateMainTF(p.T(), p.standaloneTerraformOptions, keyPath, p.terraformConfig, p.terratestConfig)
	require.NoError(p.T(), err)

	p.proxyBastion = proxyBastion
}

func (p *TfpProxyProvisioningTestSuite) TfpSetupSuite(terratestConfig *config.TerratestConfig, terraformConfig *config.TerraformConfig) {
	testSession := session.NewSession()
	p.session = testSession

	rancherConfig := new(rancher.Config)
	ranchFrame.LoadConfig(configs.Rancher, rancherConfig)

	p.rancherConfig = rancherConfig

	adminUser := &management.User{
		Username: "admin",
		Password: rancherConfig.AdminPassword,
	}

	p.adminUser = adminUser

	userToken, err := token.GenerateUserToken(adminUser, p.rancherConfig.Host)
	require.NoError(p.T(), err)

	client, err := rancher.NewClient(userToken.Token, testSession)
	require.NoError(p.T(), err)

	p.client = client

	rancherConfig.AdminToken = userToken.Token

	keyPath := rancher2.SetKeyPath()
	terraformOptions := framework.Setup(p.T(), terraformConfig, terratestConfig, keyPath)
	p.terraformOptions = terraformOptions
}

func (p *TfpProxyProvisioningTestSuite) TestTfpNoProxyProvisioning() {
	nodeRolesDedicated := []config.Nodepool{config.EtcdNodePool, config.ControlPlaneNodePool, config.WorkerNodePool}

	tests := []struct {
		name      string
		nodeRoles []config.Nodepool
		module    string
	}{
		{"RKE1", nodeRolesDedicated, "ec2_rke1"},
		{"RKE2", nodeRolesDedicated, "ec2_rke2"},
		{"K3S", nodeRolesDedicated, "ec2_k3s"},
	}

	for _, tt := range tests {
		terratestConfig := *p.terratestConfig
		terratestConfig.Nodepools = tt.nodeRoles
		terraformConfig := *p.terraformConfig
		terraformConfig.Module = tt.module
		terraformConfig.Proxy.ProxyBastion = ""

		p.TfpSetupSuite(&terratestConfig, &terraformConfig)

		provisioning.GetK8sVersion(p.T(), p.client, &terratestConfig, &terraformConfig, configs.DefaultK8sVersion)

		tt.name = tt.name + " Kubernetes version: " + terratestConfig.KubernetesVersion
		testUser, testPassword, clusterName, poolName := configs.CreateTestCredentials()

		p.Run((tt.name), func() {
			keyPath := rancher2.SetKeyPath()
			defer cleanup.Cleanup(p.T(), p.terraformOptions, keyPath)

			provisioning.Provision(p.T(), p.client, p.rancherConfig, &terraformConfig, &terratestConfig, testUser, testPassword, clusterName, poolName, p.terraformOptions, nil)
			provisioning.VerifyCluster(p.T(), p.client, clusterName, &terraformConfig, &terratestConfig)
		})
	}

	if p.terratestConfig.LocalQaseReporting {
		qase.ReportTest()
	}
}

func (p *TfpProxyProvisioningTestSuite) TestTfpProxyProvisioning() {
	nodeRolesDedicated := []config.Nodepool{config.EtcdNodePool, config.ControlPlaneNodePool, config.WorkerNodePool}

	tests := []struct {
		name      string
		nodeRoles []config.Nodepool
		module    string
	}{
		{"RKE1", nodeRolesDedicated, "ec2_rke1"},
		{"RKE2", nodeRolesDedicated, "ec2_rke2"},
		{"K3S", nodeRolesDedicated, "ec2_k3s"},
	}

	for _, tt := range tests {
		terratestConfig := *p.terratestConfig
		terratestConfig.Nodepools = tt.nodeRoles
		terraformConfig := *p.terraformConfig
		terraformConfig.Module = tt.module
		terraformConfig.Proxy.ProxyBastion = p.proxyBastion

		p.TfpSetupSuite(&terratestConfig, &terraformConfig)

		provisioning.GetK8sVersion(p.T(), p.client, &terratestConfig, &terraformConfig, configs.DefaultK8sVersion)

		tt.name = tt.name + " Kubernetes version: " + terratestConfig.KubernetesVersion
		testUser, testPassword, clusterName, poolName := configs.CreateTestCredentials()

		p.Run((tt.name), func() {
			keyPath := rancher2.SetKeyPath()
			defer cleanup.Cleanup(p.T(), p.terraformOptions, keyPath)

			provisioning.Provision(p.T(), p.client, p.rancherConfig, &terraformConfig, &terratestConfig, testUser, testPassword, clusterName, poolName, p.terraformOptions, nil)
			provisioning.VerifyCluster(p.T(), p.client, clusterName, &terraformConfig, &terratestConfig)
		})
	}

	if p.terratestConfig.LocalQaseReporting {
		qase.ReportTest()
	}
}

func TestTfpProxyProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(TfpProxyProvisioningTestSuite))
}