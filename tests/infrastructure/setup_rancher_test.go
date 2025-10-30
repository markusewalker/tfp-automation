package infrastructure

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/rancher/shepherd/clients/rancher"
	shepherdConfig "github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/rancher/tests/actions/featureflags"
	"github.com/rancher/tfp-automation/config"
	"github.com/rancher/tfp-automation/defaults/keypath"
	"github.com/rancher/tfp-automation/defaults/providers"
	"github.com/rancher/tfp-automation/framework"
	"github.com/rancher/tfp-automation/framework/set/defaults"
	"github.com/rancher/tfp-automation/framework/set/resources/rancher2"
	resources "github.com/rancher/tfp-automation/framework/set/resources/sanity"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RancherTestSuite struct {
	suite.Suite
	session          *session.Session
	terraformConfig  *config.TerraformConfig
	terratestConfig  *config.TerratestConfig
	standaloneConfig *config.Standalone
	cattleConfig     map[string]any
	rancherConfig    *rancher.Config
	terraformOptions *terraform.Options
}

func (i *RancherTestSuite) TestCreateRancher() {
	i.cattleConfig = shepherdConfig.LoadConfigFromFile(os.Getenv(shepherdConfig.ConfigEnvironmentKey))
	i.rancherConfig, i.terraformConfig, i.terratestConfig, i.standaloneConfig = config.LoadTFPConfigs(i.cattleConfig)

	_, keyPath := rancher2.SetKeyPath(keypath.SanityKeyPath, i.terratestConfig.PathToRepo, i.terraformConfig.Provider)
	terraformOptions := framework.Setup(i.T(), i.terraformConfig, i.terratestConfig, keyPath)
	i.terraformOptions = terraformOptions

	_, err := resources.CreateMainTF(i.T(), i.terraformOptions, keyPath, i.rancherConfig, i.terraformConfig, i.terratestConfig)
	require.NoError(i.T(), err)

	if i.terraformConfig.Provider != providers.Linode && i.terraformConfig.Provider != providers.Vsphere {
		testSession := session.NewSession()
		i.session = testSession

		client, err := PostRancherSetup(i.T(), i.terraformOptions, i.rancherConfig, i.session, i.terraformConfig.Standalone.RancherHostname, keyPath, false, false)
		require.NoError(i.T(), err)

		if i.standaloneConfig.FeatureFlags != nil && i.standaloneConfig.FeatureFlags.Turtles != "" {
			switch i.standaloneConfig.FeatureFlags.Turtles {
			case defaults.ToggledOff:
				featureflags.UpdateFeatureFlag(client, defaults.Turtles, false)
			case defaults.ToggledOn:
				featureflags.UpdateFeatureFlag(client, defaults.Turtles, true)
			}
		}
	}
}

func TestRancherTestSuite(t *testing.T) {
	suite.Run(t, new(RancherTestSuite))
}
