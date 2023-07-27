package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	OF = iota
	DF
	IF
	TF
	SF
	ZF
	NF
	AF
	PF
	CF
)

const (
	ES = iota
	CS
	SS
	DS
)

const (
	AL = iota
	CL
	DL
	BL
	AH
	CH
	DH
	BH
	AX
	CX
	DX
	BX
	SP
	BP
	SI
	DI
)

var (
	grname = map[int]string{
		AL: "al",
		CL: "cl",
		DL: "dl",
		BL: "bl",
		AH: "ah",
		CH: "ch",
		DH: "dh",
		BH: "bh",
		AX: "ax",
		CX: "cx",
		DX: "dx",
		BX: "bx",
		SP: "sp",
		BP: "bp",
		SI: "si",
		DI: "di",
	}
)

type CPU struct {
	FR map[int]bool
	SR map[int]uint16
	GR map[int]uint16
}

type VM struct {
	CPU    CPU
	Memory []byte
	IP     uint32
	Header executableHeader
	Data   []byte
	Text   []byte
	Debug  Debug
}

type Debug struct {
	DebugMode bool
	Buf       []byte
}

type executableHeader struct {
	Magic   uint16
	Flags   uint8
	Cpu     uint8
	Hdrlen  uint8
	Unused  uint8
	Version uint16
	Text    uint32
	Data    uint32
	Bss     uint32
	Entry   uint32
	Total   uint32
	Syms    uint32
}

func initFR(v *VM) {
	v.CPU.FR = map[int]bool{
		OF: false,
		DF: false,
		IF: false,
		TF: false,
		SF: false,
		ZF: false,
		AF: false,
		PF: false,
		CF: false,
	}
}

func initSR(v *VM) {
	v.CPU.SR = map[int]uint16{
		ES: 0,
		CS: 0,
		SS: 0,
		DS: 0,
	}
}

func initGR(v *VM) {
	v.CPU.GR = map[int]uint16{
		AX: 0,
		CX: 0,
		DX: 0,
		BX: 0,
		SP: 0,
		BP: 0,
		SI: 0,
		DI: 0,
	}
}

func setHeader(v *VM, buf []byte) {
	br := bytes.NewReader(buf)
	binary.Read(br, binary.LittleEndian, &(v.Header))
}

func setMemory(v *VM, buf []byte) {
	offset := uint32(v.Header.Hdrlen)
	for i := offset; i < uint32(v.Header.Hdrlen)+uint32(v.Header.Text); i++ {
		v.Text[i-offset] = buf[i]
	}
	offset = uint32(v.Header.Hdrlen) + v.Header.Text
	for i := uint32(offset); i < uint32(v.Header.Hdrlen)+uint32(v.Header.Text+v.Header.Data); i++ {
		v.Data[i-offset] = buf[i]
	}
}

func (v *VM) printRegister() {
	for i := AX; i <= DI; i++ {
		fmt.Fprintf(os.Stderr, "%04x", v.CPU.GR[i])
		fmt.Fprintf(os.Stderr, " ")
	}
	if v.CPU.FR[OF] {
		fmt.Fprintf(os.Stderr, "O")
	} else {
		fmt.Fprintf(os.Stderr, "-")
	}
	if v.CPU.FR[SF] {
		fmt.Fprintf(os.Stderr, "S")
	} else {
		fmt.Fprintf(os.Stderr, "-")
	}
	if v.CPU.FR[ZF] {
		fmt.Fprintf(os.Stderr, "Z")
	} else {
		fmt.Fprintf(os.Stderr, "-")
	}
	if v.CPU.FR[CF] {
		fmt.Fprintf(os.Stderr, "C")
	} else {
		fmt.Fprintf(os.Stderr, "-")
	}
	fmt.Fprintf(os.Stderr, " ")
	fmt.Fprintf(os.Stderr, "%04x:", v.IP)
}

func initStack(v *VM, args, env []string) {
	argenvlen := 0
	for _, _v := range args {
		argenvlen += len(_v) + 1
	}
	for _, _v := range env {
		argenvlen += len(_v) + 1
	}
	if argenvlen%2 == 1 {
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = 0
	}
	for i := len(env) - 1; i >= 0; i-- {
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = 0
		for j := len(env[i]) - 1; j >= 0; j-- {
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = env[i][j]
		}
	}
	envoffset := v.CPU.GR[SP]
	argaddr := make([]uint16, 0, len(env))
	for i := len(args) - 1; i >= 0; i-- {
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = 0
		for j := len(args[i]) - 1; j >= 0; j-- {
			v.CPU.GR[SP]--
			v.Data[v.CPU.GR[SP]] = args[i][j]
		}
		argaddr = append(argaddr, v.CPU.GR[SP])
	}
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = 0
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = 0
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(envoffset >> 8)
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(envoffset & 0x00ff)
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = 0
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = 0
	for _, addr := range argaddr {
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(addr >> 8)
		v.CPU.GR[SP]--
		v.Data[v.CPU.GR[SP]] = byte(addr & 0x00ff)
	}
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(uint16(len(args)) >> 8)
	v.CPU.GR[SP]--
	v.Data[v.CPU.GR[SP]] = byte(uint16(len(args)) & 0x00ff)
}

