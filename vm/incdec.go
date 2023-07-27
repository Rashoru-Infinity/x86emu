package vm

import (
	"encoding/binary"
)

func inc(v *VM, op byte) {
	v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|op&0b00000111)] < 0x8000 && v.CPU.GR[(int)(1<<3|op&0b00000111)]+1 >= 0x8000)
	v.CPU.FR[SF] = v.CPU.GR[(int)(1<<3|op&0b00000111)]+1 >= 0x8000
	v.CPU.FR[ZF] = v.CPU.GR[(int)(1<<3|op&0b00000111)]+1 == 0
	v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|op&0b00000111)]&0x0f+1 > 0x0f
	v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[(int)(1<<3|op&0b00000111)] + 1))
	v.CPU.GR[(int)(1<<3|op&0b00000111)]++
}

func dec(v *VM, op byte) {
	v.CPU.FR[OF] = (v.CPU.GR[(int)(1<<3|op&0b00000111)] >= 0x8000 && v.CPU.GR[(int)(1<<3|op&0b00000111)]-1 < 0x8000)
	v.CPU.FR[SF] = v.CPU.GR[(int)(1<<3|op&0b00000111)]-1 >= 0x8000
	v.CPU.FR[ZF] = v.CPU.GR[(int)(1<<3|op&0b00000111)]-1 == 0
	v.CPU.FR[AF] = v.CPU.GR[(int)(1<<3|op&0b00000111)]&0x0f == 0
	v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[(int)(1<<3|op&0b00000111)] - 1))
	v.CPU.GR[(int)(1<<3|op&0b00000111)]--
}

func grp4inc(v *VM, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			v.CPU.FR[OF] = (v.Data[disp] < 0x80 && v.Data[disp]+1 >= 0x80)
			v.CPU.FR[SF] = v.Data[disp]+1 >= 0x80
			v.CPU.FR[ZF] = v.Data[disp]+1 == 0
			v.CPU.FR[AF] = v.Data[disp]&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(v.Data[disp] + 1))
			v.Data[disp]++
			return
		}
		v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))] < 0x80 && v.Data[eabase(v, uint16(rm))]+1 >= 0x80)
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))]+1 >= 0x80
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))]+1 == 0
		v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))]&0x0f+1 > 0x0f
		v.CPU.FR[PF] = checkPF((uint8)(v.Data[eabase(v, uint16(rm))] + 1))
		v.Data[eabase(v, uint16(rm))]++
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+uint16(disp)] < 0x80 && v.Data[eabase(v, uint16(rm))+uint16(disp)]+1 >= 0x80)
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]+1 >= 0x80
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]+1 == 0
			v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(v.Data[eabase(v, uint16(rm))+uint16(disp)] + 1))
			v.Data[eabase(v, uint16(rm))+uint16(disp)]++
		} else {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] < 0x80 && v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]+1 >= 0x80)
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]+1 >= 0x80
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]+1 == 0
			v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] + 1))
			v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]++
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		v.CPU.FR[OF] = (v.Data[disp] < 0x80 && v.Data[disp]+1 >= 0x80)
		v.CPU.FR[SF] = v.Data[disp]+1 >= 0x80
		v.CPU.FR[ZF] = v.Data[disp]+1 == 0
		v.CPU.FR[AF] = v.Data[disp]&0x0f+1 > 0x0f
		v.CPU.FR[PF] = checkPF((uint8)(v.Data[disp] + 1))
		v.Data[disp]++
		return
	case 0b11:
		switch rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm)]&0x00ff < 0x80 && v.CPU.GR[int(1<<3|rm)]&0x00ff+1 >= 0x80)
			v.CPU.FR[SF] = v.CPU.GR[int(1<<3|rm)]&0x00ff+1 >= 0x80
			v.CPU.FR[ZF] = v.CPU.GR[int(1<<3|rm)]&0x00ff+1 == 0
			v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm)]&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[int(1<<3|rm)]&0x00ff + 1))
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | ((v.CPU.GR[int(1<<3|rm)] + 1) & 0x00ff)
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm&0x03)]>>8 < 0x80 && v.CPU.GR[int(1<<3|rm&0x03)]>>8+1 >= 0x80)
			v.CPU.FR[SF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8+1 >= 0x80
			v.CPU.FR[ZF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8+1 == 0
			v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[int(1<<3|rm&0x03)]>>8 + 1))
			v.CPU.GR[int(1<<3|rm&0x03)] = (v.CPU.GR[int(1<<3|rm&0x03)]>>8+1)<<8 | (v.CPU.GR[int(1<<3|rm&0x03)] & 0x00ff)

		}
	}
}

