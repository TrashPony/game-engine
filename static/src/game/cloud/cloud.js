import {gameStore} from "../store";
import {Scene} from "../create";
import {GetGlobalPos} from "../map/gep_global_pos";
import {MoveSprite} from "../utils/move_sprite";

function CreateCloud(cloudState) {
  if (gameStore.gameReady) {
    let cloud = gameStore.clouds[cloudState.id];
    if (!cloud) cloud = createCloud(cloudState);

    let pos = GetGlobalPos(cloudState.x, cloudState.y, cloudState.m);

    MoveSprite(cloud, pos.x, pos.y, 1000, null);
  }
}

function createCloud(cloudState) {

  let pos = GetGlobalPos(cloudState.x, cloudState.y, cloudState.m);

  let cloud = Scene.make.sprite({
    x: pos.x,
    y: pos.y,
    key: "cloud" + cloudState.type_id,
    add: true
  });

  cloud.setOrigin(0.5);
  cloud.setAngle(cloudState.r);
  cloud.setAlpha(cloudState.a);
  cloud.setDepth(1000);

  cloud.id = cloudState.id;

  gameStore.clouds[cloud.id] = cloud;

  return cloud
}

function removeCloud(id) {
  let cloud = gameStore.clouds[id];
  if (cloud) {
    cloud.destroy();
    delete gameStore.clouds[id];
  }
}

export {CreateCloud, removeCloud}