func (v *VM) Load(file string, debug bool, args []string) error {
	v.Text = make([]byte, 0x10000)
	v.Data = make([]byte, 0x10000)
	fd, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return err
	}
	buf, err := io.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return err
	}
	fd.Close()
	initFR(v)
	initSR(v)
	initGR(v)
	setHeader(v, buf)
	setMemory(v, buf)
	env := make([]string, 0, 1)
	env = append(env, "PATH=/usr:/usr/bin")
	initStack(v, args, env)
	v.Debug.DebugMode = debug
	v.Debug.Buf = make([]byte, 0, 6)
	return nil
}

func (v *VM) fetch() byte {
	ret := v.Text[v.IP]
	v.IP++
	if v.Debug.DebugMode {
		v.Debug.Buf = append(v.Debug.Buf, ret)
	}
	return ret
}

func eabase(v *VM, rm uint16) uint16 {
	switch rm {
	case 0b000:
		return v.CPU.GR[BX] + v.CPU.GR[SI]
	case 0b001:
		return v.CPU.GR[DX] + v.CPU.GR[DI]
	case 0b010:
		return v.CPU.GR[BP] + v.CPU.GR[SI]
	case 0b011:
		return v.CPU.GR[BP] + v.CPU.GR[DI]
	case 0b100:
		return v.CPU.GR[SI]
	case 0b101:
		return v.CPU.GR[DI]
	case 0b110:
		return v.CPU.GR[BP]
	}
	return v.CPU.GR[BX]
}

