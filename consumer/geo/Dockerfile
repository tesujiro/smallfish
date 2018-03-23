FROM centos
MAINTAINER tesujiro <tesujiro@gmail.com>
RUN echo "now building..."
ADD ./csgeo /
EXPOSE 80
CMD ["/csgeo","-dbserver","cockroachdb-public","-dbport","26257","-kafka_server","my-kafka-kafka","-kafka_port","9092"]
