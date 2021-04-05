REPOSITORY := "vorteil"

# build a singular container using provided environment variable
.PHONY: build-singular
build-singular:
	echo "building ${CONTAINER}"	
	docker build ${CONTAINER} --tag ${REPOSITORY}/${CONTAINER}:latest
	docker build ${CONTAINER} --tag ${REPOSITORY}/${CONTAINER}:${VERSION}
	docker push ${REPOSITORY}/${CONTAINER}:latest
	docker push ${REPOSITORY}/${CONTAINER}:${VERSION}


