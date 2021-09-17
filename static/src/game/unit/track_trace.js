import {Game, Scene} from "../create";
import {Blur} from "../blur";

let i = 0

let pipeLine = null;

let pipelines = {}

function DrawTrackTrace(unit) {

  let trackTrace = Scene.make.image({
    x: unit.sprite.x,
    y: unit.sprite.y,
    frame: unit.body.name + "_bottom",
    key: unit.body.name,
    add: true
  });
  trackTrace.setOrigin(0.5);
  trackTrace.setScale(unit.sprite.scale / 1.5, unit.sprite.scale / 1.1);
  trackTrace.setAngle(unit.sprite.angle);
  trackTrace.setDepth(0);


  let alpha = {alpha: 0.2};
  applicablePipeLine(trackTrace, alpha)
  i++

  let allTime = 1000;

  Scene.tweens.add({
    targets: alpha,
    alpha: 0,
    duration: allTime,
    ease: 'Linear',
    onUpdate: function () {
      applicablePipeLine(trackTrace, alpha.alpha.toFixed(2))
    },
    onComplete: function () {
      trackTrace.destroy();
    }
  })

  trackTrace.setBlendMode('NORMAL'); //ERASE, MULTIPLY
}

function applicablePipeLine(trackTrace, alpha) {
  if (!pipelines[String(alpha)]) {
    pipeLine = Scene.renderer.pipelines.add(String(alpha), new Blur(Game))
    pipeLine.set1f('alpha', alpha)
    pipelines[String(alpha)] = pipeLine
  } else {
    trackTrace.setPipeline(String(alpha));
  }
}

export {DrawTrackTrace}
