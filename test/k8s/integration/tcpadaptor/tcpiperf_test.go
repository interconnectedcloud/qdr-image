package tcpadaptor

import (
	"fmt"
	"testing"

	"github.com/skupperproject/skupper/test/utils/base"
)

func TestTcpIperf(t *testing.T) {
	interrupted := false
	for name, meshSize := range meshSizes {
		for _, dataSize := range dataSizes {
			if interrupted {
				break
			}
			t.Run(fmt.Sprintf("%s-%s", name, dataSize), func(t *testing.T) {
				// Starting test for given iteration
				t.Logf("Testing mesh of %d routers with %s data size", meshSize, dataSize)

				// Defining the ClusterNeeds (1 single namespace)
				needs := base.ClusterNeeds{
					NamespaceId:     namespaceId,
					PublicClusters:  1,
					PrivateClusters: 0,
				}

				var runner base.ClusterTestRunner
				runner = &base.ClusterTestRunnerBase{}
				runner.BuildOrSkip(t, needs, nil)

				// Handling interruptions
				base.HandleInterruptSignal(t, func(t *testing.T) {
					interrupted = true
					Teardown(t, runner)
				})
				defer Teardown(t, runner)

				// Setup and deferred Teardown
				Setup(t, runner)

				// Deploying
				Deploy(t, runner, meshSize)

				// Run client job with appropriate dataSize
				RunJob(t, runner, dataSize, meshSize)
			})
		}
	}
}
