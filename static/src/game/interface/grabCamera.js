function GrabCamera(scene) {
  if (scene.game.input.activePointer.rightButtonDown()) {
    scene.cameras.main.stopFollow();
    if (scene.game.origDragPoint) {
      // move the camera by the amount the mouse has moved since last update
      scene.cameras.main.scrollX +=
        scene.game.origDragPoint.x - scene.game.input.activePointer.position.x;
      scene.cameras.main.scrollY +=
        scene.game.origDragPoint.y - scene.game.input.activePointer.position.y;
    } // set new drag origin to current position
    scene.game.origDragPoint = scene.game.input.activePointer.position.clone();
  } else {
    scene.game.origDragPoint = null;
  }
}

export {GrabCamera}
