# jms-amqp-tests

This project contains Junit tests that can be executed
against an `AMQP router` that has a connection to a `broker`.

The `router` must provide `$coordinator` links to the respective `broker`, as well as link routes to a waypoint address that reaches the `broker`.

## Environment variable

`QPID_JMS_TRANSACTION_ROUTER_URL` can be defined to provide
the URL that the AMQP producer and consumer will use to communicate
with the respective router component.

The default value to be used, if it is not provided is:
`amqp://127.0.0.1:5672`.

`QPID_JMS_TRANSACTION_ADDRESS` the address to be used for sending and receiving messages. The default value is: `trx.testQueue`.

## Local sample

To run this atomic test locally, go to the `sample` directory and run the `sample-docker-run.sh` script.

