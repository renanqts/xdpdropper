// It depends on bpf_link, available in Linux kernel version 5.7 or newer.
package xdp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/cilium/ebpf/link"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf xdp.c -- -I../../includes

type XDP interface {
	AddToDrop(string) error
	RemoveFromDrop(string) error
}

type xdp struct {
	objs bpfObjects
}

func New(ifaceName string) (XDP, error) {

	// Look up the network interface by name.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("Fail to lookup the iface %s: %s", ifaceName, err.Error())
	}

	// Load pre-compiled programs into the kernel.
	objs := bpfObjects{}
	err = loadBpfObjects(&objs, nil)
	if err != nil {
		return nil, fmt.Errorf("Fail loading eBPF program into the kernel: %s", err.Error())
	}
	defer objs.Close()

	// Attach the program.
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpProgFunc,
		Interface: iface.Index,
	})
	if err != nil {
		return nil, fmt.Errorf("Fail attaching XDP program: %s", err.Error())
	}
	defer l.Close()

	log.Printf("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)

	return xdp{
		objs: objs,
	}, nil
}

func (x xdp) AddToDrop(strIP string) error {
	ip, err := ip2long(strIP)
	if err != nil {
		return err
	}
	packetCount := int32(0) // As this is a new entry, set the counter to 0
	err = x.objs.XdpStatsMap.Put(ip, packetCount)
	return err
}

func (x xdp) RemoveFromDrop(strIP string) error {
	ip, err := ip2long(strIP)
	if err != nil {
		return err
	}
	err = x.objs.XdpStatsMap.Delete(ip)
	return err
}

func ip2long(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}
