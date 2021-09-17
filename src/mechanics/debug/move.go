package debug

import (
	"sync"
	"time"
)

var Store = newStore()

func newStore() *MessageStore {
	return &MessageStore{
		messages: make([]*Message, 0),
		Move:     false,
		MoveInit: false,

		MoveEndPoint: false,

		AStartNeighbours: false,
		AStartResult:     false,

		RegionFindDebug: false,
		RegionResult:    false,

		HandAlgorithm: false,

		SearchCollisionLine:       false,
		SearchCollisionLineResult: false,
		SearchCollisionLineStep:   false,

		UnitUnitCollision: false,

		Collisions:     false,
		WeaponFirePos:  false,
		FlyBulletLevel: false,
	}
}

type MessageStore struct {
	messages                  []*Message
	mx                        sync.Mutex
	Move                      bool
	MoveInit                  bool
	AStartNeighbours          bool
	AStartResult              bool
	RegionFindDebug           bool
	RegionResult              bool
	HandAlgorithm             bool
	SearchCollisionLineResult bool
	SearchCollisionLine       bool
	MoveEndPoint              bool
	UnitUnitCollision         bool
	SearchCollisionLineStep   bool
	WeaponFirePos             bool
	Collisions                bool
	FlyBulletLevel            bool
}

type Message struct {
	Type  string
	Color string
	X     int
	Y     int
	ToX   int
	ToY   int
	Size  int
	MapID int
	MS    int64
	Text  string
}

func (s *MessageStore) AddMessage(msgType, color string, x, y, toX, toY, size, mpId int, ms int64) {
	s.mx.Lock()

	s.messages = append(s.messages, &Message{
		Type:  msgType,
		Color: color,
		X:     x,
		Y:     y,
		ToX:   toX,
		ToY:   toY,
		Size:  size,
		MapID: mpId,
		MS:    ms})

	s.mx.Unlock()

	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (s *MessageStore) AddMessageText(msgType, color string, x, y, toX, toY, size, mpId int, ms int64, text string) {
	s.mx.Lock()

	s.messages = append(s.messages, &Message{
		Type:  msgType,
		Color: color,
		X:     x,
		Y:     y,
		ToX:   toX,
		ToY:   toY,
		Size:  size,
		MapID: mpId,
		MS:    ms,
		Text:  text})

	s.mx.Unlock()

	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (s *MessageStore) GetAllMessages() []*Message {
	s.mx.Lock()
	defer s.mx.Unlock()

	result := s.messages
	s.messages = make([]*Message, 0)

	return result
}
