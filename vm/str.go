package vm

import "encoding/binary"

func movs(v *VM, op byte) {
	switch op {
	case 0xa4:
		fallthrough
	case 0xa5:
		w := op & 0b00000001
		if w == 0 && v.CPU.FR[DF] {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]]
			v.CPU.GR[SI]--
			v.CPU.GR[DI]--
		} else if w == 0 && !v.CPU.FR[DF] {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]]
			v.CPU.GR[SI]++
			v.CPU.GR[DI]++
		} else if w == 1 && v.CPU.FR[DF] {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]]
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]+1] = v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]+1]
			v.CPU.GR[SI] -= 2
			v.CPU.GR[DI] -= 2
		} else {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]]
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]+1] = v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]+1]
			v.CPU.GR[SI] += 2
			v.CPU.GR[DI] += 2
		}
	}
}

func cmps(v *VM, op byte) {
	switch op {
	case 0xa6:
		fallthrough
	case 0xa7:
		w := op & 0b00000001
		if w == 0 && v.CPU.FR[DF] {
			src1 := v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]]
			src2 := v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]]
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI]--
			v.CPU.GR[DI]--
		} else if w == 0 && !v.CPU.FR[DF] {
			src1 := v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]]
			src2 := v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]]
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI]++
			v.CPU.GR[DI]++
		} else if w == 1 && v.CPU.FR[DF] {
			src1 := binary.LittleEndian.Uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]:])
			src2 := binary.LittleEndian.Uint16(v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]:])
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI] -= 2
			v.CPU.GR[DI] -= 2
		} else {
			src1 := binary.LittleEndian.Uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]:])
			src2 := binary.LittleEndian.Uint16(v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[SI]:])
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI] += 2
			v.CPU.GR[DI] += 2
		}
	}
}

func stos(v *VM, op byte) {
	switch op {
	case 0xaa:
		fallthrough
	case 0xab:
		w := op & 0b00000001
		if w == 0 && v.CPU.FR[DF] {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = byte(v.CPU.GR[AX] & 0x00ff)
			v.CPU.GR[SI]--
			v.CPU.GR[DI]--
		} else if w == 0 && !v.CPU.FR[DF] {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = byte(v.CPU.GR[AX] & 0x00ff)
			v.CPU.GR[SI]++
			v.CPU.GR[DI]++
		} else if w == 1 && v.CPU.FR[DF] {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = byte(v.CPU.GR[AX] & 0x00ff)
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]+1] = byte(v.CPU.GR[AX] & 0xff00 >> 8)
			v.CPU.GR[SI] -= 2
			v.CPU.GR[DI] -= 2
		} else {
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]] = byte(v.CPU.GR[AX] & 0x00ff)
			v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]+1] = byte(v.CPU.GR[AX] & 0xff00 >> 8)
			v.CPU.GR[SI] += 2
			v.CPU.GR[DI] += 2
		}
	}
}

func lods(v *VM, op byte) {
	switch op {
	case 0xaa:
		fallthrough
	case 0xab:
		w := op & 0b00000001
		if w == 0 && v.CPU.FR[DF] {
			v.CPU.GR[AX] = v.CPU.GR[AX]&0xff00 | uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]])
			v.CPU.GR[SI]--
			v.CPU.GR[DI]--
		} else if w == 0 && !v.CPU.FR[DF] {
			v.CPU.GR[AX] = v.CPU.GR[AX]&0xff00 | uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]])
			v.CPU.GR[SI]++
			v.CPU.GR[DI]++
		} else if w == 1 && v.CPU.FR[DF] {
			v.CPU.GR[AX] = binary.LittleEndian.Uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]:])
			v.CPU.GR[SI] -= 2
			v.CPU.GR[DI] -= 2
		} else {
			v.CPU.GR[AX] = binary.LittleEndian.Uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]:])
			v.CPU.GR[SI] += 2
			v.CPU.GR[DI] += 2
		}
	}
}

func scas(v *VM, op byte) {
	switch op {
	case 0xae:
		fallthrough
	case 0xaf:
		w := op & 0b00000001
		if w == 0 && v.CPU.FR[DF] {
			src1 := byte(v.CPU.GR[AX] & 0x00ff)
			src2 := v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]]
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI]--
			v.CPU.GR[DI]--
		} else if w == 0 && !v.CPU.FR[DF] {
			src1 := byte(v.CPU.GR[AX] & 0x00ff)
			src2 := v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]]
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI]++
			v.CPU.GR[DI]++
		} else if w == 1 && v.CPU.FR[DF] {
			src1 := v.CPU.GR[AX]
			src2 := binary.LittleEndian.Uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]:])
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI] -= 2
			v.CPU.GR[DI] -= 2
		} else {
			src1 := v.CPU.GR[AX]
			src2 := binary.LittleEndian.Uint16(v.Data[v.CPU.SR[ES]<<4+v.CPU.GR[DI]:])
			val := src1 - src2
			v.CPU.FR[CF] = src1 < src2
			v.CPU.FR[OF] = (src1 < 0x80 && src2 >= 0x80 && src1-src2 >= 0x80) || (src1 >= 0x80 && src2 < 0x80 && src1-src2 < 0x80)
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[AF] = src1&0x0f < src2&0x0f
			v.CPU.GR[SI] += 2
			v.CPU.GR[DI] += 2
		}
	}
}
