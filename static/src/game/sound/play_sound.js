import {Scene} from "../create";
import store from "../../store/store";

let soundStore = {};

function PlayPositionSound(sounds, params, x, y, noDestroy, volumeK = 1) {

  let volume = GetVolume(x, y, volumeK)
  if (!volume) {
    return
  }

  if (!params) params = {};
  params.volume = volume;

  let soundName = sounds[getRandomInt(sounds.length)]

  let sound;
  if (soundStore.hasOwnProperty(soundName)) {
    sound = soundStore[soundName].shift()
  }

  if (!sound) sound = Scene.sound.add(soundName);

  if (!noDestroy) {
    sound.on('complete', function () {
      if (!soundStore.hasOwnProperty(soundName)) {
        soundStore[soundName] = [];
      }

      soundStore[soundName].push(sound)
    });
  }
  sound.play(params);
  return sound;
}

function GetVolume(x, y, volumeK = 1) {
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

  return (percent / 100 * store.getters.getSettings.SFXVolume) * volumeK
}

function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}

export {PlayPositionSound, GetVolume}
