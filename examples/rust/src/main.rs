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