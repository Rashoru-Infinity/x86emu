package vm

func grp2(v *VM, op byte) {
	var (
		count uint8
	)
	_v := op & 0b00000010 >> 1
	w := op & 0b00000001
	data := v.fetch()
	mod := data & 0b11000000 >> 6
	_op := data & 0b00111000 >> 3
	rm := data & 0b00000111
	if _v == 0 {
		count = 1
	} else {
		count = uint8(v.CPU.GR[CX] & 0x00ff)
	}
	switch _op {
	case 0b000:
		rol(v, w, mod, rm, count)
	case 0b001:
		ror(v, w, mod, rm, count)
	case 0b010:
		rcl(v, w, mod, rm, count)
	case 0b011:
		rcr(v, w, mod, rm, count)
	case 0b100:
		shl(v, w, mod, rm, count)
	case 0b101:
		shr(v, w, mod, rm, count)
	case 0b111:
		sar(v, w, mod, rm, count)
	}
}

func rol(v *VM, w, mod, rm byte, count uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.Data[disp] = v.Data[disp]*2 + v.Data[disp]>>7
			}
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = v.Data[disp]&0x01 == 1
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[disp]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[disp]>>7 == 1
			}
			return
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.Data[eabase(v, uint16(rm))] = v.Data[eabase(v, uint16(rm))]*2 + v.Data[eabase(v, uint16(rm))]>>7
		}
		if count&0x1f%8 != 0 {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))]&0x01 == 1
		}
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))]>>7^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))]>>7 == 1
		}
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.Data[eabase(v, uint16(rm))+uint16(disp)] = v.Data[eabase(v, uint16(rm))+uint16(disp)]*2 + v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7
			}
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x01 == 1
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
			}
		} else {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]*2 + v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7
			}
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x01 == 1
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7 == 1
			}
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.Data[eabase(v, uint16(rm))+disp] = v.Data[eabase(v, uint16(rm))+disp]*2 + v.Data[eabase(v, uint16(rm))+disp]>>7
		}
		if count&0x1f%8 != 0 {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+disp]&0x01 == 1
		}
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+disp]>>7^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+disp]>>7 == 1
		}
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			tmp := uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				tmp = tmp*2 + tmp>>7
			}
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(tmp)
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = tmp&0x01 == 1
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7 == 1
			}
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			tmp := uint8(v.CPU.GR[int(1<<3|rm&0x03)] & 0xff00 >> 8)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				tmp = tmp*2 + tmp>>7
			}
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | uint16(tmp)<<8
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = tmp&0x01 == 1
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7 == 1
			}
		default:
			for _cnt := 0; _cnt < int(count&0x1f%16); _cnt++ {
				v.CPU.GR[int(w<<3|rm)] = v.CPU.GR[int(w<<3|rm)]*2 + v.CPU.GR[int(w<<3|rm)]>>15
			}
			if count&0x1f%16 != 0 {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)]&0x01 == 1
			}
			if count&0x1f%16 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)]>>15)^1 == 1
			} else if count&0x1f%16 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
			}
		}
	}
}

