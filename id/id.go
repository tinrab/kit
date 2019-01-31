package id

import (
	"bytes"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type ID uint64

type Generator struct {
	workerID  uint16
	sequence  uint16
	startTime int64
	lastTime  int64
	mutex     sync.Mutex
}

const (
	bitLengthTimestamp = 38
	bitLengthWorkerID  = 16
	bitLengthSequence  = 10
)

func NewGenerator(workerID uint16) *Generator {
	return NewGeneratorWithStartTime(workerID, time.Unix(0, 0))
}

func NewGeneratorWithStartTime(workerID uint16, startTime time.Time) *Generator {
	return &Generator{
		workerID:  workerID,
		startTime: startTime.Unix(),
		mutex:     sync.Mutex{},
	}
}

func (g *Generator) Generate() ID {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.lastTime = time.Now().Unix() - g.startTime

	value := Encode(g.lastTime, g.workerID, g.sequence)

	g.sequence++

	if g.sequence >= (1 << bitLengthSequence) {
		g.sequence = 0

		time.Sleep(time.Second)
	}

	return value
}

func (g *Generator) GenerateList(size int) []ID {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	list := make([]ID, size)
	g.lastTime = time.Now().Unix() - g.startTime

	for i := 0; i < size; i++ {
		list[i] = Encode(g.lastTime, g.workerID, g.sequence)

		g.sequence++

		if g.sequence >= (1 << bitLengthSequence) {
			g.sequence = 0

			time.Sleep(time.Second)
			g.lastTime = time.Now().Unix() - g.startTime
		}
	}

	return list
}

func (i ID) Timestamp() int64 {
	return int64(i >> (bitLengthWorkerID + bitLengthSequence))
}

func (i ID) WorkerID() uint16 {
	return uint16((i >> bitLengthSequence) & (1<<bitLengthWorkerID - 1))
}

func (i ID) Sequence() uint16 {
	return uint16(i & (1<<bitLengthSequence - 1))
}

func (i ID) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteRune('"')
	buf.WriteString(strconv.FormatUint(uint64(i), 10))
	buf.WriteRune('"')
	return buf.Bytes(), nil
}

func (i *ID) UnmarshalJSON(data []byte) error {
	if len(data) < 3 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid JSON ID value")
	}

	value, err := Parse(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}

	*i = value

	return nil
}

func Parse(s string) (ID, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return ID(i), nil
}

func Encode(timestamp int64, workerID uint16, sequence uint16) ID {
	if timestamp >= 1<<bitLengthTimestamp {
		panic("timestamp is too big")
	}
	if sequence >= 1<<bitLengthSequence {
		panic("sequence is too big")
	}

	t := uint64(timestamp)
	w := uint64(workerID)
	s := uint64(sequence)

	return ID((t&(1<<bitLengthTimestamp-1))<<(bitLengthWorkerID+bitLengthSequence) |
		((w & (1<<bitLengthWorkerID - 1)) << bitLengthSequence) |
		(s & (1<<bitLengthSequence - 1)))
}

func GetWorkerID() (uint16, error) {
	faces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}

	for _, face := range faces {
		if face.Flags&net.FlagUp == 0 || face.Flags&net.FlagLoopback != 0 {
			continue
		}

		addresses, err := face.Addrs()
		if err != nil {
			return 0, err
		}

		for _, address := range addresses {
			var ip net.IP

			switch v := address.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			id := uint16(ip[2])<<8 | uint16(ip[3])

			return id, nil
		}
	}

	return 0, errors.New("could not generate machine ID")
}
