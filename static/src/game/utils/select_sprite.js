import {Scene} from "../create";

let structSelect = [];

function SelectSprite(x, y, callback, tint, tintOver, radius) {

  //let radius = sprite.displayHeight + 5;

  let selectSprite = Scene.make.sprite({
    x: x,
    y: y,
    key: "select_sprite",
    add: true
  });

  selectSprite.setDepth(0);
  selectSprite.setOrigin(0.5);
  selectSprite.setDisplaySize(radius, radius);
  selectSprite.setTint(tint);

  structSelect.push(selectSprite);

  Scene.tweens.add({
    targets: selectSprite,
    props: {
      displayHeight: {value: radius + 5, duration: 1000, ease: 'Cubic.easeIn'},
      displayWidth: {value: radius + 5, duration: 1000, ease: 'Cubic.easeIn'},
    },
    repeat: -1,
    yoyo: true,
  });

  selectSprite.setInteractive({
    pixelPerfect: true,
    alphaTolerance: 1
  });

  // selectSprite.on('pointerover', function () {
  //   selectSprite.setTint(tintOver);
  // });
  //
  // selectSprite.on('pointerout', function () {
  //   selectSprite.setTint(tint);
  // });
  //
  // selectSprite.on('pointerdown', function (pointer) {
  //   if (callback) callback(pointer);
  // });

  return selectSprite
}

export {SelectSprite}
