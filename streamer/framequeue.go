package streamer

import "io"

type opusFrameQueue struct {
	frames chan []byte
	closed chan struct{}
}

func newOpusFrameQueue(buffer int) *opusFrameQueue {
	return &opusFrameQueue{
		frames: make(chan []byte, buffer),
		closed: make(chan struct{}),
	}
}

func (ofq *opusFrameQueue) ProvideOpusFrame() ([]byte, error) {
	select {
	case frame, ok := <-ofq.frames:
		if !ok {
			return nil, io.EOF
		}
		return frame, nil
	case <-ofq.closed:
		return nil, io.EOF
	}
}

func (ofq *opusFrameQueue) Close() {
	close(ofq.closed)

	// Drain frames so Discord immediately hits EOF
	for {
		select {
		case <-ofq.frames:
		default:
			return
		}
	}
}

func (ofq *opusFrameQueue) Write(p []byte) (int, error) {
	frame := make([]byte, len(p))
	copy(frame, p)

	select {
	case ofq.frames <- frame:
		return len(p), nil
	case <-ofq.closed:
		return 0, io.EOF
	}
}
