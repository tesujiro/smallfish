FROM centos
MAINTAINER tesujiro <tesujiro@gmail.com>
RUN echo "now building..."
ADD ./csgeo /
ADD ./template /template
ADD ./static /static
ADD ./ssl /ssl
EXPOSE 80
CMD ["/csgeo","-kafka_server","my-kafka-kafka","-kafka_port","9092"]
