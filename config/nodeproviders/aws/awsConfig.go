package aws

type Config struct {
	AMI                   string   `json:"ami,omitempty" yaml:"ami,omitempty"`
	AWSInstanceType       string   `json:"awsInstanceType,omitempty" yaml:"awsInstanceType,omitempty"`
	AWSKeyName            string   `json:"awsKeyName,omitempty" yaml:"awsKeyName,omitempty"`
	AWSVolumeType         string   `json:"awsVolumeType,omitempty" yaml:"awsVolumeType,omitempty"`
	AWSRootSize           int64    `json:"awsRootSize,omitempty" yaml:"awsRootSize,omitempty"`
	AWSSecurityGroupNames []string `json:"awsSecurityGroupNames,omitempty" yaml:"awsSecurityGroupNames,omitempty"`
	AWSSecurityGroups     []string `json:"awsSecurityGroups,omitempty" yaml:"awsSecurityGroups,omitempty"`
	AWSSubnetID           string   `json:"awsSubnetID,omitempty" yaml:"awsSubnetID,omitempty"`
	AWSSubnets            []string `json:"awsSubnets,omitempty" yaml:"awsSubnets,omitempty"`
	AWSVpcID              string   `json:"awsVpcID,omitempty" yaml:"awsVpcID,omitempty"`
	AWSRoute53Zone        string   `json:"awsRoute53Zone,omitempty" yaml:"awsRoute53Zone,omitempty"`
	AWSZoneLetter         string   `json:"awsZoneLetter,omitempty" yaml:"awsZoneLetter,omitempty"`
	ClusterCIDR           string   `json:"clusterCIDR,omitempty" yaml:"clusterCIDR,omitempty"`
	ServiceCIDR           string   `json:"serviceCIDR,omitempty" yaml:"serviceCIDR,omitempty"`
	EnablePrimaryIPv6     bool     `json:"enablePrimaryIPv6,omitempty" yaml:"enablePrimaryIPv6,omitempty" default:"false"`
	HTTPProtocolIPv6      string   `json:"httpProtocolIPv6,omitempty" yaml:"httpProtocolIPv6,omitempty" default:"disabled"`
	IPAddressType         string   `json:"ipAddressType,omitempty" yaml:"ipAddressType,omitempty" default:"ipv4"`
	LoadBalancerType      string   `json:"loadBalancerType,omitempty" yaml:"loadBalancerType,omitempty" default:"ipv4"`
	PrivateAccess         bool     `json:"privateAccess,omitempty" yaml:"privateAccess,omitempty"`
	PublicAccess          bool     `json:"publicAccess,omitempty" yaml:"publicAccess,omitempty"`
	RegistryRootSize      int64    `json:"registryRootSize,omitempty" yaml:"registryRootSize,omitempty"`
	Region                string   `json:"region,omitempty" yaml:"region,omitempty"`
	AWSUser               string   `json:"awsUser,omitempty" yaml:"awsUser,omitempty"`
	TargetType            string   `json:"targetType,omitempty" yaml:"targetType,omitempty" default:"instance"`
	Timeout               string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	WindowsAMI            string   `json:"windowsAMI,omitempty" yaml:"windowsAMI,omitempty"`
	WindowsAWSUser        string   `json:"windowsAWSUser,omitempty" yaml:"windowsAWSUser,omitempty"`
	WindowsAWSPassword    string   `json:"windowsAWSPassword,omitempty" yaml:"windowsAWSPassword,omitempty"`
	WindowsInstanceType   string   `json:"windowsInstanceType,omitempty" yaml:"windowsInstanceType,omitempty"`
	WindowsKeyName        string   `json:"windowsKeyName,omitempty" yaml:"windowsKeyName,omitempty"`
	WindowsVolumeType     string   `json:"windowsVolumeType,omitempty" yaml:"windowsVolumeType,omitempty"`
}
