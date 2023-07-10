package vm

import (
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
	header executableHeader
	data []byte
	Debug  Debug
}

type Debug struct {
	DebugMode bool
	Buf       []byte
}

type executableHeader struct {
	magic   uint16
	flags   uint8
	cpu     uint8
	hdrlen  uint8
	unused  uint8
	version uint16
	text    uint32
	data    uint32
	bss     uint32
	entry   uint32
	total   uint32
	syms    uint32
}

func initFR(v *VM) {
	v.CPU.FR = make(map[int]bool)
	v.CPU.FR[OF] = false
	v.CPU.FR[DF] = false
	v.CPU.FR[IF] = false
	v.CPU.FR[TF] = false
	v.CPU.FR[SF] = false
	v.CPU.FR[ZF] = false
	v.CPU.FR[AF] = false
	v.CPU.FR[PF] = false
	v.CPU.FR[CF] = false
}

func initSR(v *VM) {
	v.CPU.SR = make(map[int]uint16)
	v.CPU.SR[ES] = 0
	v.CPU.SR[CS] = 0
	v.CPU.SR[SS] = 0
	v.CPU.SR[DS] = 0
}

func initGR(v *VM) {
	v.CPU.GR = make(map[int]uint16)
	v.CPU.GR[AX] = 0
	v.CPU.GR[CX] = 0
	v.CPU.GR[DX] = 0
	v.CPU.GR[BX] = 0
	v.CPU.GR[SP] = 0
	v.CPU.GR[BP] = 0
	v.CPU.GR[SI] = 0
	v.CPU.GR[DI] = 0
}

func setHeader(v *VM) {
	offset := 0
	v.header.magic = binary.LittleEndian.Uint16(v.Memory[offset : offset+2])
	offset += 2
	v.header.flags = (uint8)(v.Memory[offset])
	offset++
	v.header.cpu = (uint8)(v.Memory[offset])
	offset++
	v.header.hdrlen = (uint8)(v.Memory[offset])
	offset++
	v.header.unused = (uint8)(v.Memory[offset])
	offset++
	v.header.version = binary.LittleEndian.Uint16(v.Memory[offset : offset+2])
	offset += 2
	v.header.text = binary.LittleEndian.Uint32(v.Memory[offset : offset+4])
	offset += 4
	v.header.data = binary.LittleEndian.Uint32(v.Memory[offset : offset+4])
	offset += 4
	v.header.bss = binary.LittleEndian.Uint32(v.Memory[offset : offset+4])
	offset += 4
	v.header.entry = binary.LittleEndian.Uint32(v.Memory[offset : offset+4])
	offset += 4
	v.header.total = binary.LittleEndian.Uint32(v.Memory[offset : offset+4])
	offset += 4
	v.header.syms = binary.LittleEndian.Uint32(v.Memory[offset : offset+4])
	offset += 4
}

func (v *VM) printRegister() {
	for i := AX; i < DI; i++ {
		fmt.Printf("%04x", v.CPU.GR[i])
		fmt.Printf(" ")
	}
	fmt.Printf("----")
	fmt.Printf(" ")
	fmt.Printf("%04x:", v.IP-1)
}

func (v *VM) printInst() {
	for _, b := range v.Debug.Buf {
		fmt.Printf("%02x", b)
	}
	// padding
	for i := len(v.Debug.Buf) * 2; i < 13; i++ {
		fmt.Printf(" ")
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
		fmt.Printf("mov ")
		fmt.Printf("%s, ", grname[(int)(v.Debug.Buf[0]&0b00001000|v.Debug.Buf[0]&0b00000111)])
		if v.Debug.Buf[0]&0b00001000 != 0 {
			fmt.Printf("%02x", v.Debug.Buf[2])
		}
		fmt.Printf("%02x", v.Debug.Buf[1])
	case 0xcc:
		fmt.Printf("int 3")
	case 0xcd:
		fmt.Printf("int ")
		fmt.Printf("%02x", v.Debug.Buf[1])
	}
	fmt.Printf("\n")
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
	v.data = v.Memory[(uint32)(v.header.hdrlen) + v.header.text:]
	return nil
}

func (v *VM) fetch() byte {
	ret := v.Memory[(uint32)(v.header.hdrlen)+v.IP]
	v.IP++
	if v.Debug.DebugMode {
		v.Debug.Buf = append(v.Debug.Buf, ret)
	}
	return ret
}

func imd2regmem(v *VM, op byte) {
	w := (int)(0b00001000&op) >> 3
	reg := (int)(0b00000111 & op)
	data := (uint16)(v.fetch())
	if w != 0 {
		data |= ((uint16)(v.fetch())) << 8
	}
	switch op {
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
		if w == 0 { // 8bit data
			if reg <= BL { // low
				v.CPU.GR[1<<3|reg] &= 0xff00
				v.CPU.GR[1<<3|reg] |= data
			} else { // high
				v.CPU.GR[1<<3|reg] &= 0x00ff
				v.CPU.GR[1<<3|reg] |= data << 8
			}
		} else { // 16bit data
			v.CPU.GR[w<<3|reg] = data
		}
		return
	}
}

func interrupt(v *VM, op byte) {
	switch op {
	case 0xcc:
		return
	case 0xcd:
		data := v.fetch()
		switch data {
		}
	}
}

func (v *VM) Run() {
	if v.Debug.DebugMode {
		// print register names
		for i := AX; i < DI; i++ {
			fmt.Printf(" %s ", strings.ToUpper(grname[i]))
			fmt.Printf(" ")
		}
		fmt.Printf("FLAGS")
		fmt.Printf(" ")
		fmt.Printf("IP\n")
	}
	for {
		if v.IP == v.header.text {
			break
		}
		op := v.fetch()
		switch op {
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
			imd2regmem(v, op) //mov immediate to register/memory
			if v.Debug.DebugMode {
				v.printInst()
			}
		case 0xcc:
			fallthrough
		case 0xcd:
			if v.Debug.DebugMode {
				v.printRegister()
			}
			interrupt(v, op) //int
			if v.Debug.DebugMode {
				v.printInst()
			}
		}
		if v.Debug.DebugMode {
			v.Debug.Buf = v.Debug.Buf[:0] // set length to 0
		}
	}
}
