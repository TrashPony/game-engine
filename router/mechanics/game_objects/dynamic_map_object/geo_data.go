package dynamic_map_object

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/getlantern/deepcopy"
	"strconv"
)

func (o *Object) CalculateScale() {

	// подравниваем хп под размер туши
	percentLife := (float64(o.GetHP()) / float64(o.MaxHP)) * 100               // сохраняем хп в процентном соотношение
	o.MaxHP = int(float64(o.TypeMaxHP) * float64(o.GetScale()) / float64(100)) // смотрим размер обьекта от оригинала и высчитываем его макс хп
	o.SetHP(int(float64(o.MaxHP) * (percentLife / float64(100))))              // востанавливаем хп в % соотношение

	// подравниваем тени под тушу
	o.XShadowOffset = int(float64(o.TypeXShadowOffset) * (float64(o.GetScale()) / 100))
	o.YShadowOffset = int(float64(o.TypeYShadowOffset) * (float64(o.GetScale()) / 100))

	o.SetGeoData()
}

func (o *Object) SetGeoData() {

	if o.GetPhysicalModel().GeoData == nil {
		err := deepcopy.Copy(&o.GetPhysicalModel().GeoData, &o.TypeGeoData)
		if err != nil {
			println(err.Error())
		}
	}

	o.GetPhysicalModel().SetHeight(o.HeightType * (float64(o.GetScale()) / 100))

	for i, geoPoint := range o.GetPhysicalModel().GeoData {

		if geoPoint == nil {
			continue
		}

		geoPoint.SetParentType("object")
		geoPoint.SetKey(geoPoint.GetParentType() + strconv.Itoa(o.GetID()) + strconv.Itoa(o.TypeGeoData[i].GetX()) + strconv.Itoa(o.TypeGeoData[i].GetY()))
		geoPoint.SetParentID(o.ID)

		geoPoint.SetHeight(o.GetPhysicalModel().GetHeight())

		// применяем размер обьекта к геодате
		geoPoint.SetRadius(int(float64(o.TypeGeoData[i].GetRadius()) * (float64(o.GetScale()) / 100)))

		// получаем позицию гео точки на карте
		x := int(float64(o.TypeGeoData[i].GetX())*(float64(o.GetScale())/100)) + o.GetPhysicalModel().GetX()
		y := int(float64(o.TypeGeoData[i].GetY())*(float64(o.GetScale())/100)) + o.GetPhysicalModel().GetY()

		// поворачиваем геодату на угол обьекта
		newX, newY := game_math.RotatePoint(float64(x), float64(y), float64(o.GetPhysicalModel().GetX()), float64(o.GetPhysicalModel().GetY()), o.GetPhysicalModel().GetRotate())

		geoPoint.SetX(int(newX))
		geoPoint.SetY(int(newY))
	}

	o.CacheGeoData = make([]byte, 0)
	for _, geoPoint := range o.physicalModel.GeoData {
		o.CacheGeoData = append(o.CacheGeoData, binary_msg.GetIntBytes(int(geoPoint.X))...)
		o.CacheGeoData = append(o.CacheGeoData, binary_msg.GetIntBytes(int(geoPoint.Y))...)
		o.CacheGeoData = append(o.CacheGeoData, binary_msg.GetIntBytes(int(geoPoint.Radius))...)
	}
}
