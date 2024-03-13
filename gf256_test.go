package gf256

import (
	"fmt"
	"math/bits"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert.Equal(t, uint8(0b110), add(0b101, 0b11))
}

func TestSimpleMul(t *testing.T) {
	assert.Equal(t, 2, bits.Len8(3))
	assert.Equal(t, 3, bits.Len8(4))

	t.Run("normal", func(t *testing.T) {
		assert.Equal(t, uint8(0b11), simpleMul(0b11, 0b1))
		assert.Equal(t, uint8(0b101), simpleMul(0b101, 0b1))
		assert.Equal(t, uint8(0b101|0b1010), simpleMul(0b101, 0b11))
	})

	t.Run("special", func(t *testing.T) {
		assert.Equal(t, uint8(0b101^0b10100), simpleMul(0b101, 0b101))
		assert.Equal(t, uint8(0b1101^0b11010), simpleMul(0b1101, 0b11))
	})

	t.Run("mod", func(t *testing.T) {
		a := uint16(0b1110101)
		b := a << 2
		c := a << 3
		s := a ^ b ^ c

		assert.Equal(t,
			strings.ReplaceAll("10_0000_1001", "_", ""),
			fmt.Sprintf("%b", s),
		)

		poly := uint16(0x11b)
		assert.Equal(t,
			strings.ReplaceAll("1_0001_1011", "_", ""),
			fmt.Sprintf("%b", poly),
		)

		assert.Equal(t, 10, bits.Len16(s))
		assert.Equal(t, 9, bits.Len16(poly))

		s ^= poly << 1
		assert.Equal(t,
			strings.ReplaceAll("11_1111", "_", ""),
			fmt.Sprintf("%b", s),
		)

		assert.Equal(t, uint8(0b111111), simpleMul(0b1110101, 0b1101))
	})

	t.Run("exp", func(t *testing.T) {
		x := uint8(3)
		assert.Equal(t, uint8(5), simpleMul(x, 3))
		assert.Equal(t, uint8(5), simpleExp(x, 2))
		assert.Equal(t, uint8(1), simpleExp(x, 0))
		assert.Equal(t, uint8(15), simpleExp(x, 3))
		assert.Equal(t, uint8(129), simpleExp(x, 88))
		assert.Equal(t, uint8(203), simpleExp(x, 205))
		assert.Equal(t, uint8(1), simpleExp(x, 255))
		assert.Equal(t, uint8(246), simpleExp(x, 254))
		assert.Equal(t, uint8(121), simpleExp(x, 212))
		assert.Equal(t, uint8(129), simpleExp(x, 255+88))
	})

	t.Run("log table size", func(t *testing.T) {
		assert.Equal(t, 256, len(globalLogTable))
		assert.Equal(t, []byte{0, 0, 0x19, 1}, globalLogTable[:4])

		assert.Equal(t, 255, len(globalExpTable))
		assert.Equal(t, []byte{1, 3, 5, 15, 17, 0x33}, globalExpTable[:6])
	})
}

func TestFastMul(t *testing.T) {
	assert.Equal(t, simpleMul(1, 2), fastMul(1, 2))
	assert.Equal(t, simpleMul(3, 2), fastMul(3, 2))
	assert.Equal(t, simpleMul(3, 3), fastMul(3, 3))
	assert.Equal(t, simpleMul(3, 3), fastMul(3, 3))
}

func TestFastMul_Ex1(t *testing.T) {
	assert.Equal(t, simpleMul(2, 13), fastMul(2, 13))
}

func TestFastMul_All(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			x := uint8(a)
			y := uint8(b)

			v1 := simpleMul(x, y)
			v2 := fastMul(x, y)
			if v1 != v2 {
				t.Errorf("Mismatch x, y = %d %d", x, y)
				return
			}
		}
	}
}
