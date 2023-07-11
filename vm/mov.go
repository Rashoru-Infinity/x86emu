package vm

import "encoding/binary"

func mov(v *VM, op byte) {
	switch op {
	case 0x00:
		fallthrough
	case 0x01:
		fallthrough
	case 0x02:
		fallthrough
	case 0x03:
		// do something
	case 0x88:
		fallthrough
	case 0x89:
		fallthrough
	case 0x8a:
		fallthrough
	case 0x8b:
		d := (0b00000010 & op) >> 2
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
				case AX:
					fallthrough
				case CX:
					fallthrough
				case DX:
					fallthrough
				case BX:
					fallthrough
				case SP:
					fallthrough
				case BP:
					fallthrough
				case SI:
					fallthrough
				case DI:
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
					v.Data[eabase(v, (uint16)(rm))] = (byte)(v.CPU.GR[(int)(1<<3|reg)] >> 8)
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))]) << 8)
				}
			case AX:
				fallthrough
			case CX:
				fallthrough
			case DX:
				fallthrough
			case BX:
				fallthrough
			case SP:
				fallthrough
			case BP:
				fallthrough
			case SI:
				fallthrough
			case DI:
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
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(reg)] & 0x00ff)
					} else {
						v.Data[eabase(v, (uint16)(rm))+(uint16)(^disp+1)] = (byte)(v.CPU.GR[(int)(reg)] & 0x00ff)
					}
				} else {
					if disp < 0x80 {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(^disp+1)])
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
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8)
					} else {
						v.Data[eabase(v, (uint16)(rm))+(uint16)(^disp+1)] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8)
					}
				} else {
					if disp < 0x80 {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) << 8)
					} else {
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(^disp+1)]) << 8)
					}
				}
			case AX:
				fallthrough
			case CX:
				fallthrough
			case DX:
				fallthrough
			case BX:
				fallthrough
			case SP:
				fallthrough
			case BP:
				fallthrough
			case SI:
				fallthrough
			case DI:
				if d == 0 {
					if disp < 0x80 {
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp+1)] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0xff00 >> 8)
					} else {
						v.Data[eabase(v, (uint16)(rm))+(uint16)(^disp+1)] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(^disp+1+1)] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0xff00 >> 8)
					}
				} else {
					if disp < 0x80 {
						v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)+(uint16)(disp)):])
					} else {
						v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)+(uint16)(^disp+1)):])
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
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp))] = (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(v.Data[eabase(v, (uint16)(rm)+(uint16)(disp))])
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp))] = (byte)(v.CPU.GR[(int)(1<<3|reg)] & 0xff00 >> 8)
				} else {
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm)+(uint16)(disp))]) << 8)
				}
			case AX:
				fallthrough
			case CX:
				fallthrough
			case DX:
				fallthrough
			case BX:
				fallthrough
			case SP:
				fallthrough
			case BP:
				fallthrough
			case SI:
				fallthrough
			case DI:
				if d == 0 {
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp))] = (byte)(v.CPU.GR[(int)(w<<3|reg)] & 0xff00 >> 8)
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp+1))] = (byte)(v.CPU.GR[(int)(w<<3|reg)&0x00ff])
				} else {
					v.CPU.GR[(int)(w<<3|reg)] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)+(uint16)(disp)):])
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
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | v.CPU.GR[(int)(1<<3|rm)]&0x00ff
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | v.CPU.GR[(int)(1<<3|rm)]&0xff00>>8
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
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | v.CPU.GR[(int)(1<<3|rm)]&0x00ff<<8
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff | v.CPU.GR[(int)(1<<3|rm)]&0xff00
				}
			case AX:
				fallthrough
			case CX:
				fallthrough
			case DX:
				fallthrough
			case BX:
				fallthrough
			case SP:
				fallthrough
			case BP:
				fallthrough
			case SI:
				fallthrough
			case DI:
				v.CPU.GR[(int)(w<<3|reg)] = v.CPU.GR[(int)(w<<3|reg)]
			}
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
		return
	}
}
