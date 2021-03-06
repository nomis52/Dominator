package filesystem

import (
	"github.com/Symantec/Dominator/lib/format"
	"syscall"
	"time"
)

func (objSrv *ObjectServer) garbageCollector() (uint64, error) {
	if objSrv.gc == nil {
		return 0, nil
	}
	objSrv.rwLock.Lock()
	if time.Since(objSrv.lastGarbageCollection) < time.Second {
		objSrv.rwLock.Unlock()
		return 0, nil
	}
	objSrv.lastGarbageCollection = time.Now()
	objSrv.rwLock.Unlock()
	free, capacity, err := objSrv.getSpaceMetrics()
	if err != nil {
		return 0, err
	}
	cleanupStartPercent := sanitisePercentage(*objectServerCleanupStartPercent)
	cleanupStopPercent := sanitisePercentage(*objectServerCleanupStopPercent)
	if cleanupStopPercent >= cleanupStartPercent {
		cleanupStopPercent = cleanupStartPercent - 1
	}
	utilisation := (capacity - free) * 100 / capacity
	if utilisation < cleanupStartPercent {
		return 0, nil
	}
	bytesToDelete := (utilisation - cleanupStopPercent) * capacity / 100
	bytesDeleted, err := objSrv.gc(bytesToDelete)
	if err != nil {
		objSrv.logger.Printf("Error collecting garbage, only deleted: %s: %s\n",
			format.FormatBytes(bytesDeleted), err)
		return 0, err
	}
	return bytesDeleted, nil
}

func (t *ObjectServer) getSpaceMetrics() (uint64, uint64, error) {
	if fd, err := syscall.Open(t.baseDir, syscall.O_RDONLY, 0); err != nil {
		t.logger.Printf("error opening: %s: %s", t.baseDir, err)
		return 0, 0, err
	} else {
		defer syscall.Close(fd)
		var statbuf syscall.Statfs_t
		if err := syscall.Fstatfs(fd, &statbuf); err != nil {
			t.logger.Printf("error getting file-system stats: %s\n", err)
			return 0, 0, err
		}
		rootReservation := statbuf.Bfree - statbuf.Bavail
		return uint64(statbuf.Bavail) * uint64(statbuf.Bsize),
			uint64(statbuf.Blocks-rootReservation) * uint64(statbuf.Bsize), nil
	}
}

func sanitisePercentage(percent int) uint64 {
	if percent < 1 {
		return 1
	}
	if percent > 99 {
		return 99
	}
	return uint64(percent)
}