func (v *VM) Run() MSG {
	var (
		msg MSG
		err error
	)
	if v.Debug.DebugMode {
		// print register names
		for i := AX; i <= DI; i++ {
			fmt.Fprintf(os.Stderr, " %s ", strings.ToUpper(grname[i]))
			fmt.Fprintf(os.Stderr, " ")
		}
		fmt.Fprintf(os.Stderr, "FLAGS")
		fmt.Fprintf(os.Stderr, " ")
		fmt.Fprintf(os.Stderr, "IP\n")
	}
	//stacktop := make([]uint16, 8)
	for {
		/*
			sp := v.CPU.GR[SP]
			for i := uint32(v.CPU.GR[SP]); i <= uint32(v.CPU.GR[SP])+14 && i < 0x10000; i += 2 {
				stacktop[(i-uint32(v.CPU.GR[SP]))/2] = binary.LittleEndian.Uint16(v.Data[i:])
			}
		*/
		/*
			if v.Debug.DebugMode {
				// print register names
				for i := AX; i <= DI; i++ {
					fmt.Fprintf(os.Stderr, " %s ", strings.ToUpper(grname[i]))
					fmt.Fprintf(os.Stderr, " ")
				}
				fmt.Fprintf(os.Stderr, "FLAGS")
				fmt.Fprintf(os.Stderr, " ")
				fmt.Fprintf(os.Stderr, "IP\n")
			}
		*/
		/*
			if v.IP >= v.Header.Text {
				break
			}
		*/
		if v.Debug.DebugMode {
			v.printRegister()
		}
		op := v.fetch()
		switch op {
		case 0x00:
			fallthrough
		case 0x01:
			fallthrough
		case 0x02:
			fallthrough
		case 0x03:
			fallthrough
		case 0x04:
			fallthrough
		case 0x05:
			add(v, op)
		case 0x06:
			push(v, op)
		case 0x07:
			pop(v, op)
		case 0x08:
			fallthrough
		case 0x09:
			fallthrough
		case 0x0a:
			fallthrough
		case 0x0b:
			fallthrough
		case 0x0c:
			fallthrough
		case 0x0d:
			or(v, op)
		case 0x0e:
			fallthrough
		case 0x10:
			fallthrough
		case 0x13:
			fallthrough
		case 0x14:
			fallthrough
		case 0x15:
			adc(v, op)
		case 0x16:
			push(v, op)
		case 0x17:
			pop(v, op)
		case 0x18:
			fallthrough
		case 0x19:
			fallthrough
		case 0x1a:
			fallthrough
		case 0x1b:
			fallthrough
		case 0x1c:
			fallthrough
		case 0x1d:
			sbb(v, op)
		case 0x1e:
			push(v, op)
		case 0x1f:
			pop(v, op)
		case 0x20:
			fallthrough
		case 0x21:
			fallthrough
		case 0x22:
			fallthrough
		case 0x23:
			fallthrough
		case 0x24:
			fallthrough
		case 0x25:
			and(v, op)
		case 0x26:
			//segment override ES:
		case 0x27:
			daa(v)
		case 0x28:
			fallthrough
		case 0x29:
			fallthrough
		case 0x2a:
			fallthrough
		case 0x2b:
			fallthrough
		case 0x2c:
			fallthrough
		case 0x2d:
			sub(v, op)
		case 0x2e:
			//segment override
		case 0x2f:
			das(v)
		case 0x30:
			fallthrough
		case 0x31:
			fallthrough
		case 0x32:
			fallthrough
		case 0x33:
			fallthrough
		case 0x34:
			fallthrough
		case 0x35:
			xor(v, op)
		case 0x36:
			//segment override SS:
		case 0x37:
			aaa(v)
		case 0x38:
			fallthrough
		case 0x39:
			fallthrough
		case 0x3a:
			fallthrough
		case 0x3b:
			fallthrough
		case 0x3c:
			fallthrough
		case 0x3d:
			cmp(v, op)
		case 0x3e:
			//segment override DS:
		case 0x3f:
			aas(v)
		case 0x40:
			fallthrough
		case 0x41:
			fallthrough
		case 0x42:
			fallthrough
		case 0x43:
			fallthrough
		case 0x44:
			fallthrough
		case 0x45:
			fallthrough
		case 0x46:
			fallthrough
		case 0x47:
			inc(v, op)
		case 0x48:
			fallthrough
		case 0x49:
			fallthrough
		case 0x4a:
			fallthrough
		case 0x4b:
			fallthrough
		case 0x4c:
			fallthrough
		case 0x4d:
			fallthrough
		case 0x4e:
			fallthrough
		case 0x4f:
			dec(v, op)
		case 0x50:
			fallthrough
		case 0x51:
			fallthrough
		case 0x52:
			fallthrough
		case 0x53:
			fallthrough
		case 0x54:
			fallthrough
		case 0x55:
			fallthrough
		case 0x56:
			fallthrough
		case 0x57:
			push(v, op)
		case 0x58:
			fallthrough
		case 0x59:
			fallthrough
		case 0x5a:
			fallthrough
		case 0x5b:
			fallthrough
		case 0x5c:
			fallthrough
		case 0x5d:
			fallthrough
		case 0x5e:
			fallthrough
		case 0x5f:
			pop(v, op)
		case 0x70:
			fallthrough
		case 0x71:
			fallthrough
		case 0x72:
			fallthrough
		case 0x73:
			fallthrough
		case 0x74:
			fallthrough
		case 0x75:
			fallthrough
		case 0x76:
			fallthrough
		case 0x77:
			fallthrough
		case 0x78:
			fallthrough
		case 0x79:
			fallthrough
		case 0x7a:
			fallthrough
		case 0x7b:
			fallthrough
		case 0x7c:
			fallthrough
		case 0x7d:
			fallthrough
		case 0x7e:
			fallthrough
		case 0x7f:
			jcc8(v, op)
		case 0x80:
			fallthrough
		case 0x81:
			fallthrough
		case 0x82:
			fallthrough
		case 0x83:
			grp1(v, op)
		case 0x84:
			fallthrough
		case 0x85:
			test(v, op)
		case 0x86:
			fallthrough
		case 0x87:
			xchg(v, op)
		case 0x88:
			fallthrough
		case 0x89:
			fallthrough
		case 0x8a:
			fallthrough
		case 0x8b:
			fallthrough
		case 0x8c:
			mov(v, op)
		case 0x8d:
			lea(v, op)
		case 0x8e:
			mov(v, op)
		case 0x8f:
			pop(v, op)
		case 0x90:
			fallthrough
		case 0x91:
			fallthrough
		case 0x92:
			fallthrough
		case 0x93:
			fallthrough
		case 0x94:
			fallthrough
		case 0x95:
			fallthrough
		case 0x96:
			fallthrough
		case 0x97:
			xchg(v, op)
		case 0x98:
			cbw(v)
		case 0x99:
			cwd(v)
		case 0x9a:
			// call direct intersegment
		case 0x9b:
			// wait
		case 0x9c:
			// pushf
		case 0x9d:
			// popf
		case 0x9e:
			sahf(v)
		case 0x9f:
			lahf(v)
		case 0xa0:
			fallthrough
		case 0xa1:
			fallthrough
		case 0xa2:
			fallthrough
		case 0xa3:
			mov(v, op)
		case 0xa4:
			fallthrough
		case 0xa5:
			movs(v, op)
		case 0xa6:
			fallthrough
		case 0xa7:
			cmps(v, op)
		case 0xa8:
			fallthrough
		case 0xa9:
			test(v, op)
		case 0xaa:
			fallthrough
		case 0xab:
			stos(v, op)
		case 0xac:
			fallthrough
		case 0xad:
			lods(v, op)
		case 0xae:
			fallthrough
		case 0xaf:
			scas(v, op)
		case 0xb0:
			fallthrough
		case 0xb1:
			fallthrough
		case 0xb2:
			fallthrough
		case 0xb3:
			fallthrough
		case 0xb4:
			fallthrough
		case 0xb5:
			fallthrough
		case 0xb6:
			fallthrough
		case 0xb7:
			fallthrough
		case 0xb8:
			fallthrough
		case 0xb9:
			fallthrough
		case 0xba:
			fallthrough
		case 0xbb:
			fallthrough
		case 0xbc:
			fallthrough
		case 0xbd:
			fallthrough
		case 0xbe:
			fallthrough
		case 0xbf:
			mov(v, op)
		case 0xc2:
			fallthrough
		case 0xc3:
			ret(v, op)
		case 0xc4:
			// les
		case 0xc5:
			// lds
		case 0xc6:
			fallthrough
		case 0xc7:
			mov(v, op)
		case 0xca:
			// ret intersegment adding immediate to SP
		case 0xcb:
			// ret intersegment
		case 0xcc:
			// int 3
		case 0xcd:
			if msg, err = interrupt(v, op); err != nil {
				if v.Debug.DebugMode {
					for _, b := range v.Debug.Buf {
						fmt.Fprintf(os.Stderr, "%02x", b)
					}
					fmt.Fprintf(os.Stderr, "\n")
				}
				goto ret
			}
		case 0xd0:
			fallthrough
		case 0xd1:
			fallthrough
		case 0xd2:
			fallthrough
		case 0xd3:
			grp2(v, op)
		case 0xd4:
			aam(v)
		case 0xd5:
			aad(v)
		case 0xd7:
			xlat(v)
		case 0xe0:
			fallthrough
		case 0xe1:
			fallthrough
		case 0xe2:
			fallthrough
		case 0xe3:
			jcc8(v, op)
		case 0xe4:
			fallthrough
		case 0xe5:
			// in
		case 0xe6:
			fallthrough
		case 0xe7:
			// out
		case 0xe8:
			call(v)
		case 0xe9:
			jcc8(v, op)
		case 0xea:
			// jmp direct intersegment
		case 0xeb:
			jcc8(v, op)
		case 0xec:
			fallthrough
		case 0xed:
			// in
		case 0xee:
			fallthrough
		case 0xef:
			// out
		case 0xf0:
			// lock
		case 0xf2:
			fallthrough
		case 0xf3:
			// rep
		case 0xf4:
			// hlt
		case 0xf5:
			cmc(v)
		case 0xf6:
			fallthrough
		case 0xf7:
			grp3(v, op)
		case 0xf8:
			clc(v)
		case 0xf9:
			stc(v)
		case 0xfa:
			cli(v)
		case 0xfb:
			sti(v)
		case 0xfc:
			cld(v)
		case 0xfd:
			std(v)
		case 0xfe:
			grp4(v, op)
		case 0xff:
			grp5(v)
		}
		if v.Debug.DebugMode {
			for _, b := range v.Debug.Buf {
				fmt.Fprintf(os.Stderr, "%02x", b)
			}
		}
		if v.Debug.DebugMode {
			v.Debug.Buf = v.Debug.Buf[:0] // set length to 0
			/*
				fmt.Fprintf(os.Stderr, "---stack---\n")
				for i := 0; i < 10 && int(v.CPU.GR[SP])+i <= 0xffff; i++ {
					fmt.Fprintf(os.Stderr, "%04x:%02x\n", v.CPU.GR[SP]+uint16(i), v.Data[int(v.CPU.GR[SP])+i])
				}
				fmt.Fprintf(os.Stderr, "-----------\n")
			*/
			/*
				for _, v := range stacktop {
					fmt.Fprintf(os.Stderr, " [%x: %x]", sp, v)
					sp += 2
				}
			*/
			fmt.Fprintf(os.Stderr, "\n")
		}
	}
ret:
	return msg
}
