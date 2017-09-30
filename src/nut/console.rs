
use rocket;

pub fn server() {
    rocket::ignite()
        .mount("/", routes![super::controllers::index])
        .launch();
}

pub fn version() {
    println!("{}", super::VERSION);
}

pub fn generate_locale(name: String) {
    info!("generate file {}", name);
}

pub fn generate_migration(name: String) {
    info!("generate file {}", name);
}
