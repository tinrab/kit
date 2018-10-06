package expr

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
)

func TestBinary(t *testing.T) {
	assert.Equal(t, 2.0, New("1 + 1").MustEvaluate())
	assert.Equal(t, 2.0, New("3 - 1").MustEvaluate())
	assert.Equal(t, 4.0, New("2 * 2").MustEvaluate())
	assert.Equal(t, 2.0, New("4 / 2").MustEvaluate())
}

func TestCompound(t *testing.T) {
	assert.Equal(t, 2.0, New("4 * 2 / 4").MustEvaluate())
	assert.Equal(t, 4.0, New("(1+1)*2/(2-1)").MustEvaluate())
}

func TestFunc(t *testing.T) {
	assert.Equal(t, math.Pi, New("Pi").MustEvaluate())
	assert.Equal(t, math.E, New("E").MustEvaluate())
	assert.Equal(t, 3.0, New("floor(3.5)").MustEvaluate())
	assert.Equal(t, 3.0, New("round(3.49)").MustEvaluate())
	assert.Equal(t, 4.0, New("ceil(3.5)").MustEvaluate())
	assert.Equal(t, 8.0, New("pow(2, 3)").MustEvaluate())
}

func TestBool(t *testing.T) {
	assert.True(t, New("true").MustEvaluate().(bool))
	assert.False(t, New("false").MustEvaluate().(bool))
	assert.True(t, New("2 == 2").MustEvaluate().(bool))
	assert.False(t, New("2 == 3").MustEvaluate().(bool))
}

func TestUnary(t *testing.T) {
	assert.Equal(t, -3.0, New("-3").MustEvaluate().(float64))
	assert.True(t, New("!false").MustEvaluate().(bool))
	assert.True(t, New("3 == (2 + 1)").MustEvaluate().(bool))
}
