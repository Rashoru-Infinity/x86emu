package vm

import "encoding/binary"

func xor(v *VM, op byte) {
	switch op {
	case 0x30:
		fallthrough
	case 0x31:
		fallthrough
	case 0x32:
		fallthrough
	case 0x33:
		d := (0b00000010 & op) >> 1
		w := 0b00000001 & op
		val := v.fetch()
		mod := (0b11000000 & val) >> 6
		reg := (0b00111000 & val) >> 3
		rm := 0b00000111 & val
		switch mod {
		case 0b00:
			if rm == 0b110 {
				disp := (int)(v.fetch()) | ((int)(v.fetch()) << 8)
				switch w<<3 | reg {
				case AL:
					fallthrough
				case CL:
					fallthrough
				case DL:
					fallthrough
				case BL:
					val := v.Data[disp] ^ (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
					if d == 0 {
						v.Data[disp] = val
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					val := v.Data[disp] ^ (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF(val)
					if d == 0 {
						v.Data[disp] = val
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(val) << 8)
					}
				default:
					val := binary.LittleEndian.Uint16(v.Data[disp:]) ^ v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					if d == 0 {
						v.Data[disp] = (byte)(val & 0x00ff)
						v.Data[disp+1] = (byte)(val >> 8)
					} else {
						v.CPU.GR[(int)(w<<3|reg)] = val
					}
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
				val := v.Data[eabase(v, (uint16)(rm))] ^ (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))] = val
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				val := v.Data[eabase(v, (uint16)(rm))] ^ (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))] = val
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(val) << 8)
				}
			default:
				val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) ^ v.CPU.GR[(int)(w<<3|reg)]
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))] = (byte)(val & 0x00ff)
					v.Data[eabase(v, (uint16)(rm))+1] = (byte)(val >> 8)
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = val
				}
				return
			}
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
				if d == 0 {
					if disp < 0x80 {
						val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] ^ (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] ^ (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] = val
					}
				} else {
					if disp < 0x80 {
						val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] ^ (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] ^ (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
					}
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					if disp < 0x80 {
						val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] ^ (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] ^ (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] = val
					}
				} else {
					if disp < 0x80 {
						val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] ^ (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(val) << 8)
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] ^ (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(val) << 8)
					}
				}
			default:
				if d == 0 {
					if disp < 0x80 {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) ^ v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(val & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = (byte)(val >> 8)
					} else {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) ^ v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] = (byte)(val & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff+1] = (byte)(val >> 8)
					}
				} else {
					if disp < 0x80 {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) ^ v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.CPU.GR[(int)(w<<3|reg)] = val
					} else {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) ^ v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[CF] = false
						v.CPU.FR[OF] = false
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.CPU.GR[(int)(w<<3|reg)] = val
					}
				}
			}
		case 0b10:
			disp := (int)(v.fetch()) | ((int)(v.fetch()) << 8)
			switch w<<3 | reg {
			case AL:
				fallthrough
			case CL:
				fallthrough
			case DL:
				fallthrough
			case BL:
				val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] ^ (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] ^ (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]&0xff00>>8)
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF(val)
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] >> 8)
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) << 8)
				}
			default:
				val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) ^ v.CPU.GR[(int)(w<<3|reg)]
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(val >> 8)
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp)+1)] = (byte)(val & 0x00ff)
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = val
				}
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
					val := (v.CPU.GR[(int)(1<<3|reg)] & 0x00ff) ^ (v.CPU.GR[(int)(1<<3|rm)] & 0x00ff)
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
					v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff&0x0f+v.CPU.GR[(int)(1<<3|rm)]&0x00ff&0x0f > 0x0f
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | (val & 0x00ff)
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (val & 0x00ff)
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					val := (v.CPU.GR[(int)(1<<3|reg)] & 0x00ff) ^ (v.CPU.GR[(int)(1<<3|rm&0x03)] >> 8)
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm&0x03)]&0x00ff | val<<8
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | val&0x00ff
					}
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
					val := (v.CPU.GR[(int)(1<<3|reg&0x03)] >> 8) ^ (v.CPU.GR[(int)(1<<3|rm)] & 0x00ff)
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | val&0x00ff
					} else {
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | val&0x00ff<<8
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					val := (v.CPU.GR[(int)(1<<3|reg&0x03)] >> 8) ^ (v.CPU.GR[(int)(1<<3|rm&0x03)] >> 8)
					v.CPU.FR[ZF] = val&0x00ff == 0
					v.CPU.FR[SF] = val&0x00ff >= 0x80
					v.CPU.FR[CF] = false
					v.CPU.FR[OF] = false
					v.CPU.FR[PF] = checkPF((uint8)(val))
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm&0x03)] = v.CPU.GR[(int)(1<<3|rm&0x03)]&0x00ff | val&0x00ff<<8
					} else {
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | val&0x00ff<<8
					}
				}
			default:
				val := v.CPU.GR[(int)(w<<3|reg)] ^ v.CPU.GR[(int)(w<<3|rm)]
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				if d == 0 {
					v.CPU.GR[(int)(w<<3|rm)] = val
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = val
				}
			}
		}
	case 0x34:
		fallthrough
	case 0x35:
		w := 0b00000100 & op >> 2
		if w == 1 {
			data := (uint16)(v.fetch()) | (uint16)(v.fetch())<<8
			val := v.CPU.GR[AX] ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[AX] = val
		} else {
			data := (uint16)(v.fetch())
			val := v.CPU.GR[AX]&0x00ff ^ data
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[AX] = v.CPU.GR[AX]&0xff00 | val&0x00ff
		}
	}
}

