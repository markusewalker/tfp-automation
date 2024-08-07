package config

import (
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	"github.com/rancher/tfp-automation/config/authproviders"
	"github.com/rancher/tfp-automation/config/nodeproviders"
)

type TestClientName string
type Role string
type PSACT string

const (
	TerraformConfigurationFileKey = "terraform"
	TerratestConfigurationFileKey = "terratest"

	AdminClientName    TestClientName = "Admin User"
	StandardClientName TestClientName = "Standard User"

	ClusterOwner Role = "cluster-owner"
	ProjectOwner Role = "project-owner"

	RancherPrivileged PSACT = "rancher-privileged"
	RancherRestricted PSACT = "rancher-restricted"
)

var EtcdNodePool = Nodepool{
	Etcd:         true,
	Controlplane: false,
	Worker:       false,
	Quantity:     1,
}

var ControlPlaneNodePool = Nodepool{
	Etcd:         false,
	Controlplane: true,
	Worker:       false,
	Quantity:     1,
}

var WorkerNodePool = Nodepool{
	Etcd:         false,
	Controlplane: false,
	Worker:       true,
	Quantity:     1,
}

var ScaleUpEtcdNodePool = Nodepool{
	Etcd:         true,
	Controlplane: false,
	Worker:       false,
	Quantity:     3,
}

var ScaleUpControlPlaneNodePool = Nodepool{
	Etcd:         false,
	Controlplane: true,
	Worker:       false,
	Quantity:     2,
}

var ScaleUpWorkerNodePool = Nodepool{
	Etcd:         false,
	Controlplane: false,
	Worker:       true,
	Quantity:     3,
}

// String stringer for the TestClientName
func (c TestClientName) String() string {
	return string(c)
}

type GoogleAuthEncodedJSON struct {
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url,omitempty" yaml:"auth_provider_x509_cert_url,omitempty"`
	AuthURI                 string `json:"auth_uri,omitempty" yaml:"auth_uri,omitempty"`
	ClientEmail             string `json:"client_email,omitempty" yaml:"client_email,omitempty"`
	ClientID                string `json:"client_id,omitempty" yaml:"client_id,omitempty"`
	ClientX509CertURL       string `json:"client_x509_cert_url,omitempty" yaml:"client_x509_cert_url,omitempty"`
	PrivateKey              string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
	PrivateKeyID            string `json:"private_key_id,omitempty" yaml:"private_key_id,omitempty"`
	ProjectID               string `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	TokenURI                string `json:"token_uri,omitempty" yaml:"token_uri,omitempty"`
	Type                    string `json:"type,omitempty" yaml:"type,omitempty"`
}

type Nodepool struct {
	Quantity         int64  `json:"quantity,omitempty" yaml:"quantity,omitempty"`
	Etcd             bool   `json:"etcd,omitempty" yaml:"etcd,omitempty"`
	Controlplane     bool   `json:"controlplane,omitempty" yaml:"controlplane,omitempty"`
	Worker           bool   `json:"worker,omitempty" yaml:"worker,omitempty"`
	InstanceType     string `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
	DesiredSize      int64  `json:"desiredSize,omitempty" yaml:"desiredSize,omitempty"`
	MaxSize          int64  `json:"maxSize,omitempty" yaml:"maxSize,omitempty"`
	MinSize          int64  `json:"minSize,omitempty" yaml:"minSize,omitempty"`
	MaxPodsContraint int64  `json:"maxPodsContraint,omitempty" yaml:"maxPodsContraint,omitempty"`
}

type PrivateRegistries struct {
	EngineInsecureRegistry string `json:"engineInsecureRegistry,omitempty" yaml:"engineInsecureRegistry,omitempty"`
	Password               string `json:"password,omitempty" yaml:"password,omitempty"`
	URL                    string `json:"url,omitempty" yaml:"url,omitempty"`
	Username               string `json:"username,omitempty" yaml:"username,omitempty"`
	AuthConfigSecretName   string `json:"authConfigSecretName,omitempty" yaml:"authConfigSecretName,omitempty"`
	TLSSecretName          string `json:"tlsSecretName,omitempty" yaml:"tlsSecretName,omitempty"`
	CABundle               string `json:"caBundle,omitempty" yaml:"caBundle,omitempty"`
	Insecure               bool   `json:"insecure,omitempty" yaml:"insecure,omitempty"`
	SystemDefaultRegistry  string `json:"systemDefaultRegistry,omitempty" yaml:"systemDefaultRegistry,omitempty"`
}

