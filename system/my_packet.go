package system

type MyPacket struct {
	cyclesNeeded  int
	fmID          int
	streamNum     int
	notifier      Callable
	sender        Callable
	preferredVnet int
	preferredDest int
	preferredSrc  int
	msgSize       uint64
	readyTime     uint64
	flowID        int
	parentFlowID  int
	childFlowID   int
	channelID     int
	chunkID       int
}

func (mp *MyPacket) SetNotifier(c Callable) {
	mp.notifier = c
}

func (mp *MyPacket) SetSender(c Callable) {
	mp.sender = c
}
func NewMyPacket(preferredVnet int, preferredSrc int, preferredDest int) *MyPacket {
	return &MyPacket{
		preferredVnet: preferredVnet,
		preferredSrc:  preferredSrc,
		preferredDest: preferredDest,
	}
}