func ror(v *VM, w, mod, rm byte, count uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.Data[disp] = v.Data[disp]/2 + v.Data[disp]<<7
			}
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = v.Data[disp]>>7 == 1
			}
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = (v.Data[disp]>>7)^(v.Data[disp]&0x40>>6) == 1
			}
			return
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.Data[eabase(v, uint16(rm))] = v.Data[eabase(v, uint16(rm))]/2 + v.Data[eabase(v, uint16(rm))]<<7
		}
		if count&0x1f%8 != 0 {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))]>>7 == 1
		}
		if count&0x1f%8 == 1 {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))]>>7)^(v.Data[eabase(v, uint16(rm))]&0x40>>6) == 1
		}
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.Data[eabase(v, uint16(rm))+uint16(disp)] = v.Data[eabase(v, uint16(rm))+uint16(disp)]/2 + v.Data[eabase(v, uint16(rm))+uint16(disp)]<<7
			}
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
			}
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7)^(v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x40>>6) == 1
			}
		} else {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]/2 + v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]<<7
			}
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7 == 1
			}
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7)^(v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x40>>6) == 1
			}
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.Data[eabase(v, uint16(rm))+disp] = v.Data[eabase(v, uint16(rm))+disp]/2 + v.Data[eabase(v, uint16(rm))+disp]<<7
		}
		if count&0x1f%8 != 0 {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+disp]>>7 == 1
		}
		if count&0x1f%8 == 1 {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+disp]>>7)^(v.Data[eabase(v, uint16(rm))+disp]&0x40>>6) == 1
		}
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			tmp := uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				tmp = tmp/2 + tmp<<7
			}
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(tmp)
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = tmp>>7 == 1
			}
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = (tmp>>7)^(tmp&0x40>>6) == 1
			}
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			tmp := uint8(v.CPU.GR[int(1<<3|rm&0x03)] & 0xff00 >> 8)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				tmp = tmp/2 + tmp<<7
			}
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | uint16(tmp)<<8
			if count&0x1f%8 != 0 {
				v.CPU.FR[CF] = tmp>>7 == 1
			}
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = (tmp>>7)^(tmp&0x40>>6) == 1
			}
		default:
			for _cnt := 0; _cnt < int(count&0x1f%16); _cnt++ {
				v.CPU.GR[int(w<<3|rm)] = v.CPU.GR[int(w<<3|rm)]/2 + v.CPU.GR[int(w<<3|rm)]<<15
			}
			if count&0x1f%16 != 0 {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)]&0x01 == 1
			}
			if count&0x1f%16 == 1 {
				v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)]>>15)^(v.CPU.GR[int(w<<3|rm)]&0x4000>>14) == 1
			}
		}
	}
}

func rcl(v *VM, w, mod, rm byte, count uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[disp]>>7 == 1
				v.Data[disp] = v.Data[disp]*2 + v.Data[disp]>>7
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.Data[disp]>>7)^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.Data[disp] >> 7) == 1
			}
			return
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))]>>7 == 1
			v.Data[eabase(v, uint16(rm))] = v.Data[eabase(v, uint16(rm))]*2 + v.Data[eabase(v, uint16(rm))]>>7
		}
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))]>>7)^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))] >> 7) == 1
		}
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
				v.Data[eabase(v, uint16(rm))+uint16(disp)] = v.Data[eabase(v, uint16(rm))+uint16(disp)]*2 + v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7)^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+uint16(disp)] >> 7) == 1
			}
		} else {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7 == 1
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]*2 + v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7)^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] >> 7) == 1
			}
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+disp]>>7 == 1
			v.Data[eabase(v, uint16(rm))+disp] = v.Data[eabase(v, uint16(rm))+disp]*2 + v.Data[eabase(v, uint16(rm))+disp]>>7
		}
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+disp]>>7)^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = (v.Data[eabase(v, uint16(rm))+disp] >> 7) == 1
		}
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			tmp := uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp>>7 == 1
				tmp = tmp*2 + tmp>>7
			}
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(tmp)
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (tmp>>7)^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = (tmp >> 7) == 1
			}
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			tmp := uint8(v.CPU.GR[int(1<<3|rm&0x03)] & 0xff00 >> 8)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp>>7 == 1
				tmp = tmp*2 + tmp>>7
			}
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | uint16(tmp)<<8
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (tmp>>7)^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7 == 1
			}
		default:
			for _cnt := 0; _cnt < int(count&0x1f%16); _cnt++ {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
				v.CPU.GR[int(w<<3|rm)] = v.CPU.GR[int(w<<3|rm)]*2 + v.CPU.GR[int(w<<3|rm)]&0x8000>>15
			}
			if count&0x1f%16 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)]>>15)^1 == 1
			} else if count&0x1f%16 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
			}
		}
	}
}

