package vm

func grp1(v *VM, op byte) {
	s := 0b00000010 & op >> 1
	w := 0b00000001 & op
	op = v.fetch()
	mod := op & 0b11000000 >> 6
	rm := op & 0b00000111
	op = op & 0b00111000 >> 3
	switch op {
	case 0b000:
		grp1add(v, s, w, mod, rm)
	case 0b001:
		grp1or(v, w, mod, rm)
	case 0b010:
		grp1adc(v, s, w, mod, rm)
	case 0b011:
		grp1sbb(v, s, w, mod, rm)
	case 0b100:
		grp1and(v, w, mod, rm)
	case 0b101:
		grp1sub(v, s, w, mod, rm)
	case 0b110:
		grp1xor(v, w, mod, rm)
	case 0b111:
		grp1cmp(v, s, w, mod, rm)
	}
}
