package vm

func grp4(v *VM, op byte) {
	data := v.fetch()
	mod := data & 0b11000000 >> 6
	_op := data & 0b00111000 >> 4
	rm := data & 0b00000011
	switch _op {
	case 0b000:
		grp4inc(v, mod, rm)
	case 0b001:
		grp4dec(v, mod, rm)
	}
}
