package vm

func lea(v *VM, op byte) {
	switch op {
	case 0x8d:
		data := v.fetch()
		mod := (0b11000000 & data) >> 6
		reg := (0b00111000 & data) >> 3
		rm := 0b00000111 & data
		switch mod {
		case 0b00:
			if rm == 0b110 {
				disp := (int)(v.fetch()) | ((int)(v.fetch()) << 8)
				v.CPU.GR[(int)(1<<3|reg)] = uint16(disp)
				return
			}
			v.CPU.GR[(int)(1<<3|reg)] = eabase(v, (uint16)(rm))
		case 0b01:
			disp := v.fetch()
			if disp < 0x80 {
				v.CPU.GR[(int)(1<<3|reg)] = eabase(v, (uint16)(rm)) + (uint16)(disp)
			} else {
				v.CPU.GR[(int)(1<<3|reg)] = eabase(v, (uint16)(rm)) - (uint16)(^disp+1)&0x00ff
			}
		case 0b10:
			disp := (int)(v.fetch()) | ((int)(v.fetch()) << 8)
			v.CPU.GR[(int)(1<<3|reg)] = eabase(v, (uint16)(rm)) + (uint16)(disp)
		case 0b11:
			v.CPU.GR[(int)(1<<3|reg)] = v.CPU.GR[(int)(1<<3|rm)]
		}
	}
}
