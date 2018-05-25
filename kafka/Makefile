BROKER_NAME=my-kafka
ZOOKEEPER_NAME=myzk

all: 

install: install_kafka

#install_zk:
	#helm install --name $(ZOOKEEPER_NAME) incubator/zookeeper
install_kafka:
	-helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
	helm install --name $(BROKER_NAME) -f values.yaml incubator/kafka
	#helm install --name $(BROKER_NAME) incubator/kafka

delete: 
	- helm delete --purge $(BROKER_NAME)
	#- helm delete --purge $(ZOOKEEPER_NAME)

testclient:
	kubectl apply -f testclient.yaml

delete_testclient:
	kubectl delete po testclient

conn:
	kubectl exec -ti testclient -- /usr/bin/kafka-console-consumer --zookeeper $(BROKER_NAME)-zookeeper:2181 --topic new_topic
	#kubectl exec -ti testclient -- /usr/bin/kafka-topics --zookeeper $(BROKER_NAME)-zookeeper:2181 --list
	#kubectl exec -ti testclient -- bash


topics:
	kubectl exec -ti testclient -- /usr/bin/kafka-topics --zookeeper $(BROKER_NAME)-zookeeper:2181 --describe
