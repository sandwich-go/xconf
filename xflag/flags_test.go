package xflag

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type logging struct {
	Interval int64
	Path     string
}

type socket struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type tcp struct {
	ReadTimeout time.Duration
	socket
}

type network struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	TestMap         map[string]string
	TestStringSlice []string

	TestIntSlice   []int
	TestInt8Slice  []int8
	TestInt16Slice []int16
	TestInt32Slice []int32
	TestInt64Slice []int64

	TestUIntSlice   []uint
	TestUInt8Slice  []uint8
	TestUInt16Slice []uint16
	TestUInt32Slice []uint32
	TestUInt64Slice []uint64
	tcp
}

type Cfg1 struct {
	logging
	network
}

func TestFlagMakerExample(t *testing.T) {
	cfg := Cfg1{}

	{
		f := NewMaker()
		err := f.Set(&cfg)
		assert.True(t, err == nil)
		f.PrintDefaults()
	}

	args := []string{
		"--network.tcp.socket.read_timeout", "5ms",
		"--network.tcp.read_timeout", "3ms",
		"-logging.path", "/var/log",
	}
	args, err := ParseArgs(cfg, args)
	assert.False(t, err == nil)
	args, err = ParseArgs(&cfg, args)
	assert.True(t, err == nil)
	assert.Equal(t, 0, len(args))

	expected := Cfg1{
		network: network{
			tcp: tcp{
				ReadTimeout: time.Duration(3) * time.Millisecond,
				socket: socket{
					ReadTimeout: time.Duration(5) * time.Millisecond,
				},
			},
		},
		logging: logging{
			Path: "/var/log",
		},
	}
	assert.Equal(t, expected, cfg)
}
