package vm

import "encoding/binary"

func grp3mul(v *VM, w, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if w == 1 {
				val := uint32(v.CPU.GR[AX]) * uint32(binary.LittleEndian.Uint16(v.Data[disp:]))
				v.CPU.FR[CF] = val&0xffff0000 != 0
				v.CPU.FR[OF] = val&0xffff0000 != 0
				v.CPU.GR[DX] = uint16(val >> 16)
				v.CPU.GR[AX] = uint16(val & 0x0000ffff)
			} else {
				val := v.CPU.GR[AX] & 0x00ff * uint16(v.Data[disp])
				v.CPU.FR[CF] = val&0xff00 != 0
				v.CPU.FR[OF] = val&0xff00 != 0
				v.CPU.GR[AX] = uint16(val)
			}
			return
		}
		if w == 1 {
			val := uint32(v.CPU.GR[AX]) * uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]))
			v.CPU.FR[CF] = val&0xffff0000 != 0
			v.CPU.FR[OF] = val&0xffff0000 != 0
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
		} else {
			val := v.CPU.GR[AX] & 0x00ff * uint16(v.Data[eabase(v, uint16(rm))])
			v.CPU.FR[CF] = val&0xff00 != 0
			v.CPU.FR[OF] = val&0xff00 != 0
			v.CPU.GR[AX] = uint16(val)
		}
	case 0b01:
		disp := v.fetch()
		if w == 1 && disp < 0x80 {
			val := uint32(v.CPU.GR[AX]) * uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.FR[CF] = val&0xffff0000 != 0
			v.CPU.FR[OF] = val&0xffff0000 != 0
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
		} else if w == 1 && disp >= 0x80 {
			val := uint32(v.CPU.GR[AX]) * uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]))
			v.CPU.FR[CF] = val&0xffff0000 != 0
			v.CPU.FR[OF] = val&0xffff0000 != 0
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
		} else if w == 0 && disp < 0x80 {
			val := v.CPU.GR[AX] & 0x00ff * uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.FR[CF] = val&0xff00 != 0
			v.CPU.FR[OF] = val&0xff00 != 0
			v.CPU.GR[AX] = uint16(val)
		} else {
			val := v.CPU.GR[AX] & 0x00ff * uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)])
			v.CPU.FR[CF] = val&0xff00 != 0
			v.CPU.FR[OF] = val&0xff00 != 0
			v.CPU.GR[AX] = uint16(val)
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if w == 0 {
			val := v.CPU.GR[AX] & 0x00ff * uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.FR[CF] = val&0xff00 != 0
			v.CPU.FR[OF] = val&0xff00 != 0
			v.CPU.GR[AX] = uint16(val)
		} else {
			val := uint32(v.CPU.GR[AX]) * uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.FR[CF] = val&0xffff0000 != 0
			v.CPU.FR[OF] = val&0xffff0000 != 0
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
		}
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			val := v.CPU.GR[AX] & 0x00ff * (v.CPU.GR[(int)(1<<3|rm)] & 0x00ff)
			v.CPU.FR[CF] = val&0xff00 != 0
			v.CPU.FR[OF] = val&0xff00 != 0
			v.CPU.GR[AX] = uint16(val)
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			val := v.CPU.GR[AX] & 0x00ff * (v.CPU.GR[(int)(1<<3|rm)] >> 8)
			v.CPU.FR[CF] = val&0xff00 != 0
			v.CPU.FR[OF] = val&0xff00 != 0
			v.CPU.GR[AX] = uint16(val)
		default:
			val := uint32(v.CPU.GR[AX]) * uint32(v.CPU.GR[(int)(w<<3|rm)])
			v.CPU.FR[CF] = val&0xffff0000 != 0
			v.CPU.FR[OF] = val&0xffff0000 != 0
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
		}
	}
}

