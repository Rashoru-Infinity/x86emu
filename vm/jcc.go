package vm

import "encoding/binary"

func jcc8(v *VM, op byte) {
	const (
		opsize = 2
	)
	offset := v.IP - 1
	data := uint32(v.fetch())
	switch op {
	case 0x70:
		if v.CPU.FR[OF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[OF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x71:
		if !v.CPU.FR[OF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if !v.CPU.FR[OF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x72:
		if v.CPU.FR[CF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[CF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x73:
		if !v.CPU.FR[CF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if !v.CPU.FR[CF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x74:
		if v.CPU.FR[ZF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[ZF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x75:
		if !v.CPU.FR[ZF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if !v.CPU.FR[ZF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x76:
		if (v.CPU.FR[CF] || v.CPU.FR[ZF]) && data < 0x80 {
			v.IP = data + offset + opsize
		} else if (v.CPU.FR[CF] || v.CPU.FR[ZF]) && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x77:
		if (!v.CPU.FR[CF] && !v.CPU.FR[ZF]) && data < 0x80 {
			v.IP = data + offset + opsize
		} else if (!v.CPU.FR[CF] && !v.CPU.FR[ZF]) && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x78:
		if v.CPU.FR[SF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[SF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x79:
		if !v.CPU.FR[SF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if !v.CPU.FR[SF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x7a:
		if v.CPU.FR[PF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[PF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x7b:
		if !v.CPU.FR[PF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if !v.CPU.FR[PF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x7c:
		if v.CPU.FR[SF] != v.CPU.FR[OF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[SF] != v.CPU.FR[OF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x7d:
		if v.CPU.FR[SF] == v.CPU.FR[OF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[SF] == v.CPU.FR[OF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x7e:
		if (v.CPU.FR[SF] != v.CPU.FR[OF] || v.CPU.FR[ZF]) && data < 0x80 {
			v.IP = data + offset + opsize
		} else if (v.CPU.FR[SF] != v.CPU.FR[OF] || v.CPU.FR[ZF]) && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0x7f:
		if v.CPU.FR[SF] == v.CPU.FR[OF] && !v.CPU.FR[ZF] && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[SF] == v.CPU.FR[OF] && !v.CPU.FR[ZF] && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0xe0:
		count := v.CPU.GR[CX] - 1
		if !v.CPU.FR[ZF] && count != 0 && data < 0x80 {
			v.IP = data + offset + opsize
		} else if !v.CPU.FR[ZF] && count != 0 && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0xe1:
		count := v.CPU.GR[CX] - 1
		if v.CPU.FR[ZF] && count != 0 && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.FR[ZF] && count != 0 && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0xe2:
		count := v.CPU.GR[CX] - 1
		if count != 0 && data < 0x80 {
			v.IP = data + offset + opsize
		} else if count != 0 && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0xe3:
		if v.CPU.GR[CX] == 0 && data < 0x80 {
			v.IP = data + offset + opsize
		} else if v.CPU.GR[CX] == 0 && data >= 0x80 {
			v.IP = (^data + 1) + offset + opsize
		}
	case 0xe9:
		data |= uint32(v.fetch()) << 8
		if data < 0x8000 {
			v.IP += uint32(data)
		} else {
			v.IP -= uint32(^data + 1)
		}
	case 0xeb:
		if data < 0x80 {
			v.IP += data
		} else if data >= 0x80 {
			v.IP += 0xffffff00 | data
		}
	}
	v.IP &= 0x0000ffff
}

func grp5jmp(v *VM, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			v.IP = uint32(binary.LittleEndian.Uint16(v.Data[disp:]))
			return
		}
		v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]))
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+uint16(disp):]))
		} else {
			v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-uint16(^disp+1):]))
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]))
	case 0b11:
		v.IP = uint32(v.CPU.GR[int(1<<3|rm)])
	}
}
