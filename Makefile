REPOSITORY := direktiv
# dependencies
.PHONY: dependencies
dependencies:
	echo "fetching dependencies"
	mkdir -p /home/runner/.docker/cli-plugins
	wget https://github.com/christian-korneck/docker-pushrm/releases/download/v1.7.0/docker-pushrm_linux_amd64
	mv ./docker-pushrm_linux_amd64 /home/runner/.docker/cli-plugins/docker-pushrm
	chmod +x ~/.docker/cli-plugins/docker-pushrm
	chmod +x ./getversion.sh
	chmod +x ./pushall.sh
# build a singular container using provided environment variable
.PHONY: build-singular
build-singular:
	echo "building ${CONTAINER}";
	echo $(shell REPOSITORY=${REPOSITORY} CONTAINER=${CONTAINER} ./getversion.sh) > ${CONTAINER}/VERSION
	DOCKER_BUILDKIT=1 docker build ${CONTAINER} -t ${REPOSITORY}/${CONTAINER}:latest -t ${REPOSITORY}/${CONTAINER}:$(shell REPOSITORY=${REPOSITORY} CONTAINER=${CONTAINER} ./getversion.sh) --no-cache
	docker push ${REPOSITORY}/${CONTAINER}:$(shell REPOSITORY=${REPOSITORY} CONTAINER=${CONTAINER} ./getversion.sh)
	docker push ${REPOSITORY}/${CONTAINER}:latest
	cd ${CONTAINER} && docker pushrm docker.io/${REPOSITORY}/${CONTAINER}
# build all containers using provided version variable
.PHONY: all
all:
	@for f in $(shell ls ${MYDIR} -I LICENSE -I examples -I build-local.sh -I pushall.sh -I getversion.sh -I pkg -I Makefile -I readme.md -I docker-pushrm_linux_amd64 -I cli); do \
		REPOSITORY=${REPOSITORY} CONTAINER=$${f} ./pushall.sh; \
		cd $${f}; \
		docker pushrm docker.io/${REPOSITORY}/$${f}; \
		cd ../; \
	done

.PHONY: reinit
reinit:
	@for f in $(shell ls ${MYDIR} -I LICENSE -I examples -I build-local.sh -I pushall.sh -I getversion.sh -I pkg -I Makefile -I readme.md -I docker-pushrm_linux_amd64 -I cli); do \
		cd $${f}; \
		rm go.mod; \
		rm go.sum; \
		go mod init github.com/direktiv/direktiv-apps/$${f}; \
		go mod tidy; \
		cd ../; \
	done