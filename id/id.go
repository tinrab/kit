package id

import (
	"bytes"
	"github.com/tinrab/kit"
	"strconv"
	"sync"
	"time"
)

var (
	ErrInvalidJSONValue = kit.NewErrorWithMessage("invalid JSON ID value")
)

type ID uint64

type Generator struct {
	workerID    uint16
	counter     uint16
	startTime   int64
	lastTime    int64
	advanceTime int64
	mutex       sync.Mutex
}

func NewGenerator(workerID uint16) *Generator {
	return NewGeneratorWithStartTime(workerID, time.Unix(0, 0))
}

func NewGeneratorWithStartTime(workerID uint16, startTime time.Time) *Generator {
	return &Generator{
		workerID:    workerID,
		startTime:   startTime.UTC().UnixNano(),
		lastTime:    startTime.UTC().UnixNano(),
		advanceTime: 0,
		mutex:       sync.Mutex{},
	}
}

func (g *Generator) Generate() ID {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	t := g.timeOffset()

	value := ID((t/int64(time.Second))<<32) |
		ID(g.workerID)<<16 |
		ID(g.counter)

	g.counter++
	if t-g.lastTime > int64(time.Second) {
		g.counter = 0
	}
	g.lastTime = t

	return value
}

func (g *Generator) GenerateList(size int) []ID {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	t := g.timeOffset()

	list := make([]ID, size)
	for i := 0; i < size; i++ {
		t = g.timeOffset()
		list[i] = ID((t/int64(time.Second))<<32) |
			ID(g.workerID)<<16 |
			ID(g.counter)

		g.counter++
		if g.counter == 0 {
			g.advanceTime += int64(time.Second)
		}
	}

	g.lastTime = t

	return list
}

func (g Generator) timeOffset() int64 {
	return time.Now().UTC().UnixNano() - g.startTime + g.advanceTime
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
		return ErrInvalidJSONValue
	}

	value, err := ParseID(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}

	*i = value

	return nil
}

func (i ID) Decode() (int64, uint16, uint16) {
	t := int64(i >> 32)
	w := uint16((i >> 16) & (1<<16 - 1))
	c := uint16(i & (1<<16 - 1))

	return t, w, c
}

func ParseID(s string) (ID, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return ID(i), nil
}
