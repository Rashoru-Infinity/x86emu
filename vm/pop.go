package vm

import (
	"encoding/binary"
)

func pop(v *VM, op byte) {
	switch op {
	case 0x07:
		v.CPU.SR[ES] = binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:])
	case 0x17:
		v.CPU.SR[SS] = binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:])
	case 0x1f:
		v.CPU.SR[DS] = binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:])
	case 0x58:
		fallthrough
	case 0x59:
		fallthrough
	case 0x5a:
		fallthrough
	case 0x5b:
		fallthrough
	case 0x5c:
		fallthrough
	case 0x5d:
		fallthrough
	case 0x5e:
		fallthrough
	case 0x5f:
		reg := op & 0b00000111
		v.CPU.GR[(int)(1<<3|reg)] = binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:])
	case 0x8f:
		data := v.fetch()
		mod := data & 0b11000000 >> 6
		rm := data & 00000111
		switch mod {
		case 0b00:
			if rm == 0b110 {
				disp := uint16(v.fetch()) | uint16(v.fetch())<<8
				v.Data[disp] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) & 0x00ff)
				v.Data[disp+1] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) >> 8)
				break
			}
			v.Data[eabase(v, uint16(rm))] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) & 0x00ff)
			v.Data[eabase(v, uint16(rm))] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) >> 8)
		case 0b01:
			disp := v.fetch()
			if disp < 0x80 {
				v.Data[eabase(v, uint16(rm))+uint16(disp)] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) & 0x00ff)
				v.Data[eabase(v, uint16(rm))+uint16(disp)+1] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) >> 8)
			} else {
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) & 0x00ff)
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)+1] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) >> 8)
			}
		case 0b10:
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			v.Data[disp] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) & 0x00ff)
			v.Data[disp+1] = byte(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]) >> 8)
		case 0b11:
			v.CPU.GR[int(1<<3|rm)] = binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:])
		}
	}
	v.CPU.GR[SP] += 2
}
