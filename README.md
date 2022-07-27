# XDP network packet dropper via API

This is a small implementation of a dropper via eBPF XDP.   
It provides a simple JSON API where an IP can be `add/remove` to be dropped,   
meaning that all incoming packets in a certain interface with the source IP will be dropped.

## API
The current APIs available are:

| context | method | content | return code |
|---|---|---|---|---|
| `/health` | `GET` | | 200 |
| `/` | `POST` | `{"ip":"1.1.1.1"}` | 201 |
| `/` | `DELETE` | `{"ip":"1.1.1.1"}` | 204 |

**Examples**:

Dropping incoming packets with the source IP `1.1.1.1`
```shell
curl -d '{"ip":"1.1.1.1"}' -H "Content-Type: application/json" http://localhost:8080/ -v
```

To remove the IP from the list:
```shell
curl -d '{"ip":"1.1.1.1"}' -X DELETE -H "Content-Type: application/json" http://localhost:8080/ -v
```

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