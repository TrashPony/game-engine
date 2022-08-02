package behavior_rule

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	"sync"
)

type BehaviorRules struct {
	Rules []*BehaviorRule `json:"rules"`
	Meta  *Meta           `json:"meta"`
	Key   string          `json:"key"`
}

type BehaviorRule struct {
	Action   string        `json:"action"`
	Meta     *Meta         `json:"meta"`
	PassRule *BehaviorRule `json:"access_rule"`
	StopRule *BehaviorRule `json:"stop_rule"`
	Exit     bool          `json:"exit"`
}

type Patrol struct {
	Path       []*coordinate.Coordinate // маршрут состоящий из x,y и радиуса
	ToIDIndex  int                      // текущая цель куда двигается группа, по достижение происходит ++, если -1
	AutoChange bool                     // true, ToIDIndex меняется автоматически как только достигнута цель
}

type SubGroup struct {
	members map[int]bool
	mx      sync.RWMutex
}

func (s *SubGroup) AddMember(id int) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.members == nil {
		s.members = make(map[int]bool)
	}

	s.members[id] = true
}

func (s *SubGroup) RemoveMember(id int) {
	s.mx.Lock()
	defer s.mx.Unlock()

	delete(s.members, id)
}

func (s *SubGroup) GetMembers() []int {
	s.mx.RLock()
	defer s.mx.RUnlock()

	ids := make([]int, 0, len(s.members))
	for id := range s.members {
		ids = append(ids, id)
	}

	return ids
}

func (s *SubGroup) GetCountMembers() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return len(s.members)
}

type Meta struct {
	ID       int       `json:"ID"`
	Type     string    `json:"type"`
	BaseID   int       `json:"base_id"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
	Radius   int       `json:"radius"`
	Patrol   *Patrol   `json:"patrol"`
	Role     string    `json:"role"`
	SubGroup *SubGroup `json:"sub_group"`
}
