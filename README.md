## Haskell

```sh
$ cabal init

$ cabal install --only-dependencies
...
setup: Missing dependencies on foreign libraries:
* Missing (or bad) C libraries: ssl, crypto
This problem can usually be solved by installing the system packages that
provide these libraries (you may need the "-dev" versions). If the libraries
are already installed but in a non-standard location then you can use the
flags --extra-include-dirs= and --extra-lib-dirs= to specify where they are.If
the library files do exist, it may contain errors that are caught by the C
compiler at the preprocessing stage. In this case you can re-run configure
with the verbosity flag -v3 to see the error messages.

cabal: Failed to build mysql-0.2.1 (which is required by
mysql-simple-0.4.8.1). See the build log above for details

$ brew install icu4c

$ brew install openssl

$ cabal install --only-dependencies
```
