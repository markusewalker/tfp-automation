package ranchers

import (
	"os"
	"testing"

	shepherdConfig "github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/config/operations"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/rancher/tests/actions/features"
	infraConfig "github.com/rancher/tests/validation/recurring/infrastructure/config"
	"github.com/rancher/tfp-automation/config"
	"github.com/rancher/tfp-automation/defaults/keypath"
	"github.com/rancher/tfp-automation/framework"
	featureDefaults "github.com/rancher/tfp-automation/framework/set/defaults/features"
	"github.com/rancher/tfp-automation/framework/set/resources/airgap"
	"github.com/rancher/tfp-automation/framework/set/resources/rancher2"
	"github.com/rancher/tfp-automation/tests/extensions/ssh"
	"github.com/stretchr/testify/require"
)

// CreateAirgapRancher is a function that creates an airgap Rancher setup, either via CLI or web application
func CreateAirgapRancher(t *testing.T, provider string) error {
	os.Getenv("CLOUD_PROVIDER_VERSION")

	configPath := os.Getenv("CATTLE_TEST_CONFIG")
	cattleConfig := shepherdConfig.LoadConfigFromFile(configPath)
	rancherConfig, terraformConfig, terratestConfig, standaloneConfig := config.LoadTFPConfigs(cattleConfig)

	if provider != "" {
		terraformConfig.Provider = provider
	}

	_, keyPath := rancher2.SetKeyPath(keypath.AirgapKeyPath, terratestConfig.PathToRepo, terraformConfig.Provider)
	terraformOptions := framework.Setup(t, terraformConfig, terratestConfig, keyPath)

	_, bastion, err := airgap.CreateMainTF(t, terraformOptions, keyPath, rancherConfig, terraformConfig, terratestConfig)
	if err != nil {
		return err
	}

	_, err = operations.ReplaceValue([]string{"terraform", "airgapBastion"}, terraformConfig.AirgapBastion, cattleConfig)
	require.NoError(t, err)

	_, err = operations.ReplaceValue([]string{"terraform", "privateRegistries", "systemDefaultRegistry"}, terraformConfig.PrivateRegistries.SystemDefaultRegistry, cattleConfig)
	require.NoError(t, err)

	_, err = operations.ReplaceValue([]string{"terraform", "privateRegistries", "url"}, terraformConfig.PrivateRegistries.URL, cattleConfig)
	require.NoError(t, err)

	rancherConfig, terraformConfig, terratestConfig, _ = config.LoadTFPConfigs(cattleConfig)
	infraConfig.WriteConfigToFile(os.Getenv(configEnvironmentKey), cattleConfig)

	sshKey, err := os.ReadFile(terraformConfig.PrivateKeyPath)
	require.NoError(t, err)

	_, err = ssh.StartBastionSSHTunnel(bastion, terraformConfig.Standalone.OSUser, sshKey, "8443", standaloneConfig.RancherHostname, "443")
	require.NoError(t, err)

	testSession := session.NewSession()

	client, err := PostRancherSetup(t, terraformOptions, rancherConfig, testSession, terraformConfig.Standalone.RancherHostname, keyPath, false)
	require.NoError(t, err)

	_, err = operations.ReplaceValue([]string{"rancher", "adminToken"}, client.RancherConfig.AdminToken, cattleConfig)
	require.NoError(t, err)

	rancherConfig, terraformConfig, terratestConfig, _ = config.LoadTFPConfigs(cattleConfig)
	infraConfig.WriteConfigToFile(os.Getenv(configEnvironmentKey), cattleConfig)

	if standaloneConfig.FeatureFlags != nil && standaloneConfig.FeatureFlags.Turtles != "" {
		switch standaloneConfig.FeatureFlags.Turtles {
		case featureDefaults.ToggledOff:
			features.UpdateFeatureFlag(client, featureDefaults.Turtles, false)
		case featureDefaults.ToggledOn:
			features.UpdateFeatureFlag(client, featureDefaults.Turtles, true)
		}
	}

	return nil
}
