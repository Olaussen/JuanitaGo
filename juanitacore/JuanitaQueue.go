package juanitacore

import (
	"juanitaGo/structs"
	"math/rand"
)

type JuanitaQueue struct {
	tracks []structs.JuanitaSearch
}

func NewJuanitaQueue() JuanitaQueue {
	return JuanitaQueue{tracks: make([]structs.JuanitaSearch, 0)}
}

func NewJuanitaQueueWithTracks(tracks []structs.JuanitaSearch) JuanitaQueue {
	return JuanitaQueue{tracks: tracks}
}

func (queue *JuanitaQueue) EnqueueBack(track structs.JuanitaSearch) {
	queue.tracks = append(queue.tracks, track)
}

func (queue *JuanitaQueue) EnqueueFirst(track structs.JuanitaSearch) {
	queue.tracks = append([]structs.JuanitaSearch{track}, queue.tracks...)
}

func (queue *JuanitaQueue) Dequeue() *structs.JuanitaSearch {
	track := queue.tracks[0]
	queue.tracks = queue.tracks[1:]
	return &track
}

func (queue *JuanitaQueue) SkipTo(index int) *structs.JuanitaSearch {
	track := queue.tracks[index]
	queue.tracks = queue.tracks[index+1:]
	return &track
}

func (queue *JuanitaQueue) Shuffle() {
	for i := range queue.tracks {
		j := rand.Intn(i + 1)
		queue.tracks[i], queue.tracks[j] = queue.tracks[j], queue.tracks[i]
	}
}

func (queue *JuanitaQueue) IsEmpty() bool {
	return len(queue.tracks) == 0
}
