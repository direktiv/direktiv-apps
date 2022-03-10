BASE_IMAGE := "gerke74/base"

clean:
	docker rmi `docker images --filter label=direktiv-base -q` -f

.PHONY: builder
builder: 
	@if [ ! `docker images --filter label=direktiv-base -q` ]; then \
		docker build --label direktiv-base -t ${BASE_IMAGE} --target base . ; \
		docker push ${BASE_IMAGE}; \
	fi
	docker build --build-arg BASE_IMAGE=${BASE_IMAGE} -t builder .


.PHONY: build-local
build-local: 
	docker build --target base -t base .
	rm -Rf /tmp/base && vorteil projects convert-container local.docker/base /tmp/base
	vorteil images build --vm.disk-size "+30GiB" --vm.kernel=22.3.1 --format=vmdk /tmp/base/ -f