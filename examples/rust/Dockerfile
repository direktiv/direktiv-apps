# Dockerfile for creating a statically-linked Rust application using docker's
# multi-stage build feature. This also leverages the docker build cache to avoid
# re-downloading dependencies if they have not changed.
FROM rust:latest AS build
WORKDIR /usr/src/myapp

COPY . .

RUN cargo install --path .
# Copy the statically-linked binary into a scratch container.
FROM debian:buster-slim
COPY --from=build /usr/local/cargo/bin/greeting /greeting

USER 1000
CMD ["./greeting"]