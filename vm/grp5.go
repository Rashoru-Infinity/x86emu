package vm

func grp5(v *VM) {
	op := v.fetch()
	mod := 0b11000000 & op >> 6
	rm := 0b00000111 & op
	switch 0b00111000 & op >> 3 {
	case 0b000:
		grp5inc(v, mod, rm)
	case 0b001:
		grp5dec(v, mod, rm)
	case 0b010:
		grp5call(v, mod, rm)
	case 0b100:
		grp5jmp(v, mod, rm)
	case 0b110:
		grp5push(v, mod, rm)
	}
}
