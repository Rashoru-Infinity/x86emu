package vm

func cbw(v *VM) {
	if v.CPU.GR[AX]&0x0080 != 0 {
		v.CPU.GR[AX] = 0xff00 | v.CPU.GR[AX]&0x00ff
	} else {
		v.CPU.GR[AX] = v.CPU.GR[AX] & 0x00ff
	}
}

func cwd(v *VM) {
	if v.CPU.GR[AX]&0x8000 != 0 {
		v.CPU.GR[DX] = 0xffff
	} else {
		v.CPU.GR[DX] = 0x0000
	}
}
