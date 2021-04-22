package tcpadaptor

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/apps/router"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/apps/router/mgmt"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/apps/router/mgmt/entities"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/k8s"
	"github.com/skupperproject/skupper/pkg/kube"
	"github.com/skupperproject/skupper/pkg/utils"
	"github.com/skupperproject/skupper/test/utils/base"
	skconstants "github.com/skupperproject/skupper/test/utils/constants"
	k8s2 "github.com/skupperproject/skupper/test/utils/k8s"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
)

var (
	// meshSizes and dataSizes will be used to compose a matrix for testing
	meshSizes = map[string]int{"one-router": 1, "two-routers": 2, "three-routers": 3}
	dataSizes = []string{"100M", "500M", "1G"}

	iperfServerLabels = map[string]string{"app": "iperf3-server"}
	iperfClientLabels = map[string]string{"app": "iperf3-client"}
)

const (
	iperfImage      = "quay.io/fgiorgetti/iperf3"
	namespaceId     = "tcpadaptor"
	iperfServer     = "iperf3-server"
	iperfRouterName = "iperf3-router"
	timeout         = 120 * time.Second
)

func Setup(t *testing.T, runner base.ClusterTestRunner) {
	// creating namespace
	ctx, err := runner.GetPublicContext(1)
	t.Logf("creating namespace: %s", ctx.Namespace)
	assert.Assert(t, err)
	assert.Assert(t, ctx.CreateNamespace())
}

func Deploy(t *testing.T, runner base.ClusterTestRunner, meshSize int) base.ClusterTestRunner {

	// retrieving cluster context
	ctx, _ := runner.GetPublicContext(1)

	// k8s clients
	deployments := ctx.VanClient.KubeClient.AppsV1().Deployments(ctx.Namespace)
	services := ctx.VanClient.KubeClient.CoreV1().Services(ctx.Namespace)
	configMaps := ctx.VanClient.KubeClient.CoreV1().ConfigMaps(ctx.Namespace)

	/*
	 * iPerf3 server deployment
	 */
	t.Logf("deploying iperf3-server")

	// preparing iperf server deployment
	dep, err := k8s.NewDeployment(iperfServer, ctx.Namespace, k8s.DeploymentOpts{
		Image:         iperfImage,
		Labels:        map[string]string{"app": iperfServer},
		RestartPolicy: v1.RestartPolicyAlways,
		Args:          []string{"-s"},
	})
	assert.Assert(t, err, "error generating iperf3-server deployment")
	_, err = deployments.Create(dep)
	assert.Assert(t, err, "error deploying iperf3-server")

	// creating an iperf3-server service
	svc := k8s.NewServiceClusterIP(iperfServer, ctx.Namespace, []int{5201}, iperfServerLabels, iperfServerLabels)
	_, err = services.Create(svc)
	assert.Assert(t, err, "error creating iperf3-server service")

	// wait for iperf3-server to be running
	waitPodRunning(t, iperfServer, ctx)

	/*
	 * Router mesh deployment
	 */
	t.Logf("deploying router mesh")

	// creating router mesh
	for i := 1; i <= meshSize; i++ {
		createTcpListener := i == 1
		createTcpConnector := i == meshSize
		createRouterConnector := i != meshSize
		qdr := routerConfig(i, createTcpListener, createTcpConnector, createRouterConnector)
		// app=router-1 (or router-2, ...)
		routerLabels := map[string]string{"app": qdr.Id}

		// creating config map
		cm := router.NewConfigMap(qdr, ctx.Namespace, routerLabels)
		_, err = configMaps.Create((*v1.ConfigMap)(cm))
		assert.Assert(t, err)

		// Deploying the router instance
		t.Logf("deploying %s", qdr.Id)
		routerDep, err := router.NewDeployment(ctx.Namespace, qdr, router.QpidDispatchDeploymentOpts{
			DeploymentOpts: k8s.DeploymentOpts{
				Labels: routerLabels,
			},
			ConfigMap: cm,
		})
		assert.Assert(t, err)
		_, err = deployments.Create(routerDep)
		assert.Assert(t, err)

		// Creating the router service
		routerSvc := k8s.NewServiceClusterIP(qdr.Id, ctx.Namespace, []int{55672}, routerLabels, routerLabels)
		_, err = services.Create(routerSvc)
		assert.Assert(t, err)

		// Create the iperf3-router service
		if createTcpListener {
			iperfRouterSvc := k8s.NewServiceClusterIP(iperfRouterName, ctx.Namespace, []int{5201}, routerLabels, routerLabels)
			_, err = services.Create(iperfRouterSvc)
			assert.Assert(t, err)
		}
	}

	// Waiting on routers to be running
	for i := 1; i <= meshSize; i++ {
		routerId := fmt.Sprintf("router-%d", i)
		// Wait for router to be up and running
		waitPodRunning(t, routerId, ctx)
	}
	ctx.KubectlExec("get pods")

	// Wait for mesh size to match
	t.Logf("Waiting for router mesh to have %d nodes", meshSize)
	retryCtx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()
	for i := 1; i <= meshSize; i++ {
		routerLabel := fmt.Sprintf("app=router-%d", i)
		pods, err := kube.GetDeploymentPods("", routerLabel, ctx.Namespace, ctx.VanClient.KubeClient)
		assert.Assert(t, err)
		for _, pod := range pods {
			t.Logf("Validating router network on %s", pod.Name)
			err = utils.RetryWithContext(retryCtx, skconstants.DefaultTick, func() (bool, error) {
				nodes, err := mgmt.QdmanageQuery(ctx.VanClient, ctx.Namespace, pod.Name, pod.Spec.Containers[0].Name, entities.Node{}, nil)
				if err != nil {
					t.Logf("qdmanage query failed with: %v", err)
					return false, nil
				}
				if len(nodes) != meshSize {
					return false, nil
				}
				return true, nil
			})
			assert.Assert(t, err)
		}
	}
	return runner
}

