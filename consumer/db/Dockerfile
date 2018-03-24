FROM centos
MAINTAINER tesujiro <tesujiro@gmail.com>
RUN echo "now building..."
ADD ./csdb /
EXPOSE 80
CMD ["/csdb","-dbserver","cockroachdb-public","-dbport","26257","-kafka_server","my-kafka-kafka","-kafka_port","9092"]
