package com.redhat.qpid.dispatch.jms.transaction.basic.clients;

import javax.jms.JMSException;
import javax.jms.MessageConsumer;
import javax.jms.TextMessage;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * A very basic JMS Consumer that consumes one ore more messages
 * from a pre-defined queue and allows clients to commit or rollback.
 *
 */
public class JmsConsumer extends JmsBaseClient {

    private MessageConsumer consumer;
    private AtomicInteger attemptCounter  = new AtomicInteger();
    private AtomicInteger receivedCounter = new AtomicInteger();

    public JmsConsumer(String url, String queueName) {
        super(url, queueName);
    }

    /**
     * Connects and creates a consumer
     */
    public void start() {
        try {
            super.connect();
            this.consumer = this.sess.createConsumer(this.queue);
        } catch (JMSException e) {
            e.printStackTrace();
        }
    }

    /**
     * Perform the given number of attempts to receive messages.
     * @param attempts
     * @return
     */
    public int receiveMessages(int attempts) {
        int count = 0;
        for ( int i = 0; i < attempts; i++ ) {
            count += receiveMessage();
        }
        return count;
    }

    /**
     * Attempt to receive a message within the pre-defined TIMEOUT.
     * @return
     */
    public int receiveMessage() {

        try {
            attemptCounter.incrementAndGet();
            TextMessage msg = (TextMessage) consumer.receive(TIMEOUT);
            if ( msg != null ) {
                this.receivedCounter.incrementAndGet();
                return 1;
            }
        } catch (JMSException e) {
            e.printStackTrace();
        }

        return 0;

    }

}
