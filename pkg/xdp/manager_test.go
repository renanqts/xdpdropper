package xdp

import (
	"net"
	"testing"

	"github.com/renanqts/xdpdropper/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	logConfig := logger.NewDefaultConfig()
	logConfig.Level = "debug"
	_ = logger.Init(logConfig)
}

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
		key   []byte
		value uint32
	)
	iter := xdp.objs.DropMap.Iterate()
	for iter.Next(&key, &value) {
		actualIP := net.IP(key) // IPv4 source address in network byte order.
		actualCounter := value
		assert.Equal(t, expectedIP, actualIP.String())
		assert.Equal(t, uint32(0), actualCounter)
	}
	assert.Nil(t, iter.Err())

	err = xdp.RemoveFromDrop(expectedIP)
	assert.Nil(t, err)
	assert.Equal(t, false, iter.Next(&key, &value))
}
