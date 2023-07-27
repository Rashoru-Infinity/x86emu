package vm

import "encoding/binary"

func sbb(v *VM, op byte) {
	switch op {
	case 0x18:
		fallthrough
	case 0x19:
		fallthrough
	case 0x1a:
		fallthrough
	case 0x1b:
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
					if d == 0 {
						val := v.Data[disp] - (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) <= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[disp]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg)])&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[disp]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)])&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val < 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.Data[disp] = val
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[disp]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) >= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[disp]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg)])&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) > v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[disp]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)])&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					if d == 0 {
						val := v.Data[disp] - (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[disp] <= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[disp] < (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val < 0x80) || (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.Data[disp] = val
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8) - v.Data[disp]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) >= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[disp] >= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[disp]) > v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[disp] > (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(val) << 8)
					}
				default:
					if d == 0 {
						val := binary.LittleEndian.Uint16(v.Data[disp:]) - v.CPU.GR[(int)(w<<3|reg)]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[disp:]) <= v.CPU.GR[(int)(1<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[disp:])&0x0f <= v.CPU.GR[(int)(w<<3|reg)]&0x0f
						} else {
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[disp:]) < v.CPU.GR[(int)(1<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[disp:])&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[disp:]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[disp:]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.Data[disp] = (byte)(val & 0x00ff)
						v.Data[disp+1] = (byte)(val >> 8)
					} else {
						val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[disp:])
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[disp:]) >= v.CPU.GR[(int)(1<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[disp:])&0x0f >= v.CPU.GR[(int)(w<<3|reg)]&0x0f
						} else {
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[disp:]) > v.CPU.GR[(int)(1<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[disp:])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[disp:]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[disp:]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
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
				if d == 0 {
					val := v.Data[eabase(v, (uint16)(rm))] - (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) <= v.CPU.GR[(int)(1<<3|reg)]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) < v.CPU.GR[(int)(1<<3|reg)]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.Data[eabase(v, (uint16)(rm))] = val
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[eabase(v, (uint16)(rm))]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) >= v.CPU.GR[(int)(1<<3|reg)]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) > v.CPU.GR[(int)(1<<3|reg)]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					val := v.Data[eabase(v, (uint16)(rm))] - (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.Data[eabase(v, (uint16)(rm))] = val
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8) - v.Data[eabase(v, (uint16)(rm))]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) >= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) > v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(val) << 8)
				}
			default:
				if d == 0 {
					val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) - v.CPU.GR[(int)(w<<3|reg)]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) <= v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])&0x0f <= v.CPU.GR[(int)(w<<3|reg)]&0x0f
					} else {
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) < v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					v.Data[eabase(v, (uint16)(rm))] = (byte)(val & 0x00ff)
					v.Data[eabase(v, (uint16)(rm))+1] = (byte)(val >> 8)
				} else {
					val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) >= v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])&0x0f >= v.CPU.GR[(int)(w<<3|reg)]&0x0f
					} else {
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) > v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
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
						val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] - (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) <= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] - (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) <= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&+0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = val
					}
				} else {
					if disp < 0x80 {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) >= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) >= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) > v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val >= 0x80)
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
						val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] - (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] - (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = val
					}
				} else {
					if disp < 0x80 {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) >= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(val) << 8)
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8) - v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) >= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						} else {
							v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]) > v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(val) << 8)
					}
				}
			default:
				if d == 0 {
					if disp < 0x80 {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(disp):]) - v.CPU.GR[(int)(w<<3|reg)]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) <= v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f <= v.CPU.GR[(int)(w<<3|reg)]&0x0f
						} else {
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(val & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)+1] = (byte)(val & 0xff00 >> 8)
					} else {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) - v.CPU.GR[(int)(w<<3|reg)]
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) <= v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])&0x0f <= v.CPU.GR[(int)(w<<3|reg)]&0x0f
						} else {
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) < v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)] = (byte)(val & 0x00ff)
						v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)+1] = (byte)(val & 0xff00 >> 8)
					}
				} else {
					if disp < 0x80 {
						val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f >= v.CPU.GR[(int)(w<<3|reg)]&0x0f
						} else {
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) > v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
						v.CPU.GR[(int)(w<<3|reg)] = val
					} else {
						val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) >= v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])&0x0f >= v.CPU.GR[(int)(w<<3|reg)]&0x0f
						} else {
							v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) > v.CPU.GR[(int)(w<<3|reg)]
							v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
						}
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
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
				if d == 0 {
					val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] - (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) <= v.CPU.GR[(int)(1<<3|reg)&0x00ff]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) < v.CPU.GR[(int)(1<<3|reg)&0x00ff]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = val
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff]) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) >= v.CPU.GR[(int)(1<<3|reg)&0x00ff]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg)&0x00ff]
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (uint16)(val)
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] - (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)] & 0xff00 >> 8)
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) >= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f >= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					} else {
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
					v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | ((uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) << 8)
				}
			default:
				if d == 0 {
					val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) - v.CPU.GR[(int)(w<<3|reg)]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) <= v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f <= (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])&0x0f
					} else {
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f < (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] = (byte)(val >> 8)
					v.Data[eabase(v, (uint16)(rm)+(uint16)(disp)+1)] = (byte)(val & 0x00ff)
				} else {
					val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f >= (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])&0x0f
					} else {
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) > v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f > (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
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
						val := v.CPU.GR[(int)(1<<3|rm)]&0x00ff - v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff <= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff&0x0f <= v.CPU.GR[(int)(1<<3|reg)]&0x00ff&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff&0x0f < v.CPU.GR[(int)(1<<3|reg)]&0x00ff&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | (val & 0x00ff)
					} else {
						val := v.CPU.GR[(int)(1<<3|reg)]&0x00ff - v.CPU.GR[(int)(1<<3|rm)]&0x00ff
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff <= v.CPU.GR[(int)(1<<3|rm)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff&0x0f <= v.CPU.GR[(int)(1<<3|rm)]&0x00ff&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff < v.CPU.GR[(int)(1<<3|rm)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff&0x0f < v.CPU.GR[(int)(1<<3|rm)]&0x00ff&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
						v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|reg)]&0xff00 | (val & 0x00ff)
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					if d == 0 {
						val := v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 - v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 <= v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f <= v.CPU.GR[(int)(1<<3|reg)]&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f < v.CPU.GR[(int)(1<<3|reg)]&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
						v.CPU.GR[(int)(1<<3|rm&0x03)] = v.CPU.GR[(int)(1<<3|rm&0x03)]&0x00ff | val<<8
					} else {
						val := v.CPU.GR[(int)(1<<3|reg)]&0x00ff - v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff <= v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]&0x0f <= v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff < v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]&0x0f < v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
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
					if d == 0 {
						val := v.CPU.GR[(int)(1<<3|rm)]&0x00ff - v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff <= v.CPU.GR[(int)(1<<3|reg&0x03&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm)]&0x0f <= v.CPU.GR[(int)(1<<3|reg)]>>8&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm)]&0x0f < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
						v.CPU.GR[(int)(1<<3|rm)] = v.CPU.GR[(int)(1<<3|rm)]&0xff00 | val&0x00ff
					} else {
						val := v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 - v.CPU.GR[(int)(1<<3|rm)]&0x00ff
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 <= v.CPU.GR[(int)(1<<3|rm)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f <= v.CPU.GR[(int)(1<<3|rm)]&0x00ff&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < v.CPU.GR[(int)(1<<3|rm)]&0x00ff
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f < v.CPU.GR[(int)(1<<3|rm)]&0x00ff&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | val&0x00ff<<8
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					if d == 0 {
						val := v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 - v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[CF] = val > 0xff
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
						v.CPU.GR[(int)(1<<3|rm&0x03)] = v.CPU.GR[(int)(1<<3|rm&0x03)]&0x00ff | val&0x00ff<<8
					} else {
						val := v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 - v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
						if v.CPU.FR[CF] {
							val--
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 <= v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f <= v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f
						} else {
							v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
							v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f < v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f
						}
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[CF] = val > 0xff
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
						v.CPU.GR[(int)(1<<3|reg&0x03)] = v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff | val&0x00ff<<8
					}
				}
			default:
				if d == 0 {
					val := v.CPU.GR[(int)(w<<3|rm)] - v.CPU.GR[(int)(w<<3|reg)]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = v.CPU.GR[(int)(w<<3|rm)] <= v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = v.CPU.GR[(int)(w<<3|rm)]&0x0f <= v.CPU.GR[(int)(w<<3|reg)]&0x0f
					} else {
						v.CPU.FR[CF] = v.CPU.GR[(int)(w<<3|rm)] < v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = v.CPU.GR[(int)(w<<3|rm)]&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (v.CPU.GR[(int)(w<<3|rm)] < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (v.CPU.GR[(int)(w<<3|rm)] >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					v.CPU.GR[(int)(w<<3|rm)] = val
				} else {
					val := v.CPU.GR[(int)(w<<3|reg)] - v.CPU.GR[(int)(w<<3|rm)]
					if v.CPU.FR[CF] {
						val--
						v.CPU.FR[CF] = v.CPU.GR[(int)(w<<3|reg)] <= v.CPU.GR[(int)(w<<3|rm)]
						v.CPU.FR[AF] = v.CPU.GR[(int)(w<<3|reg)]&0x0f <= v.CPU.GR[(int)(w<<3|rm)]&0x0f
					} else {
						v.CPU.FR[CF] = v.CPU.GR[(int)(w<<3|reg)] < v.CPU.GR[(int)(w<<3|rm)]
						v.CPU.FR[AF] = v.CPU.GR[(int)(w<<3|reg)]&0x0f < v.CPU.GR[(int)(w<<3|rm)]&0x0f
					}
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && v.CPU.GR[(int)(w<<3|rm)] >= 0x8000 && val >= 0x8000) || (v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && v.CPU.GR[(int)(w<<3|rm)] < 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					v.CPU.GR[(int)(w<<3|reg)] = val
				}
			}
		}
	case 0x1c:
		fallthrough
	case 0x1d:
		w := 0b00000100 & op >> 2
		if w == 1 {
			data := (uint16)(v.fetch()) | (uint16)(v.fetch())<<8
			val := v.CPU.GR[AX] - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = v.CPU.GR[AX] <= data
				v.CPU.FR[AF] = v.CPU.GR[AX]&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = v.CPU.GR[AX] < data
				v.CPU.FR[AF] = v.CPU.GR[AX]&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (v.CPU.GR[AX] < 0x8000 && data >= 0x8000 && val >= 0x8000) || (v.CPU.GR[AX] >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[AX] = val
		} else {
			data := (uint16)(v.fetch())
			val := v.CPU.GR[AX]&0x00ff - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[AF] = v.CPU.GR[AX]&0x00ff&0x0f <= data&0x0f
			} else {
				v.CPU.FR[AF] = v.CPU.GR[AX]&0x00ff&0x0f < data&0x0f
			}
			v.CPU.FR[CF] = v.CPU.GR[AX]&0x00ff < data
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = (v.CPU.GR[AX]&0x00ff < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.CPU.GR[AX]&0x00ff >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[AX] = v.CPU.GR[AX]&0xff00 | val&0x00ff
		}
	}
}

