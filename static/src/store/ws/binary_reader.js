import {intFromBytes} from "../../utils";
import {ServerGunDesignator} from "../../game/interface/server_gun_designator";
import {AddUnitMoveBufferData} from "../../game/unit/move";
import {RotateUnitGun} from "../../game/unit/rotate";
import {ExplosionBullet} from "../../game/weapon/explosion";
import {AddBulletMoveBufferData} from "../../game/bullet/fly";
import {FireWeapon} from "../../game/weapon/fire";
import {gameStore} from "../../game/store";
import {ObjectDead} from "../../game/map/object_dead";
import {RotateTurretGun} from "../../game/map/structures/turret";
import {CreateRadarObject} from "../../game/watch/create";
import {UpdateObject} from "../../game/watch/update";
import {AddObjectMoveBufferData} from "../../game/map/move_object";
import {RemoveRadarObject} from "../../game/watch/remove";

export const SourceItemBin = {
  1: "squadInventory", 2: "Constructor", 3: "empty", 4: "box", 5: "storage",
}

const MarksTypes = {
  1: "fly", 2: "ground", 3: "structure", 4: "resource", 5: "bullet",
}

function parseMegaPackData(data, store) {

  // [1[eventID],
  //      4[data_size], data[data],
  //      ...
  //      4[data_size], data[data],

  let unitMoveSize = intFromBytes(data.slice(1, 5))
  let stopByte = 5 + unitMoveSize;
  BinaryReader(data.subarray(5, 5 + unitMoveSize), store)

  for (; stopByte < data.length;) {
    let subData = intFromBytes(data.slice(stopByte, stopByte + 4))
    stopByte = 4 + stopByte;
    BinaryReader(data.subarray(stopByte, stopByte + subData), store)
    stopByte = subData + stopByte;
  }
}

