# blockchain-go

Simple and experimental Blockchain implementation with PoW.

## Usage

Grab blockchain binary;

```
$ go get github.com/azer/blockchain-go/blockchain
```

Create an `.env` file in the working directory;

```
BC_DB_FILE = /tmp/bc-db
BC_TARGET_BITS = 24
BC_BLOCKS_BUCKET = blocks
```

Run the app;

```
$ blockchain print-chain | less
```