func rcr(v *VM, w, mod, rm byte, count uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[disp]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[disp]>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[disp]&0x01 == 1
				v.Data[disp] = v.Data[disp]/2 + v.Data[disp]<<7
			}
			return
		}
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))]>>7^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))]>>7 == 1
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))]&0x01 == 1
			v.Data[eabase(v, uint16(rm))] = v.Data[eabase(v, uint16(rm))]/2 + v.Data[eabase(v, uint16(rm))]<<7
		}
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x01 == 1
				v.Data[eabase(v, uint16(rm))+uint16(disp)] = v.Data[eabase(v, uint16(rm))+uint16(disp)]/2 + v.Data[eabase(v, uint16(rm))+uint16(disp)]<<7
			}
		} else {
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x01 == 1
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]/2 + v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]<<7
			}
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+disp]&0x01 == 1
			v.Data[eabase(v, uint16(rm))+disp] = v.Data[eabase(v, uint16(rm))+disp]/2 + v.Data[eabase(v, uint16(rm))+disp]<<7
		}
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			tmp := uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff)
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7^1 == 1
			} else {
				v.CPU.FR[OF] = tmp>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp&0x01 == 1
				tmp = tmp/2 + tmp<<7
			}
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(tmp)
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			tmp := uint8(v.CPU.GR[int(1<<3|rm&0x03)] & 0xff00 >> 8)
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp>>7 == 1
				tmp = tmp/2 + tmp<<7
			}
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | uint16(tmp)<<8
		default:
			if count&0x1f%16 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.CPU.GR[int(w<<3|rm)]>>15^1 == 1
			} else if count&0x1f%16 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%16); _cnt++ {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
				v.CPU.GR[int(w<<3|rm)] = v.CPU.GR[int(w<<3|rm)]/2 + v.CPU.GR[int(w<<3|rm)]<<15
			}
		}
	}
}

func shl(v *VM, w, mod, rm byte, count uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[disp]>>7 == 1
				v.Data[disp] *= 2
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[disp]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[disp]>>7 == 1
			}
			v.CPU.FR[PF] = checkPF(v.Data[disp])
			v.CPU.FR[ZF] = v.Data[disp] == 0
			v.CPU.FR[SF] = v.Data[disp]&0x80 != 0
			return
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))]>>7 == 1
			v.Data[eabase(v, uint16(rm))] *= 2
		}
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))]>>7^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))]>>7 == 1
		}
		v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))])
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))] == 0
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))]&0x80 != 0
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
				v.Data[eabase(v, uint16(rm))+uint16(disp)] *= 2
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
			}
			v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))+uint16(disp)])
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+uint16(disp)] == 0
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x80 != 0
		} else {
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7 == 1
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] *= 2
			}
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7 == 1
			}
			v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))-uint16(^disp+1)])
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] == 0
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x80 != 0
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
			v.Data[eabase(v, uint16(rm))+disp] *= 2
		}
		if count&0x1f%8 == 1 && v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+disp]>>7^1 == 1
		} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+disp]>>7 == 1
		}
		v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))+disp])
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+disp] == 0
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+disp]&0x80 != 0
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			tmp := uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp>>7 == 1
				tmp *= 2
			}
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(tmp)
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7 == 1
			}
			v.CPU.FR[PF] = checkPF(tmp)
			v.CPU.FR[ZF] = tmp == 0
			v.CPU.FR[SF] = tmp&0x80 != 0
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			tmp := uint8(v.CPU.GR[int(1<<3|rm&0x03)] & 0xff00 >> 8)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp>>7 == 1
				tmp = tmp * 2
			}
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | uint16(tmp)<<8
			if count&0x1f%8 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7^1 == 1
			} else if count&0x1f%8 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = tmp>>7 == 1
			}
			v.CPU.FR[PF] = checkPF(tmp)
			v.CPU.FR[ZF] = tmp == 0
			v.CPU.FR[SF] = tmp&0x80 != 0
		default:
			for _cnt := 0; _cnt < int(count&0x1f%16); _cnt++ {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
				v.CPU.GR[int(w<<3|rm)] = v.CPU.GR[int(w<<3|rm)] * 2
			}
			if count&0x1f%16 == 1 && v.CPU.FR[CF] {
				v.CPU.FR[OF] = (v.CPU.GR[int(w<<3|rm)]>>15)^1 == 1
			} else if count&0x1f%16 == 1 && !v.CPU.FR[CF] {
				v.CPU.FR[OF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
			}
			v.CPU.FR[PF] = checkPF(uint8(v.CPU.GR[int(w<<3|rm)] & 0x00ff))
			v.CPU.FR[ZF] = v.CPU.GR[int(w<<3|rm)] == 0
			v.CPU.FR[SF] = v.CPU.GR[int(w<<3|rm)]&0x8000 != 0
		}
	}
}

