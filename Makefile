REPOSITORY := "vorteil"

# build a singular container using provided environment variable
.PHONY: build-singular
build-singular:
	echo "building ${CONTAINER}"	
	docker build ${CONTAINER} --tag ${REPOSITORY}/${CONTAINER}
	docker push ${REPOSITORY}/${CONTAINER}


