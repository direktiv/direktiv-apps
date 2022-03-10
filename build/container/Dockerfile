# FROM golang:1.17-buster as builder
# COPY cmd src/cmd/
# COPY go.mod src/
# COPY go.sum src/
# RUN cd src/ && go mod download
# RUN cd src && GO_ENABLED=0 go build -o /manager -ldflags=" -s -w" cmd/manager/*.go

# FROM ubuntu:21.10 as runner
# RUN apt-get update && apt-get install ca-certificates net-tools iproute2 buildah -y
# RUN echo "unqualified-search-registries = ['docker.io']" >> /etc/containers/registries.conf 
# COPY --from=builder /manager /runner

# FROM ubuntu:21.10 as linux
# ARG LINUX=5.16.11
# RUN apt-get update && apt-get install wget curl build-essential make libelf-dev \
#     flex bison bc libvdeplug-dev libssl-dev -y
# RUN wget https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-$LINUX.tar.xz
# RUN tar -xvf linux-$LINUX.tar.xz
# COPY linux.config /
# RUN cp /linux.config /linux-$LINUX/.config
# RUN cd /linux-$LINUX && make ARCH=um -j `nproc`
# RUN mkdir /out && mkdir /build && cp /linux-$LINUX/linux /out/linux
# CMD ["cp", "/out/linux", "/build/linux"]

# FROM golang:1.17-buster as coredns
# RUN apt-get update 
# RUN apt-get install git make -y
# RUN git clone -b v1.9.0 --depth 1 https://github.com/coredns/coredns
# RUN cd coredns && echo fanout:github.com/networkservicemesh/fanout >> plugin.cfg && \ 
#   make

# FROM ubuntu:21.10 as application
# RUN apt-get update
# RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -y install tzdata
# RUN apt-get install tar gzip git net-tools iproute2 ca-certificates -y
# # install docker
# RUN apt-get install ca-certificates curl iptables gnupg lsb-release -y
# RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
# RUN echo \
#   "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
#   $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
# RUN apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io
# RUN mkdir /disk

# RUN update-alternatives --set iptables /usr/sbin/iptables-legacy
# RUN update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy

# RUN curl -L -o /slirp4netns https://github.com/rootless-containers/slirp4netns/releases/download/v1.2.0-beta.0/slirp4netns-x86_64
# RUN chmod 755 /slirp4netns

# COPY --from=linux /out/linux /usr/local/bin/linux
# COPY --from=coredns /go/bin/coredns /coredns
# COPY --from=builder /manager /usr/local/bin/manager

# ENV DOCKER_HOST=tcp://127.0.0.1:2375
# COPY daemon.json /etc/docker/daemon.json

# RUN mkdir /buildtmp
# ENV TMPDIR /buildtmp

# CMD ["/usr/local/bin/manager"]

FROM golang:1.17-buster as builder
COPY cmd src/cmd/
COPY go.mod src/
COPY go.sum src/
RUN cd src/ && go mod download
RUN cd src && GO_ENABLED=0 go build -o /runner -ldflags=" -s -w" cmd/runner/*.go
RUN cd src && GO_ENABLED=0 go build -o /builder -ldflags=" -s -w" cmd/builder/*.go

FROM ubuntu:21.10 as base
RUN apt-get update && apt-get install net-tools qemu slirp
RUN apt-get install ca-certificates curl iptables gnupg lsb-release -y
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
RUN echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io
COPY --from=builder /runner /runner
CMD ["/runner"]

FROM base
ARG BASE_IMAGE=gerke74/base
RUN apt-get update && apt-get install wget -y
RUN wget https://github.com/vorteil/vorteil/releases/latest/download/vorteil_linux-x86.tar.gz && \
      tar -xzf vorteil_linux-x86.tar.gz && \
        ./vorteil version

RUN ./vorteil projects convert-container $BASE_IMAGE /tmp/base
RUN  ./vorteil images build --vm.disk-size "+30GiB" \
           --vm.kernel=22.3.1 --format=vmdk /tmp/base/ -f

ENV DOCKER_HOST=tcp://127.0.0.1:2375
RUN mkdir -p /tmp/tars

RUN apt-get install qemu-system-x86 -y
COPY --from=builder /builder /builder
CMD ["/builder"]
