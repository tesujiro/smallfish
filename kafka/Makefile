BROKER_NAME=my-kafka

all: 

install:
	helm install --name $(BROKER_NAME) -f values.yaml incubator/kafka

delete:
	helm delete --purge $(BROKER_NAME)

