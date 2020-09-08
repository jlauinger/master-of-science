package go_fuse

import (
	"log"
	"reflect"
	"unsafe"
)

// doBatchForget - forget a list of NodeIds
func doBatchForget(server *Server, req *request) {
	in := (*_BatchForgetIn)(req.inData)
	wantBytes := uintptr(in.Count) * unsafe.Sizeof(_ForgetOne{})
	safeCount := uintptr(in.Count)

	if uintptr(len(req.arg)) < wantBytes {
		// We have no return value to complain, so log an error.
		safeCount = uintptr(len(req.arg)) / unsafe.Sizeof(_ForgetOne{})
		log.Printf("Too few bytes for batch forget. Got %d bytes (enough for %d entries), want %d (%d entries)",
			len(req.arg), safeCount, wantBytes, in.Count)
	}

	forgets := make([]_ForgetOne, safeCount, safeCount)
	h := (*reflect.SliceHeader)(unsafe.Pointer(&forgets))
	h.Data = uintptr(unsafe.Pointer(&req.arg[0]))

	for i, f := range forgets {
		if server.opts.Debug {
			log.Printf("doBatchForget: rx %d %d/%d: FORGET i%d {Nlookup=%d}",
				req.inHeader.Unique, i+1, len(forgets), f.NodeId, f.Nlookup)
		}
		if f.NodeId == pollHackInode {
			continue
		}
		server.fileSystem.Forget(f.NodeId, f.Nlookup)
	}
}
