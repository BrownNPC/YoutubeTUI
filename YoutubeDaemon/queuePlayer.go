package daemon

func playQueue() {
	if Player.queue == nil {
		panic("assert: asked to play queue with nil queue")
	}
	if len(Player.queue) == 0 {
		return
	}
	t := &Player.queue[0]
	go PlayTrack(t)
}