// waitPodRunning waits for pods with selector: "app=${appId}" to be in Running
// state, or until it times out.
func waitPodRunning(t *testing.T, appId string, cluster *base.ClusterContext) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	t.Logf("waiting on %s pod to be running", appId)
	utils.RetryWithContext(ctx, skconstants.DefaultTick, func() (bool, error) {
		pods, err := kube.GetDeploymentPods("", "app="+appId, cluster.Namespace, cluster.VanClient.KubeClient)
		assert.Assert(t, err)
		if len(pods) == 0 {
			return false, nil
		}
		for _, pod := range pods {
			podReady, err := kube.WaitForPodStatus(cluster.Namespace, cluster.VanClient.KubeClient, pod.Name, v1.PodRunning, timeout, skconstants.DefaultTick)
			assert.Assert(t, err)

			// Pod must be ready
			err = utils.RetryWithContext(ctx, skconstants.DefaultTick, func() (bool, error) {
				return kube.IsPodReady(podReady), nil
			})
			assert.Assert(t, err)
		}
		return true, nil
	})
}

func routerConfig(id int, tcpListener, tcpConnector, routerConnector bool) router.QpidDispatch {
	qpidDispatch := router.QpidDispatch{
		Id:   fmt.Sprintf("router-%d", id),
		Role: router.RouterRoleInterior,
		Listeners: []router.Listener{
			{Host: "0.0.0.0", Port: 5672},
		},
		InterRouterListeners: []router.Listener{
			{Host: "0.0.0.0", Port: 55672},
		},
		Logs: []router.Log{
			{
				Module:           "DEFAULT",
				Enable:           "trace+",
				IncludeTimestamp: true,
				IncludeSource:    true,
			},
		},
	}

	if tcpListener {
		qpidDispatch.TcpListeners = []router.TcpListener{
			{Host: "0.0.0.0", Port: "5201", Address: "iperf3"},
		}
	}

	if tcpConnector {
		qpidDispatch.TcpConnectors = []router.TcpConnector{
			{Host: "iperf3-server", Port: "5201", Address: "iperf3"},
		}
	}

	if routerConnector {
		qpidDispatch.InterRouterConnectors = []router.Connector{
			{Host: fmt.Sprintf("router-%d", id+1), Port: 55672},
		}
	}

	return qpidDispatch
}

func Teardown(t *testing.T, runner base.ClusterTestRunner) {
	ctx, _ := runner.GetPublicContext(1)
	t.Logf("deleting namespace: %s", ctx.Namespace)
	ctx.DeleteNamespace()
}

