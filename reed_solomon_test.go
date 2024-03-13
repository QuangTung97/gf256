package gf256

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReedSolomon(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		data := reedSolomonEncode([2]byte{10, 20})
		assert.Equal(t, [6]byte{
			10, 20,
			0x36, 0x28, 0x72, 0x6c,
		}, data)

		input := reedSolomonDecode(data, [6]bool{true, true})
		assert.Equal(t, [2]byte{10, 20}, input)

		input = reedSolomonDecode(data, [6]bool{false, true, true})
		assert.Equal(t, [2]byte{10, 20}, input)

		input = reedSolomonDecode(data, [6]bool{false, false, true, true})
		assert.Equal(t, [2]byte{10, 20}, input)

		input = reedSolomonDecode(data, [6]bool{true, false, false, true})
		assert.Equal(t, [2]byte{10, 20}, input)
	})
}

func TestReedSolomon_Random(t *testing.T) {
	seed := time.Now().UnixNano()
	fmt.Println("SEED =", seed)
	r := rand.New(rand.NewSource(seed))

	for i := 0; i < 100000; i++ {
		a := uint8(r.Intn(256))
		b := uint8(r.Intn(256))

		input := [2]byte{a, b}
		data := reedSolomonEncode(input)

		true1 := r.Intn(5)
		true2 := r.Intn(5-true1) + true1 + 1

		var trueLabel [6]bool
		trueLabel[true1] = true
		trueLabel[true2] = true

		newInput := reedSolomonDecode(data, trueLabel)
		assert.Equal(t, input, newInput)
	}
}
