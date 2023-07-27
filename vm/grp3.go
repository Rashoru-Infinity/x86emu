package vm

func grp3(v *VM, op byte) {
	w := 0b00000001 & op
	_op := v.fetch()
	mod := _op & 0b11000000 >> 6
	rm := _op & 0b00000111
	_op = _op & 0b00111000 >> 3
	switch _op {
	case 0b000:
		grp3test(v, w, mod, rm)
	case 0b010:
		grp3not(v, w, mod, rm)
	case 0b011:
		grp3neg(v, w, mod, rm)
	case 0b100:
		grp3mul(v, w, mod, rm)
	case 0b101:
		grp3imul(v, w, mod, rm)
	case 0b110:
		grp3div(v, w, mod, rm)
	case 0b111:
		grp3idiv(v, w, mod, rm)
	}
}
