package conf

import "time"

type Conf struct {
	MaxX    float64
	MaxY    float64
	Address string
}

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const MaxBufferSize = 1024
const TimeOut = time.Second * 30
