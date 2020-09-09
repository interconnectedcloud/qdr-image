package router_broker

import (
	"os"
	"testing"
	"time"

	"github.com/interconnectedcloud/qdr-image/test/k8s/integration/router-broker/common"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/apps/broker"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/apps/router"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/k8s"
	"github.com/skupperproject/skupper/pkg/kube"
	"github.com/skupperproject/skupper/test/utils/base"
	"github.com/skupperproject/skupper/test/utils/constants"
	"gotest.tools/assert"
	core "k8s.io/api/core/v1"
)

// TestMain helps parsing the common test flags and running package level tests
func TestMain(m *testing.M) {
	base.ParseFlags()
	os.Exit(m.Run())
}

// Setup deploys a single router->broker topology
// under the provided namespaceId
func Setup(t *testing.T, namespaceId string) {
	var err error
	var ctx *base.ClusterContext

	// Use the provided namespace id
	needs := base.ClusterNeeds{
		NamespaceId:     namespaceId,
		PublicClusters:  1,
		PrivateClusters: 0,
	}

	t.Logf("Building ClusterTestRunner for %s", needs.NamespaceId)
	common.TestRunner.BuildOrSkip(t, needs, nil)
	ctx, err = common.TestRunner.GetPublicContext(1)
	assert.Assert(t, err)

	//
	// - Creating the namespace
	//
	err = ctx.CreateNamespace()
	assert.Assert(t, err)

	//
	// - Deploying topology
	//
	t.Logf("%s - starting topology setup", time.Now().String())
	deployments := ctx.VanClient.KubeClient.AppsV1().Deployments(ctx.Namespace)
	services := ctx.VanClient.KubeClient.CoreV1().Services(ctx.Namespace)
	configMaps := ctx.VanClient.KubeClient.CoreV1().ConfigMaps(ctx.Namespace)

	//
	// - Deploying the Broker
	//
	brokerLabels := map[string]string{"app": "broker"}
	brokerQueues := []string{"trx.testQueue"}

	// Preparing the Deployment
	brokerDep, err := broker.NewDeployment(ctx.Namespace, broker.ActiveMQArtemisDeploymentOpts{
		Name:   "broker",
		Labels: brokerLabels,
		Queues: brokerQueues,
	})
	_, err = deployments.Create(brokerDep)
	assert.Assert(t, err)

	// Preparing the Service
	brokerSvc := k8s.NewServiceClusterIP("broker", ctx.Namespace, []int{5672}, brokerLabels, brokerLabels)
	_, err = services.Create(brokerSvc)
	assert.Assert(t, err)

	//
	// - Deploying the Router
	//
	routerLabels := map[string]string{"app": "router"}

	// Preparing the Router Configuration
	routerConfig := router.QpidDispatch{
		Id:   "router",
		Role: router.RouterRoleInterior,
		Listeners: []router.Listener{
			{Host: "0.0.0.0", Port: 5672},
		},
		Connectors: []router.Connector{
			{Name: "broker", Host: "broker", Port: 5672, RouteContainer: true},
		},
		Addresses: []router.Address{
			{Prefix: "trx", Waypoint: true},
		},
		LinkRoutes: []router.LinkRoute{
			{Prefix: "$coordinator", Direction: "in", Connection: "broker"},
			{Prefix: "$coordinator", Direction: "out", Connection: "broker"},
			{Prefix: "trx", Direction: "in", Connection: "broker"},
			{Prefix: "trx", Direction: "out", Connection: "broker"},
		},
	}
	routerConfigMap := router.NewConfigMap(routerConfig, ctx.Namespace, routerLabels)
	_, err = configMaps.Create((*core.ConfigMap)(routerConfigMap))
	assert.Assert(t, err)

	// Deploying the router instance
	routerDep, err := router.NewDeployment(ctx.Namespace, routerConfig, router.QpidDispatchDeploymentOpts{
		Labels:    routerLabels,
		ConfigMap: routerConfigMap,
	})
	assert.Assert(t, err)
	_, err = deployments.Create(routerDep)
	assert.Assert(t, err)

	// Creating the router service
	routerSvc := k8s.NewServiceClusterIP("router", ctx.Namespace, []int{5672}, routerLabels, routerLabels)
	_, err = services.Create(routerSvc)
	assert.Assert(t, err)

	// Waiting on both broker and router pods
	brokerPods, err := kube.GetDeploymentPods(brokerDep.Name, "app=broker", ctx.Namespace, ctx.VanClient.KubeClient)
	routerPods, err := kube.GetDeploymentPods(routerDep.Name, "app=router", ctx.Namespace, ctx.VanClient.KubeClient)
	assert.Assert(t, err)
	for _, pod := range append(brokerPods, routerPods...) {
		_, err = kube.WaitForPodStatus(ctx.Namespace, ctx.VanClient.KubeClient, pod.Name, core.PodRunning, constants.ImagePullingAndResourceCreationTimeout, constants.DefaultTick)
		assert.Assert(t, err, "pod not running %s", pod.Name)
	}

	t.Logf("%s - setup is complete", time.Now().String())
}

// Teardown deletes the namespace created earlier during Setup
func Teardown(t *testing.T) {
	t.Logf("%s - starting topology teardown", time.Now().String())
	ctx, err := common.TestRunner.GetPublicContext(1)
	assert.Assert(t, err)

	err = ctx.DeleteNamespace()
	assert.Assert(t, err)
	t.Logf("%s - teardown is complete", time.Now().String())
}
