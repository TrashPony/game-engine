function AddMoveBufferData(data, object) {
  data.serverPos = {x: data.x, y: data.y}

  if (!object.bufferMoveTick) {
    object.bufferMoveTick = [];
  }

  while (object.bufferMoveTick.length > 1) {
    object.bufferMoveTick.shift();
  }

  object.bufferMoveTick.push(data);
  // показываем прошлый тик что бы компенсировать сетевые лаги
  if (!object.updaterPos && object.bufferMoveTick.length >= 0) {
    object.updaterPos = true;
  }
}

export {AddMoveBufferData}
