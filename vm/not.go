package vm

import "encoding/binary"

func grp3not(v *VM, w, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if w == 1 {
				val := ^binary.LittleEndian.Uint16(v.Data[disp:])
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
				v.Data[disp] = byte(val & 0x00ff)
				v.Data[disp+1] = byte(val >> 8)
			} else {
				val := ^v.Data[disp]
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
				v.Data[disp] = val
			}
			return
		}
		if w == 1 {
			val := ^binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+1] = byte(val >> 8)
		} else {
			val := ^v.Data[eabase(v, uint16(rm))]
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
			v.Data[eabase(v, uint16(rm))] = val
		}
	case 0b01:
		disp := v.fetch()
		if w == 1 && disp < 0x80 {
			val := ^binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
			v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = byte(val & 0x00ff)
			v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = byte(val >> 8)
		} else if w == 1 && disp >= 0x80 {
			val := ^binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
			v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = byte(val & 0x00ff)
			v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)+1] = byte(val >> 8)
		} else if w == 0 && disp < 0x80 {
			val := ^v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
			v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
		} else {
			val := ^v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
			v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = val
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if w == 0 {
			val := ^v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
			v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
		} else {
			val := ^binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
			v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = byte(val & 0x00ff)
			v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = byte(val >> 8)
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
			val := ^(v.CPU.GR[(int)(1<<3|rm)] & 0x00ff)
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
			v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | val&0x00ff
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			val := ^(v.CPU.GR[(int)(1<<3|rm&0x03)] >> 8)
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
			v.CPU.GR[(int)(1<<3|rm&0x03)] = v.CPU.GR[(int)(1<<3|rm&0x03)]&0x00ff | val<<8
		default:
			val := ^v.CPU.GR[(int)(w<<3|rm)]
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
			v.CPU.GR[(int)(w<<3|rm)] = val
		}
	}
}
