package vm

func xlat(v *VM) {
	v.CPU.GR[AX] = v.CPU.GR[AX]&0xff00 | uint16(v.Data[v.CPU.SR[DS]<<4+v.CPU.GR[BX]+v.CPU.GR[AX]&0x00ff])
}
