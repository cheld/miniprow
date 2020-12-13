ifeq (,$(GITHUB_SHA))
	GITHUB_SHA=local
endif
ifeq (,$(GCP_PROJECT))
	GCP_PROJECT=$(shell gcloud config get-value project 2> /dev/null)
endif

clean: 
	rm -Rf ./bin
	
build: clean
	go build -v -o ./bin/miniprow -ldflags="-X github.com/cheld/miniprow/pkg/piper/server.version=${GITHUB_SHA}" cmd/miniprow/miniprow.go
	go build -v -o ./bin/boskosctl cmd/boskosctl/boskosctl.go 

docker-clean:
	docker rmi -f cheld/miniprow:latest 2>/dev/null 
	docker rmi -f eu.gcr.io/${GCP_PROJECT}/cicd-bot:latest 2>/dev/null 

docker: docker-clean build 
	docker build . --file Dockerfile --tag cheld/miniprow:latest
	docker build examples/hello --file examples/hello/Dockerfile --tag eu.gcr.io/${GCP_PROJECT}/cicd-bot:latest

docker-push:
	docker push cheld/miniprow:latest
	docker push eu.gcr.io/${GCP_PROJECT}/cicd-bot:latest