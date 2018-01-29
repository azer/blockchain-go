# blockchain-go

Simple and experimental Blockchain implementation with PoW. I'm working on it for learning the basic concepts, it's still under development. 

## Usage

Grab blockchain binary;

```
$ go get github.com/azer/blockchain-go/blockchain
```

Create an `.env` file in the working directory;

```bash
BC_TARGET_BITS = 24 # higher makes mining more difficult.
BC_DB_FILE = /tmp/bc-db
BC_BLOCKS_BUCKET = blocks
```

Run the app;

```
$ blockchain print-chain | less
```
