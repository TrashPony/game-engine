#### Сетевой игровой движок для игр в жанре `action`, `rpg`, `rts`, `tbs`, `tower defense` с видом сверху/изометрией.

Движок не является коробочным решением готовым к использованию и не является конкурентом нормальным движкам. :)

Движок является серверной частью игры решающую вопросы:

- [Сессия, игровое поле, GameLoop](#create-world)
- [Реализация игровых обьектов: "юниты", "строения", "пули"](#game-objects)
- [Бинарный протокол связи между серверов и клиентом](#binary-protocol)
- [Ввод пользователя](#input)
- [Движение обьектов в мире](#move)
- [Баллистика пуль/стрельба](#weapon)
- [Физическую модель обьектов](#physical-model)
- [Обнаружение коллизий, реакция на коллизию](#collision)
- [Реализация "Обзора"](#watch)
- [Создание и обновление игровых обьектов на стороне клиента](#create-update-client)
- [ИИ](#ai)
- [Поиск пути, управление движением ИИ](#ai-support)

Движок не решает вопросы графики и звука, это лежит на плечах клиента. В данном случае реализован клиент на игровом
движке [Phaser3](https://phaser.io/phaser3), в роли клиента может выступать все что умеет в вебсокеты. (или вообще любые
сокеты с небольшой переделкой api). Например хорошим клиентом может стать [Ebitengine](https://ebiten.org/).

#### Игра ["Veliri"](https://yandex.ru/games/app/184316?lang=ru) (Session/MMO/Action) разработаная с помощью этого движка:

([ссылка на видео](https://www.youtube.com/watch?v=7D_ILFRQG2MQ))
[![Watch the video](/readme_assets/img_1.png)](https://www.youtube.com/watch?v=7D_ILFRQG2MQ)

### Как работать с движком:

1) форкаете
2) работаете

### Немного об архитектуре:

Сервер состоит из 3х частей:

- `Роутер` - является посредником между клиентом и нодой. Хранит в себе все игровые сущности и библиотеку с математикой,
  поэтому испортируется как библиотека в ноду.
- `Нода` - это сервис где происходит сама "игра", но 1 ноде может быть много игор.
- `Клиент` - та часть что видит игрок, является "терминалом" который просто выводит происходящие на экран с помощью
  графического движка.

Для работы сервера необходимо поднять 1 роутер и хотя бы 1 ноду. К одному роутеру может быть подключено много нод что
позволяет горизонтально масштабировать игру.

- Роутер и нода общаются между собой по rpc.
- Роутер и клиент общаются через websocket.

![This is an image](./readme_assets/img_2.png)

##### Как подключается нода:

1) когда нода запускается
   она [стучится по rpc](https://github.com/TrashPony/game-engine/blob/master/node/rpc/rpc.go#L56)
   на `veliriURL` указаный в [main.ini](./main.ini)
2) дополнительным аргументом передает параметр `nodeUrl` указаный в [main.ini](./main.ini) (этот аргумент передает ip:
   port машины на которой запущена нода)
3) если соеденение успешно
   то [роутер ловит запрос](https://github.com/TrashPony/game-engine/blob/master/router/rpc/server.go#L65), и
   регистрирует ноду
   в [сторедже нод](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/factories/nodes/nodes.go#L8),
   дает ему uuid и отсылает как ответ.
4) при добавление ноды в сторедж, роутер
   поднимает [обратный rpc канал](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/factories/nodes/nodes.go#L32)
   на ноду для полнодуплексной связи.
5) когда оба канала открыты, на ноде будут запускатся новые сессии. В случае если на ноде происходит ошибка, или ошибка
   при передаче, или отвал по таймауту то нода удаляется из стореджа.

##### API ноды

Для каждой ноды на роутере есть свое апи выраженое как методы ноды. Посмотреть
можно [тут](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/node/node.go#L74)
, нода ловит эти запросы [тут](https://github.com/TrashPony/game-engine/blob/master/node/rpc/rpc.go#L86).

- `CreateBattle` - создает новую сессию
- `FindBattle` - ищет бой по uuid
- `InitBattle` - когда бой создан, игрок запрашивает его состояние при загрузке на клиенте
- `StartLoad` - состояние было успешно получено и игрок запрашивает все обьекты в игре.
- `Input` - ввод игрока (клава/мышь)
- `CreateUnit` - создает юнита игроку инициатору
- `CreateBot` - создает бота в конкретную команду
- `CreateObj` - создает обьект (турель), в команду

##### Что посылает нода

Каждая нода [отправляет](https://github.com/TrashPony/game-engine/blob/master/node/rpc/rpc.go#L163) на роутер
данные обновление мира для каждого активного боя. Роутер
их [перенаправляет](https://github.com/TrashPony/game-engine/blob/master/router/rpc/server.go#L81) клиентам по
вебсокетам.

Маршрутизация сообщений для игроков на сокетах происходит
вот [тут](https://github.com/TrashPony/game-engine/blob/master/router/web_socket/sender.go#L75)

### Как запустить:

- поднимаем сервер

```bash
go mod tidy;
go run ./router/main.go;
go run ./node/main.go;
```

- запускаем статику

```bash
cd .\static\;
npm run dev;
#заходим http://localhost:8083/

## ИЛИ
cd .\static\;
npm run build;
#заходим http://localhost:8086/
```

настрока сети на стороне клиента находится
вот [тут](https://github.com/TrashPony/game-engine/blob/master/static/src/const.js) <br>
роутер ловит сообщения сокета вот [тут](https://github.com/TrashPony/game-engine/blob/master/router/main.go#L17)

### ----

<h3 id="create-world">
Сессия, игровое поле, GameLoop
</h3>

Сессию определяет
[обьект`Battle`](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/battle/battle.go),
но все игровые обьекты прикрепляются к
обьекту [карты (`Map`)](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/map/map.go#L10)
который встроен в `Battle`
.

Все игровые обьекты находящиеся на карте (юниты, пули, строения) привязываются к ней
полем `MapID` ([например](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/unit/unit.go#L20))
.

[Юниты](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/factories/units/units.go#L8)
и [пули](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/factories/bullets/bullets.go#L9) хранятся в
отдельных стореджах, при добавление в сторедж,
он [смотрит](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/factories/units/units.go#L56) на
поле `MapID` и кладет в соотвествующий массив.

Строения же находятся в
карте [карте (`Map`)](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/map/map.go)
.

```go
type Map struct {
// Размер карты в пикселях (в дальности в игре в писелях), в идиале квадрат, прямоугольник не тестировался)
XSize        int     `json:"x_size"`
YSize        int     `json:"y_size"`
// текстуры земли, не на что не влияют, нужны только для отрисовке на клиенте.
Flore map[int]map[int]*dynamic_map_object.Flore `json:"flore"`
// Не изменяемые обьекты (например горы и овраги), игрок видит эти обьекты всегда независимо от радара/обзора и они никогда не изменяются
StaticObjects   map[int]*dynamic_map_object.Object `json:"-"`
// Тут находятся игровые обьекты с которыми можно взаимодествовать (убить, построить, передвинуть и тд.), та же эти обьекты могут иметь поведение (например турели).
// Эти обьекты игрок видит только когда открыл их в тумане войны. Когда обьект ушел обратно в туман игрок запоминает его расположение и состояние.
// Игрок не видит измененияесли с обьектом вне поле его зрения.
DynamicObjects   []*dynamic_map_object.Object `json:"-"`
// Карта не плоская и у каждой клетке 16x16px есть своя высота (в текущей реализации это влияет только на пули), если она указана то хранится тут если нет то используется DefaultLevel
LevelMap [][]*LvlMap `json:"level_map"`
// высота карты по умолчанию
DefaultLevel float64 `json:"default_level"`
// кеширования непроходимых участков из за обьектов на карте для ускорения расчета коллизий. Подробнее смотри раздел коллизий.
GeoZones [][]*Zone `json:"-"`
}
```

#### Создание сессии и GameLoop

Создание сессии происходит на ноде при соотвествующем запросе
по [rpc](https://github.com/TrashPony/game-engine/blob/master/node/rpc/requests.go#L32).

Когда сессия успешно создана, она попадает
в [сторедж](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/factories/quick_battles/quick_battles.go#L24)
всех сессий на ноде.

Что бы сессия попала в `GameLoop`,
специалий [метод](https://github.com/TrashPony/game-engine/blob/master/node/game_loop/game_loop.go#L29) смотрит все
сессии и если она не инициализиорована то запускает её.

`GameLoop` - это игровой цикл который отслеживает и запускает все игровые механизмы (движения обьектов, расчеты физики,
обзора,
нанесения урона, пользовательский ввод и все такое). Одна итерация `GameLoop` это 1 кадр на стороне сервера, время этого
кадра ровняется [_const.ServerTick](https://github.com/TrashPony/game-engine/blob/master/router/const/const.go#L9) в мс.

Если итерация отработала
быстрее [_const.ServerTick](https://github.com/TrashPony/game-engine/blob/master/router/const/const.go#L9) то она
вычесляет дельту и спит это время, если дольше то все лагают.

В время работы итерации собираются сообщения об изменениях или событиях в специальный
обьект [web_socket.MessagesStore{}](https://github.com/TrashPony/game-engine/blob/master/node/game_loop/game_loop.go#L53)
.

В конце работы итерации всем
игрокам [отсылаются сообщения](https://github.com/TrashPony/game-engine/blob/master/node/game_loop/send_game_loop_data.go#L16)
предворительно попуская через фильтр видимости.

Подробнее о протоколе и способе обмена данными в
разделе [Бинарный протокол связи между серверов и клиентом](#binary-protocol)

<h3 id="game-objects">
Реализация игровых обьектов: "юниты", "строения", "пули"
</h3>

"Юниты", "строения", "пули" - это обьекты которые находятся на карте, имеют физическую модель и обрабатываются
в [`GameLoop`](https://github.com/TrashPony/game-engine/blob/master/node/game_loop/game_loop.go).

#### Юниты

Юниты - это обьекты которые перемещаются по карте своим ходом. В игровом представление может быть как танк так и
человек.

Характиристики юнита определяет
его ["тело"](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/body/body.go#L3)
тело представляет собой куб (длинна, ширина, высота). Параметр `Radius` для быстрый расчетов где не нужна
точность. Тело юнита имеет дальность обзора, слоты под оружие и тип шасси который определяет механику движения.

![This is an image](./readme_assets/img_3.png)<br>
<sub>красный квадрат это длинна, ширина. Синий круг `Radius`. Высота влияет на пули.</sub>

```go
type Body struct {
Texture       string              `json:"texture"` // текстура для отрисовки на клиенте.
MaxHP         int                 `json:"max_hp"`
Scale         int                 `json:"scale"`  // сервер заточен на текстуры размером 128х128px, это параметр говорит % от этого размера. Этот параметр влияет на все параметры габаритов и якорей для оружия (подробнее в разделе оружия).
Length        int                 `json:"length"` // длинна полигона (Length * (Scale/100))
Width         int                 `json:"width"`  // ширина полигона (Width * (Scale/100))
Height        int                 `json:"height"` // высота полигона (Height * (Scale/100))
Radius        int                 `json:"radius"` // радиус корпуса для простого расчета колизи в поиске пути
RangeView     int                 `json:"range_view"`  // дальность видимости
RangeRadar    int                 `json:"range_radar"` // дальность радара
Weapons       map[int]*WeaponSlot `json:"weapons"`      // слоты и якоря для оружия (подробнее в разделе оружия)
ChassisType   string              `json:"chassis_type"` // тип шасси, в текущий реализации есть только гусеницы. Можно реализовать свою механику перемещания. (например антиграв, колеса, полет)
Speed         float64             `json:"speed"` // максимальная скорость движения (W)
ReverseSpeed  float64             `json:"-"`     // максимальная скорость движения назад (S)
PowerFactor   float64             `json:"power_factor"`// сила ускорения
ReverseFactor float64             `json:"-"`           // сила ускорения назад (тормоз)
TurnSpeed     float64             `json:"turn_speed"`  // скорость поворота
MoveDrag      float64             `json:"move_drag"` // сила трения при движение
AngularDrag   float64             `json:"-"`         // сила трения при повороте
Weight        float64             `json:"-"`           // вес (влияет на столкновения)
}
```

Параметры связанные с
физикой [перетекают](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/unit/physical_model.go#L21)
в
структуру [физической модели](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/physical_model/physical_model.go)
при первом обращение к ней.

#### Строения

Строения- это обьекты которые не могут самостоятельно перемещатся или не могут перемещатся вообще. У строений нет "тела"
как у юнитов. Их физическую модель
представляет [массив](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/dynamic_map_object/object.go#L44)
из
непроходимых [точек](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/obstacle_point/obstacle_point.go#L7)
.

![This is an image](./readme_assets/img_4.png)<br>

При инициализации или движение обьекта геодата собирается в
методе [SetGeoData](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/dynamic_map_object/geo_data.go)
и после регестрируется
в [Map.GeoZones](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/map/geo_zone.go)

Некоторые строения как и юниты могут иметь дальность обзора, оружие
и [поведение](https://github.com/TrashPony/game-engine/blob/master/node/game_loop/game_loop_structure/game_loop_structure.go)
.

#### Общее (Юниты/Строения)

В строение и юнитов встроены общие структуры и интерфейсы:

```go 
{
   physicalModel    *physical_model.PhysicalModel          // физическое представление объекта (именно с этой структурой происходит движение, столкновения, пападания пуль)
   gunner           *gunner.Gunner                         // интерфейс для использования оружия
   BurstOfShots     *burst_of_shots.BurstOfShots           // очередь пуль которые вытеют из ствола, у оружия может быть несколько стволов и стрельба очередью
   weaponTarget     *target.Target                         // цель для атаки
   visibleObjects   *visible_objects.VisibleObjectsStore   // обьекты которые "видит", в текущей реализации встраиваается экземляр "комманды"
}
```

Движение юнитов и строений обрабатывается
в `GameLoop` [одним методом](https://github.com/TrashPony/game-engine/blob/master/node/game_loop/game_loop_move/game_loop_move.go)
, [стрельба](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/attack/fire.go#L25)
и [поворот](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/attack/rotate_gun.go#L19) орудий тоже.

#### Пули

[Пуля](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_objects/bullet/bullet.go) в игре
представляет из себя шарик
радиусом [_const.AmmoRadius](https://github.com/TrashPony/game-engine/blob/master/router/const/const.go#L7), физика пули
работает на
формулах ["Полета снаряда"](https://en.wikipedia.org/wiki/Projectile_motion) ([реализация](https://github.com/TrashPony/game-engine/blob/master/router/mechanics/game_math/bullet_ballistic.go))
. В зависимости от типа
снаряда ([пуля](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/fly_bullets/fly_bullet.go)
, [лазер](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/fly_bullets/fly_laser.go),
[ракеты](https://github.com/TrashPony/game-engine/blob/master/node/mechanics/fly_bullets/fly_missile.go)) может немного
менять поведение.

Пули в `GameLoop`
обрабатываются [тут](https://github.com/TrashPony/game-engine/blob/master/node/game_loop/game_loop_bullets/game_loop_bullets.go)

<h3 id="binary-protocol">
Бинарный протокол связи между серверов и клиентом
</h3>

<h3 id="input">
Ввод пользователя
</h3>

<h3 id="move">
Движение обьектов в мире
</h3>

<h3 id="weapon">
Баллистика пуль/стрельба
</h3>

<h3 id="physical-model">
Физическую модель обьектов
</h3>

<h3 id="collision">
Обнаружение коллизий, реакция на коллизию
</h3>

<h3 id="watch">
Реализация "Обзора" игроков
</h3>

<h3 id="create-update-client">
Создание и обновление игровых обьектов на стороне клиента
</h3>

<h3 id="ai">
ИИ
</h3>

<h3 id="ai-support">
Поиск пути, управление движением ИИ
</h3>