func shr(v *VM, w, mod, rm byte, count uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			if count&0x1f%8 != 0 {
				v.CPU.FR[OF] = v.Data[disp]>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[disp]&0x01 == 1
				v.Data[disp] /= 2
			}
			v.CPU.FR[PF] = checkPF(v.Data[disp])
			v.CPU.FR[ZF] = v.Data[disp] == 0
			v.CPU.FR[SF] = v.Data[disp]&0x80 != 0
			return
		}
		if count&0x1f%8 == 1 {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))]>>7 == 1
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))]&0x01 == 1
			v.Data[eabase(v, uint16(rm))] /= 2
		}
		v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))])
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))] == 0
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))]&0x80 != 0
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x01 == 1
				v.Data[eabase(v, uint16(rm))+uint16(disp)] /= 2
			}
			v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))+uint16(disp)])
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+uint16(disp)] == 0
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x80 != 0
		} else {
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x01 == 1
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] /= 2
			}
			v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))-uint16(^disp+1)])
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] == 0
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x80 != 0
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		if count&0x1f%8 == 1 {
			v.CPU.FR[OF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]>>7 == 1
		}
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x01 == 1
			v.Data[eabase(v, uint16(rm))+disp] /= 2
		}
		v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))+disp])
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+disp] == 0
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+disp]&0x80 != 0
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			tmp := uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff)
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = uint8(v.CPU.GR[int(1<<3|rm)]&0x00ff)>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp&0x01 == 1
				tmp /= 2
			}
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(tmp)
			v.CPU.FR[PF] = checkPF(tmp)
			v.CPU.FR[ZF] = tmp == 0
			v.CPU.FR[SF] = tmp&0x80 != 0
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			tmp := uint8(v.CPU.GR[int(1<<3|rm&0x03)] & 0xff00 >> 8)
			if count&0x1f%8 == 1 {
				v.CPU.FR[OF] = uint8(v.CPU.GR[int(1<<3|rm&0x03)]&0xff00>>8)>>7 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp&0x01 == 1
				tmp /= 2
			}
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | uint16(tmp)<<8
			v.CPU.FR[PF] = checkPF(tmp)
			v.CPU.FR[ZF] = tmp == 0
			v.CPU.FR[SF] = tmp&0x80 != 0
		default:
			if count&0x1f%16 == 1 {
				v.CPU.FR[OF] = v.CPU.GR[int(1<<3|rm)]>>15 == 1
			}
			for _cnt := 0; _cnt < int(count&0x1f%16); _cnt++ {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
				v.CPU.GR[int(w<<3|rm)] /= 2
			}
			v.CPU.FR[PF] = checkPF(uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff))
			v.CPU.FR[ZF] = v.CPU.GR[int(w<<3|rm)] == 0
			v.CPU.FR[SF] = v.CPU.GR[int(w<<3|rm)]&0x8000 != 0
		}
	}
}

