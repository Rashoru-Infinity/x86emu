package vm

func interrupt(v *VM, op byte) (MSG, error) {
	var (
		msg MSG
		err error
	)
	switch op {
	/*
		case 0xcc:
			return msg, err
	*/
	case 0xcd:
		v.fetch()
		/*
			if v.Debug.DebugMode {
				for _, b := range v.Debug.Buf {
					fmt.Fprintf(os.Stderr, "%02x", b)
				}
				// padding
				for i := len(v.Debug.Buf) * 2; i < 13; i++ {
					fmt.Fprintf(os.Stderr, " ")
				}
				switch op {
				case 0xcc:
					fmt.Fprintf(os.Stderr, "int 3")
				case 0xcd:
					fmt.Fprintf(os.Stderr, "int ")
					fmt.Fprintf(os.Stderr, "%02x", v.Debug.Buf[1])
				}
				fmt.Fprintf(os.Stderr, "\n")
			}
		*/
		msg, err = X86syscall(v)
	}
	return msg, err
}
