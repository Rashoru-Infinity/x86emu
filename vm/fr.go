package vm

func checkPF(val uint8) bool {
	bit := 0
	for mask := 0b00000001; mask <= 0b10000000; mask <<= 1 {
		bit += mask & (int)(val) / (int)(mask)
	}
	return bit%2 == 0
}

func sahf(v *VM) {
	v.CPU.FR[SF] = v.CPU.GR[AX]>>8&0b10000000 == 0b10000000
	v.CPU.FR[ZF] = v.CPU.GR[AX]>>8&0b01000000 == 0b01000000
	v.CPU.FR[AF] = v.CPU.GR[AX]>>8&0b00010000 == 0b00010000
	v.CPU.FR[PF] = v.CPU.GR[AX]>>8&0b00000100 == 0b00000100
	v.CPU.FR[CF] = v.CPU.GR[AX]>>8&0b00000001 == 0b00000001
}

func lahf(v *VM) {
	v.CPU.GR[AX] = 0
	if v.CPU.FR[SF] {
		v.CPU.GR[AX] |= 0b10000000
	}
	if v.CPU.FR[ZF] {
		v.CPU.GR[AX] |= 0b01000000
	}
	if v.CPU.FR[AF] {
		v.CPU.GR[AX] |= 0b00010000
	}
	if v.CPU.FR[PF] {
		v.CPU.GR[AX] |= 0b00000100
	}
	if v.CPU.FR[CF] {
		v.CPU.GR[AX] |= 0b00000001
	}
}

func cmc(v *VM) {
	v.CPU.FR[CF] = !v.CPU.FR[CF]
}

func clc(v *VM) {
	v.CPU.FR[CF] = false
}

func stc(v *VM) {
	v.CPU.FR[CF] = true
}

func cli(v *VM) {
	v.CPU.FR[IF] = false
}

func sti(v *VM) {
	v.CPU.FR[IF] = true
}

func cld(v *VM) {
	v.CPU.FR[DF] = false
}

func std(v *VM) {
	v.CPU.FR[DF] = true
}
