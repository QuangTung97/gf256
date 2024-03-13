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

func computeLogTable() map[uint8]uint8 {
	table := map[uint8]uint8{}
	x := uint8(3)
	for i := 0; i < 255; i++ {
		v := simpleExp(x, uint16(i))
		table[v] = uint8(i)
	}
	return table
}

func computeAntiLogTable() map[uint8]uint8 {
	table := map[uint8]uint8{}
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
	sum := uint8(sum16 % 255)
	return globalExpTable[sum]
}