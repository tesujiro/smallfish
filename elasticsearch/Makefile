BROKER_NAME=smallfish
define SSD_YAML
kind: StorageClass
apiVersion: extensions/v1beta1
metadata:
  name: ssd
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-ssd
endef
export SSD_YAML

all: 

install: 
	-helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
	helm install --name $(BROKER_NAME) -f values.yaml incubator/elasticsearch 
	#helm install --name $(BROKER_NAME) incubator/elasticsearch --set client.heapSize=64m,data.heapSize=256m,data.persistence.size=1Gi,master.heapSize=64m,master.persistence.size=64Mi,master.replicas=1
	#echo "$$SSD_YAML" |helm install --name $(BROKER_NAME) incubator/elasticsearch --set client.heapSize=128M,data.heapSize=128M,data.persistence.size=1Gi,master.heapSize=128M,master.persistence.size=64Mi,data.storageClass=ssd,data.storage=1Gi

delete: 
	- helm delete --purge $(BROKER_NAME)

port_forward:
	export POD_NAME=$(kubectl get pods --namespace default -l "app=elasticsearch,component=client,release=smallfish" -o jsonpath="{.items[0].metadata.name}")
	echo "Visit http://127.0.0.1:9200 to use Elasticsearch"
	kubectl port-forward --namespace default $POD_NAME 9200:9200


