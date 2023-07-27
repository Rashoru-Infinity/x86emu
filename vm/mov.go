package vm

import (
	"encoding/binary"
)

func mov(v *VM, op byte) {
	switch op {
	case 0x88:
		fallthrough
	case 0x89:
		fallthrough
	case 0x8a:
		fallthrough
	case 0x8b:
		d := (0b00000010 & op) >> 1
		w := 0b00000001 & op
		data := v.fetch()
		mod := (0b11000000 & data) >> 6
		reg := (0b00111000 & data) >> 3
		rm := 0b00000111 & data
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
					if d == 0 {
						v.Data[disp] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[AX]&0xff00 | (uint16)(v.Data[disp])
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					if d == 0 {
						v.Data[disp] = (byte)((v.CPU.GR[(int)(1<<3|reg)] & 0xff00) >> 8)
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(v.Data[disp]) << 8)
					}
				default:
					if d == 0 {
						v.Data[disp] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
						v.Data[disp+1] = (byte)(v.CPU.GR[(int)(w<<3|reg)] >> 8)
					} else {
						v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[disp:])
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
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm))])
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))] = (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)] >> 8)
				} else {
					v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))]) << 8)
				}
			default:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
					v.Data[eabase(v, (uint16)(rm))+1] = (byte)(v.CPU.GR[(int)(w<<3|reg)] >> 8)
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])
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
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
					} else {
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
					}
				} else {
					if disp < 0x80 {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff])
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
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)] & 0xff00 >> 8)
					} else {
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] = (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)] & 0xff00 >> 8)
					}
				} else {
					if disp < 0x80 {
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) << 8)
					} else {
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]) << 8)
					}
				}
			default:
				if d == 0 {
					if disp < 0x80 {
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0xff00 >> 8)
					} else {
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff+1] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0xff00 >> 8)
					}
				} else {
					if disp < 0x80 {
						v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
					} else {
						v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:])
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
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)] & 0xff00 >> 8)
				} else {
					v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) << 8)
				}
			default:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0xff00 >> 8)
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp)+1)] = (byte)(v.CPU.GR[(int)(w<<3|reg)&0x00ff])
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
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
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | v.CPU.GR[(int)(1<<3|reg)]*0x00ff
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | v.CPU.GR[(int)(1<<3|rm)]&0x00ff
					}

				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | v.CPU.GR[(int)(1<<3|reg)]&0xff00>>8
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | v.CPU.GR[(int)(1<<3|rm)]&0xff00>>8
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
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff | v.CPU.GR[(int)(1<<3|reg)]&0x00ff<<8
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | v.CPU.GR[(int)(1<<3|rm)]&0x00ff<<8
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					if d == 0 {
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff | v.CPU.GR[(int)(1<<3|reg)]&0xff00
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | v.CPU.GR[(int)(1<<3|rm)]&0xff00
					}
				}
			default:
				if d == 0 {
					v.CPU.GR[(int)(w<<3|rm)] = v.CPU.GR[(int)(w<<3|reg)]
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = v.CPU.GR[(int)(w<<3|rm)]
				}
			}
		}
	case 0x8c:
	case 0x8e:
		d := op & 0b00000010 >> 1
		w := op & 0b00000001
		data := v.fetch()
		mod := data & 0b11000000 >> 6
		reg := data & 0b00011000 >> 3
		rm := data & 0b00000111
		switch mod {
		case 0b00:
			if rm == 0b110 {
				disp := (int)(v.fetch()) | ((int)(v.fetch()) << 8)
				if d == 0 {
					v.Data[disp] = byte(v.CPU.SR[(int)(reg)] & 0x00ff)
					v.Data[disp+1] = byte(v.CPU.SR[(int)(reg)] >> 8)
				} else {
					v.CPU.SR[(int)(reg)] = binary.LittleEndian.Uint16(v.Data[disp:])
				}
				return
			}
			if d == 0 {
				v.Data[eabase(v, uint16(rm))] = byte(v.CPU.SR[int(reg)] & 0x00ff)
				v.Data[eabase(v, uint16(rm))+1] = byte(v.CPU.SR[int(reg)] >> 8)
			} else {
				v.CPU.SR[int(reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])
			}
		case 0b01:
			disp := v.fetch()
			if d == 0 {
				if disp < 0x80 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.SR[(int)(reg)] & 0x00ff)
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = (byte)(v.CPU.SR[(int)(reg)] & 0xff00 >> 8)
				} else {
					v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] = (byte)(v.CPU.SR[(int)(reg)] & 0x00ff)
					v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff+1] = (byte)(v.CPU.SR[(int)(reg)] & 0xff00 >> 8)
				}
			} else {
				if disp < 0x80 {
					v.CPU.SR[(int)(reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
				} else {
					v.CPU.SR[(int)(reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:])
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
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8)
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) << 8)
				}
			default:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0xff00 >> 8)
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp)+1)] = (byte)(v.CPU.GR[(int)(w<<3|reg)&0x00ff])
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
				}
			}
		case 0b11:
			if d == 0 {
				v.CPU.GR[int(1<<3|rm)] = v.CPU.SR[int(reg)]
			} else {
				v.CPU.SR[int(reg)] = v.CPU.GR[int(1<<3|rm)]
			}
		}
	case 0xa0:
		fallthrough
	case 0xa1:
		fallthrough
	case 0xa2:
		fallthrough
	case 0xa3:
		dw := op & 0b00000011
		addr := uint16(v.fetch()) | uint16(v.fetch())<<8
		switch dw {
		case 00:
			v.CPU.GR[AX] = v.CPU.GR[AX]&0xff00 | uint16(v.Data[addr])
		case 01:
			v.Data[addr] = byte(v.CPU.GR[AX] & 0x00ff)
		case 10:
			v.CPU.GR[AX] = binary.LittleEndian.Uint16(v.Data[addr:])
		case 11:
			v.Data[addr] = byte(v.CPU.GR[AX] & 0x00ff)
			v.Data[addr+1] = byte(v.CPU.GR[AX] & 0xff00 >> 8)
		}
	case 0xb0:
		fallthrough
	case 0xb1:
		fallthrough
	case 0xb2:
		fallthrough
	case 0xb3:
		fallthrough
	case 0xb4:
		fallthrough
	case 0xb5:
		fallthrough
	case 0xb6:
		fallthrough
	case 0xb7:
		fallthrough
	case 0xb8:
		fallthrough
	case 0xb9:
		fallthrough
	case 0xba:
		fallthrough
	case 0xbb:
		fallthrough
	case 0xbc:
		fallthrough
	case 0xbd:
		fallthrough
	case 0xbe:
		fallthrough
	case 0xbf:
		w := (int)(0b00001000&op) >> 3
		reg := (int)(0b00000111 & op)
		data := (uint16)(v.fetch())
		if w == 0 { // 8bit data
			if reg <= BL { // low
				v.CPU.GR[1<<3|reg] &= 0xff00
				v.CPU.GR[1<<3|reg] |= data
			} else { // high
				v.CPU.GR[1<<3|reg] &= 0x00ff
				v.CPU.GR[1<<3|reg] |= data << 8
			}
		} else { // 16bit data
			data |= ((uint16)(v.fetch())) << 8
			v.CPU.GR[w<<3|reg] = data
		}
	/*
		if v.Debug.DebugMode {
			for _, b := range v.Debug.Buf {
				fmt.Fprintf(os.Stderr, "%02x", b)
			}
			// padding
			for i := len(v.Debug.Buf) * 2; i < 13; i++ {
				fmt.Fprintf(os.Stderr, " ")
			}
			fmt.Fprintf(os.Stderr, "mov ")
			fmt.Fprintf(os.Stderr, "%s, ", grname[(int)(v.Debug.Buf[0]&0b00001000|v.Debug.Buf[0]&0b00000111)])
			if v.Debug.Buf[0]&0b00001000 != 0 {
				fmt.Fprintf(os.Stderr, "%02x", v.Debug.Buf[2])
			}
			fmt.Fprintf(os.Stderr, "%02x", v.Debug.Buf[1])
			fmt.Fprintf(os.Stderr, "\n")
		}
	*/
	case 0xc6:
		fallthrough
	case 0xc7:
		w := op & 0b00000001
		data := v.fetch()
		mod := data & 0b11000000 >> 6
		rm := data & 0b00000111
		switch mod {
		case 0b00:
			if rm == 0b110 {
				disp := uint16(v.fetch()) | uint16(v.fetch())<<8
				v.Data[disp] = v.fetch()
				if w == 1 {
					v.Data[disp+1] = v.fetch()
				}
				return
			}
			v.Data[eabase(v, uint16(rm))] = v.fetch()
			if w == 1 {
				v.Data[eabase(v, uint16(rm))+1] = v.fetch()
			}
		case 0b01:
			disp := v.fetch()
			if disp < 0x80 {
				v.Data[eabase(v, uint16(rm))+uint16(disp)] = v.fetch()
				if w == 1 {
					v.Data[eabase(v, uint16(rm))+uint16(disp)+1] = v.fetch()
				}
			} else {
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = v.fetch()
				if w == 1 {
					v.Data[eabase(v, uint16(rm))-uint16(^disp+1)+1] = v.fetch()
				}
			}
		case 0b10:
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			v.Data[eabase(v, uint16(rm))+disp] = v.fetch()
			if w == 1 {
				v.Data[eabase(v, uint16(rm))+disp+1] = v.fetch()
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
				v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(v.fetch())
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0x00ff | uint16(v.fetch())<<8
			default:
				v.CPU.GR[int(w<<3|rm)] = uint16(v.fetch()) | uint16(v.fetch())<<8
			}
		}
	}
}
