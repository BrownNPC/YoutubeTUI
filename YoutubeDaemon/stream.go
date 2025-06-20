package daemon

// func beginStreaming(t Track, streamingURL string) {
// 	mu.Lock()
// 	if Player.close != nil {
// 		Player.close()
// 	}
// 	mu.Unlock()

// 	resp, err := http.Get(streamingURL)
// 	if err != nil {
// 		Player.Events <- EventErr(errors.Join(errors.New("failed to fetch streaming url"), err))
// 		return
// 	}
// 	f := httprs.NewHttpReadSeeker(resp)

// 	reader, _, err := newWebMReader(f)
// 	if err != nil {
// 		Player.Events <- EventErr(errors.Join(errors.New("failed to decode webm stream"), err))
// 		go beginStreaming(t, streamingURL)
// 		return
// 	}
// 	plr := otoCtx.NewPlayer(reader)
// 	plr.Play()
// 	mu.Lock()
// 	defer mu.Unlock()
// 	Player.plr = plr
// 	Player.close = func() {
// 		plr.Close()
// 		f.Close()
// 		reader.Close()
// 	}
// 	Player.Events <- EventTrackStarted(t)
// }
