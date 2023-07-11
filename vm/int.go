package vm

func interrupt(v *VM, op byte) (MSG, error) {
	var (
		msg MSG
		err error
	)
	switch op {
	case 0xcc:
		return msg, err
	case 0xcd:
		v.fetch()
		if v.Debug.DebugMode {
			v.printRegister()
		}
		if v.Debug.DebugMode {
			v.printInst()
		}
		msg, err = X86syscall(v)
	}
	return msg, err
}
