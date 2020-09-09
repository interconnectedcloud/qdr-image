#
# To run this sample, you must have a running router -> broker topology
#
export QPID_JMS_TRANSACTION_ROUTER_URL='amqp://127.0.0.1:5672'
mvn test
