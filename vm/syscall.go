package vm

import (
	"encoding/binary"
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
	m1i1 int16
	m1i2 int16
	m1i3 int16
	m1p1 uint16
	m1p2 uint16
	m1p3 uint16
}

type MS2 struct {
	m2i1 int16
	m2i2 int16
	m2i3 int16
	m2l1 int32
	m2l2 int32
	m2p1 uint16
}

type MS3 struct {
	m3i1 int16
	m3i2 int16
	m3p1 uint16
	m3ca1 uint16
}

type MS4 struct {
	m4l1 int32
	m4l2 int32
	m4l3 int32
	m4l4 int32
	m4l5 int32
}

type MS5 struct {
	m5c1 int8
	m5c2 int8
	m5i1 int16
	m5i2 int16
	m5l1 int32
	m5l2 int32
	m5l3 int32
}

type MS6 struct {
	m6i1 int16
	m6i2 int16
	m6i3 int16
	m6l1 int32
	m6f1 int16
}

type MSG struct {
	m_source int16
	m_type int16
	ms1 MS1
	ms2 MS2
	ms3 MS3
	ms4 MS4
	ms5 MS5
	ms6 MS6
}

func syscall(v *VM, who, syscallnr int, msg *MSG) {

	switch (syscallnr) {
	case EXIT:
	case FORK:
	case READ:
	case WRITE:
		msg.ms1.m1i1 = (int16)(binary.LittleEndian.Uint16(v.data[v.CPU.GR[BX]:v.CPU.GR[BX]+2]))
		msg.ms1.m1i2 = (int16)(binary.LittleEndian.Uint16(v.data[v.CPU.GR[BX]+2:v.CPU.GR[BX]+4]))
		msg.ms1.m1i3 = (int16)(binary.LittleEndian.Uint16(v.data[v.CPU.GR[BX]+4:v.CPU.GR[BX]+6]))
		// m1i1~3のポインタのポインタが挿すデータを使う
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
}