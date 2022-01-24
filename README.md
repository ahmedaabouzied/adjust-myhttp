# MyHTTP

Parallel HTTP request initiator

## Building

### Requirements

1. Go. Preferably go 1.7. This was the version used to write this code.
2. Make. You can install make by installing `build-essential`
   On Ubuntu, you can use apt

```bash
sudo apt install build-essential
```

### Building stpes

Ensuring `make` is installed just run `make` in the root
directory of the cloned repo.

## Running

```bash
./myhttp -parallel <concurrency limit> [URLs]
```

### Arguments

`parallel` (int): The limit of prallel requests to make

### Example

```bash
./myhttp -parallel 2 adjust.com google.com facebook.com
```

### Example output

```bash
./myhttp -parallel 2 adjust.com google.com facebook.com
Limit = 2
URL count = 3
Workers 2
https://google.com 427bf71d4420c73b169d39933bf6652e
https://adjust.com 16f21af9ecae9734c575a2c772715733
https://facebook.com df490e48106bbe2e32cb1756bb45b7c8
```

## Testing

### Unit tests

```bash
go test
```

### Benchmarks

```bash
go test -bench=.
```