func grp1sbb(v *VM, s, w, mod, rm uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if s == 0 && w == 1 {
				data := uint16(v.fetch()) | uint16(v.fetch())<<8
				dst := binary.LittleEndian.Uint16(v.Data[disp:])
				val := dst - data
				if v.CPU.FR[CF] {
					val--
					v.CPU.FR[AF] = dst&0x0f <= data&0x0f
					v.CPU.FR[CF] = dst <= data
				} else {
					v.CPU.FR[AF] = dst&0x0f < data&0x0f
					v.CPU.FR[CF] = dst < data
				}
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
				v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
				v.Data[disp] = byte(val & 0x00ff)
				v.Data[disp+1] = byte(val >> 8)
				return
			} else if s == 1 && w == 1 {
				data := uint16(v.fetch())
				if data >= 0x80 {
					data = 0xff00 | data
				}
				dst := binary.LittleEndian.Uint16(v.Data[disp:])
				val := dst - data
				if v.CPU.FR[CF] {
					val--
					v.CPU.FR[AF] = dst&0x0f <= data&0x0f
					v.CPU.FR[CF] = dst <= data
				} else {
					v.CPU.FR[AF] = dst&0x0f < data&0x0f
					v.CPU.FR[CF] = dst < data
				}
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
				v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
				v.Data[disp] = byte(val & 0x00ff)
				v.Data[disp+1] = byte(val >> 8)
				return
			}
			data := uint16(v.fetch())
			val := uint16(v.Data[disp]) - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[AF] = v.Data[disp]&0x0f <= byte(data&0x0f)
				v.CPU.FR[CF] = uint16(v.Data[disp]) <= data
			} else {
				v.CPU.FR[AF] = v.Data[disp]&0x0f < byte(data&0x0f)
				v.CPU.FR[CF] = uint16(v.Data[disp]) < data
			}
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = (v.Data[disp] < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.Data[disp] >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
			v.Data[disp] = byte(val)
			return
		}
		if s == 0 && w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) <= data
				v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) < data
				v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+1] = byte(val >> 8)
			return
		} else if s == 1 && w == 1 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			dst := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])
			val := dst - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = dst <= data
				v.CPU.FR[AF] = dst&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = dst < data
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+1] = byte(val >> 8)
			return
		}
		data := v.fetch()
		val := v.Data[eabase(v, uint16(rm))] - data
		if v.CPU.FR[CF] {
			val--
			v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))]&0x0f <= data&0x0f
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))] <= data
		} else {
			v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))]&0x0f < data&0x0f
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))] < data
		}
		v.CPU.FR[ZF] = val == 0
		v.CPU.FR[SF] = val >= 0x80
		v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))] < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.Data[eabase(v, uint16(rm))] >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
		v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		v.Data[eabase(v, uint16(rm))] = byte(val & 0x00ff)
		return
	case 0b01:
		disp := uint16(v.fetch())
		if s == 0 && w == 1 && disp < 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) <= data
				v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < data
				v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+disp+1] = byte(val >> 8)
			return
		} else if s == 0 && w == 1 && disp >= 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			dst := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1):])
			val := dst - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = dst <= data
				v.CPU.FR[AF] = dst&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = dst < data
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))-(^disp+1)] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))-(^disp+1)+1] = byte(val >> 8)
			return
		} else if s == 1 && w == 1 && disp < 0x80 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			dst := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])
			val := dst - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = dst <= data
				v.CPU.FR[AF] = dst&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = dst < data
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+disp+1] = byte(val >> 8)
			return
		} else if s == 1 && w == 1 && disp >= 0x80 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			dst := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1):])
			val := dst - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = dst <= data
				v.CPU.FR[AF] = dst&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = dst < data
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))-(^disp+1)] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))-(^disp+1)+1] = byte(val >> 8)
			return
		}
		data := uint16(v.fetch())
		if disp < 0x80 {
			val := uint16(v.Data[eabase(v, uint16(rm))+disp]) - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))+disp]) <= data
				v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))+disp])&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))+disp]) < data
				v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))+disp])&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+disp] < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.Data[eabase(v, uint16(rm))+disp] >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
		} else if disp >= 0x80 {
			val := uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]) - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]) <= data
				v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]&0x0f) <= data&0x0f
			} else {
				v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]) < data
				v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]&0x0f) < data&0x0f
			}
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff] < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff] >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff] = byte(val & 0x00ff)
		}
		return
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if s == 0 && w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			dst := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])
			val := dst - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = dst <= data
				v.CPU.FR[AF] = dst&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = dst < data
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+disp+1] = byte(val >> 8)
			return
		} else if s == 1 && w == 1 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			dst := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])
			val := dst - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = dst <= data
				v.CPU.FR[AF] = dst&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = dst < data
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
			v.Data[eabase(v, uint16(rm))+disp+1] = byte(val >> 8)
			return
		}
		data := uint16(v.fetch())
		val := uint16(v.Data[eabase(v, uint16(rm))+disp]) - data
		if v.CPU.FR[CF] {
			val--
			v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))+disp]) <= data
			v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))+disp])&0x0f <= data&0x0f
		} else {
			v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))+disp]) < data
			v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))+disp])&0x0f < data&0x0f
		}
		v.CPU.FR[ZF] = val&0x00ff == 0
		v.CPU.FR[SF] = val&0x00ff >= 0x80
		v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+disp] < 0x80 && data >= 0x80 && val >= 0x80) || (v.Data[eabase(v, uint16(rm))+disp] >= 0x80 && data < 0x80 && val < 0x80)
		v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		v.Data[eabase(v, uint16(rm))+disp] = byte(val & 0x00ff)
		return
	case 0b11:
		if s == 0 && w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := v.CPU.GR[int(w<<3|rm)] - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)] <= data
				v.CPU.FR[AF] = v.CPU.GR[int(w<<3|rm)]&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)] < data
				v.CPU.FR[AF] = v.CPU.GR[int(w<<3|rm)]&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)] < 0x8000 && data >= 0x8000 && val >= 0x8000) || (v.CPU.GR[int(w<<3|rm)] >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[int(w<<3|rm)] = val
			return
		} else if s == 1 && w == 1 {
			data := uint16(v.fetch())
			val := v.CPU.GR[int(w<<3|rm)] - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)] <= data
				v.CPU.FR[AF] = v.CPU.GR[int(w<<3|rm)]&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)] < data
				v.CPU.FR[AF] = v.CPU.GR[int(w<<3|rm)]&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)] < 0x8000 && data >= 0x8000 && val >= 0x8000) || (v.CPU.GR[int(w<<3|rm)] >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[int(w<<3|rm)] = val
			return
		}
		switch rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			data := uint16(v.fetch())
			val := v.CPU.GR[int(1<<3|rm)]&0x00ff - data
			if v.CPU.FR[CF] {
				val--
				v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm)]&0x00ff&0x0f <= data&0x0f
				v.CPU.FR[CF] = v.CPU.GR[int(1<<3|rm)]&0x00ff <= data
			} else {
				v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm)]&0x00ff&0x0f < data&0x0f
				v.CPU.FR[CF] = v.CPU.GR[int(1<<3|rm)]&0x00ff < data
			}
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm)]&0x00ff < 0x80 && data >= 0x80 && val >= 0x80) || (v.CPU.GR[int(1<<3|rm)]&0x00ff >= 0x80 && data < 0x80 && val < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | val&0x00ff
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			data := uint16(v.fetch())
			val := v.CPU.GR[int(1<<3|rm&0x03)]>>8 - data
			if v.CPU.FR[CF] {
				val++
				v.CPU.FR[CF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8 <= data
				v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8&0x0f <= data&0x0f
			} else {
				v.CPU.FR[CF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8 < data
				v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8&0x0f < data&0x0f
			}
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm&0x03)]>>8 < 0x80 && data >= 0x80 && val >= 0x80) || (v.CPU.GR[int(1<<3|rm&0x03)]>>8 >= 0x80 && data < 0x80 && val < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | val&0x00ff<<8
		}
	}
}