func grp1xor(v *VM, w, mod, rm uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			data := v.fetch()
			val := v.Data[disp] ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(val))
			v.Data[disp] = val
			v.Data[disp] = byte(val)
			return
		}
		if w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+1] = byte(val >> 8)
			return
		}
		data := v.fetch()
		val := v.Data[eabase(v, uint16(rm))] ^ data
		v.CPU.FR[ZF] = val == 0
		v.CPU.FR[SF] = val >= 0x80
		v.CPU.FR[CF] = false
		v.CPU.FR[OF] = false
		v.CPU.FR[PF] = checkPF((uint8)(val))
		v.Data[eabase(v, uint16(rm))] = byte(val)
		return
	case 0b01:
		disp := uint16(v.fetch())
		if w == 1 && disp < 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+disp+1] = byte(val >> 8)
			return
		} else if w == 1 && disp >= 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+uint16(^disp+1)&0x00ff:]) ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+uint16(^disp+1)&0x00ff] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+uint16(^disp+1)&0x00ff+1] = byte(val >> 8)
			return
		}
		data := v.fetch()
		if disp < 0x80 {
			val := v.Data[eabase(v, uint16(rm))+disp] ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = val
			return
		} else if disp >= 0x80 {
			val := v.Data[eabase(v, uint16(rm))+uint16(^disp+1)&0x00ff] ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+uint16(^disp+1)&0x00ff] = val
			return
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+disp+1] = byte(val >> 8)
			return
		} else {
			data := v.fetch()
			val := v.Data[eabase(v, uint16(rm))+disp] ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val))
			v.Data[eabase(v, uint16(rm))+disp] = val
			return
		}
	case 0b11:
		if w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := v.CPU.GR[int(w<<3|rm)] ^ data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = false
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[int(w<<3|rm)] = val
			return
		} else {
			switch rm {
			case AL:
				fallthrough
			case CL:
				fallthrough
			case DL:
				fallthrough
			case BL:
				data := v.fetch()
				val := byte(v.CPU.GR[int(1<<3|rm)]&0x00ff) ^ data
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val))
				v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(val)
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				data := v.fetch()
				val := byte(v.CPU.GR[int(1<<3|rm&0x03)]&0xff00>>8) ^ data
				v.CPU.FR[ZF] = val&0xff == 0
				v.CPU.FR[SF] = val&0xff >= 0x80
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0x00ff | uint16(val)<<8
			default:
				data := uint16(v.fetch()) | uint16(v.fetch())<<8
				val := v.CPU.GR[int(1<<3|rm)]&0xff00 ^ data
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = false
				v.CPU.FR[OF] = false
				v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				v.CPU.GR[int(1<<3|rm)] = val
			}
		}
	}
}
