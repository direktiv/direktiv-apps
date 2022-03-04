FROM golang:1.17-buster as builder
COPY cmd src/cmd/
COPY go.mod src/
COPY go.sum src/
RUN cd src/ && go mod download
RUN cd src && GO_ENABLED=0 go build -o /manager -ldflags=" -s -w" cmd/manager/*.go

# FROM ubuntu:21.10 as runner
# RUN apt-get update && apt-get install ca-certificates net-tools iproute2 buildah -y
# RUN echo "unqualified-search-registries = ['docker.io']" >> /etc/containers/registries.conf 
# COPY --from=builder /manager /runner

FROM ubuntu:21.10 as linux
ARG LINUX=5.16.11
RUN apt-get update && apt-get install wget curl build-essential make libelf-dev \
    flex bison bc libvdeplug-dev libssl-dev -y
COPY linux.config /
RUN wget https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-$LINUX.tar.xz
RUN tar -xvf linux-$LINUX.tar.xz
RUN cp /linux.config /linux-$LINUX/.config
RUN cd /linux-$LINUX && make ARCH=um -j `nproc`
RUN mkdir /out && mkdir /build && cp /linux-$LINUX/linux /out/linux
CMD ["cp", "/out/linux", "/build/linux"]
 
FROM ubuntu:21.10 as application
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -y install tzdata
RUN apt-get install tar gzip git net-tools iproute2 buildah  ca-certificates slirp -y
RUN echo "unqualified-search-registries = ['docker.io']" >> /etc/containers/registries.conf 

COPY --from=linux /out/linux /usr/local/bin/linux
COPY --from=builder /manager /usr/local/bin/manager
RUN mkdir /buildtmp
ENV TMPDIR /buildtmp
CMD ["/usr/local/bin/manager"]