func RunJob(t *testing.T, runner base.ClusterTestRunner, dataSize string, meshSize int) {
	t.Logf("running iperf3-client job")
	ctx, _ := runner.GetPublicContext(1)
	jobs := ctx.VanClient.KubeClient.BatchV1().Jobs(ctx.Namespace)
	job := k8s.NewJob("iperf3-client", ctx.Namespace, k8s.JobOpts{
		Image:        iperfImage,
		BackoffLimit: 0,
		Restart:      v1.RestartPolicyNever,
		Labels:       iperfClientLabels,
		Args:         []string{"-c", iperfRouterName, "-n", dataSize},
	})

	// Retrieving last topology change for all the routers
	topoMapBefore, err := retrieveTopologyChangeMap(ctx, meshSize)
	assert.Assert(t, err)

	// Create iperf3 client Job
	_, err = jobs.Create(job)
	assert.Assert(t, err, "error creating job")

	// Wait for job to finish
	job, jobErr := k8s2.WaitForJob(ctx.Namespace, ctx.VanClient.KubeClient, job.Name, skconstants.ImagePullingAndResourceCreationTimeout)

	// Job logs
	pods, err := kube.GetDeploymentPods("", "app=iperf3-client", ctx.Namespace, ctx.VanClient.KubeClient)
	assert.Assert(t, err)
	logs, err := kube.GetPodContainerLogs(pods[0].Name, "", ctx.Namespace, ctx.VanClient.KubeClient)
	assert.Assert(t, err)
	t.Logf("job logs: %s", logs)

	// verifying router pods (and eventually logs)
	verifyRouterPods(t, runner, meshSize)

	// Retrieving last topology change for all the routers
	topoMapAfter, err := retrieveTopologyChangeMap(ctx, meshSize)
	assert.Assert(t, err)

	// Comparing topology change times
	t.Logf("Validating topology has NOT changed after test completed")
	assert.Assert(t, reflect.DeepEqual(topoMapBefore, topoMapAfter), "topology has changed after client execution")

	// Before verifying job has passed or failed, we need to verify router logs
	assert.Assert(t, jobErr)

	// Assert job completed successfully
	k8s2.AssertJob(t, job)
}

func verifyRouterPods(t *testing.T, runner base.ClusterTestRunner, meshSize int) {
	ctx, _ := runner.GetPublicContext(1)
	for i := 1; i <= meshSize; i++ {
		pods, err := kube.GetDeploymentPods("", fmt.Sprintf("app=router-%d", i), ctx.Namespace, ctx.VanClient.KubeClient)
		assert.Assert(t, err)
		for _, pod := range pods {
			restarts := pod.Status.ContainerStatuses[0].RestartCount
			t.Logf("restart count on pod %s = %d", pod.Name, int(restarts))
			// just adding more verbosity in case of router restarts
			if restarts != int32(0) {
				t.Logf("ERROR - router has restarted. Logs:")
				linesToTail := int64(50)
				logs, _ := kube.GetPodContainerLogsWithOpts(pod.Name, "", ctx.Namespace, ctx.VanClient.KubeClient, v1.PodLogOptions{TailLines: &linesToTail, Previous: true})
				t.Logf(logs)
			}
			assert.Assert(t, restarts == int32(0), "pod %s has been restarted %d times", pod.Name, int(restarts))
		}
	}
}

// retrieveTopologyChangeMap returns a map indexed by router id (1..meshSize) containing
// the lastTopoChange value returned by qdmanage for each respective router.
func retrieveTopologyChangeMap(ctx *base.ClusterContext, meshSize int) (map[int]int, error) {
	topoMap := map[int]int{}
	for i := 1; i <= meshSize; i++ {
		routerLabel := fmt.Sprintf("app=router-%d", i)
		pods, err := kube.GetDeploymentPods("", routerLabel, ctx.Namespace, ctx.VanClient.KubeClient)
		if err != nil {
			return topoMap, err
		}
		for _, pod := range pods {
			nodes, err := mgmt.QdmanageQuery(ctx.VanClient, ctx.Namespace, pod.Name, pod.Spec.Containers[0].Name, entities.Node{}, func(entity entities.Entity) bool {
				node := entity.(entities.Node)
				return node.NextHop == "(self)"
			})
			if err != nil {
				return topoMap, err
			}
			node := nodes[0].(entities.Node)
			topoMap[i] = node.LastTopoChange
		}
	}
	return topoMap, nil
}
