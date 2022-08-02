package generate_ids

import "sync"

var mx sync.Mutex

var generateUnitID int

func GetBotID() int {
	mx.Lock()
	defer mx.Unlock()

	generateUnitID -= 1
	return generateUnitID
}

var markIDGenerate = 0

func GetMarkID() int {
	mx.Lock()
	defer mx.Unlock()

	markIDGenerate++
	return markIDGenerate
}

var playerID int

func GetPlayerFakeID() int {
	mx.Lock()
	defer mx.Unlock()
	playerID++
	return playerID
}

var mapID int

func GetMapID() int {
	mx.Lock()
	defer mx.Unlock()
	mapID--
	return mapID
}
