import {CreateRadarObject, RemoveRadarObject} from './object'
import {UpdateObject} from "./update_objects";

function RadarWork(events) {
  // если обьект создается то метка не нужна
  for (let event of events) {
    if (event.ao === "createObj") {
      CreateRadarObject(event.rm, event.o)
    }

    if (event.ao === "updateObj") {
      UpdateObject(event.rm, event.o)
    }

    if (event.ao === "removeObj") {
      RemoveRadarObject(event.rm)
    }
  }
}

export {RadarWork}
