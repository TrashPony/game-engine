import {gameStore} from "../store";
import {GetGlobalPos} from "../map/gep_global_pos";
import {MoveSprite} from "../utils/move_sprite";
import {SetBodyAngle} from "./rotate";
import {Scene} from "../create";
import {DrawTrackTrace} from "./track_trace";
import {GetVolume, PlayPositionSound} from "../sound/play_sound";

function AddUnitMoveBufferData(data) {

  let unit = gameStore.units[data.id];

  if (unit) {
    if (!unit.bufferMoveTick) {
      unit.bufferMoveTick = [];
    }

    while (unit.bufferMoveTick.length > 5) {
      unit.bufferMoveTick.shift();
    }

    unit.bufferMoveTick.push(data);

    // показываем прошлый тик что бы компенсировать сетевые лаги
    if (!unit.updaterPos && unit.bufferMoveTick.length > 0) {
      unit.updaterPos = true;
    }
  }
}

function MoveTo(unit, ms) {

  if (!gameStore.units.hasOwnProperty(unit.id)) return;

  let path = unit.bufferMoveTick.shift();
  if (!path) {
    return;
  }

  unit.x = path.x;
  unit.y = path.y;

  let pos = GetGlobalPos(path.x, path.y, path.m);
  MoveSprite(unit.sprite, pos.x, pos.y, ms);
  SetBodyAngle(unit, path.r, ms, true, Scene.shadowXOffset * 3);

  unit.speedDirection = path.d;
  unit.speed = path.s;
  unit.angularVelocity = path.av;
  unit.animateSpeed = path.a;
  unit.rotate = path.r;

  AnimationMove(unit);
}

let engagesSound = {
  apd_light: "engine_6",
};

function AnimationMove(unit) {
  if (unit.speed && unit.speed > 0 && unit.animateSpeed) {

    unit.sprite.bodyBottomLeft.anims.resume();
    unit.sprite.bodyBottomRight.anims.resume();

    // направление гусей
    if (unit.speedDirection && !unit.sprite.bodyBottomLeft.anims.forward) {
      unit.sprite.bodyBottomLeft.anims.reverse()
    }
    if (!unit.speedDirection && unit.sprite.bodyBottomLeft.anims.forward) {
      unit.sprite.bodyBottomLeft.anims.reverse()
    }

    if (unit.speedDirection && !unit.sprite.bodyBottomRight.anims.forward) {
      unit.sprite.bodyBottomRight.anims.reverse()
    }
    if (!unit.speedDirection && unit.sprite.bodyBottomRight.anims.forward) {
      unit.sprite.bodyBottomRight.anims.reverse()
    }

    speedBody(unit.speed, unit.sprite.bodyBottomLeft.anims);
    let animSpeed = speedBody(unit.speed, unit.sprite.bodyBottomRight.anims);
    DrawTrackTrace(unit)

    if (!unit.engigeSound) {
      unit.engigeSound = PlayPositionSound([engagesSound[unit.body.name]], {loop: true}, unit.sprite.x, unit.sprite.y);
    } else {
      let volume = GetVolume(unit.sprite.x, unit.sprite.y);
      unit.engigeSound.setVolume(volume * (1 - (animSpeed / 75)));
      unit.engigeSound.resume()
    }

  } else {
    if (unit.angularVelocity !== 0) {

      unit.sprite.bodyBottomLeft.anims.resume();
      unit.sprite.bodyBottomRight.anims.resume();

      if (unit.angularVelocity > 0) {
        if (!unit.sprite.bodyBottomLeft.anims.forward) {
          unit.sprite.bodyBottomLeft.anims.reverse()
        }
        if (unit.sprite.bodyBottomRight.anims.forward) {
          unit.sprite.bodyBottomRight.anims.reverse()
        }
      } else {
        if (unit.sprite.bodyBottomLeft.anims.forward) {
          unit.sprite.bodyBottomLeft.anims.reverse()
        }
        if (!unit.sprite.bodyBottomRight.anims.forward) {
          unit.sprite.bodyBottomRight.anims.reverse()
        }
      }

      speedBody(unit.angularVelocity, unit.sprite.bodyBottomLeft.anims);
      let animSpeed = speedBody(unit.angularVelocity, unit.sprite.bodyBottomRight.anims);
      DrawTrackTrace(unit)

      if (!unit.engigeSound) {
        unit.engigeSound = PlayPositionSound([engagesSound[unit.body.name]], {loop: true}, unit.sprite.x, unit.sprite.y);
      } else {
        let volume = GetVolume(unit.sprite.x, unit.sprite.y);
        unit.engigeSound.setVolume(volume * (1 - (animSpeed / 200)));
        unit.engigeSound.resume()
      }

    } else {
      if (unit.engigeSound) unit.engigeSound.pause()
      unit.sprite.bodyBottomLeft.anims.pause();
      unit.sprite.bodyBottomRight.anims.pause();
    }
  }
}

function speedBody(speed, anim) {

  if (speed < 0) {
    speed = speed * -1;
  }

  let ms = 110;

  if (speed < 10) {
    ms = 110
  } else if (speed >= 10 && speed < 15) {
    ms = 105
  } else if (speed >= 15 && speed < 20) {
    ms = 100
  } else if (speed >= 20 && speed < 25) {
    ms = 95
  } else if (speed >= 25 && speed < 30) {
    ms = 90
  } else if (speed >= 30 && speed < 40) {
    ms = 85
  } else if (speed >= 40 && speed < 50) {
    ms = 80
  } else if (speed >= 50 && speed < 60) {
    ms = 75
  } else if (speed >= 60 && speed < 70) {
    ms = 70
  } else if (speed >= 70 && speed < 80) {
    ms = 65
  } else if (speed >= 80 && speed < 90) {
    ms = 40
  } else if (speed >= 90 && speed < 100) {
    ms = 35
  } else if (speed >= 100 && speed < 110) {
    ms = 30
  } else if (speed >= 110 && speed < 120) {
    ms = 25
  } else if (speed >= 120 && speed < 130) {
    ms = 20
  } else if (speed >= 130 && speed < 140) {
    ms = 15
  } else if (speed >= 140) {
    ms = 10
  }

  anim.msPerFrame = ms
  return ms
}

export {MoveTo, AddUnitMoveBufferData}
