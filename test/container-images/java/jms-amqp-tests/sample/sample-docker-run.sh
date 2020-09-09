function cleanup() {
    echo "Stopping router and broker containers"
    docker stop broker-jms-amqp-tests router-jms-amqp-tests test-jms-amqp-tests| true

    echo "Removing docker network"
    docker network rm jms-amqp-tests
}

MAX_ATTEMPTS=120
DELAY=1

# clean up the containers and network
trap cleanup EXIT

echo "Creating docker network"
docker network create jms-amqp-tests | true

echo "Starting router and broker containers"
docker run --rm --name broker-jms-amqp-tests -d --net=jms-amqp-tests \
       -e 'AMQ_USER=admin' -e 'AMQ_PASSWORD=admin' \
       -e 'AMQ_EXTRA_ARGS=--queues trx/queue1' \
       quay.io/artemiscloud/activemq-artemis-broker:latest | true

docker run --rm --name router-jms-amqp-tests -d --net=jms-amqp-tests \
       -e 'QDROUTERD_CONF=/opt/router/qdrouterd.conf' \
       -v `pwd`/qdrouterd.conf:/opt/router/qdrouterd.conf:Z \
       quay.io/interconnectedcloud/qdrouterd:latest | true

echo "Waiting for router connection to broker to become active"
attempts=0
while [[ `docker exec router-jms-amqp-tests qdstat -c | grep broker | wc -l` -eq 0 ]]; do
    echo -n "."
    let attempts+=1
    [[ $attempts -eq ${MAX_ATTEMPTS} ]] && echo && echo "Connection not established, exiting..." && exit 1
    sleep ${DELAY}
done
echo

echo "Running test container"
docker run --rm --name test-jms-amqp-tests --net=jms-amqp-tests \
       -e 'QPID_JMS_TRANSACTION_ROUTER_URL=amqp://router-jms-amqp-tests:5672' \
       -v `pwd`/result:/opt/jms-amqp-tests/target/surefire-reports/:Z \
       fgiorgetti/jms-amqp-tests

