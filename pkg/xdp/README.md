# eBPF XDP

## How to generate
`bpf_`*` files are generated by bpf2go, please, don't touch them.
To generate those, from the repository root, use:
```shell
make
```

## How to run
Inside the vagrant machine, use:
```
sudo go test -v
```