package gf256

import (
	"math/bits"
)

func add(a, b uint8) uint8 {
	return a ^ b
}

func simpleExp(n uint8, k uint16) uint8 {
	res := uint8(1)
	for i := uint16(0); i < k; i++ {
		res = simpleMul(res, n)
	}
	return res
}

func simpleMul(a, b uint8) uint8 {
	sum := uint16(0)
	count := bits.Len8(b)
	for offset := 0; offset < count; offset++ {
		mask := uint8(1 << offset)
		if b&mask == 0 {
			continue
		}
		sum ^= uint16(a) << offset
	}
	poly := uint16(0x11b)
	polyLen := 9
	sumLen := bits.Len16(sum)

	shift := sumLen - polyLen
	for ; shift >= 0; shift-- {
		mask := uint16(1 << (shift + polyLen - 1))
		if sum&mask == 0 {
			continue
		}
		sum ^= poly << shift
	}
	return uint8(sum)
}

func computeLogTable() []uint8 {
	table := make([]uint8, 256)
	x := uint8(3)
	for i := 0; i < 255; i++ {
		v := simpleExp(x, uint16(i))
		table[v] = uint8(i)
	}
	return table
}

func computeAntiLogTable() []uint8 {
	table := make([]uint8, 255)
	x := uint8(3)
	for i := 0; i < 255; i++ {
		v := simpleExp(x, uint16(i))
		table[uint8(i)] = v
	}
	return table
}

var globalLogTable = computeLogTable()
var globalExpTable = computeAntiLogTable()

func fastMul(a, b uint8) uint8 {
	if a == 0 || b == 0 {
		return 0
	}
	aLog := globalLogTable[a]
	bLog := globalLogTable[b]
	sum16 := uint16(aLog) + uint16(bLog)
	if sum16 >= 255 {
		sum16 -= 255
	}
	return globalExpTable[sum16]
}

func computeMulTable() []uint8 {
	res := make([]byte, 256*256)
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			res[i*256+j] = simpleMul(uint8(i), uint8(j))
		}
	}
	return res
}

var globalMulTable = computeMulTable()

func tableMul(a, b uint8) uint8 {
	i := uint16(a)
	j := uint16(b)
	return globalMulTable[i<<8|j]
}

func tableMul4(a, b uint32) uint32 {
	x0 := uint8(a & 0xff)
	a >>= 8
	x1 := uint8(a & 0xff)
	a >>= 8
	x2 := uint8(a & 0xff)
	a >>= 8
	x3 := uint8(a & 0xff)

	y0 := uint8(b & 0xff)
	b >>= 8
	y1 := uint8(b & 0xff)
	b >>= 8
	y2 := uint8(b & 0xff)
	b >>= 8
	y3 := uint8(b & 0xff)

	res := uint32(tableMul(x3, y3))
	res <<= 8
	res |= uint32(tableMul(x2, y2))
	res <<= 8
	res |= uint32(tableMul(x1, y1))
	res <<= 8
	res |= uint32(tableMul(x0, y0))

	return res
}

func simpleInv(x uint8) uint8 {
	logX := globalLogTable[x]
	if logX == 0 {
		return 1
	}
	return globalExpTable[255-logX]
}

func computeInvTable() []uint8 {
	result := make([]uint8, 256)
	for i := 1; i < 256; i++ {
		x := uint8(i)
		result[i] = simpleInv(x)
	}
	return result
}

var invTable = computeInvTable()

func tableInv(x uint8) uint8 {
	return invTable[x]
}
