PROJECT_NAME=ssugt-project
APP_NAME=web-backend

DOCKER_IMAGE_NAME=$(PROJECT_NAME)/$(APP_NAME)

WORK_DIR_LINUX=./cmd/ssugt-projects
CONFIG_DIR_LINUX=./cmd/ssugt-projects/config

WORK_DIR_WINDOWS=.\cmd\ssugt-projects
CONFIG_DIR_WINDOWS=.\cmd\ssugt-projects\config

docker.linux.local: build.linux
	docker image rm -f $(DOCKER_IMAGE_NAME) || (echo "Image $(DOCKER_IMAGE_NAME) didn't exist so not removed."; exit 0)
	docker build -t $(DOCKER_IMAGE_NAME) -f cmd/ssugt-projects/Dockerfile .

build.linux: build.clean
	mkdir -p cmd/ssugt-projects/build
	go build -o cmd/ssugt-projects/build/main cmd/ssugt-projects/*.go
	cp -R cmd/ssugt-projects/config/* cmd/ssugt-projects/build

build.linux.local: build.clean
	mkdir -p cmd/ssugt-projects/build
	go build -o cmd/ssugt-projects/build/main cmd/ssugt-projects/*.go
	cp -R cmd/ssugt-projects/config/* cmd/ssugt-projects/build
	@echo "build.local: OK"

run.linux: build.linux
	go run $(WORK_DIR_LINUX)/*.go \
		-config.files $(CONFIG_DIR_LINUX)/application.yaml \
		-env.vars.file $(CONFIG_DIR_LINUX)/application.env \

run.windows:
	go run $(WORK_DIR_WINDOWS)\. \
		-config.files $(CONFIG_DIR_WINDOWS)\application.yaml \
		-env.vars.file $(CONFIG_DIR_WINDOWS)\application.env

build.clean:
	rm -rf ./cmd/ssugt-projects/build