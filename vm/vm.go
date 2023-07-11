package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func setHeader(v *VM) {
	br := bytes.NewReader(v.Memory)
	binary.Read(br, binary.LittleEndian, &(v.Header))
}

func (v *VM) printRegister() {
	for i := AX; i < DI; i++ {
		fmt.Fprintf(os.Stderr, "%04x", v.CPU.GR[i])
		fmt.Fprintf(os.Stderr, " ")
	}
	fmt.Fprintf(os.Stderr, "----")
	fmt.Fprintf(os.Stderr, " ")
	fmt.Fprintf(os.Stderr, "%04x:", v.IP-1)
}

func (v *VM) printInst() {
	for _, b := range v.Debug.Buf {
		fmt.Fprintf(os.Stderr, "%02x", b)
	}
	// padding
	for i := len(v.Debug.Buf) * 2; i < 13; i++ {
		fmt.Fprintf(os.Stderr, " ")
	}
	switch v.Debug.Buf[0] {
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
		fmt.Fprintf(os.Stderr, "mov ")
		fmt.Fprintf(os.Stderr, "%s, ", grname[(int)(v.Debug.Buf[0]&0b00001000|v.Debug.Buf[0]&0b00000111)])
		if v.Debug.Buf[0]&0b00001000 != 0 {
			fmt.Fprintf(os.Stderr, "%02x", v.Debug.Buf[2])
		}
		fmt.Fprintf(os.Stderr, "%02x", v.Debug.Buf[1])
	case 0xcc:
		fmt.Fprintf(os.Stderr, "int 3")
	case 0xcd:
		fmt.Fprintf(os.Stderr, "int ")
		fmt.Fprintf(os.Stderr, "%02x", v.Debug.Buf[1])
	}
	fmt.Fprintf(os.Stderr, "\n")
}

func (v *VM) Load(file string, debug bool) error {
	var err error
	v.Memory, err = os.ReadFile(file)
	if err != nil {
		return err
	}
	v.Debug.DebugMode = debug
	v.Debug.Buf = make([]byte, 0, 6)
	initFR(v)
	initSR(v)
	initGR(v)
	setHeader(v)
	v.Data = v.Memory[(uint32)(v.Header.Hdrlen)+v.Header.Text:]
	return nil
}

func (v *VM) fetch() byte {
	ret := v.Memory[(uint32)(v.Header.Hdrlen)+v.IP]
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
		for i := AX; i < DI; i++ {
			fmt.Fprintf(os.Stderr, " %s ", strings.ToUpper(grname[i]))
			fmt.Fprintf(os.Stderr, " ")
		}
		fmt.Fprintf(os.Stderr, "FLAGS")
		fmt.Fprintf(os.Stderr, " ")
		fmt.Fprintf(os.Stderr, "IP\n")
	}
	for {
		if v.IP == v.Header.Text {
			break
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
		case 0x88:
			fallthrough
		case 0x89:
			fallthrough
		case 0x8a:
			fallthrough
		case 0x8b:
			if v.Debug.DebugMode {
				v.printRegister()
			}
			mov(v, op)
			if v.Debug.DebugMode {
				v.printInst()
			}
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
			if v.Debug.DebugMode {
				v.printRegister()
			}
			mov(v, op) //mov immediate to register/memory
			if v.Debug.DebugMode {
				v.printInst()
			}
		case 0xcc:
			fallthrough
		case 0xcd:
			if msg, err = interrupt(v, op); err != nil {
				goto ret
			}
		}
		if v.Debug.DebugMode {
			v.Debug.Buf = v.Debug.Buf[:0] // set length to 0
		}
	}
ret:
	return msg
}
