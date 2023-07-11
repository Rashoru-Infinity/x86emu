package vm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"syscall"
)

const (
	MM = iota
	FS
)

const (
	_UNUSED_0 = iota
	EXIT
	FORK
	READ
	WRITE
	OPEN
	CLOSE
	WAIT
	CREAT
	LINK
	UNLINK
	WAITPID
	CHDIR
	TIME
	MKNOD
	CHMOD
	CHOWN
	BRK
	STAT
	LSEEK
	GETPID
	MOUNT
	UMOUNT
	SETUID
	GETUID
	STIME
	PTRACE
	ALARM
	FSTAT
	PAUSE
	UTIME
	_UNUSED_31
	_UNUSED_32
	ACCESS
	_UNUSED_34
	_UNUSED_35
	SYNC
	KILL
	RENAME
	MKDIR
	RMDIR
	DUP
	PIPE
	TIMES
	_UNUSED_44
	_UNUSED_45
	SETGID
	GETGID
	SIGNAL
	_UNUSED_49
	_UNUSED_50
	_UNUSED_51
	_UNUSED_52
	_UNUSED_53
	IOCTL
	FCNTL
	_UNUSED_56
	_UNUSED_57
	_UNUSED_58
	EXEC
	UMASK
	CHROOT
	SETSID
	GETPGRP
	KSIG
	UNPAUSE
	_UNUSED_66
	REVIVE
	TASK_REPLY
	_UNUSED_69
	_UNUSED_70
	SIGACTION
	SIGSUSPEND
	SIGPENDING
	SIGPROCMASK
	SIGRETURN
	REBOOT
	SVRCTL
)

type MS1 struct {
	M1i1 int16
	M1i2 int16
	M1i3 int16
	M1p1 uint16
	M1p2 uint16
	M1p3 uint16
}

type MS2 struct {
	M2i1 int16
	M2i2 int16
	M2i3 int16
	M2l1 int32
	M2l2 int32
	M2p1 uint16
}

type MS3 struct {
	M3i1  int16
	M3i2  int16
	M3p1  uint16
	M3ca1 uint16
}

type MS4 struct {
	M4l1 int32
	M4l2 int32
	M4l3 int32
	M4l4 int32
	M4l5 int32
}

type MS5 struct {
	M5c1 int8
	M5c2 int8
	M5i1 int16
	M5i2 int16
	M5l1 int32
	M5l2 int32
	M5l3 int32
}

type MS6 struct {
	M6i1 int16
	M6i2 int16
	M6i3 int16
	M6l1 int32
	M6f1 int16
}

type MSG struct {
	M_source int16
	M_type   int16
}

func X86syscall(v *VM) (MSG, error) {
	msg := MSG{}
	br := bytes.NewReader(v.Data[v.CPU.GR[BX]:])
	var err error = nil
	binary.Read(br, binary.LittleEndian, &msg)
	switch msg.M_type {
	case EXIT:
		{
			msg1 := MS1{}
			arg := v.Data[v.CPU.GR[BX]+4:]
			br = bytes.NewReader(arg)
			binary.Read(br, binary.LittleEndian, &msg1)
			if v.Debug.DebugMode {
				fmt.Fprintf(os.Stderr, "<exit(%d)>\n", msg1.M1i1)
			}
			err = errors.New("end of program")
		}
	case FORK:
	case READ:
	case WRITE:
		{
			msg1 := MS1{}
			arg := v.Data[v.CPU.GR[BX]+4:]
			br = bytes.NewReader(arg)
			binary.Read(br, binary.LittleEndian, &msg1)
			if v.Debug.DebugMode {
				fmt.Fprintf(os.Stderr, "<write(%d, 0x%04x, %d)", msg1.M1i1, v.CPU.GR[BX]+msg1.M1p1, msg1.M1i2)
			}
			ret, err := syscall.Write((int)(msg1.M1i1), v.Data[v.CPU.GR[BX]+msg1.M1p1:v.CPU.GR[BX]+msg1.M1p1+(uint16)(msg1.M1i2)])
			if err != nil {
				msg.M_type = -(int16)(err.(syscall.Errno))
				break
			}
			msg.M_type = (int16)(ret)
			if v.Debug.DebugMode {
				fmt.Fprintf(os.Stderr, " ==> %d>\n", msg.M_type)
			}
		}
	case OPEN:
	case CLOSE:
	case WAIT:
	case CREAT:
	case LINK:
	case UNLINK:
	case WAITPID:
	case CHDIR:
	case TIME:
	case MKNOD:
	case CHMOD:
	case CHOWN:
	case BRK:
	case STAT:
	case LSEEK:
	case GETPID:
	case MOUNT:
	case UMOUNT:
	case SETUID:
	case GETUID:
	case STIME:
	case PTRACE:
	case ALARM:
	case FSTAT:
	case PAUSE:
	case UTIME:
	case ACCESS:
	case SYNC:
	case KILL:
	case RENAME:
	case MKDIR:
	case RMDIR:
	case DUP:
	case PIPE:
	case TIMES:
	case SETGID:
	case GETGID:
	case SIGNAL:
	case IOCTL:
	case FCNTL:
	case EXEC:
	case UMASK:
	case CHROOT:
	case SETSID:
	case GETPGRP:
	case KSIG:
	case UNPAUSE:
	case REVIVE:
	case TASK_REPLY:
	case SIGACTION:
	case SIGSUSPEND:
	case SIGPENDING:
	case SIGPROCMASK:
	case SIGRETURN:
	case REBOOT:
	case SVRCTL:
	}
	return msg, err
}
