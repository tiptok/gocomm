package xbyte

// BytePoolCap  sync.Pool []byte 实现（性能相差不多）
type BytePoolCap struct {
	c    chan []byte
	w    int // initialize []byte len
	wcap int // initialize []byte cap len
}

func (bp *BytePoolCap) Get() (b []byte) {
	select {
	case b = <-bp.c:
	default:
		if bp.wcap > 0 {
			b = make([]byte, bp.w, bp.wcap)
		} else {
			b = make([]byte, bp.w)
		}
	}
	return
}

func (bp *BytePoolCap) Put(b []byte) {
	select {
	case bp.c <- b:
		// buffer went back into pool
	default:
		// buffer didn't go back into pool, just discard
	}
}

func NewBytePoolCap(maxSize int, width int, capwidth int) (bp *BytePoolCap) {
	return &BytePoolCap{
		c:    make(chan []byte, maxSize),
		w:    width,
		wcap: capwidth,
	}
}
