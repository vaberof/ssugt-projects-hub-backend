PROJECT_NAME=ssugt-projects-hub
APP_NAME=web-backend

DOCKER_IMAGE_NAME=$(PROJECT_NAME)/$(APP_NAME)

WORK_DIR_LINUX=./cmd/ssugt-projects-hub
CONFIG_DIR_LINUX=./cmd/ssugt-projects-hub/config

WORK_DIR_WINDOWS=.\cmd\ssugt-projects-hub
CONFIG_DIR_WINDOWS=.\cmd\ssugt-projects-hub\config

docker.linux.local: build.linux
	docker image rm -f $(DOCKER_IMAGE_NAME) || (echo "Image $(DOCKER_IMAGE_NAME) didn't exist so not removed."; exit 0)
	docker build -t $(DOCKER_IMAGE_NAME) -f cmd/ssugt-projects-hub/Dockerfile .

build.linux: build.linux.clean
	mkdir -p cmd/ssugt-projects-hub/build
	go build -o cmd/ssugt-projects-hub/build/main cmd/ssugt-projects-hub/*.go
	cp -R cmd/ssugt-projects-hub/config/* cmd/ssugt-projects-hub/build

build.linux.local: build.linux.clean
	mkdir -p cmd/ssugt-projects-hub/build
	go build -o cmd/ssugt-projects-hub/build/main cmd/ssugt-projects-hub/*.go
	cp -R cmd/ssugt-projects-hub/config/* cmd/ssugt-projects-hub/build
	@echo "build.local: OK"

run.linux: build.linux
	go run $(WORK_DIR_LINUX)/*.go \
		-config.files $(CONFIG_DIR_LINUX)/application.yaml \
		-env.vars.file $(CONFIG_DIR_LINUX)/application.env \

build.linux.clean:
	rm -rf ./cmd/ssugt-projects-hub/build

run.windows:
	go run $(WORK_DIR_WINDOWS)\. \
		-config.files $(CONFIG_DIR_WINDOWS)\application.yaml \
		-env.vars.file $(CONFIG_DIR_WINDOWS)\application.env
