package xdp

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperations(t *testing.T) {
	objs := bpfObjects{}
	err := loadBpfObjects(&objs, nil)
	assert.Nil(t, err)
	defer objs.Close()

	xdp := xdp{
		objs: objs,
	}

	expectedIP := "1.2.3.5"
	err = xdp.AddToDrop(expectedIP)
	assert.Nil(t, err)

	var (
		key   uint32
		value uint32
	)
	iter := objs.XdpStatsMap.Iterate()
	for iter.Next(&key, &value) {
		actualIP := int2ip(key) // IPv4 source address in network byte order.
		actualCounter := value
		assert.Equal(t, expectedIP, actualIP)
		assert.Equal(t, uint32(0), actualCounter)
	}
	assert.Nil(t, iter.Err())

	err = xdp.RemoveFromDrop(expectedIP)
	assert.Nil(t, err)
	assert.Equal(t, false, iter.Next(&key, &value))

}

func int2ip(nn uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip.String()
}
