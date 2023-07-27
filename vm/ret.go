package vm

import "encoding/binary"

func ret(v *VM, op byte) {
	switch op {
	case 0xc2:
		next := binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:])
		v.CPU.GR[SP] += uint16(v.fetch()) | uint16(v.fetch())<<8 + 2
		v.IP = uint32(next)
	case 0xc3:
		v.IP = uint32(binary.LittleEndian.Uint16(v.Data[v.CPU.GR[SP]:]))
		v.CPU.GR[SP] += 2
	}
}
