REPOSITORY := vorteil
# dependencies
.PHONY: dependencies
dependencies:
	echo "fetching dependencies"
	mkdir -p /home/runner/.docker/cli-plugins
	wget https://github.com/christian-korneck/docker-pushrm/releases/download/v1.7.0/docker-pushrm_linux_amd64
	mv ./docker-pushrm_linux_amd64 /home/runner/.docker/cli-plugins/docker-pushrm
	chmod +x ~/.docker/cli-plugins/docker-pushrm
	chmod +x ./getversion.sh
# build a singular container using provided environment variable
.PHONY: build-singular
build-singular:
	echo "building ${CONTAINER}";
	echo v$(shell REPOSITORY=${REPOSITORY} CONTAINER=${CONTAINER} ./getversion.sh) > ${CONTAINER}/VERSION
	docker build ${CONTAINER} -t ${REPOSITORY}/${CONTAINER}:latest -t ${REPOSITORY}/${CONTAINER}:v$(shell REPOSITORY=${REPOSITORY} CONTAINER=${CONTAINER} ./getversion.sh)
	docker push ${REPOSITORY}/${CONTAINER}:v$(shell REPOSITORY=${REPOSITORY} CONTAINER=${CONTAINER} ./getversion.sh)
	docker push ${REPOSITORY}/${CONTAINER}:latest
	cd ${CONTAINER} && docker pushrm docker.io/${REPOSITORY}/${CONTAINER}
# build all containers using provided version variable
.PHONY: all
all:
	echo "building all containers with version ${VERSION}"
	@for f in $(shell ls ${MYDIR} -I pkg -I Makefile -I readme.md -I docker-pushrm_linux_amd64 -I cli); do \
		docker build $${f} --tag ${REPOSITORY}/$${f}:${VERSION}; \
		docker push ${REPOSITORY}/$${f}:${VERSION}; \
		cd $${f}; \
		docker pushrm docker.io/${REPOSITORY}/$${f}; \
		cd ../; \
	done
