# XDPDropper

This is a small implementation of a dropper via eBPF XDP.   
It provides a simple JSON API where an IP can be `add/remove` to be dropped,   
meaning that all incoming packets in a certain interface with the source IP will be dropped.

## What and Why eBPF XDP
XDP is a technology that allows to attach eBPF programs to low-level hooks, implemented by network device drivers in the Linux kernel, as well as generic hooks that run after the device driver.

XDP can be used to achieve high-performance packet processing in an eBPF architecture, in this case here, used to drop packets.   
It can be very useful in case of DDoS attacks or others.   

## Configuration
This is a simple program, but some configuration is possible.   
This is done via environment variables:

| variable | description | type | default | required |
|---|---|---|---|---|
| `XDPDROPPER_ADDRESS` | `[address]:[port]` of the HTTP API | `string` | `0.0.0.0:8080` | `no` |
| `XDPDROPPER_IFACE` | Name of the interface in which the packets should be dropped. The eBPF XDP program will be attached to this interface |  `string` | | `yes` |
| `XDPDROPPER_LOGLEVEL` | Available log levels `INFO`, `DEBUG` | `string` | `INFO` | `yes` |

## API
The current APIs available are:

| context | method | content | return code |
|---|---|---|---|
| `/health` | `GET` | `NA` | 200 |
| `/drop` | `POST` | `{"ip":"1.1.1.1"}` | 201 |
| `/drop` | `DELETE` | `{"ip":"1.1.1.1"}` | 204 |

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

## Reference
- [eBPF XDP: The Basics and a Quick Tutorial](https://www.tigera.io/learn/guides/ebpf/ebpf-xdp/)
- [BPF and XDP Reference Guide](https://docs.cilium.io/en/stable/bpf/)
