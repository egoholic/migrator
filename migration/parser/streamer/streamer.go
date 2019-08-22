package streamer

type Streamer []rune

func New(runes []rune) Streamer {
	return Streamer(runes)
}

func (runes Streamer) Stream(ch chan<- rune) {
	for _, r := range runes {
		ch <- r
	}
}
