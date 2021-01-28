#
# Variables defined here can be used by the test suites
# or either by the client container images that will be
# executed by them.
#
# Note that this is just to reference all environment
# variables that can be defined. The test suite must
# set them in the respective pod/jobs when running 
# against a Kubernetes cluster.
#

# general variables
QPID_DISPATCH_IMAGE="quay.io/interconnectedcloud/qdrouterd:latest"
ACTIVEMQ_ARTEMIS_IMAGE="quay.io/artemiscloud/activemq-artemis-broker:latest"

# clients/java/jms-amqp-tests
JMS_AMQP_TESTS_IMAGE="quay.io/atomictests/jms-amqp-tests:latest"
QPID_JMS_TRANSACTION_ROUTER_URL="amqp://127.0.0.1:5672"
QPID_JMS_TRANSACTION_ADDRESS="trx.testQueue"
