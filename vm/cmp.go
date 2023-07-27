package vm

import (
	"encoding/binary"
)

func cmp(v *VM, op byte) {
	switch op {
	case 0x38:
		fallthrough
	case 0x39:
		fallthrough
	case 0x3a:
		fallthrough
	case 0x3b:
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
						v.CPU.FR[CF] = (uint16)(v.Data[disp]) < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.Data[disp]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)])&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val < 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[disp]
						v.CPU.FR[CF] = (uint16)(v.Data[disp]) > v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.Data[disp]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)])&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
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
						v.CPU.FR[CF] = (uint16)(v.Data[disp]) < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[disp] < (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val < 0x80) || (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8) - v.Data[disp]
						v.CPU.FR[CF] = (uint16)(v.Data[disp]) > v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[disp] > (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && (uint16)(v.Data[disp]) >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && (uint16)(v.Data[disp]) < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
					}
				default:
					if d == 0 {
						val := binary.LittleEndian.Uint16(v.Data[disp:]) - v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[disp:]) < v.CPU.GR[(int)(1<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[disp:])&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[disp:]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[disp:]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					} else {
						val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[disp:])
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[disp:]) > v.CPU.GR[(int)(1<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[disp:])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[disp:]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[disp:]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
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
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) < v.CPU.GR[(int)(1<<3|reg)]
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[eabase(v, (uint16)(rm))]
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) > v.CPU.GR[(int)(1<<3|reg)]
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
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
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8) - v.Data[eabase(v, (uint16)(rm))]
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))]) > v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]&0x00ff < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
				}
			default:
				if d == 0 {
					val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) - v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) < v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				} else {
					val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])
					v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) > v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm)):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
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
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] - (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]) < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&+0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
					}
				} else {
					if disp < 0x80 {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff) - v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]) > v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]&0x00ff)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
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
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) <= v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f <= (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
					} else {
						val := v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] - (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]) < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg&0x03)]>>8)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF(val)
					}
				} else {
					if disp < 0x80 {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
					} else {
						val := (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8) - v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]
						v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]) < v.CPU.GR[(int)(1<<3|reg)]>>8
						v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8)&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x80
						v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 < 0x80 && val >= 0x80)
						v.CPU.FR[PF] = checkPF(val)
					}
				}
			default:
				if d == 0 {
					if disp < 0x80 {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(disp):]) - v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					} else {
						val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) - v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) < v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))&0x0f-(uint16)(^disp+1)&0x00ff:]) < v.CPU.GR[(int)(w<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					}
				} else {
					if disp < 0x80 {
						val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) > v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
					} else {
						val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:])
						v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) > v.CPU.GR[(int)(w<<3|reg)]
						v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:])&0x0f > v.CPU.GR[(int)(w<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val == 0
						v.CPU.FR[SF] = val >= 0x8000
						v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))-(uint16)(^disp+1)&0x00ff:]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
						v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
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
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) < v.CPU.GR[(int)(1<<3|reg)&0x00ff]
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff]) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg)&0x00ff]
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)&0x00ff])&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)&0x00ff] < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
				}
			case AH:
				fallthrough
			case CH:
				fallthrough
			case DH:
				fallthrough
			case BH:
				if d == 0 {
					val := v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] - (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8)
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) < v.CPU.GR[(int)(1<<3|reg)]>>8
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f < (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8)&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 >= 0x80 && val >= 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 < 0x80 && val < 0x80)
					v.CPU.FR[PF] = checkPF(val)
				} else {
					val := (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8) - v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]
					v.CPU.FR[CF] = (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]) > v.CPU.GR[(int)(1<<3|reg)]>>8
					v.CPU.FR[AF] = v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)]&0x0f > (byte)(v.CPU.GR[(int)(1<<3|reg)]>>8)&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x80
					v.CPU.FR[OF] = (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] < 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 >= 0x80 && val < 0x80) || (v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)] >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 < 0x80 && val >= 0x80)
					v.CPU.FR[PF] = checkPF(val)
				}
			default:
				if d == 0 {
					val := binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) - v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f < (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				} else {
					val := v.CPU.GR[(int)(w<<3|reg)] - binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])
					v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) > v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):])&0x0f > (uint16)(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp)])&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val < 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, (uint16)(rm))+(uint16)(disp):]) >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val >= 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
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
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm)]&0x0f < v.CPU.GR[(int)(1<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
					} else {
						val := v.CPU.GR[(int)(1<<3|reg)]&0x00ff - v.CPU.GR[(int)(1<<3|rm)]&0x00ff
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff < v.CPU.GR[(int)(1<<3|rm)]&0x00ff
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]&0x0f < v.CPU.GR[(int)(1<<3|rm)]&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
					}
				case AH:
					fallthrough
				case CH:
					fallthrough
				case DH:
					fallthrough
				case BH:
					if d == 0 {
						val := v.CPU.GR[(int)(1<<3|rm)]>>8 - v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm)]>>8 < v.CPU.GR[(int)(1<<3|reg)]&0x00ff
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm)]>>8&0x0f < v.CPU.GR[(int)(1<<3|reg)]&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
					} else {
						val := v.CPU.GR[(int)(1<<3|reg)]&0x00ff - v.CPU.GR[(int)(1<<3|rm)]>>8
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg)]&0x00ff < v.CPU.GR[(int)(1<<3|rm)]>>8
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]&0x0f < v.CPU.GR[(int)(1<<3|rm)]>>8&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|rm)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|rm)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
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
						val := v.CPU.GR[(int)(1<<3|rm)]&0x00ff - v.CPU.GR[(int)(1<<3|reg)]>>8
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm)]&0x00ff < v.CPU.GR[(int)(1<<3|reg)]>>8
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm)]>>8&0x0f < v.CPU.GR[(int)(1<<3|reg)]>>8&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && v.CPU.GR[(int)(1<<3|reg)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
					} else {
						val := v.CPU.GR[(int)(1<<3|reg)]>>8 - v.CPU.GR[(int)(1<<3|rm)]&0x00ff
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg)]>>8 < v.CPU.GR[(int)(1<<3|rm)]&0x00ff
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg)]>>8&0x0f < v.CPU.GR[(int)(1<<3|rm)]&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|rm)]&0x00ff < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
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
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f < v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[CF] = val > 0xff
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
					} else {
						val := v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 - v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
						v.CPU.FR[CF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < v.CPU.GR[(int)(1<<3|rm&0x03)]>>8
						v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|reg&0x03)]>>8&0x0f < v.CPU.GR[(int)(1<<3|rm&0x03)]>>8&0x0f
						v.CPU.FR[ZF] = val&0x00ff == 0
						v.CPU.FR[SF] = val&0x00ff >= 0x80
						v.CPU.FR[CF] = val > 0xff
						v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 < 0x80 && v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 >= 0x80 && val >= 0x80) || (v.CPU.GR[(int)(1<<3|reg&0x03)]>>8 >= 0x80 && v.CPU.GR[(int)(1<<3|rm&0x03)]>>8 < 0x80 && val < 0x80)
						v.CPU.FR[PF] = checkPF((uint8)(val))
					}
				}
			default:
				if d == 0 {
					val := v.CPU.GR[(int)(w<<3|rm)] - v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[CF] = v.CPU.GR[(int)(w<<3|rm)] < v.CPU.GR[(int)(w<<3|reg)]
					v.CPU.FR[AF] = v.CPU.GR[(int)(w<<3|rm)]&0x0f < v.CPU.GR[(int)(w<<3|reg)]&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (v.CPU.GR[(int)(w<<3|rm)] < 0x8000 && v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && val >= 0x8000) || (v.CPU.GR[(int)(w<<3|rm)] >= 0x8000 && v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				} else {
					val := v.CPU.GR[(int)(w<<3|reg)] - v.CPU.GR[(int)(w<<3|rm)]
					v.CPU.FR[CF] = v.CPU.GR[(int)(w<<3|reg)] < v.CPU.GR[(int)(w<<3|rm)]
					v.CPU.FR[AF] = v.CPU.GR[(int)(w<<3|reg)]&0x0f < v.CPU.GR[(int)(w<<3|rm)]&0x0f
					v.CPU.FR[ZF] = val == 0
					v.CPU.FR[SF] = val >= 0x8000
					v.CPU.FR[OF] = (v.CPU.GR[(int)(w<<3|reg)] < 0x8000 && v.CPU.GR[(int)(w<<3|rm)] >= 0x8000 && val >= 0x8000) || (v.CPU.GR[(int)(w<<3|reg)] >= 0x8000 && v.CPU.GR[(int)(w<<3|rm)] < 0x8000 && val < 0x8000)
					v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
				}
			}
		}
	case 0x3c:
		fallthrough
	case 0x3d:
		w := 0b00000100 & op >> 2
		if w == 1 {
			data := (uint16)(v.fetch()) | (uint16)(v.fetch())<<8
			val := v.CPU.GR[AX] - data
			v.CPU.FR[CF] = v.CPU.GR[AX] < data
			v.CPU.FR[AF] = v.CPU.GR[AX]&0x0f < data&0x0f
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[OF] = (v.CPU.GR[AX] < 0x8000 && data >= 0x8000 && val >= 0x8000) || (v.CPU.GR[AX] >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		} else {
			data := (uint16)(v.fetch())
			val := v.CPU.GR[AX]&0x00ff - data
			v.CPU.FR[AF] = v.CPU.GR[AX]&0x0f < data&0x0f
			v.CPU.FR[CF] = val > 0xff
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[OF] = (v.CPU.GR[AX]&0x00ff < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.CPU.GR[AX]&0x00ff >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		}
	}
}

