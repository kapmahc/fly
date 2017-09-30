#![feature(plugin)]
#![plugin(rocket_codegen)]

extern crate rocket;
extern crate fly;


fn main() {
    rocket::ignite().mount("/", routes![fly::nut::controllers::index]).launch();
}
