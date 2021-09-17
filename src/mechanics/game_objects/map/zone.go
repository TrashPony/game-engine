package _map

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/obstacle_point"
	"sync"
)

type Zone struct {
	Size      int                                      `json:"size"`
	DiscreteX int                                      `json:"discrete_x"`
	DiscreteY int                                      `json:"discrete_y"`
	Obstacle  map[string]*obstacle_point.ObstaclePoint `json:"obstacle"`
	Rect      *game_math.Polygon                       `json:"rect"`
	mx        sync.RWMutex                             `json:"-"`
}

func (mp *Map) getZone(x, y int) *Zone {

	if x < 0 || y < 0 {
		return nil
	}

	if len(mp.GeoZones)-1 < x {
		return nil
	}

	if len(mp.GeoZones[x])-1 < y {
		return nil
	}

	// записывать в зоны в игре не надо тк что чтение безопасно
	return mp.GeoZones[x][y]
}

func (mp *Map) GetObstaclesByZone(x, y int) <-chan *obstacle_point.ObstaclePoint {
	zone := mp.getZone(x/_const.DiscreteSize, y/_const.DiscreteSize)
	if zone == nil {
		oChan := make(chan *obstacle_point.ObstaclePoint, 1)
		close(oChan)
		return oChan
	} else {
		return zone.GetObstaclesChan()
	}
}

func (mp *Map) GetObstaclesByZoneUnsafe(x, y int) (map[string]*obstacle_point.ObstaclePoint, *sync.RWMutex) {
	zone := mp.getZone(x/_const.DiscreteSize, y/_const.DiscreteSize)
	if zone == nil {
		return nil, nil
	} else {
		return zone.GetFastUnsafeRange()
	}
}

type Obstacler interface {
	GetX() int
	GetY() int
	GetGeoData() []*obstacle_point.ObstaclePoint
}

func (mp *Map) AddGeoDataObjectsToZone(obj Obstacler) {

	zone := mp.getZone(obj.GetX()/_const.DiscreteSize, obj.GetY()/_const.DiscreteSize)
	if zone == nil {
		return
	}

	zones := zone.GetNeighboursZone(mp)
	zones = append(zones, zone)

	for _, zone := range zones {
		for _, gPoint := range obj.GetGeoData() {
			if zone.Rect.DetectCollisionRectToCircle(&game_math.Point{X: float64(gPoint.GetX()), Y: float64(gPoint.GetY())}, gPoint.GetRadius()) {
				zone.AddObstacle(gPoint)
			}
		}
	}
}

func (mp *Map) RemoveGeoDataObjectsToZone(obj Obstacler) {

	zone := mp.getZone(obj.GetX()/_const.DiscreteSize, obj.GetY()/_const.DiscreteSize)
	if zone == nil {
		return
	}

	zones := zone.GetNeighboursZone(mp)
	zones = append(zones, zone)

	for _, zone := range zones {
		for _, obstacle := range obj.GetGeoData() {
			zone.RemoveObstacle(obstacle)
		}
	}
}

func (z *Zone) GetNeighboursZone(mp *Map) []*Zone {

	neighboursZones := make([]*Zone, 0)
	checkRegion := func(x, y int) {
		if x >= 0 && x < len(mp.GeoZones) && y >= 0 && y < len(mp.GeoZones[x]) {
			if mp.GeoZones[x][y] != nil {
				neighboursZones = append(neighboursZones, mp.GeoZones[x][y])
			}
		}
	}

	//строго лево
	checkRegion(z.DiscreteX-1, z.DiscreteY)
	//строго право
	checkRegion(z.DiscreteX+1, z.DiscreteY)
	//верх центр
	checkRegion(z.DiscreteX, z.DiscreteY-1)
	//низ центр
	checkRegion(z.DiscreteX, z.DiscreteY+1)

	//верх лево
	checkRegion(z.DiscreteX-1, z.DiscreteY-1)
	//верх право
	checkRegion(z.DiscreteX+1, z.DiscreteY-1)
	//низ лево
	checkRegion(z.DiscreteX-1, z.DiscreteY+1)
	//низ паво
	checkRegion(z.DiscreteX+1, z.DiscreteY+1)

	return neighboursZones
}

func (z *Zone) AddObstacle(obstacle *obstacle_point.ObstaclePoint) {
	z.mx.Lock()
	defer z.mx.Unlock()

	if z.Obstacle == nil {
		z.Obstacle = make(map[string]*obstacle_point.ObstaclePoint, 0)
	}

	z.Obstacle[obstacle.GetKey()] = obstacle
}

func (z *Zone) RemoveObstacle(obstacle *obstacle_point.ObstaclePoint) {
	z.mx.Lock()
	defer z.mx.Unlock()

	if z.Obstacle == nil {
		return
	}

	delete(z.Obstacle, obstacle.GetKey())
}

func (z *Zone) GetObstaclesChan() <-chan *obstacle_point.ObstaclePoint {

	z.mx.RLock()

	obstacles := make(chan *obstacle_point.ObstaclePoint, len(z.Obstacle))

	go func() {
		defer func() {
			z.mx.RUnlock()
			close(obstacles)
		}()

		for _, obstacle := range z.Obstacle {
			obstacles <- obstacle
		}
	}()

	return obstacles
}

func (z *Zone) GetFastUnsafeRange() (map[string]*obstacle_point.ObstaclePoint, *sync.RWMutex) {
	z.mx.RLock()
	return z.Obstacle, &z.mx
}
