package vm

import "encoding/binary"

func push(v *VM, op byte) {
	switch op {
	case 0x06:
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[ES] >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[ES] & 0x00ff)
	case 0x0e:
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[CS] >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[CS] & 0x00ff)
	case 0x16:
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[SS] >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[SS] & 0x00ff)
	case 0x1e:
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[DS] >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.SR[DS] & 0x00ff)
	case 0x50:
		fallthrough
	case 0x51:
		fallthrough
	case 0x52:
		fallthrough
	case 0x53:
		fallthrough
	case 0x54:
		fallthrough
	case 0x55:
		fallthrough
	case 0x56:
		fallthrough
	case 0x57:
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.GR[(int)(1<<3|op-0x50)] >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = (byte)(v.CPU.GR[(int)(1<<3|op-0x50)] & 0x00ff)
	}
}

func grp5push(v *VM, mod, rm byte) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			data := uint32(binary.LittleEndian.Uint16(v.Data[disp:]))
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = byte(data & 0x0000ff00 >> 8)
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = byte(data & 0x000000ff)
			return
		}
		data := uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]))
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(data & 0x0000ff00 >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(data & 0x000000ff)
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			data := uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+uint16(disp):]))
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = byte(data & 0x0000ff00 >> 8)
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = byte(data & 0x000000ff)
		} else {
			data := uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-uint16(^disp+1):]))
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = byte(data & 0x0000ff00 >> 8)
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = byte(data & 0x000000ff)
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		data := uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]))
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(data & 0x0000ff00 >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(data & 0x000000ff)
	case 0b11:
		data := uint32(v.CPU.GR[int(1<<3|rm)])
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(data & 0x0000ff00 >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(data & 0x000000ff)
	}
}
