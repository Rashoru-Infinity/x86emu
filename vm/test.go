package vm

import "encoding/binary"

func test(v *VM, op byte) {
	switch op {
	case 0x84:
		fallthrough
	case 0x85:
		w := 0b00000001 & op
		val := v.fetch()
		mod := (0b11000000 & val) >> 6
		reg := (0b00111000 & val) >> 3
		rm := 0b00000111 & val
		switch mod {
		case 0b00:
			if rm == 0b110 {
				disp := uint16(v.fetch()) | uint16(v.fetch())<<8
				switch w<<3 | reg {
				case AL:
					fallthrough
				case CL:
					fallthrough
				case DL:
					fallthrough
				case BL:
					val := v.Data[disp] & (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					val := v.Data[disp] & (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
				default:
					val := binary.LittleEndian.Uint16(v.Data[disp:]) & v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				}
				return
			}
			switch w<<3 | reg {
			case AL:
				fallthrough
			case CL:
				fallthrough
			case DL:
				fallthrough
			case BL:
				val := v.Data[eabase(v, (uint16)(rm))] & (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				val := v.Data[eabase(v, (uint16)(rm))] & (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
			default:
				val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) & v.CPU.GR[(int)(w<<3|reg)]
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			}
			return
		case 0b01:
			disp := v.fetch()
			switch w<<3 | reg {
			case AL:
				fallthrough
			case CL:
				fallthrough
			case DL:
				fallthrough
			case BL:
				if disp < 0x80 {
					val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] & (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
				} else {
					val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] & (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if disp < 0x80 {
					val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] & (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
				} else {
					val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] & (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
				}
			default:
				if disp < 0x80 {
					val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) & v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				} else {
					val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) & v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				}
			}
		case 0b10:
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			switch w<<3 | reg {
			case AL:
				fallthrough
			case CL:
				fallthrough
			case DL:
				fallthrough
			case BL:
				val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] & (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] & (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
			default:
				val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) & v.CPU.GR[(int)(w<<3|reg)]
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			}
		case 0b11:
			switch w<<3 | reg {
			case AL:
				fallthrough
			case CL:
				fallthrough
			case DL:
				fallthrough
			case BL:
				switch w<<3 | rm {
				case AL:
					fallthrough
				case CL:
					fallthrough
				case DL:
					fallthrough
				case BL:
					val := v.CPU.GR[(int)(1<<3|reg)] & 0x00ff & v.CPU.GR[(int)(1<<3|rm)] & 0x00ff
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					val := v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8 & v.CPU.GR[(int)(1<<3|rm&0x03)] & 0xff00 >> 8
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				switch w<<3 | rm {
				case AL:
					fallthrough
				case CL:
					fallthrough
				case DL:
					fallthrough
				case BL:
					val := v.CPU.GR[(int)(1<<3|reg&0x03)] & 0xff00 >> 8 & v.CPU.GR[(int)(1<<3|rm)] & 0x00ff
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					val := v.CPU.GR[(int)(1<<3|reg&0x03)] & 0xff00 >> 8 & v.CPU.GR[(int)(1<<3|rm&0x03)] & 0xff00 >> 8
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
				}
			default:
				val := v.CPU.GR[(int)(w<<3|reg)] & v.CPU.GR[(int)(w<<3|rm)]
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			}
		}
	case 0xa8:
		fallthrough
	case 0xa9:
		w := 0b00000100 & op >> 2
		if w == 1 {
			data := (uint16)(v.fetch()) | (uint16)(v.fetch())<<8
			val := v.CPU.GR[AX] & data
			v.CPU.FR[CF] = false
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		} else {
			data := (uint16)(v.fetch())
			val := v.CPU.GR[AX] & 0x00ff & data
			v.CPU.FR[CF] = false
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		}
	}
}

func grp3test(v *VM, w, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if w == 1 {
				data := uint16(v.fetch()) | uint16(v.fetch())<<8
				val := binary.LittleEndian.Uint16(v.Data[disp:]) & data
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
			} else {
				val := v.Data[disp] & v.fetch()
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
			}
			return
		}
		if w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) & data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
		} else {
			val := v.Data[eabase(v, uint16(rm))] & v.fetch()
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
		}
	case 0b01:
		disp := v.fetch()
		if w == 1 && disp < 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) & data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
		} else if w == 1 && disp >= 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) & data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
		} else if w == 0 && disp < 0x80 {
			val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] & v.fetch()
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
		} else {
			val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] & v.fetch()
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if w == 0 {
			val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(val)
		} else {
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
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
			data := v.fetch()
			val := v.CPU.GR[(int)(1<<3|rm)] & 0x00ff & uint16(data)
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			data := v.fetch()
			val := v.CPU.GR[(int)(1<<3|rm&0x03)] & 0xff00 >> 8 & uint16(data)
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
		default:
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := v.CPU.GR[(int)(w<<3|rm)] & data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
		}
	}
}