func grp3imul(v *VM, w, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if w == 1 {
				val := int32(v.CPU.GR[AX]) * int32(binary.LittleEndian.Uint16(v.Data[disp:]))
				v.CPU.GR[DX] = uint16(val >> 16)
				v.CPU.GR[AX] = uint16(val & 0x0000ffff)
				v.CPU.FR[CF] = int32(int8(val&0x0000ffff)) != val
				v.CPU.FR[OF] = int32(int8(val&0x0000ffff)) != val
			} else {
				val := int16(v.CPU.GR[AX]&0x00ff) * int16(v.Data[disp])
				v.CPU.GR[AX] = uint16(val)
				v.CPU.FR[CF] = int16(int8(val&0x00ff)) != val
				v.CPU.FR[OF] = int16(int8(val&0x00ff)) != val
			}
			return
		}
		if w == 1 {
			val := int32(v.CPU.GR[AX]) * int32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]))
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
			v.CPU.FR[CF] = int32(int8(val&0x0000ffff)) != val
			v.CPU.FR[OF] = int32(int8(val&0x0000ffff)) != val
		} else {
			val := int16(v.CPU.GR[AX]&0x00ff) * int16(v.Data[eabase(v, uint16(rm))])
			v.CPU.GR[AX] = uint16(val)
			v.CPU.FR[CF] = int16(int8(val&0x00ff)) != val
			v.CPU.FR[OF] = int16(int8(val&0x00ff)) != val
		}
	case 0b01:
		disp := v.fetch()
		if w == 1 && disp < 0x80 {
			val := int32(v.CPU.GR[AX]) * int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
			v.CPU.FR[CF] = int32(int8(val&0x0000ffff)) != val
			v.CPU.FR[OF] = int32(int8(val&0x0000ffff)) != val
		} else if w == 1 && disp >= 0x80 {
			val := int32(v.CPU.GR[AX]) * int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]))
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
			v.CPU.FR[CF] = int32(int8(val&0x0000ffff)) != val
			v.CPU.FR[OF] = int32(int8(val&0x0000ffff)) != val
		} else if w == 0 && disp < 0x80 {
			val := int16(v.CPU.GR[AX]) & 0x00ff * int16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.GR[AX] = uint16(val)
			v.CPU.FR[CF] = int16(int8(val&0x00ff)) != val
			v.CPU.FR[OF] = int16(int8(val&0x00ff)) != val
		} else {
			val := int16(v.CPU.GR[AX]&0x00ff) * int16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)])
			v.CPU.GR[AX] = uint16(val)
			v.CPU.FR[CF] = int16(int8(val&0x00ff)) != val
			v.CPU.FR[OF] = int16(int8(val&0x00ff)) != val
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if w == 0 {
			val := int16(v.CPU.GR[AX]&0x00ff) * int16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.GR[AX] = uint16(val)
			v.CPU.FR[CF] = int16(int8(val&0x00ff)) != val
			v.CPU.FR[OF] = int16(int8(val&0x00ff)) != val
		} else {
			val := int32(v.CPU.GR[AX]) * int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
			v.CPU.FR[CF] = int32(int8(val&0x0000ffff)) != val
			v.CPU.FR[OF] = int32(int8(val&0x0000ffff)) != val
		}
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			val := int16(v.CPU.GR[AX]&0x00ff) * int16(v.CPU.GR[(int)(1<<3|rm)]&0x00ff)
			v.CPU.GR[AX] = uint16(val)
			v.CPU.FR[CF] = int16(int8(val&0x00ff)) != val
			v.CPU.FR[OF] = int16(int8(val&0x00ff)) != val
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			val := int16(v.CPU.GR[AX]&0x00ff) * int16(v.CPU.GR[(int)(1<<3|rm&0x03)]&0xff00>>8)
			v.CPU.GR[AX] = uint16(val)
			v.CPU.FR[CF] = int16(int8(val&0x00ff)) != val
			v.CPU.FR[OF] = int16(int8(val&0x00ff)) != val
		default:
			val := int32(v.CPU.GR[AX]) * int32(v.CPU.GR[(int)(w<<3|rm)])
			v.CPU.GR[DX] = uint16(val >> 16)
			v.CPU.GR[AX] = uint16(val & 0x0000ffff)
			v.CPU.FR[CF] = int32(int8(val&0x0000ffff)) != val
			v.CPU.FR[OF] = int32(int8(val&0x0000ffff)) != val
		}
	}
}