func grp1cmp(v *VM, s, w, mod, rm uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if s == 0 && w == 1 {
				data := uint16(v.fetch()) | uint16(v.fetch())<<8
				dst := binary.LittleEndian.Uint16(v.Data[disp:])
				val := dst - data
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = dst > 0xffff-data
				v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
				v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
				return
			} else if s == 1 && w == 1 {
				data := uint16(v.fetch())
				if data >= 0x80 {
					data = 0xff00 | data
				}
				dst := binary.LittleEndian.Uint16(v.Data[disp:])
				val := dst - data
				v.CPU.FR[ZF] = val == 0
				v.CPU.FR[SF] = val >= 0x8000
				v.CPU.FR[CF] = dst > 0xffff-data
				v.CPU.FR[OF] = (dst < 0x8000 && data >= 0x8000 && val >= 0x8000) || (dst >= 0x8000 && data < 0x8000 && val < 0x8000)
				v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
				v.CPU.FR[AF] = dst&0x0f < data&0x0f
				return
			}
			data := v.fetch()
			val := v.Data[disp] - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x80
			v.CPU.FR[CF] = v.Data[disp] < data
			v.CPU.FR[OF] = (v.Data[disp] < 0x80 && data >= 0x80 && val >= 0x80) || (v.Data[disp] >= 0x80 && data < 0x80 && val < 0x80)
			v.CPU.FR[PF] = checkPF(uint8(val & 0x00ff))
			v.CPU.FR[AF] = v.Data[disp]&0x0f < byte(data&0x0f)
			return
		}
		if s == 0 && w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])&0x0f < data&0x0f
			return
		} else if s == 1 && w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			if data >= 0x80 {
				data = 0xff00 | data
			}
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])&0x0f < data&0x0f
			return
		}
		data := v.fetch()
		val := v.Data[eabase(v, uint16(rm))] - data
		v.CPU.FR[ZF] = val == 0
		v.CPU.FR[SF] = val >= 0x80
		v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))] < data
		v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))] < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.Data[eabase(v, uint16(rm))] >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
		v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))]&0x0f < data&0x0f
		return
	case 0b01:
		disp := uint16(v.fetch())
		if s == 0 && w == 1 && disp < 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])&0x0f < data&0x0f
			return
		} else if s == 0 && w == 1 && disp >= 0x80 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1):]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:])&0x0f < data&0x0f
			return
		} else if s == 1 && w == 1 && disp < 0x80 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])&0x0f < data&0x0f
			return
		} else if s == 1 && w == 1 && disp >= 0x80 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff:])&0x0f < data&0x0f
			return
		}
		data := uint16(v.fetch())
		if disp < 0x80 {
			val := uint16(v.Data[eabase(v, uint16(rm))+disp]) - data
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))+disp]) < data
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+disp] < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.Data[eabase(v, uint16(rm))+disp] >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))+disp])&0x0f < data&0x0f
			return
		} else if disp >= 0x80 {
			val := uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]) - data
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]) < data
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff] < 0x80 && data >= 0x80 && val&0x00ff >= 0x80) || (v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff] >= 0x80 && data < 0x80 && val&0x00ff < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))-(^disp+1)&0x00ff]&0x0f) < data&0x0f
			return
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if s == 0 && w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])&0x0f < data&0x0f
			return
		} else if s == 1 && w == 0 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < data
			v.CPU.FR[OF] = (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) < 0x8000 && data >= 0x8000 && val >= 0x8000) || (binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]) >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:])&0x0f < data&0x0f
			return
		}
		data := uint16(v.fetch())
		val := uint16(v.Data[eabase(v, uint16(rm))+disp]) - data
		v.CPU.FR[ZF] = val&0x00ff == 0
		v.CPU.FR[SF] = val&0x00ff >= 0x80
		v.CPU.FR[CF] = uint16(v.Data[eabase(v, uint16(rm))+disp]) < data
		v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+disp] < 0x80 && data >= 0x80 && val >= 0x80) || (v.Data[eabase(v, uint16(rm))+disp] >= 0x80 && data < 0x80 && val < 0x80)
		v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
		v.CPU.FR[AF] = uint16(v.Data[eabase(v, uint16(rm))+disp])&0x0f < data&0x0f
		return
	case 0b11:
		if s == 0 && w == 1 {
			data := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := v.CPU.GR[int(w<<3|rm)] - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)] < data
			v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)] < 0x8000 && data >= 0x8000 && val >= 0x8000) || (v.CPU.GR[int(w<<3|rm)] >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = v.CPU.GR[int(w<<3|rm)]&0x0f < data&0x0f
			return
		} else if s == 1 && w == 1 {
			data := uint16(v.fetch())
			if data >= 0x80 {
				data = 0xff00 | data
			}
			val := v.CPU.GR[int(w<<3|rm)] - data
			v.CPU.FR[ZF] = val == 0
			v.CPU.FR[SF] = val >= 0x8000
			v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)] < data
			v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)] < 0x8000 && data >= 0x8000 && val >= 0x8000) || (v.CPU.GR[int(w<<3|rm)] >= 0x8000 && data < 0x8000 && val < 0x8000)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = v.CPU.GR[int(w<<3|rm)]&0x0f < data&0x0f
			return
		}
		data := uint16(v.fetch())
		switch 1<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			val := v.CPU.GR[int(1<<3|rm)]&0x00ff - data
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = v.CPU.GR[int(1<<3|rm)]&0x00ff < data
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm)]&0x00ff < 0x80 && data >= 0x80 && val >= 0x80) || (v.CPU.GR[int(1<<3|rm)]&0x00ff >= 0x80 && data < 0x80 && val < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm)]&0x0f < data&0x0f
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			val := v.CPU.GR[int(1<<3|rm&0x03)]>>8 - data
			v.CPU.FR[ZF] = val&0x00ff == 0
			v.CPU.FR[SF] = val&0x00ff >= 0x80
			v.CPU.FR[CF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8 < data
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm&0x03)]>>8 < 0x80 && data >= 0x80 && val >= 0x80) || (v.CPU.GR[int(1<<3|rm&0x03)]>>8 >= 0x80 && data < 0x80 && val < 0x80)
			v.CPU.FR[PF] = checkPF((uint8)(val & 0x00ff))
			v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8&0x0f < data&0x0f
		}
	}
}
