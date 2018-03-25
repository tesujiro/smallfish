
install:
	helm install stable/prometheus --name my-prometheus -f values.yaml

delete:
	helm delete my-prometheus --purge

sportforward_server:
	kubectl port-forward `kubectl get pods -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}"` 9090

sportforward_alertmanager:
	kubectl port-forward `kubectl get pods -l "app=prometheus,component=alertmanager" -o jsonpath="{.items[0].metadata.name}"` 9093

