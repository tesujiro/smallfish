POD="my-grafana"

install:
	helm install stable/grafana --name $(POD) -f values.yaml

delete:
	helm delete $(POD) --purge

portforward:
	kubectl port-forward `kubectl get pods | grep $(POD) | awk '{print $$1}'` 3000


