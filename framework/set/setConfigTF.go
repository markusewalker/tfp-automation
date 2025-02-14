package set

import (
	"os"
	"strings"

	"github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/tfp-automation/config"
	"github.com/rancher/tfp-automation/defaults/clustertypes"
	"github.com/rancher/tfp-automation/defaults/configs"
	"github.com/rancher/tfp-automation/defaults/keypath"
	"github.com/rancher/tfp-automation/defaults/modules"
	"github.com/rancher/tfp-automation/framework/set/defaults"
	"github.com/rancher/tfp-automation/framework/set/provisioning/airgap"
	custom "github.com/rancher/tfp-automation/framework/set/provisioning/custom/rke1"
	customV2 "github.com/rancher/tfp-automation/framework/set/provisioning/custom/rke2k3s"
	"github.com/rancher/tfp-automation/framework/set/provisioning/hosted"
	"github.com/rancher/tfp-automation/framework/set/provisioning/imported"
	"github.com/rancher/tfp-automation/framework/set/provisioning/multiclusters"
	nodedriver "github.com/rancher/tfp-automation/framework/set/provisioning/nodedriver/rke1"
	nodedriverV2 "github.com/rancher/tfp-automation/framework/set/provisioning/nodedriver/rke2k3s"
	"github.com/rancher/tfp-automation/framework/set/resources/rancher2"
	resources "github.com/rancher/tfp-automation/framework/set/resources/rancher2"
	"github.com/sirupsen/logrus"
)

// ConfigTF is a function that will set the main.tf file based on the module type.
func ConfigTF(client *rancher.Client, rancherConfig *rancher.Config, terraformConfig *config.TerraformConfig, terratestConfig *config.TerratestConfig,
	testUser, testPassword, clusterName, poolName string, rbacRole config.Role, configMap []map[string]any) ([]string, error) {
	module := terraformConfig.Module

	clusterNames := []string{clusterName}
	var file *os.File
	keyPath := rancher2.SetKeyPath(keypath.RancherKeyPath)

	file, err := os.Create(keyPath + configs.MainTF)
	if err != nil {
		logrus.Infof("Failed to reset/overwrite main.tf file. Error: %v", err)
		return nil, err
	}

	defer file.Close()

	newFile, rootBody := resources.SetProvidersAndUsersTF(rancherConfig, terraformConfig, testUser, testPassword, false, configMap)

	rootBody.AppendNewline()

	if terraformConfig.MultiCluster {
		clusterNames, err = multiclusters.SetMultiCluster(client, rancherConfig, configMap, clusterName, newFile, rootBody, file, rbacRole, poolName)
		return clusterNames, err
	} else {
		switch {
		case module == clustertypes.AKS:
			_, err = hosted.SetAKS(terraformConfig, clusterName, terratestConfig.KubernetesVersion, terratestConfig.Nodepools, newFile, rootBody, file)
			return clusterNames, err
		case module == clustertypes.EKS:
			_, err = hosted.SetEKS(terraformConfig, clusterName, terratestConfig.KubernetesVersion, terratestConfig.Nodepools, newFile, rootBody, file)
			return clusterNames, err
		case module == clustertypes.GKE:
			_, err = hosted.SetGKE(terraformConfig, clusterName, terratestConfig.KubernetesVersion, terratestConfig.Nodepools, newFile, rootBody, file)
			return clusterNames, err
		case strings.Contains(module, clustertypes.RKE1) && !strings.Contains(module, defaults.Custom) && !strings.Contains(module, defaults.Airgap) && !strings.Contains(module, defaults.Import):
			_, err = nodedriver.SetRKE1(terraformConfig, clusterName, poolName, terratestConfig.KubernetesVersion, terratestConfig.PSACT, terratestConfig.Nodepools,
				terratestConfig.SnapshotInput, newFile, rootBody, file, rbacRole)
			return clusterNames, err
		case (strings.Contains(module, clustertypes.RKE2) || strings.Contains(module, clustertypes.K3S)) && !strings.Contains(module, defaults.Custom) && !strings.Contains(module, defaults.Airgap) && !strings.Contains(module, defaults.Import):
			_, err = nodedriverV2.SetRKE2K3s(client, terraformConfig, clusterName, poolName, terratestConfig.KubernetesVersion, terratestConfig.PSACT, terratestConfig.Nodepools,
				terratestConfig.SnapshotInput, newFile, rootBody, file, rbacRole)
			return clusterNames, err
		case module == modules.CustomEC2RKE1:
			_, err = custom.SetCustomRKE1(rancherConfig, terraformConfig, terratestConfig, nil, clusterName, newFile, rootBody, file)
			return clusterNames, err
		case module == modules.CustomEC2RKE2 || module == modules.CustomEC2K3s:
			_, err = customV2.SetCustomRKE2K3s(rancherConfig, terraformConfig, terratestConfig, nil, clusterName, newFile, rootBody, file)
			return clusterNames, err
		case module == modules.AirgapRKE1:
			_, err = airgap.SetAirgapRKE1(rancherConfig, terraformConfig, terratestConfig, nil, clusterName, newFile, rootBody, file)
			return clusterNames, err
		case module == modules.AirgapRKE2 || module == modules.AirgapK3S:
			_, err = airgap.SetAirgapRKE2K3s(rancherConfig, terraformConfig, terratestConfig, nil, clusterName, newFile, rootBody, file)
			return clusterNames, err
		case module == modules.ImportRKE1:
			_, err = imported.SetImportedRKE1(rancherConfig, terraformConfig, terratestConfig, clusterName, newFile, rootBody, file)
			return clusterNames, err
		case module == modules.ImportRKE2 || module == modules.ImportK3s:
			_, err = imported.SetImportedRKE2K3s(rancherConfig, terraformConfig, terratestConfig, clusterName, newFile, rootBody, file)
			return clusterNames, err

		default:
			logrus.Errorf("Unsupported module: %v", module)
		}

		return clusterNames, nil
	}
}
