package com.redhat.qpid.dispatch.jms.transaction.basic.clients;

import javax.jms.JMSException;
import javax.jms.MessageProducer;
import javax.jms.TextMessage;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Basic JMS producer implementation that can be used to send one ore more
 * static (small) message(s) to a pre-defined queue and offers ability
 * to commit or rollback current session.
 */
public class JmsProducer extends JmsBaseClient {

    private AtomicInteger messageCounter = new AtomicInteger();
    private MessageProducer producer;

    public JmsProducer(String url, String queueName) {
        super(url, queueName);
    }

    /**
     * Connects and creates a producer to the pre-defined queue.
     */
    public void start() {
        try {
            this.connect();
            this.producer = this.sess.createProducer(this.queue);
        } catch (JMSException e) {
            e.printStackTrace();
        }
    }

    /**
     * Sends the given number of messages (small and static)
     * @param messageCount
     */
    public void sendMessages(int messageCount) {
        for (int i = 0; i < messageCount; i++) {
            sendMessage();
        }
    }

    /**
     * Send a single message to the defined queue.
     */
    public void sendMessage() {
        TextMessage msg = createMessage();
        try {
            this.producer.send(msg);
        } catch (JMSException e) {
            e.printStackTrace();
        }
    }

    /**
     * Creates the static message that will be sent.
     * @return
     */
    private TextMessage createMessage() {

        int i = messageCounter.incrementAndGet();

        TextMessage msg = null;
        try {
            msg = this.sess.createTextMessage("Text message " + i);
            msg.setIntProperty("MESSAGE_NUMBER", i);
        } catch (JMSException e) {
            e.printStackTrace();
        }

        return msg;

    }

}
