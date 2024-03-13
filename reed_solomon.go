package gf256

func reedSolomonEncode(data [2]byte) [6]byte {
	y0 := data[0]
	y1 := data[1]

	return [6]byte{
		y0,
		y1,
		tableMul(y0, 2^1) ^ tableMul(y1, 2),
		tableMul(y0, 3^1) ^ tableMul(y1, 3),
		tableMul(y0, 4^1) ^ tableMul(y1, 4),
		tableMul(y0, 5^1) ^ tableMul(y1, 5),
	}
}

func reedSolomonDecode(data [6]byte, trueLabel [6]bool) [2]byte {
	x0 := uint8(255)
	x1 := uint8(255)

	for i := uint8(0); i < 6; i++ {
		if trueLabel[i] {
			x0 = i
			break
		}
	}
	for i := x0 + 1; i < 6; i++ {
		if trueLabel[i] {
			x1 = i
			break
		}
	}

	dx := x0 ^ x1
	dxInv := tableInv(dx)

	y0 := data[x0]
	y1 := data[x1]

	sum1 := tableMul(y0, x1) ^ tableMul(y1, x0)
	sum2 := tableMul(y0, x1^1) ^ tableMul(y1, x0^1)

	return [2]byte{
		tableMul(sum1, dxInv),
		tableMul(sum2, dxInv),
	}
}
