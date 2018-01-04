ORGANIZATION=hyperpilot
IMAGE=snap-ddagent
TAG=test
GLIDE=$(which glide)
GO_EXECUTABLE ?= go
# For windows developer, use $(go list ./... | grep -v /vendor/)
PACKAGES=$(glide novendor)

init:
	glide install

test:
	${GO_EXECUTABLE} test ${PACKAGES} -v

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o ${IMAGE}

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o ${IMAGE}

complete-build-linux: init build-linux

docker-build:
	docker build --no-cache . -t ${ORGANIZATION}/${IMAGE}:${TAG}

docker-push:
	docker push ${ORGANIZATION}/${IMAGE}:${TAG}
