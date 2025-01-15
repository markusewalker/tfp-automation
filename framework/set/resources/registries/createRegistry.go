package registries

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/rancher/tfp-automation/config"
	"github.com/rancher/tfp-automation/framework/set/defaults"
	"github.com/rancher/tfp-automation/framework/set/resources/sanity/rke2"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

const (
	authRegistry    = "auth_registry"
	nonAuthRegistry = "non_auth_registry"
)

// CreateAuthenticatedRegistry is a helper function that will create an authenticated registry.
func CreateAuthenticatedRegistry(file *os.File, newFile *hclwrite.File, rootBody *hclwrite.Body, terraformConfig *config.TerraformConfig,
	rke2AuthRegistryPublicDNS string) (*os.File, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	registryScriptPath := filepath.Join(userDir, "go/src/github.com/rancher/tfp-automation/framework/set/resources/registries/auth-registry.sh")

	registryScriptContent, err := os.ReadFile(registryScriptPath)
	if err != nil {
		return nil, err
	}

	_, provisionerBlockBody := rke2.CreateNullResource(rootBody, terraformConfig, rke2AuthRegistryPublicDNS, authRegistry)

	provisionerBlockBody.SetAttributeValue(defaults.Inline, cty.ListVal([]cty.Value{
		cty.StringVal("echo '" + string(registryScriptContent) + "' > /tmp/auth-registry.sh"),
		cty.StringVal("chmod +x /tmp/auth-registry.sh"),
		cty.StringVal("bash -c '/tmp/auth-registry.sh " + terraformConfig.StandaloneRegistry.RegistryUsername + " " +
			terraformConfig.StandaloneRegistry.RegistryPassword + " " + terraformConfig.StandaloneRegistry.RegistryName + " " +
			rke2AuthRegistryPublicDNS + " " + terraformConfig.Standalone.RancherTagVersion + " " +
			terraformConfig.StandaloneRegistry.AssetsPath + " " + terraformConfig.Standalone.OSUser + "'"),
	}))

	_, err = file.Write(newFile.Bytes())
	if err != nil {
		logrus.Infof("Failed to append configurations to main.tf file. Error: %v", err)
		return nil, err
	}

	return file, nil
}

// CreateNonAuthenticatedRegistry is a helper function that will create a non-authenticated registry.
func CreateNonAuthenticatedRegistry(file *os.File, newFile *hclwrite.File, rootBody *hclwrite.Body, terraformConfig *config.TerraformConfig,
	rke2NonAuthRegistryPublicDNS string) (*os.File, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	registryScriptPath := filepath.Join(userDir, "go/src/github.com/rancher/tfp-automation/framework/set/resources/registries/non-auth-registry.sh")

	registryScriptContent, err := os.ReadFile(registryScriptPath)
	if err != nil {
		return nil, err
	}

	_, provisionerBlockBody := rke2.CreateNullResource(rootBody, terraformConfig, rke2NonAuthRegistryPublicDNS, nonAuthRegistry)

	provisionerBlockBody.SetAttributeValue(defaults.Inline, cty.ListVal([]cty.Value{
		cty.StringVal("echo '" + string(registryScriptContent) + "' > /tmp/non-auth-registry.sh"),
		cty.StringVal("chmod +x /tmp/non-auth-registry.sh"),
		cty.StringVal("bash -c '/tmp/non-auth-registry.sh " + terraformConfig.StandaloneRegistry.RegistryName + " " +
			rke2NonAuthRegistryPublicDNS + " " + terraformConfig.Standalone.RancherTagVersion + " " +
			terraformConfig.StandaloneRegistry.AssetsPath + " " + terraformConfig.Standalone.OSUser + "'"),
	}))

	_, err = file.Write(newFile.Bytes())
	if err != nil {
		logrus.Infof("Failed to append configurations to main.tf file. Error: %v", err)
		return nil, err
	}

	return file, nil
}