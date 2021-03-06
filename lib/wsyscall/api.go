package wsyscall

import "syscall"

const (
	MS_BIND = 1 << iota

	RUSAGE_CHILDREN = iota
	RUSAGE_SELF
	RUSAGE_THREAD
)

type Rusage struct {
	Utime    Timeval
	Stime    Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
}

type Stat_t struct {
	Dev     uint64
	Ino     uint64
	Nlink   uint64
	Mode    uint32
	Uid     uint32
	Gid     uint32
	Rdev    uint64
	Size    int64
	Blksize int64
	Blocks  int64
	Atim    syscall.Timespec
	Mtim    syscall.Timespec
	Ctim    syscall.Timespec
}

type Timeval struct {
	Sec  int64
	Usec int64
}

func Lstat(path string, statbuf *Stat_t) error {
	var rawStatbuf syscall.Stat_t
	if err := syscall.Lstat(path, &rawStatbuf); err != nil {
		return err
	}
	convertStat(statbuf, &rawStatbuf)
	return nil
}

func Mount(source string, target string, fstype string, flags uintptr,
	data string) error {
	return mount(source, target, fstype, flags, data)
}

func Getrusage(who int, rusage *Rusage) error {
	return getrusage(who, rusage)
}

func SetAllGid(gid int) error {
	return setAllGid(gid)
}

func SetAllUid(uid int) error {
	return setAllUid(uid)
}

func Stat(path string, statbuf *Stat_t) error {
	var rawStatbuf syscall.Stat_t
	if err := syscall.Stat(path, &rawStatbuf); err != nil {
		return err
	}
	convertStat(statbuf, &rawStatbuf)
	return nil
}

func UnshareMountNamespace() error {
	return unshareMountNamespace()
}
