import {Scene} from "../create";
import store from "../../store/store";

function PlayPositionSound(sounds, params, x, y, noDestroy) {

  let volume = GetVolume(x, y)
  if (!volume) {
    return
  }

  if (!params) params = {};
  params.volume = volume;

  let soundName = sounds[getRandomInt(sounds.length)]
  let sound = Scene.sound.add(soundName);

  if (!noDestroy) {
    sound.on('complete', function () {
      sound.destroy();
    });
  }
  sound.play(params);

  return sound;
}

function GetVolume(x, y) {
  let dist = Phaser.Math.Distance.Between(Scene.cameras.main.worldView.x + Scene.cameras.main.worldView.width / 2, Scene.cameras.main.worldView.y + Scene.cameras.main.worldView.height / 2, x, y);
  let percent = 100 - ((dist / (Scene.cameras.main.worldView.width / 2)) * 100);

  if (percent < 0) {
    return 0;
  }

  if (percent + 30 > 100) {
    percent = 100;
  } else {
    percent = percent + 30;
  }

  return percent / 100 * store.getters.getSettings.SFXVolume
}

function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}

export {PlayPositionSound, GetVolume}