type TerraformConfig struct {
	AWSConfig                           nodeproviders.AWSConfig      `json:"awsConfig,omitempty" yaml:"awsConfig,omitempty"`
	AzureConfig                         nodeproviders.AzureConfig    `json:"azureConfig,omitempty" yaml:"azureConfig,omitempty"`
	GoogleConfig                        nodeproviders.GoogleConfig   `json:"googleConfig,omitempty" yaml:"googleConfig,omitempty"`
	LinodeConfig                        nodeproviders.LinodeConfig   `json:"linodeConfig,omitempty" yaml:"linodeConfig,omitempty"`
	VsphereConfig                       nodeproviders.VsphereConfig  `json:"vsphereConfig,omitempty" yaml:"vsphereConfig,omitempty"`
	ADConfig                            authproviders.ADConfig       `json:"adConfig,omitempty" yaml:"adConfig,omitempty"`
	AzureADConfig                       authproviders.AzureADConfig  `json:"azureADConfig,omitempty" yaml:"azureADConfig,omitempty"`
	GithubConfig                        authproviders.GithubConfig   `json:"githubConfig,omitempty" yaml:"githubConfig,omitempty"`
	OktaConfig                          authproviders.OktaConfig     `json:"oktaConfig,omitempty" yaml:"oktaConfig,omitempty"`
	OpenLDAPConfig                      authproviders.OpenLDAPConfig `json:"openLDAPConfig,omitempty" yaml:"openLDAPConfig,omitempty"`
	AuthProvider                        string                       `json:"authProvider,omitempty" yaml:"authProvider,omitempty"`
	CloudCredentialName                 string                       `json:"cloudCredentialName,omitempty" yaml:"cloudCredentialName,omitempty"`
	DefaultClusterRoleForProjectMembers string                       `json:"defaultClusterRoleForProjectMembers,omitempty" yaml:"defaultClusterRoleForProjectMembers,omitempty"`
	EnableNetworkPolicy                 bool                         `json:"enableNetworkPolicy,omitempty" yaml:"enableNetworkPolicy,omitempty"`
	ETCD                                *rkev1.ETCD                  `json:"etcd,omitempty" yaml:"etcd,omitempty"`
	ETCDRKE1                            *management.ETCDService      `json:"etcdRKE1,omitempty" yaml:"etcdRKE1,omitempty"`
	HostnamePrefix                      string                       `json:"hostnamePrefix,omitempty" yaml:"hostnamePrefix,omitempty"`
	MachineConfigName                   string                       `json:"machineConfigName,omitempty" yaml:"machineConfigName,omitempty"`
	Module                              string                       `json:"module,omitempty" yaml:"module,omitempty"`
	NetworkPlugin                       string                       `json:"networkPlugin,omitempty" yaml:"networkPlugin,omitempty"`
	NodeTemplateName                    string                       `json:"nodeTemplateName,omitempty" yaml:"nodeTemplateName,omitempty"`
	PrivateRegistries                   *PrivateRegistries           `json:"privateRegistries,omitempty" yaml:"privateRegistries,omitempty"`
}

type Scaling struct {
	ScaledDownNodeCount int64      `json:"scaledDownNodeCount,omitempty" yaml:"scaledDownNodeCount,omitempty"`
	ScaledDownNodepools []Nodepool `json:"scaledDownNodepools,omitempty" yaml:"scaledDownNodepools,omitempty"`
	ScaledUpNodeCount   int64      `json:"scaledUpNodeCount,omitempty" yaml:"scaledUpNodeCount,omitempty"`
	ScaledUpNodepools   []Nodepool `json:"scaledUpNodepools,omitempty" yaml:"scaledUpNodepools,omitempty"`
}

type Snapshots struct {
	CreateSnapshot               bool   `json:"createSnapshot,omitempty" yaml:"createSnapshot,omitempty"`
	RestoreSnapshot              bool   `json:"restoreSnapshot,omitempty" yaml:"restoreSnapshot,omitempty"`
	SnapshotName                 string `json:"snapshotName,omitempty" yaml:"snapshotName,omitempty"`
	UpgradeKubernetesVersion     string `json:"upgradeKubernetesVersion,omitempty" yaml:"upgradeKubernetesVersion,omitempty"`
	SnapshotRestore              string `json:"snapshotRestore,omitempty" yaml:"snapshotRestore,omitempty"`
	ControlPlaneConcurrencyValue string `json:"controlPlaneConcurrencyValue,omitempty" yaml:"controlPlaneConcurrencyValue,omitempty"`
	WorkerConcurrencyValue       string `json:"workerConcurrencyValue,omitempty" yaml:"workerConcurrencyValue,omitempty"`
}

type TerratestConfig struct {
	KubernetesVersion         string     `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	UpgradedKubernetesVersion string     `json:"upgradedKubernetesVersion,omitempty" yaml:"upgradedKubernetesVersion,omitempty"`
	NodeCount                 int64      `json:"nodeCount,omitempty" yaml:"nodeCoun,omitemptyt"`
	Nodepools                 []Nodepool `json:"nodepools,omitempty" yaml:"nodepools,omitempty"`
	ScalingInput              Scaling    `json:"scalingInput,omitempty" yaml:"scalingInput,omitempty"`
	PSACT                     string     `json:"psact,omitempty" yaml:"psact,omitempty"`
	SnapshotInput             Snapshots  `json:"snapshotInput,omitempty" yaml:"snapshotInput,omitempty"`
}