func sar(v *VM, w, mod, rm byte, count uint8) {
	switch mod {
	case 0b00:
		if rm == 0b110 {
			disp := uint16(v.fetch()) | uint16(v.fetch())<<8
			sign := v.Data[disp] & 0x80
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[disp]&0x01 == 1
				v.Data[disp] = (v.Data[disp] & 0x7f / 2) | sign
			}
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(v.Data[disp])
			v.CPU.FR[ZF] = v.Data[disp] == 0
			v.CPU.FR[SF] = v.Data[disp]&0x80 != 0
			return
		}
		sign := v.Data[eabase(v, uint16(rm))] & 0x80
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))]&0x01 == 1
			v.Data[eabase(v, uint16(rm))] = (v.Data[eabase(v, uint16(rm))] & 0x7f / 2) | sign
		}
		v.CPU.FR[OF] = false
		v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))])
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))] == 0
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))]&0x80 != 0
	case 0b01:
		disp := v.fetch()
		if disp < 0x80 {
			sign := v.Data[disp] & 0x80
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x01 == 1
				v.Data[eabase(v, uint16(rm))+uint16(disp)] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x7f/2 | sign
			}
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))+uint16(disp)])
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+uint16(disp)] == 0
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x80 != 0
		} else {
			sign := v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] & 0x80
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x01 == 1
				v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x7f/2 | sign
			}
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))-uint16(^disp+1)])
			v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)] == 0
			v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))-uint16(^disp+1)]&0x80 != 0
		}
	case 0b10:
		disp := uint16(v.fetch()) | uint16(v.fetch())<<8
		for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
			v.CPU.FR[CF] = v.Data[eabase(v, uint16(rm))+uint16(disp)]&0x01 == 1
			v.Data[eabase(v, uint16(rm))+disp] = v.Data[eabase(v, uint16(rm))+disp] & 0x7f / 2
		}
		v.CPU.FR[OF] = false
		v.CPU.FR[PF] = checkPF(v.Data[eabase(v, uint16(rm))+disp])
		v.CPU.FR[ZF] = v.Data[eabase(v, uint16(rm))+disp] == 0
		v.CPU.FR[SF] = v.Data[eabase(v, uint16(rm))+disp]&0x80 != 0
	case 0b11:
		switch w<<3 | rm {
		case AL:
			fallthrough
		case CL:
			fallthrough
		case DL:
			fallthrough
		case BL:
			tmp := uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp&0x01 == 1
				tmp = tmp & 0x7f / 2
			}
			v.CPU.GR[int(1<<3|rm)] = v.CPU.GR[int(1<<3|rm)]&0xff00 | uint16(tmp)
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(tmp)
			v.CPU.FR[ZF] = tmp == 0
			v.CPU.FR[SF] = tmp&0x80 != 0
		case AH:
			fallthrough
		case CH:
			fallthrough
		case DH:
			fallthrough
		case BH:
			tmp := uint8(v.CPU.GR[int(1<<3|rm&0x03)] & 0xff00 >> 8)
			for _cnt := 0; _cnt < int(count&0x1f%8); _cnt++ {
				v.CPU.FR[CF] = tmp&0x01 == 1
				tmp = tmp & 0x7f / 2
			}
			v.CPU.GR[int(1<<3|rm&0x03)] = v.CPU.GR[int(1<<3|rm&0x03)]&0x00ff | uint16(tmp)<<8
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(tmp)
			v.CPU.FR[ZF] = tmp == 0
			v.CPU.FR[SF] = tmp&0x80 != 0
		default:
			for _cnt := 0; _cnt < int(count&0x1f%16); _cnt++ {
				v.CPU.FR[CF] = v.CPU.GR[int(w<<3|rm)]>>15 == 1
				v.CPU.GR[int(w<<3|rm)] = v.CPU.GR[int(w<<3|rm)] & 0x7f / 2
			}
			v.CPU.FR[OF] = false
			v.CPU.FR[PF] = checkPF(uint8(v.CPU.GR[int(1<<3|rm)] & 0x00ff))
			v.CPU.FR[ZF] = v.CPU.GR[int(w<<3|rm)] == 0
			v.CPU.FR[SF] = v.CPU.GR[int(w<<3|rm)]&0x8000 != 0
		}
	}
}
