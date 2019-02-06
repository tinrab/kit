package id

import (
	"bytes"
	"encoding/base64"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type ID int64

type Generator struct {
	workerID  uint16
	sequence  uint16
	startTime int64
	lastTime  int64
	mutex     sync.Mutex
}

const (
	bitLengthTimestamp = 37
	bitLengthWorkerID  = 16
	bitLengthSequence  = 10
)

var (
	ErrGenerateMachineID  = errors.New("could not generate machine ID")
	ErrInvalidBase62Value = errors.New("invalid base62 value")
	ErrInvalidValue       = errors.New("invalid value")

	base62Alphabet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
	buf.WriteString(strconv.FormatInt(int64(i), 10))
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

func (i ID) Base64() string {
	if i <= 0 {
		panic(ErrInvalidValue.Error())
	}

	var data []byte
	v := i
	for v > 0 {
		data = append(data, byte(v&255))
		v >>= 8
	}

	return base64.URLEncoding.EncodeToString(data)
}

func (i ID) Base62() string {
	if i <= 0 {
		panic(ErrInvalidValue.Error())
	}

	v := int64(i)
	e := ""

	for v > 0 {
		e = string(base62Alphabet[v%int64(len(base62Alphabet))]) + e
		v /= int64(len(base62Alphabet))
	}

	return string(e)
}

func Parse(s string) (ID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if i <= 0 {
		return 0, ErrInvalidValue
	}
	return ID(i), nil
}

func ParseBase64(s string) (ID, error) {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}

	id := int64(0)
	for i, b := range data {
		id += int64(b) << (8 * uint64(i))
	}

	if id <= 0 {
		return 0, ErrInvalidValue
	}

	return ID(id), nil
}

func ParseBase62(s string) (ID, error) {
	data := []rune(s)
	id := int64(0)

	for i := len(data) - 1; i >= 0; i-- {
		c := data[i]
		j := -1
		for k := 0; k < len(base62Alphabet); k++ {
			if base62Alphabet[k] == c {
				j = k
				break
			}
		}

		if j == -1 {
			return 0, ErrInvalidBase62Value
		}

		m := int64(len(data) - 1 - i)
		p := int64(1)
		for ; m > 0; m-- {
			p *= int64(len(base62Alphabet))
		}

		id += int64(j) * p
	}

	if id <= 0 {
		return 0, ErrInvalidValue
	}

	return ID(id), nil
}

func Encode(timestamp int64, workerID uint16, sequence uint16) ID {
	if timestamp >= 1<<bitLengthTimestamp {
		panic("timestamp is too big")
	}
	if sequence >= 1<<bitLengthSequence {
		panic("sequence is too big")
	}

	t := int64(timestamp)
	w := int64(workerID)
	s := int64(sequence)

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

	return 0, ErrGenerateMachineID
}
