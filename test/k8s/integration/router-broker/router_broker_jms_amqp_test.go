//+build router_broker transaction integration

package router_broker

import (
	"os"

	"github.com/interconnectedcloud/qdr-image/test/k8s/integration/router-broker/common"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/constants"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/k8s"
	"github.com/skupperproject/skupper/test/utils/base"
	skconstants "github.com/skupperproject/skupper/test/utils/constants"
	k8s2 "github.com/skupperproject/skupper/test/utils/k8s"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"

	"testing"
)

const JmsAmqpTestsImageEnvVar = "JMS_AMQP_TESTS_IMAGE"

func TestJmsAmqp(t *testing.T) {

	defer Teardown(t)
	base.HandleInterruptSignal(t, func(t *testing.T) {
		Teardown(t)
	})
	Setup(t, "jms-amqp")

	// Cluster context
	ctx, err := common.TestRunner.GetPublicContext(1)
	assert.Assert(t, err)

	// Preparing the jms-amqp-tests job
	jmsAmqpTests := k8s.NewJob("jms-amqp-tests", ctx.Namespace, k8s.JobOpts{
		Image:        utils.StrDefault(os.Getenv(JmsAmqpTestsImageEnvVar), constants.JmsAmqpTestsImage),
		BackoffLimit: 1,
		Restart:      v1.RestartPolicyNever,
		Env: map[string]string{
			"QPID_JMS_TRANSACTION_ROUTER_URL": "amqp://router:5672",
		},
		Labels: map[string]string{
			"app": "jms-amqp-tests",
		},
	})

	// Running the job
	_, err = ctx.VanClient.KubeClient.BatchV1().Jobs(ctx.Namespace).Create(jmsAmqpTests)
	assert.Assert(t, err)

	// Waiting for job to complete
	job, err := k8s2.WaitForJob(ctx.Namespace, ctx.VanClient.KubeClient, jmsAmqpTests.Name, skconstants.ImagePullingAndResourceCreationTimeout)
	assert.Assert(t, err)
	k8s2.AssertJob(t, job)

}
