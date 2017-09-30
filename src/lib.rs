#![feature(plugin)]
#![plugin(rocket_codegen)]
#![plugin(docopt_macros)]

extern crate docopt;
extern crate env_logger;
#[macro_use]
extern crate log;
extern crate rocket;
#[macro_use]
extern crate serde_derive;


use std::env;


pub mod nut;
pub mod forum;
pub mod survey;
pub mod reading;
pub mod ops;
pub mod mall;
pub mod erp;
pub mod pos;




docopt!(Args derive Debug, "
A complete open source e-commerce solution(by rust language).

Usage:
  fly generate (migration | locale) [--name=<fn>]
  fly db (create | migrate | rollback | reset | status | drop)
  fly cache (list | clear)
  fly (-h | --help)
  fly (-v | --version)

Options:
  --name=<fn>   File name.
  -h --help     Show this screen.
  -v --version  Show version.
");

use std::borrow::Borrow;
pub fn run() {
    if env::args().len() == 1 {
        // TODO
        nut::console::server();
        return;
    }


    // let cfg = rocket::config::Config::build(
    //     env::var("ROCKET_ENV")
    //         .unwrap()
    //         .parse::<Environment>()
    //         .unwrap(),
    // ).finalize()
    //     .unwrap();
    // println!("{:?}", cfg.root());
    let ins = rocket::ignite();
    let cfg = ins.config().borrow();
    println!("{:?}", cfg);

    let args: Args = Args::docopt().deserialize().unwrap_or_else(|e| e.exit());
    if args.cmd_generate {
        let name = args.flag_name.to_string();
        if name == "" {
            panic!("name must empty")
        }
        if args.cmd_migration {
            nut::console::generate_migration(name);
            return;
        }
        if args.cmd_locale {
            nut::console::generate_locale(name);
            return;
        }
    }
    if args.flag_version {
        nut::console::version();
        return;
    }
    println!("{:?}", args);
}
