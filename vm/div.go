package vm

import "encoding/binary"

func grp3div(v *VM, w, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if w == 1 {
				val := (uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) / uint32(binary.LittleEndian.Uint16(v.Data[disp:]))
				v.CPU.GR[DX] = uint16((uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) % uint32(binary.LittleEndian.Uint16(v.Data[disp:])))
				v.CPU.GR[AX] = uint16(val)
			} else {
				val := v.CPU.GR[AX] / uint16(v.Data[disp])
				v.CPU.GR[AX] = (v.CPU.GR[AX]%uint16(v.Data[disp]))<<8 | val
			}
			return
		}
		if w == 1 {
			val := (uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) / uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]))
			v.CPU.GR[DX] = uint16((uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) % uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])))
			v.CPU.GR[AX] = uint16(val)
		} else {
			val := v.CPU.GR[AX] / uint16(v.Data[eabase(v, uint16(rm))])
			v.CPU.GR[AX] = v.CPU.GR[AX]%uint16(v.Data[eabase(v, uint16(rm))])<<8 | uint16(val)
		}
	case 0b01:
		disp := v.fetch()
		if w == 1 && disp < 0x80 {
			val := (uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) / uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.GR[DX] = uint16((uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) % uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])))
			v.CPU.GR[AX] = uint16(val)
		} else if w == 1 && disp >= 0x80 {
			val := (uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) / uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]))
			v.CPU.GR[DX] = uint16((uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) % uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])))
			v.CPU.GR[AX] = uint16(val)
		} else if w == 0 && disp < 0x80 {
			val := v.CPU.GR[AX] / uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.GR[AX] = v.CPU.GR[AX]%uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])<<8 | uint16(val)
		} else {
			val := v.CPU.GR[AX] / uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)])
			v.CPU.GR[AX] = v.CPU.GR[AX]%uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)])<<8 | uint16(val)
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if w == 0 {
			val := v.CPU.GR[AX] / uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.GR[AX] = (v.CPU.GR[AX]%uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]))<<8 | uint16(val)
		} else {
			val := uint32(v.CPU.GR[AX]) / uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.GR[DX] = uint16(uint32(v.CPU.GR[AX]) % uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])))
			v.CPU.GR[AX] = uint16(val)
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
			val := v.CPU.GR[AX] / (v.CPU.GR[(int)(1<<3|rm)] & 0x00ff)
			v.CPU.GR[AX] = (v.CPU.GR[AX]%(v.CPU.GR[(int)(1<<3|rm)]&0x00ff))<<8 | uint16(val)
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			val := v.CPU.GR[AX] & 0x00ff / (v.CPU.GR[(int)(1<<3|rm&0x03)] >> 8)
			v.CPU.GR[AX] = (v.CPU.GR[AX]&0x00ff%(v.CPU.GR[(int)(1<<3|rm&0x03)]&0xff00>>8))<<8 | uint16(val)
		default:
			val := (uint32(v.CPU.GR[DX])<<16 | uint32(v.CPU.GR[AX])) / uint32(v.CPU.GR[(int)(w<<3|rm)])
			v.CPU.GR[DX] = uint16(uint32(v.CPU.GR[AX]) % uint32(v.CPU.GR[(int)(w<<3|rm)]))
			v.CPU.GR[AX] = uint16(val)
		}
	}
}

func grp3idiv(v *VM, w, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if w == 1 {
				val := (int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) / int32(binary.LittleEndian.Uint16(v.Data[disp:]))
				v.CPU.GR[DX] = uint16((int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) % int32(binary.LittleEndian.Uint16(v.Data[disp:])))
				v.CPU.GR[AX] = uint16(val)
			} else {
				val := int16(v.CPU.GR[AX]) / int16(v.Data[disp])
				v.CPU.GR[AX] = uint16((int16(v.CPU.GR[AX])%int16(v.Data[disp]))<<8 | val)
			}
			return
		}
		if w == 1 {
			val := (int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) / int32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]))
			v.CPU.GR[DX] = uint16((int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) % int32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])))
			v.CPU.GR[AX] = uint16(val)
		} else {
			val := int16(v.CPU.GR[AX]) / int16(v.Data[eabase(v, uint16(rm))])
			v.CPU.GR[AX] = uint16(int16(v.CPU.GR[AX])%int16(v.Data[eabase(v, uint16(rm))])<<8) | uint16(val)
		}
	case 0b01:
		disp := v.fetch()
		if w == 1 && disp < 0x80 {
			val := (int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) / int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.GR[DX] = uint16((int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) % int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])))
			v.CPU.GR[AX] = uint16(val)
		} else if w == 1 && disp >= 0x80 {
			val := (int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) / int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]))
			v.CPU.GR[DX] = uint16((int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) % int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])))
			v.CPU.GR[AX] = uint16(val)
		} else if w == 0 && disp < 0x80 {
			val := int16(v.CPU.GR[AX]) / int16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.GR[AX] = uint16(int16(v.CPU.GR[AX])%int16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])<<8) | uint16(val)
		} else {
			val := int16(v.CPU.GR[AX]) / int16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)])
			v.CPU.GR[AX] = uint16(int16(v.CPU.GR[AX])%int16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)])<<8) | uint16(val)
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if w == 0 {
			val := int16(v.CPU.GR[AX]) / int16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
			v.CPU.GR[AX] = uint16(int16(v.CPU.GR[AX])%int16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]))<<8 | uint16(val)
		} else {
			val := int32(v.CPU.GR[AX]) / int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]))
			v.CPU.GR[DX] = uint16(int32(v.CPU.GR[AX]) % int32(binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])))
			v.CPU.GR[AX] = uint16(val)
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
			val := int16(v.CPU.GR[AX]) / int16(v.CPU.GR[(int)(1<<3|rm)]&0x00ff)
			v.CPU.GR[AX] = uint16(int16(v.CPU.GR[AX])%int16(v.CPU.GR[(int)(1<<3|rm)]&0x00ff))<<8 | uint16(val)
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			val := int16(v.CPU.GR[AX]&0x00ff) / int16(v.CPU.GR[(int)(1<<3|rm&0x03)]&0xff00>>8)
			v.CPU.GR[AX] = uint16(int16(v.CPU.GR[AX]&0x00ff)%int16(v.CPU.GR[(int)(1<<3|rm&0x03)]&0xff00>>8))<<8 | uint16(val)
		default:
			val := (int32(v.CPU.GR[DX])<<16 | int32(v.CPU.GR[AX])) / int32(v.CPU.GR[(int)(w<<3|rm)])
			v.CPU.GR[DX] = uint16(int32(v.CPU.GR[AX]) % int32(v.CPU.GR[(int)(w<<3|rm)]))
			v.CPU.GR[AX] = uint16(val)
		}
	}
}
