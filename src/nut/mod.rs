pub mod controllers;
pub mod console;


use std::io;
use std::fmt;
use std::error;
use docopt;

pub const VERSION: &'static str = "2017.09.30";

#[derive(Debug)]
pub enum Error {
    Io(io::Error),
    Docopt(docopt::Error),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            Error::Io(ref err) => err.fmt(f),
            Error::Docopt(ref err) => err.fmt(f),
        }
    }
}

impl error::Error for Error {
    fn description(&self) -> &str {
        match *self {
            Error::Io(ref err) => err.description(),
            Error::Docopt(ref err) => err.description(),
        }
    }

    fn cause(&self) -> Option<&error::Error> {
        match *self {
            Error::Io(ref err) => Some(err),
            Error::Docopt(ref err) => Some(err),
        }
    }
}

impl From<io::Error> for Error {
    fn from(err: io::Error) -> Error {
        Error::Io(err)
    }
}

impl From<docopt::Error> for Error {
    fn from(err: docopt::Error) -> Error {
        Error::Docopt(err)
    }
}
