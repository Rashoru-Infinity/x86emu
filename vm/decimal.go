package vm

func daa(v *VM) {
	oldAL := v.CPU.GR[AX] & 0x00ff
	oldCF := v.CPU.FR[CF]
	v.CPU.FR[CF] = false
	if (v.CPU.GR[AX]&0x00ff)&0x0f > 9 || v.CPU.FR[AF] {
		v.CPU.FR[CF] = oldCF || (v.CPU.GR[AX]&0x00ff+6) > 0xff
		v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | ((v.CPU.GR[AX]&0x00ff + 6) & 0x00ff)
		v.CPU.FR[AF] = true
	} else {
		v.CPU.FR[AF] = false
	}
	if oldAL > 0x99 || oldCF {
		v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | ((v.CPU.GR[AX]&0x00ff + 0x60) & 0x00ff)
		v.CPU.FR[CF] = true
	} else {
		v.CPU.FR[CF] = false
	}
}

func das(v *VM) {
	oldAL := v.CPU.GR[AX] & 0x00ff
	oldCF := v.CPU.FR[CF]
	v.CPU.FR[CF] = false
	if (v.CPU.GR[AX]&0x00ff)&0x0f > 9 || v.CPU.FR[AF] {
		v.CPU.FR[CF] = oldCF || (v.CPU.GR[AX]&0x00ff-6) > 0xff
		v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | ((v.CPU.GR[AX]&0x00ff - 6) & 0x00ff)
		v.CPU.FR[AF] = true
	} else {
		v.CPU.FR[AF] = false
	}
	if oldAL > 0x99 || oldCF {
		v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | ((v.CPU.GR[AX]&0x00ff - 0x60) & 0x00ff)
		v.CPU.FR[CF] = true
	} else {
		v.CPU.FR[CF] = false
	}
}

func aaa(v *VM) {
	if v.CPU.GR[AX]&0x00ff&0x0f > 9 || v.CPU.FR[AF] {
		v.CPU.GR[AX] += 0x106
		v.CPU.FR[AF] = true
		v.CPU.FR[CF] = true
	} else {
		v.CPU.FR[AF] = false
		v.CPU.FR[CF] = false
	}
	v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | (v.CPU.GR[AX] & 0x00ff & 0x0f)
}

func aas(v *VM) {
	if v.CPU.GR[AX]&0x00ff&0x0f > 9 || v.CPU.FR[AF] {
		v.CPU.GR[AX] -= 6
		v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | (v.CPU.GR[AX]&0x00ff - 1)
		v.CPU.FR[AF] = true
		v.CPU.FR[CF] = true
		v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | (v.CPU.GR[AX] & 0x00ff & 0x0f)
	} else {
		v.CPU.FR[AF] = false
		v.CPU.FR[CF] = false
		v.CPU.GR[AX] = (v.CPU.GR[AX] & 0xff00) | (v.CPU.GR[AX] & 0x00ff & 0x0f)
	}
}

func aam(v *VM) {
	tmpAL := v.CPU.GR[AX] & 0x00ff
	v.CPU.GR[AX] = tmpAL/0x0a<<8 | v.CPU.GR[AX]&0x00ff
	v.CPU.GR[AX] = v.CPU.GR[AX]&0xff00 | tmpAL%0x0a
	v.CPU.FR[SF] = v.CPU.GR[AX]&0x8000>>7 == 1
	v.CPU.FR[ZF] = v.CPU.GR[AX]&0x00ff == 0
	v.CPU.FR[PF] = checkPF(uint8(v.CPU.GR[AX] & 0x00ff))
}

func aad(v *VM) {
	tmpAL := v.CPU.GR[AX] & 0x00ff
	tmpAH := v.CPU.GR[AX] & 0xff00
	v.CPU.GR[AX] = (tmpAL + tmpAH*10) & 0x00ff
	v.CPU.FR[SF] = v.CPU.GR[AX]&0x0080>>7 == 1
	v.CPU.FR[ZF] = v.CPU.GR[AX]&0x00ff == 0
	v.CPU.FR[PF] = checkPF(uint8(v.CPU.GR[AX] & 0x00ff))
}
