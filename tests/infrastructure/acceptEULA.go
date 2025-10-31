package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/shepherd/extensions/defaults"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/rancher/tests/actions/pipeline"
	"github.com/rancher/tfp-automation/framework/cleanup"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	kwait "k8s.io/apimachinery/pkg/util/wait"
)

// PostRancherSetup is a helper function that creates a Rancher client and accepts the EULA, if needed
func PostRancherSetup(t *testing.T, terraformOptions *terraform.Options, rancherConfig *rancher.Config, session *session.Session, host,
	keyPath string, isAirgap, isUpgrade bool) (*rancher.Client, error) {
	adminToken, err := CreateAdminToken(t, terraformOptions, rancherConfig, keyPath)
	if err != nil && *rancherConfig.Cleanup {
		cleanup.Cleanup(t, terraformOptions, keyPath)
	}

	rancherConfig.AdminToken = adminToken.Token

	var client *rancher.Client
	err = kwait.PollUntilContextTimeout(context.TODO(), 5*time.Second, defaults.FiveMinuteTimeout, true, func(ctx context.Context) (done bool, err error) {
		client, err = rancher.NewClient(rancherConfig.AdminToken, session)
		if err != nil {
			logrus.Warnf("Failed to create Rancher client: %v. Retrying...", err)
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	client.RancherConfig.AdminToken = rancherConfig.AdminToken
	client.RancherConfig.AdminPassword = rancherConfig.AdminPassword
	client.RancherConfig.Host = host

	if !isUpgrade {
		err = pipeline.PostRancherInstall(client, client.RancherConfig.AdminPassword)
		require.NoError(t, err)
	}

	// The FQDN needs to be set back to the Rancher host URL and not the internal FQDN
	if isAirgap {
		client.RancherConfig.Host = rancherConfig.Host
	}

	return client, nil
}
