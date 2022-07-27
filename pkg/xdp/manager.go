// It depends on bpf_link, available in Linux kernel version 5.7 or newer.
package xdp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/renanqts/xdpdropper/pkg/logger"
	"go.uber.org/zap"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf xdp.c -- -I../../includes

type XDP interface {
	AddToDrop(string) error
	RemoveFromDrop(string) error
	Close()
}

type xdp struct {
	objs bpfObjects
	link link.Link
}

func New(ifaceName string) (XDP, error) {
	// Look up the network interface by name.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("Fail to lookup the iface %s: %s", ifaceName, err.Error())
	}

	// Load pre-compiled programs into the kernel.
	logger.Log.Debug("xdp load program into the kernel")
	objs := bpfObjects{}
	err = loadBpfObjects(&objs, nil)
	if err != nil {
		return nil, fmt.Errorf("Fail loading eBPF program into the kernel: %s", err.Error())
	}

	// Attach the program.
	logger.Log.Debug("xdp attach program", zap.String("iface", ifaceName))
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpDropFunc,
		Interface: iface.Index,
	})
	if err != nil {
		return nil, fmt.Errorf("Fail attaching XDP program: %s", err.Error())
	}

	logger.Log.Info("XDP program attached", zap.String("iface", iface.Name), zap.Int("index", iface.Index))

	return xdp{
		objs: objs,
		link: l,
	}, nil
}

func (x xdp) Close() {
	x.link.Close()
	x.objs.Close()
}

func (x xdp) AddToDrop(strIP string) error {
	logger.Log.Debug("xdp dropper add", zap.String("ip", strIP))
	ip, err := ip2long(strIP)
	if err != nil {
		return err
	}
	packetCount := int32(0) // As this is a new entry, set the counter to 0
	err = x.objs.DropMap.Put(ip, packetCount)
	return err
}

func (x xdp) RemoveFromDrop(strIP string) error {
	logger.Log.Debug("xdp dropper remove", zap.String("ip", strIP))
	ip, err := ip2long(strIP)
	if err != nil {
		return err
	}
	err = x.objs.DropMap.Delete(ip)
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
