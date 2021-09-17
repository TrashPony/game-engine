import {Scene} from "../create";

function ShortDirectionRotateTween(sprite, needAngle, time) {
  // эта функция ищет оптимальный угол для поворота
  Scene.tweens.add({
    targets: sprite,
    angle: '+=' + Phaser.Math.Angle.ShortestBetween(sprite.angle, Phaser.Math.Angle.WrapDegrees(needAngle)),
    ease: 'Linear',
    duration: time,
    repeat: 0,
  });
}

/*

 tank 360 0вой градус орудия
 97 горизонталь
 40 вертикаль

*/

function PositionAttachSprite(angle, a) {
  // взятие координат угла элипса на изометричной окружности по радиусу
  let b = a; // тут было a/2 но у нас круг :)

  let psi = angle * Math.PI / 180.0;
  let fi = Math.atan2(a * Math.sin(psi), b * Math.cos(psi));
  let x = a * Math.cos(fi);
  let y = b * Math.sin(fi);

  return {x: x, y: y};
}

function VectorToPointBySpeed(x1, y1, x2, y2, speed) {
  // метод находит точку которая находится между точкой x1y1 и x2y2 на дистанции speed(в пикселях) от точки x1y1
  // возвращает инты потому что пикселы не могут быть дробным числом
  let angle = Phaser.Math.angleBetween(x2, y2, x1, y1);
  return {x: x1 + speed * Math.cos(angle), y: y1 + speed * Math.sin(angle)}
}

export {ShortDirectionRotateTween, PositionAttachSprite, VectorToPointBySpeed}
