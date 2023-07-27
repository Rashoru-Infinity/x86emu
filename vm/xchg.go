package vm

import "encoding/binary"

func xchg(v *VM, op byte) {
	switch op {
	case 0x86:
		fallthrough
	case 0x87:
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
					tmp := v.Data[disp]
					v.Data[disp] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | uint16(tmp)
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					tmp := v.Data[disp]
					v.Data[disp] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | uint16(tmp)<<8
				default:
					tmp := binary.LittleEndian.Uint16(v.Data[disp:])
					v.Data[disp] = byte(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
					v.Data[disp+1] = byte(v.CPU.GR[(int)(w<<3|reg)] >> 8)
					v.CPU.GR[(int)(w<<3|reg)] = tmp
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
				tmp := v.Data[eabase(v, (uint16)(rm))]
				v.Data[eabase(v, (uint16)(rm))] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
				v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | uint16(tmp)
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				tmp := v.Data[eabase(v, (uint16)(rm))]
				v.Data[eabase(v, (uint16)(rm))] = (byte)(v.CPU.GR[(int)(1<<3|reg)] >> 8)
				v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | uint16(tmp)<<8
			default:
				tmp := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])
				v.Data[eabase(v, (uint16)(rm))] = byte(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
				v.Data[eabase(v, (uint16)(rm))+1] = byte(v.CPU.GR[(int)(w<<3|reg)] >> 8)
				v.CPU.GR[(int)(w<<3|reg)] = tmp
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
					tmp := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = byte(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | uint16(tmp)
				} else {
					tmp := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]
					v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = byte(v.CPU.GR[(int)(1<<3|reg)] & 0x00ff)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | uint16(tmp)
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if disp < 0x80 {
					tmp := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | uint16(tmp)<<8
				} else {
					tmp := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]
					v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | uint16(tmp)<<8
				}
			default:
				if disp < 0x80 {
					tmp := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = byte(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = byte(v.CPU.GR[(int)(w<<3|reg)] >> 8)
					v.CPU.GR[(int)(w<<3|reg)] = tmp
				} else {
					tmp := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
					v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = byte(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
					v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)+1] = byte(v.CPU.GR[(int)(w<<3|reg)] >> 8)
					v.CPU.GR[(int)(w<<3|reg)] = tmp
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
				tmp := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
				v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])
				v.CPU.GR[(int)(1<<3|reg)&0x00ff] = uint16(tmp)
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				tmp := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
				v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] >> 8)
				v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | uint16(tmp)<<8
			default:
				tmp := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
				v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = byte(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
				v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = byte(v.CPU.GR[(int)(w<<3|reg)] >> 8)
				v.CPU.GR[(int)(w<<3|reg)] = tmp
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
					tmp := v.CPU.GR[(int)(1<<3|reg)] & 0x00ff
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | v.CPU.GR[(int)(1<<3|rm)]&0x00ff
					v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | tmp
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					tmp := v.CPU.GR[(int)(1<<3|reg)] >> 8
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | v.CPU.GR[(int)(1<<3|rm)]&0xff00
					v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff | tmp<<8
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
					tmp := v.CPU.GR[(int)(1<<3|reg)] >> 8
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | v.CPU.GR[(int)(1<<3|rm)]&0x00ff
					v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | tmp
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					tmp := v.CPU.GR[(int)(1<<3|reg)] >> 8
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | v.CPU.GR[(int)(1<<3|rm)]&0xff00
					v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff | tmp<<8
				}
			default:
				tmp := v.CPU.GR[(int)(w<<3|reg)]
				v.CPU.GR[(int)(w<<3|reg)] = v.CPU.GR[(int)(w<<3|rm)]
				v.CPU.GR[(int)(w<<3|rm)] = tmp
			}
		}
	case 0x90:
		fallthrough
	case 0x91:
		fallthrough
	case 0x92:
		fallthrough
	case 0x93:
		fallthrough
	case 0x94:
		fallthrough
	case 0x95:
		fallthrough
	case 0x96:
		fallthrough
	case 0x97:
		tmp := v.CPU.GR[AX]
		v.CPU.GR[AX] = v.CPU.GR[int(1<<3|op-0x90)]
		v.CPU.GR[int(1<<3|op-0x90)] = tmp
	}
}