func grp4dec(v *VM, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			v.CPU.FR[OF] = (v.Data[disp] >= 0x80 && v.Data[disp]-1 < 0x80)
			v.CPU.FR[SF] = v.Data[disp]-1 >= 0x80
			v.CPU.FR[ZF] = v.Data[disp]-1 == 0
			v.CPU.FR[AF] = v.Data[disp]&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(v.Data[disp] - 1))
			v.Data[disp]--
			return
		}
		v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))] >= 0x80 && v.Data[eabase(v, uint16(rm))]-1 < 0x80)
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))]-1 >= 0x80
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))]-1 == 0
		v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))]&0x0f == 0
		v.CPU.FR[PF] = checkPF((uint8)(v.Data[eabase(v, uint16(rm))] - 1))
		v.Data[eabase(v, uint16(rm))]--
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+uint16(disp)] >= 0x80 && v.Data[eabase(v, uint16(rm))+uint16(disp)]-1 < 0x80)
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]-1 >= 0x80
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]-1 == 0
			v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(v.Data[eabase(v, uint16(rm))+uint16(disp)] - 1))
			v.Data[eabase(v, uint16(rm))+uint16(disp)]--
		} else {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] >= 0x80 && v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]-1 < 0x80)
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]-1 >= 0x80
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]-1 == 0
			v.CPU.FR[AF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] - 1))
			v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]--
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		v.CPU.FR[OF] = (v.Data[disp] >= 0x80 && v.Data[disp]-1 < 0x80)
		v.CPU.FR[SF] = v.Data[disp]-1 >= 0x80
		v.CPU.FR[ZF] = v.Data[disp]-1 == 0
		v.CPU.FR[AF] = v.Data[disp]&0x0f == 0
		v.CPU.FR[PF] = checkPF((uint8)(v.Data[disp] - 1))
		v.Data[disp]--
	case 0b11:
		switch rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm)]&0x00ff >= 0x80 && v.CPU.GR[int(1<<3|rm)]&0x00ff-1 < 0x80)
			v.CPU.FR[SF] = v.CPU.GR[int(1<<3|rm)]&0x00ff-1 >= 0x80
			v.CPU.FR[ZF] = v.CPU.GR[int(1<<3|rm)]&0x00ff-1 == 0
			v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm)]&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[int(1<<3|rm)] - 1))
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | ((v.CPU.GR[int(1<<3|rm)] - 1) & 0x00ff)
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm&0x03)]>>8 >= 0x80 && v.CPU.GR[int(1<<3|rm&0x03)]>>8-1 < 0x80)
			v.CPU.FR[SF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8-1 >= 0x80
			v.CPU.FR[ZF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8-1 == 0
			v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm&0x03)]>>8&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[int(1<<3|rm&0x03)]>>8 - 1))
			v.CPU.GR[int(1<<3|rm&0x03)] = (v.CPU.GR[int(1<<3|rm&0x03)]>>8-1)<<8 | (v.CPU.GR[int(1<<3|rm&0x03)] & 0x00ff)
		}
	}
}

func grp5inc(v *VM, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[disp:])
			v.CPU.FR[OF] = (val < 0x8000 && val+1 >= 0x8000)
			v.CPU.FR[SF] = val+1 >= 0x8000
			v.CPU.FR[ZF] = val+1 == 0
			v.CPU.FR[AF] = val&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(val + 1))
			v.Data[disp] = byte((val + 1) & 0x00ff)
			v.Data[disp+1] = byte((val + 1) >> 8)
			return
		}
		val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])
		v.CPU.FR[OF] = val < 0x8000 && val+1 >= 0x8000
		v.CPU.FR[SF] = val+1 >= 0x8000
		v.CPU.FR[ZF] = val+1 == 0
		v.CPU.FR[AF] = val&0x0f+1 > 0x0f
		v.CPU.FR[PF] = checkPF((uint8)(val + 1))
		v.Data[eabase(v, uint16(rm))] = byte((val + 1) & 0x00ff)
		v.Data[eabase(v, uint16(rm))+1] = byte((val + 1) >> 8)
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+uint16(disp):])
			v.CPU.FR[OF] = val < 0x8000 && val+1 >= 0x8000
			v.CPU.FR[SF] = val+1 >= 0x8000
			v.CPU.FR[ZF] = val+1 == 0
			v.CPU.FR[AF] = val&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(val + 1))
			v.Data[eabase(v, uint16(rm))+uint16(disp)] = byte((val + 1) & 0x00ff)
			v.Data[eabase(v, uint16(rm))+uint16(disp)+1] = byte((val + 1) >> 8)
		} else {
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-uint16(^disp+1)&0x00ff:])
			v.CPU.FR[OF] = val < 0x8000 && val+1 >= 0x8000
			v.CPU.FR[SF] = val+1 >= 0x8000
			v.CPU.FR[ZF] = val+1 == 0
			v.CPU.FR[AF] = val&0x0f+1 > 0x0f
			v.CPU.FR[PF] = checkPF((uint8)(val + 1))
			v.Data[eabase(v, uint16(rm))-uint16(^disp+1)&0x00ff] = byte((val + 1) & 0x00ff)
			v.Data[eabase(v, uint16(rm))-uint16(^disp+1)&0x00ff+1] = byte((val + 1) >> 8)
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		val := binary.LittleEndian.Uint16(v.Data[disp:])
		v.CPU.FR[OF] = (val < 0x8000 && val+1 >= 0x8000)
		v.CPU.FR[SF] = val+1 >= 0x8000
		v.CPU.FR[ZF] = val+1 == 0
		v.CPU.FR[AF] = val&0x0f+1 > 0x0f
		v.CPU.FR[PF] = checkPF((uint8)(val + 1))
		v.Data[disp] = byte((val + 1) & 0x00ff)
		v.Data[disp+1] = byte((val + 1) & 0x00ff)
	case 0b11:
		v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm)] < 0x8000 && v.CPU.GR[int(1<<3|rm)]+1 >= 0x8000)
		v.CPU.FR[SF] = v.CPU.GR[int(1<<3|rm)]+1 >= 0x8000
		v.CPU.FR[ZF] = v.CPU.GR[int(1<<3|rm)]+1 == 0
		v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm)]&0x0f+1 > 0x0f
		v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[int(1<<3|rm)] + 1))
		v.CPU.GR[int(1<<3|rm)]++
	}
}

