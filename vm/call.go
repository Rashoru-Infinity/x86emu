package vm

import "encoding/binary"

func call(v *VM) {
	data := uint16(v.fetch()) | uint16(v.fetch())<<8
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(v.IP & 0x0000ff00 >> 8)
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(v.IP & 0x000000ff)
	if data < 0x8000 {
		v.IP += uint32(data)
	} else {
		v.IP -= uint32(^data + 1)
	}
}

func grp5call(v *VM, mod, rm byte) {
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(v.IP & 0x0000ff00 >> 8)
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(v.IP & 0x000000ff)
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			v.IP = uint32(binary.LittleEndian.Uint16(v.Data[disp:]))
			return
		}
		v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm)):]))
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+uint16(disp):]))
		} else {
			v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))-uint16(^disp+1):]))
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		v.IP = uint32(binary.LittleEndian.Uint16(v.Data[eabase(v, uint16(rm))+disp:]))
	case 0b11:
		v.IP = uint32(v.CPU.GR[int(1<<3|rm)])
	}
}
