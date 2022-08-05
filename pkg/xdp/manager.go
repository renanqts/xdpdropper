// It depends on bpf_link, available in Linux kernel version 5.7 or newer.
package xdp

import (
	"fmt"
	"net"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/renanqts/xdpdropper/pkg/logger"
	"github.com/renanqts/xdpdropper/pkg/metrics"
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

	xdp := xdp{
		objs: objs,
		link: l,
	}

	metric := metrics.NewGaugeVec("dropped_packets", "xdp", "number of droppped packets", []string{"ip"})
	go xdp.measures(metric)

	return &xdp, nil
}

func (x xdp) Close() {
	x.link.Close()
	x.objs.Close()
}

func (x xdp) AddToDrop(strIP string) error {
	logger.Log.Debug("xdp dropper add", zap.String("ip", strIP))
	ip := ip2Byte(strIP)
	packetCount := int32(0) // As this is a new entry, set the counter to 0
	err := x.objs.DropMap.Put(ip, packetCount)
	return err
}

func (x xdp) RemoveFromDrop(strIP string) error {
	logger.Log.Debug("xdp dropper remove", zap.String("ip", strIP))
	ip := ip2Byte(strIP)
	err := x.objs.DropMap.Delete(ip)
	return err
}

func (x xdp) measures(g *prometheus.GaugeVec) {
	var (
		key         []byte
		packetCount uint32
	)
	// get counters each 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		iter := x.objs.DropMap.Iterate()
		for iter.Next(&key, &packetCount) {
			src_ip := net.IP(key)
			g.WithLabelValues(src_ip.String()).Add(float64(packetCount))
		}
	}
}

func ip2Byte(ip string) []byte {
	return net.ParseIP(ip).To4()
}
