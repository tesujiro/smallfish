RELEASE=smallfish-db
ADMIN_PORT=38080
EXPOSE_PORT=36257
install:
	helm install --name $(RELEASE) -f values.yaml stable/cockroachdb
delete:
	helm delete --name $(RELEASE)
expose_service:
	echo access the database by visiting localhost:$(EXPOSE_PORT) with PostgresSQL Client
	kubectl expose service $(RELEASE)-cockroachdb-public --port=$(EXPOSE_PORT) --target-port=26257 --name=$(RELEASE)-cockroachdb-public-client-port --type=LoadBalancer
admin_port:
	echo access the admin UI by visiting http://localhost:$(ADMIN_PORT)/ with web browser
	kubectl port-forward $(RELEASE)-cockroachdb-0 $(ADMIN_PORT):8080
connect:
	kubectl run -it --rm cockroach-client \
		--image=cockroachdb/cockroach \
		--restart=Never \
		--command -- ./cockroach sql --insecure --host $(RELEASE)-cockroachdb-public
build: *.go
	go build *.go
