# FLY

A complete open source e-commerce solution by rust language.

## Install rust

```
curl https://sh.rustup.rs -sSf | sh
rustup default nightly
rustup update
cargo install rustfmt-nightly --force
cargo install racer --force
```

add to ~/.zshrd

```
export RUST_SRC_PATH=$HOME/src/rust/src
export CARGO_PATH=$HOME/.cargo
export PATH=$CARGO_PATH/bin:$PATH
export LD_LIBRARY_PATH=$(rustc --print sysroot)/lib
```

## Build

```
git clone https://github.com/kapmahc/fly.git
cd fly
cargo update
cargo build --release
```

## Atom plugins

- language-rust

## Documents

- [rust](https://doc.rust-lang.org/book/second-edition/)
- [cargo](https://crates.io/)
- [rocket](https://rocket.rs/guide/)
- [bootstrap](https://getbootstrap.com/docs/4.0/getting-started/introduction/)
- [Material Design](http://materializecss.com/getting-started.html)
