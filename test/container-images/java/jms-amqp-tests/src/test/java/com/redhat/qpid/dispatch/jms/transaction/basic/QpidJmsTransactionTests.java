package com.redhat.qpid.dispatch.jms.transaction.basic;

import com.redhat.qpid.dispatch.jms.transaction.basic.clients.JmsConsumer;
import com.redhat.qpid.dispatch.jms.transaction.basic.clients.JmsProducer;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Test;

import static org.junit.Assert.assertEquals;

/**
 * Validates JMS Producers and Consumers ability to commit or rollback
 * messages within a session.
 */
public class QpidJmsTransactionTests {

    public static final String URL = System.getenv().getOrDefault(
                                    "QPID_JMS_TRANSACTION_ROUTER_URL",
                                    "amqp://127.0.0.1:5672");
    public static final String QUEUE_NAME = System.getenv().getOrDefault(
                                            "QPID_JMS_TRANSACTION_ADDRESS",
                                            "trx.testQueue");

    public static JmsProducer producer;
    public static JmsConsumer consumer;

    /**
     * Creates and initializes a consumer and a producer, which one
     * using its own JMS Session.
     */
    @BeforeClass
    public static void setUpClass() {

        // Creating producer and consumer
        producer = new JmsProducer(URL, QUEUE_NAME);
        consumer = new JmsConsumer(URL, QUEUE_NAME);

        // Connecting both
        producer.start();
        consumer.start();

    }

    /**
     * Closes the client's JMS connections.
     */
    @AfterClass
    public static void tearDownClass() {

        // Stopping both clients
        producer.stop();
        consumer.stop();

    }

    /**
     * Producer sends messages on a transactional context
     * without committing. Receivers attempt to receive messages
     * and rollbacks. Next the producer also rollbacks and this test
     * expect that consumers do not get anything.
     */
    @Test
    public void testSendRollback() {

        int msgCount = 2;

        // Attempt to send 2 messages without committing
        producer.sendMessages(msgCount);

        // Expect nothing
        assertEquals(0, consumer.receiveMessage());
        consumer.rollback();

        // Rollback producer and see what consumer gets
        producer.rollback();

        // Attempt to receive again and should not see anything
        assertEquals(0, consumer.receiveMessage());
        consumer.rollback();

    }

    /**
     * Producer sends messages without committing. An attempt
     * to consume messages (prior to commit) is made and it expects
     * that consumer does not see any message.
     *
     * Next the producer commits and consumer attempts to get one message.
     * Test expects consumer to get one message. Then the consumer will
     * rollback and next consume all messages sent by producer, and test
     * will expect all messages to be available (so the first one that was
     * rolled back should be received again). And then consumer will commit.
     *
     * After that it is expected that no more messages are available in
     * the queue.
     *
     */
    @Test
    public void testSendCommitReceiveRollback() {

        int msgCount = 2;

        // Sending messages without committing
        producer.sendMessages(msgCount);

        // Attempt to receive - should not get any
        assertEquals(0, consumer.receiveMessage());
        consumer.rollback();

        // Committing (messages should be available now)
        producer.commit();

        // Receiving messages - expect to see them
        assertEquals(1, consumer.receiveMessage());

        // Rollback - messages should stay in queue
        consumer.rollback();

        // Receive them again and commit this time
        assertEquals(msgCount, consumer.receiveMessages(msgCount));
        consumer.commit();

        // Attempt to receive one last time - should not get any
        assertEquals(0, consumer.receiveMessage());
        consumer.rollback();

    }
}
