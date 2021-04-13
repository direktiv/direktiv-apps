## Rust

## Code

```rust
// src/main.rs
use actix_web::{App, HttpResponse, HttpServer, Responder, post};
use actix_web::web::Json;

use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
struct Input {
    name: String,
}

#[derive(Serialize)]
struct Output {
    greeting: String,
}

#[post("/")]
async fn index(info: Json<Input>) -> impl Responder {
    HttpResponse::Ok().json(Output { greeting: format!("Welcome to Direktiv, {}!", info.name) })
}

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(index)
    })
        .bind("0.0.0.0:8080")?
        .run()
        .await
}
```

## Dockerfile

```dockerfile
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
```
