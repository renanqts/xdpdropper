# XDP network packet dropper via API

This is a small implementation of a dropper via eBPF XDP.   
It provides an API where an IP can be add/remove to be dropped.

## Local tests
For the time being, tests can be performed using vagrant VM.   
For that, use:
```
vagrant up  # wait until is up and running
vagrant ssh # access the VM
sudo make test # to execute the tests
```

It also counts with a test via `docker-compose.yml`, it can be done:
```
docker-compose up
```
*it must be performed inside the VM.*