function BinaryReader(msgBytes, store) {

  if (msgBytes[0] === 100) {
    parseMegaPackData(msgBytes, store)
    return
  }

  if (!gameStore.gameReady) return;

  if (msgBytes[0] === 1) {
    // [eventID| 4[unitID], 4[speed], 4[x], 4[y], 4[z], 4[ms], 4[rotate], 4[angularVelocity], 4[mpID] 1[animate], 1[direction]]
    //      [1 | 0, 0, 2, 96, 0, 0, 0, 0, 0, 0, 7, 251, 0, 0, 4, 136, 0, 0, 0, 100, 0, 0, 0, 64, 0, 0, 0, 87, 0, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0]
    for (let i = 0; i < msgBytes.length; i += 37) {
      let event = msgBytes.slice(i, i + 37);

      AddUnitMoveBufferData({
        id: intFromBytes(event.slice(1, 5)),
        s: intFromBytes(event.slice(5, 9)),
        x: intFromBytes(event.slice(9, 13)),
        y: intFromBytes(event.slice(13, 17)),
        z: intFromBytes(event.slice(17, 21)),
        ms: intFromBytes(event.slice(21, 22)),
        r: intFromBytes(event.slice(22, 26)),
        av: intFromBytes(event.slice(26, 30)) * 0.001,
        a: intFromBytes(event.slice(30, 31)),
        d: intFromBytes(event.slice(31, 32)),
        sky: intFromBytes(event.slice(32, 33)),
        ka: intFromBytes(event.slice(33, 34)),
        kd: intFromBytes(event.slice(34, 35)),
        kw: intFromBytes(event.slice(35, 36)),
        ksp: intFromBytes(event.slice(36, 37)),
      })
    }
  }

  if (msgBytes[0] === 3) {
    for (let i = 0; i < msgBytes.length; i += 21) {
      let event = msgBytes.slice(i, i + 21);
      let msg = {
        id: intFromBytes(event.slice(1, 5)),
        x: intFromBytes(event.slice(5, 9)),
        y: intFromBytes(event.slice(9, 13)),
        ms: intFromBytes(event.slice(13, 17)),
        r: intFromBytes(event.slice(17, 21)),
      }

      AddObjectMoveBufferData(msg)
    }
  }

  if (msgBytes[0] === 6) {
    // [1[eventID] 4[typeID], 4[id], 4[x], 4[y], 4[z], 4[ms], 4[rotate], 4[mpID]]
    for (let i = 0; i < msgBytes.length; i += 23) {
      let event = msgBytes.slice(i, i + 23);

      let msg = {
        type_id: intFromBytes(event.slice(1, 2)),
        id: intFromBytes(event.slice(2, 6)),
        x: intFromBytes(event.slice(6, 10)),
        y: intFromBytes(event.slice(10, 14)),
        z: intFromBytes(event.slice(14, 18)),
        ms: intFromBytes(event.slice(18, 19)),
        r: intFromBytes(event.slice(19, 23)),
      }

      AddBulletMoveBufferData(msg)
    }
  }

  if (msgBytes[0] === 7) {
    // [1[eventID] 4[typeID], 4[x], 4[y], 4[z], 4[mpID]]
    for (let i = 0; i < msgBytes.length; i += 14) {
      let event = msgBytes.slice(i, i + 14);

      let msg = {
        type_id: intFromBytes(event.slice(1, 2)),
        x: intFromBytes(event.slice(2, 6)),
        y: intFromBytes(event.slice(6, 10)),
        z: intFromBytes(event.slice(10, 14)),
      }

      ExplosionBullet(msg);
    }
  }

  if (msgBytes[0] === 8) {
    // [1[eventID], 4[ID], 4[ms], 4[rotate]
    for (let i = 0; i < msgBytes.length; i += 11) {
      let event = msgBytes.slice(i, i + 11);
      RotateUnitGun(intFromBytes(event.slice(1, 5)), intFromBytes(event.slice(5, 6)), intFromBytes(event.slice(6, 10)), intFromBytes(event.slice(10, 11)),)
    }
  }

  if (msgBytes[0] === 9) {
    // [1[eventID] 4[typeID], 4[x], 4[y], 4[z], 4[rotate], 4[mpID]]
    for (let i = 0; i < msgBytes.length; i += 19) {
      let event = msgBytes.slice(i, i + 19);
      let msg = {
        type_id: intFromBytes(event.slice(1, 2)),
        x: intFromBytes(event.slice(2, 6)),
        y: intFromBytes(event.slice(6, 10)),
        z: intFromBytes(event.slice(10, 14)),
        r: intFromBytes(event.slice(14, 18)),
        ap: intFromBytes(event.slice(18, 19)),
      }

      FireWeapon(msg)
    }
  }

  if (msgBytes[0] === 10) {
    for (let i = 0; i < msgBytes.length; i += 10) {
      let event = msgBytes.slice(i, i + 10);
      // [1[eventID], 4[ID], 4[ms], 4[rotate]
      let msg = {
        id: intFromBytes(event.slice(1, 5)), ms: intFromBytes(event.slice(5, 6)), r: intFromBytes(event.slice(6, 10)),
      }

      RotateTurretGun(msg)
    }
  }

  if (msgBytes[0] === 13) {
    // [1[eventID], 4[id]
    let msg = {}
    if (msgBytes.length > 1) {
      msg = {
        x: intFromBytes(msgBytes.slice(1, 5)),
        y: intFromBytes(msgBytes.slice(5, 9)),
        rd: intFromBytes(msgBytes.slice(9, 10)),
        ac: intFromBytes(msgBytes.slice(10, 11)),
        aa: intFromBytes(msgBytes.slice(11, 12)),
        ap: intFromBytes(msgBytes.slice(12, 13)),
        r: intFromBytes(msgBytes.slice(13, 14)),
        chase: intFromBytes(msgBytes.slice(14, 15)),
        t: gameStore.mapBinItems[intFromBytes(msgBytes.slice(15, 16))],
        id: intFromBytes(msgBytes.slice(16, 20)),
        fireX: intFromBytes(msgBytes.slice(20, 24)),
        fireY: intFromBytes(msgBytes.slice(24, 28)),
      }
    } else {
      msg.hide = true
    }

    gameStore.fireState.target.x = msg.x;
    gameStore.fireState.target.y = msg.y;
    gameStore.fireState.firePosition.x = msg.fireX;
    gameStore.fireState.firePosition.y = msg.fireY;

    ServerGunDesignator(msg);
  }

  if (msgBytes[0] === 18) {
    // [1[eventID], 4[id], 4[x], 4[y], 4[m], 1[t]]
    for (let i = 0; i < msgBytes.length; i += 14) {
      let event = msgBytes.slice(i, i + 14);
      let msg = {
        id: intFromBytes(event.slice(1, 5)),
        x: intFromBytes(event.slice(5, 9)),
        y: intFromBytes(event.slice(9, 13)),
        t: gameStore.mapBinItems[intFromBytes(event.slice(13, 14))]
      }

      ObjectDead(msg)
    }
  }

  if (msgBytes[0] === 41) {
    for (let i = 0; i < msgBytes.length; i += 6) {
      let event = msgBytes.slice(i, i + 6);

      let mark = {id: intFromBytes(event.slice(1, 5)), to: gameStore.mapBinItems[intFromBytes(event.slice(5, 6))]}
      RemoveRadarObject(mark)
    }
  }

  if (msgBytes[0] === 42) {
    let stopByte = 1
    for (; stopByte < msgBytes.length;) {

      let typeObj = gameStore.mapBinItems[intFromBytes(msgBytes.slice(stopByte, stopByte + 1))]
      stopByte = 1 + stopByte;

      let subData = intFromBytes(msgBytes.slice(stopByte, stopByte + 4))
      stopByte = 4 + stopByte;

      CreateRadarObject(typeObj, msgBytes.subarray(stopByte, stopByte + subData))

      stopByte = subData + stopByte + 1;
    }
  }

  if (msgBytes[0] === 43) {
    let stopByte = 1
    for (; stopByte < msgBytes.length;) {

      let typeObj = gameStore.mapBinItems[intFromBytes(msgBytes.slice(stopByte, stopByte + 1))]
      stopByte = 1 + stopByte;

      let idObj = intFromBytes(msgBytes.slice(stopByte, stopByte + 4))
      stopByte = 4 + stopByte;

      let subData = intFromBytes(msgBytes.slice(stopByte, stopByte + 4))
      stopByte = 4 + stopByte;

      let data = msgBytes.subarray(stopByte, stopByte + subData)

      if (data.length > 0) {
        UpdateObject(typeObj, idObj, data)
      }

      stopByte = subData + stopByte + 1;
    }
  }

  if (msgBytes[0] === 49) {
    gameStore.radarWork = true
  }
}

export {BinaryReader}
