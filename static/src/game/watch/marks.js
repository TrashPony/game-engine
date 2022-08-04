import {gameStore} from "../store"
import {Scene} from "../create"
import {GetGlobalPos} from "../map/gep_global_pos";
import {MoveSprite} from "../utils/move_sprite";

function CreateMark(id_mark, markType, x, y) {
  if (!gameStore.gameReady) return;

  let oldMark = gameStore.radar_marks[id_mark];
  if (oldMark) {
    RemoveMark(id_mark)
  }

  let mark = {uuid: id_mark, type: markType};

  let pos = GetGlobalPos(x, y, gameStore.map.id);
  mark.sprite = Scene.make.image({
    x: pos.x,
    y: pos.y,
    key: "sprites",
    frame: markType,
    add: true
  });
  mark.sprite.setOrigin(0.5);
  mark.sprite.setScale(0.2);
  mark.sprite.setDepth(999);

  gameStore.radar_marks[id_mark] = mark;

  return mark
}

function RemoveAllMark() {
  for (let id_mark in gameStore.radar_marks) {
    RemoveMark(id_mark);
  }
}

function RemoveMark(id_mark) {
  let mark = gameStore.radar_marks[id_mark];
  if (mark) {
    mark.sprite.destroy();
    delete gameStore.radar_marks[id_mark];
  }
}

function MoveMark(data) {
  let path = data;

  let mark = gameStore.radar_marks[data.mu];
  if (!mark) {
    return
  }

  let pos = GetGlobalPos(path.x, path.y, gameStore.map.id);
  MoveSprite(mark.sprite, pos.x, pos.y, path.ms, null);
}

export {CreateMark, RemoveMark, MoveMark, RemoveAllMark}
