stage ?= local
platform ?= $(shell uname -m)
flags ?=
ecr ?= 475910951137.dkr.ecr.us-west-2.amazonaws.com
OSNAME := $(shell uname)
# islinux := "Linux" == $(OSNAME)

##################
# Help and setup #
##################
install:
	@echo "Download go.mod dependencies"
	go mod download

######################
# Development and CI #
######################
lint:
	# Disable fix flag to see the errors
	golangci-lint run --fix
tidy:
	go mod tidy
clean:
	go clean -testcache
test.all:
	CGO_ENABLED=1 gotestsum --format pkgname --no-summary=skipped -- ./... $(flags)
test.integ: stage=test
test.integ:
	@STAGE=$(stage) PRESERVE_DB=false gotestsum -- -tags 'integration' ./...
ci:
	pwd
	go env GOCACHE
	golangci-lint run --timeout 5m
	CGO_ENABLED=1 GOOS=linux go build ./...
	CGO_ENABLED=1 CI=true gotestsum --format testname --no-summary=skipped -- ./... -v --timeout=30m
	# CI=true go test ./... -v --timeout=30m
build.coreumservice:
	CGO_ENABLED=1 go build -o build/bin/coreumservice ./
start.coreumservice: build.coreumservice
start.coreumservice:
	./build/bin/coreumservice

#####################
# Deployment and CD #
#####################
require-stage:
ifndef STAGE
	$(error STAGE is undefined)
endif
docker.login:
	cat ./env/docker-stably-pass.txt | docker login -u stably --password-stdin
aws.ecr.login:
	@aws ecr get-login-password | docker login --username AWS --password-stdin $(ecr)
docker.build.service.coreumservice:
	# Need to build from the parent directory to include coreumservicemsg
	docker build --platform $(platform) -f ./build/coreumservice.Dockerfile --build-arg STAGE=${STAGE} -t stably/coreumservice:latest .
docker.build.service.coreumservice-current-platform:
	docker build -f ./build/coreumservice.Dockerfile --build-arg STAGE=${STAGE} -t stably/coreumservice:latest .

docker.run.service.coreumservice:
	@docker run -t -i -p 5011:5011 \
	-e AWS_DEFAULT_REGION=us-west-2 \
	-e AWS_ACCESS_KEY_ID=$$(aws configure get aws_access_key_id) \
	-e AWS_SECRET_ACCESS_KEY=$$(aws configure get aws_secret_access_key) \
	-e AWS_SESSION_TOKEN=$$(aws configure get aws_session_token) \
	--rm stably/coreumservice:latest

####################
# Beta deployments #
####################
docker.buildanddeploy.coreumservice.beta: require-stage
docker.buildanddeploy.coreumservice.beta: platform=linux/amd64 # Assuming that our ECS instance is amd64
docker.buildanddeploy.coreumservice.beta: ecr=475910951137.dkr.ecr.us-west-2.amazonaws.com
docker.buildanddeploy.coreumservice.beta: aws.ecr.login
docker.buildanddeploy.coreumservice.beta: docker.build.service.coreumservice
	docker tag stably/coreumservice:latest $(ecr)/stably/coreumservice:latest
	docker push $(ecr)/stably/coreumservice:latest
	aws ecs update-service --cluster PrivateServices --service coreumservice --force-new-deployment

####################
# Prod deployments #
####################
docker.buildanddeploy.coreumservice.prod: require-stage
docker.buildanddeploy.coreumservice.prod: platform=linux/amd64 # Assuming that our ECS instance is amd64
docker.buildanddeploy.coreumservice.prod: ecr=808011740389.dkr.ecr.us-west-2.amazonaws.com
docker.buildanddeploy.coreumservice.prod: aws.ecr.login
docker.buildanddeploy.coreumservice.prod: docker.build.service.coreumservice
	docker tag stably/coreumservice:latest $(ecr)/stably/coreumservice:latest
	docker push $(ecr)/stably/coreumservice:latest
	aws ecs update-service --cluster PrivateServices --service coreumservice --force-new-deployment
