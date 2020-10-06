package com.redhat.qpid.dispatch.jms.transaction.basic.clients;

import org.apache.qpid.jms.JmsConnectionFactory;

import javax.jms.Connection;
import javax.jms.JMSException;
import javax.jms.Queue;
import javax.jms.Session;

/**
 * Base JMS client class that provides basic and common behavior
 * to producer and consumers.
 */
public class JmsBaseClient {
    public static final long TIMEOUT = 5000L;

    protected String url;
    protected String queueName;

    protected Connection conn;
    protected Session sess;
    protected Queue queue;

    public JmsBaseClient(String url, String queueName) {
        this.url = url;
        this.queueName = queueName;
    }

    /**
     * Connects using the initial url and queue names, then
     * creates a session that can be used by consumers or producers.
     */
    public void connect() {
        try {
            this.conn = new JmsConnectionFactory(this.url).createConnection();
            this.sess = this.conn.createSession(true, Session.SESSION_TRANSACTED);
            this.queue = this.sess.createQueue(this.queueName);
        } catch (JMSException e) {
            e.printStackTrace();
        }
    }

    /**
     * If implementing client is connected, then it closes the
     * underlying JMS connection.
     */
    public void stop() {

        if ( this.conn == null ) {
            return;
        }

        try {
            this.conn.close();
        } catch (JMSException e) {
            e.printStackTrace();
        }

    }

    /**
     * Commits the current session
     */
    public void commit() {
        try {
            this.sess.commit();
        } catch (JMSException e) {
            e.printStackTrace();
        }
    }

    /**
     * Rollbacks current session
     */
    public void rollback() {
        try {
            this.sess.rollback();
        } catch (JMSException e) {
            e.printStackTrace();
        }
    }

}
