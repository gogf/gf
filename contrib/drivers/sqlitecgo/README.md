## sqlitecgo

- Using GCC to connect to sqlite.
- go-sqlite does not support compiling for win32 bit. This one supports it..

You need to set the environment variable CGO_ENABLED=1 and make sure that GCC is installed on your path.

windows gcc: https://jmeubank.github.io/tdm-gcc/
