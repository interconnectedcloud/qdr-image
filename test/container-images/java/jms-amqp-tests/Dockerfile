FROM maven:slim

WORKDIR /opt/jms-amqp-tests

COPY pom.xml ./
COPY src ./src/

RUN mvn -fn clean package && \
    mv /root/.m2 /.m2 && chmod 777 -R /opt/jms-amqp-tests /.m2 && \
    ln -s /.m2 /root/.m2

ENTRYPOINT ["mvn", "test"]
