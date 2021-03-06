MODULE=csgeo
GOOS=linux
SSL_KEY=ssl/server.key
SSL_CSR=ssl/server.csr
SSL_CRT=ssl/server.crt

install:
	kubectl apply -f deploy.yaml

delete:
	kubectl delete -f deploy.yaml

all: build deploy
#all: test build deploy

build: build_module build_container

build_module: ./*.go ./cmd/*.go
	GOOS=$(GOOS) go build -o $(MODULE) ./cmd/csgeo.go
	
build_container:
	docker build -t $(MODULE) .

deploy: deploy_container restart_container

deploy_container:
	docker tag $(MODULE) localhost:5000/$(MODULE)
	docker push localhost:5000/$(MODULE)
	docker rmi localhost:5000/$(MODULE)
	docker pull localhost:5000/$(MODULE)

restart_container:
	for p in `kubectl get pod | grep $(MODULE) | awk '{print $$1}'`;do\
		kubectl delete pod $$p;\
	done
	docker rmi $$(docker images -a --filter "dangling=true" -q) -f

create_secret:
ifndef APIKEY
        $(error APIKEY not set)
endif
	echo $$APIKEY > ./apikey.txt
	kubectl create secret generic google-apikey --from-file=apikey=./apikey.txt
	rm ./apikey.txt
	
delete_secret:
	kubectl delete secret google-apikey
ssh:
	kubectl exec -it `kubectl get pod | grep $(MODULE) | head -1 | awk '{print $$1}'` /bin/bash
test:
	GOCACHE=off go test . -v
logs:
	kubectl logs `kubectl get pod | grep $(MODULE) | awk '{print $$1}'` --follow

ssl_certificate:
	openssl genrsa 2048 > $(SSL_KEY)
	openssl req -new -key $(SSL_KEY) > $(SSL_CSR)
	openssl x509 -days 3650 -req -signkey $(SSL_KEY) < $(SSL_CSR) > $(SSL_CRT)
