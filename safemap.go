package util

type smap struct {
	m            map[interface{}]interface{}
	readSig      chan *readReq
	writeSig     chan *writeReq
	lenSig       chan *lenReq
	terminateSig chan bool
	delSig       chan *delReq
	scanSig      chan *scanReq
}

type readReq struct {
	key   interface{}
	value interface{}
	ok    chan bool
}

type writeReq struct {
	key   interface{}
	value interface{}
	ok    chan bool
}

type lenReq struct {
	len chan int
}

type delReq struct {
	key interface{}
	ok  chan bool
}

type scanReq struct {
	do          func(interface{}, interface{})
	doWithBreak func(interface{}, interface{}) bool
	brea        int
	done        chan bool
}

// NewSmap returns an instance of the pointer of safemap
func NewSmap() *smap {
	var mp smap
	mp.m = make(map[interface{}]interface{})
	mp.readSig = make(chan *readReq)
	mp.writeSig = make(chan *writeReq)
	mp.lenSig = make(chan *lenReq)
	mp.delSig = make(chan *delReq)
	mp.scanSig = make(chan *scanReq)
	go mp.run()
	return &mp
}

//background function to operate map in one goroutine
//this can ensure that the map is  Concurrent security.
func (s *smap) run() {
	for {
		select {
		case read := <-s.readSig:
			if value, ok := s.m[read.key]; ok {
				read.value = value
				read.ok <- true
			} else {
				read.ok <- false
			}
		case write := <-s.writeSig:
			s.m[write.key] = write.value
			write.ok <- true
		case l := <-s.lenSig:
			l.len <- len(s.m)
		case sc := <-s.scanSig:
			if sc.brea == 0 {
				for k, v := range s.m {
					sc.do(k, v)
				}
			} else {
				for k, v := range s.m {
					ret := sc.doWithBreak(k, v)
					if ret {
						break
					}
				}
			}
			sc.done <- true
		case d := <-s.delSig:
			delete(s.m, d.key)
			d.ok <- true
		case <-s.terminateSig:
			return
		}
	}
}

//Get returns the value of  key which provided.
//if the key not found in map, ok will be false.
func (s *smap) Get(key interface{}) (interface{}, bool) {
	req := &readReq{
		key: key,
		ok:  make(chan bool),
	}
	s.readSig <- req
	ok := <-req.ok
	return req.value, ok
}

//Set set the key and value to map
//ok returns true indicates that key and value is successfully added to map
func (s *smap) Set(key interface{}, value interface{}) bool {
	req := &writeReq{
		key:   key,
		value: value,
		ok:    make(chan bool),
	}
	s.writeSig <- req
	return <-req.ok //TODO 暂时先是同步的，异步的可能存在使用方面的问题。
}

//Clear clears all the key and value in map.
func (s *smap) Clear() {
	s.m = make(map[interface{}]interface{})
}

//Size returns the size of map.
func (s *smap) Size() int {
	req := &lenReq{
		len: make(chan int),
	}
	s.lenSig <- req
	return <-req.len
}

//terminate s.Run function. this function is usually called for debug.
//after this do NOT use smap again, because it can make your program block.
func (s *smap) TerminateBackGoroutine() {
	s.terminateSig <- true
}

//Del delete the key in map
func (s *smap) Del(key interface{}) bool {
	req := &delReq{
		key: key,
		ok:  make(chan bool),
	}
	s.delSig <- req
	return <-req.ok
}

//scan the map. do is a function which operate all of the key and value in map
func (s *smap) EachItem(do func(interface{}, interface{})) {
	req := &scanReq{
		do:   do,
		brea: 0,
		done: make(chan bool),
	}
	s.scanSig <- req
	<-req.done
}

//scan the map util function 'do' returns true. do is a function which operate all of the key and value in map
func (s *smap) EachItemBreak(do func(interface{}, interface{}) bool, condition bool) {
	req := &scanReq{
		doWithBreak: do,
		brea:        1,
		done:        make(chan bool),
	}
	s.scanSig <- req
	<-req.done
}

//Exists checks whether the key which provided is exists in map
func (s *smap) Exists(key interface{}) bool {
	if _,found := s.Get(key); found {
		return true
	}
	return false
}