func grp5dec(v *VM, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			val := binary.LittleEndian.Uint16(v.Data[disp:])
			v.CPU.FR[OF] = (val >= 0x8000 && val-1 < 0x8000)
			v.CPU.FR[SF] = val-1 >= 0x8000
			v.CPU.FR[ZF] = val-1 == 0
			v.CPU.FR[AF] = val&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(val - 1))
			v.Data[disp] = byte((val - 1) & 0x00ff)
			v.Data[disp+1] = byte((val - 1) >> 8)
			return
		}
		val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):])
		v.CPU.FR[OF] = val >= 0x8000 && val-1 < 0x8000
		v.CPU.FR[SF] = val-1 >= 0x8000
		v.CPU.FR[ZF] = val-1 == 0
		v.CPU.FR[AF] = val&0x0f == 0
		v.CPU.FR[PF] = checkPF((uint8)(val - 1))
		v.Data[eabase(v, uint16(rm))] = byte((val - 1) & 0x00ff)
		v.Data[eabase(v, uint16(rm))+1] = byte((val - 1) >> 8)
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+uint16(disp):])
			v.CPU.FR[OF] = val >= 0x8000 && val-1 < 0x8000
			v.CPU.FR[SF] = val-1 >= 0x8000
			v.CPU.FR[ZF] = val-1 == 0
			v.CPU.FR[AF] = val&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(val - 1))
			v.Data[eabase(v, uint16(rm))+uint16(disp)] = byte((val - 1) & 0x00ff)
			v.Data[eabase(v, uint16(rm))+uint16(disp)+1] = byte((val - 1) >> 8)
		} else {
			val := binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-uint16(^disp+1):])
			v.CPU.FR[OF] = val >= 0x8000 && val-1 < 0x8000
			v.CPU.FR[SF] = val-1 >= 0x8000
			v.CPU.FR[ZF] = val-1 == 0
			v.CPU.FR[AF] = val&0x0f == 0
			v.CPU.FR[PF] = checkPF((uint8)(val - 1))
			v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = byte((val - 1) & 0x00ff)
			v.Data[eabase(v, uint16(rm))-uint16(^disp+1)+1] = byte((val - 1) >> 8)
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		val := binary.LittleEndian.Uint16(v.Data[disp:])
		v.CPU.FR[OF] = (val >= 0x8000 && val-1 < 0x8000)
		v.CPU.FR[SF] = val-1 >= 0x8000
		v.CPU.FR[ZF] = val-1 == 0
		v.CPU.FR[AF] = val&0x0f == 0
		v.CPU.FR[PF] = checkPF((uint8)(val - 1))
		v.Data[disp] = byte((val - 1) & 0x00ff)
		v.Data[disp+1] = byte((val - 1) & 0x00ff)
	case 0b11:
		v.CPU.FR[OF] = (v.CPU.GR[int(1<<3|rm)] >= 0x8000 && v.CPU.GR[int(1<<3|rm)]-1 < 0x8000)
		v.CPU.FR[SF] = v.CPU.GR[int(1<<3|rm)]-1 >= 0x8000
		v.CPU.FR[ZF] = v.CPU.GR[int(1<<3|rm)]-1 == 0
		v.CPU.FR[AF] = v.CPU.GR[int(1<<3|rm)]&0x0f == 0
		v.CPU.FR[PF] = checkPF((uint8)(v.CPU.GR[int(1<<3|rm)] - 1))
		v.CPU.GR[int(1<<3|rm)]--
	}
}
