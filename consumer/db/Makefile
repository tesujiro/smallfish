MODULE=csdb
GOOS=linux

install:
	kubectl apply -f deploy.yaml

delete:
	kubectl delete -f deploy.yaml

all: build deploy
#all: test build deploy

build: build_module build_container

build_module: ./*.go 
	GOOS=$(GOOS) go build -o $(MODULE) .
	chmod +x $(MODULE)
	
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

private_registry:
	docker pull registry:2.3.0
	docker run -d -p 5000:5000 -v /var/opt:/var/lib/registry --restart=always --name private_registry registry:2.3.0

log:
	kubectl logs `kubectl get pod | grep $(MODULE) | awk '{print $$1}'` --follow

test:
	GOCACHE=off go test . -v
