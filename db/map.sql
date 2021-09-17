CREATE TABLE maps
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(64),
    x_size        INT, /* размер карты по Х */
    y_size        INT, /* размер карты по Y */
    level         INT, /* определяет основной уровень координат на карте еще он не перепределен конструктором */
    specification VARCHAR(255) /* описание карты */
);

CREATE TABLE maps_spawn
(
    id     SERIAL PRIMARY KEY,
    id_map INT REFERENCES maps (id), /* ид карты к которой принадлежит координата */
    x      int,
    y      int,
    radius int,
    rotate int
);

-- таблица высота мапы, основная высота карты указана в таблице maps->level если она отличается то указана тут по x/y = lvl
-- высота в игре это абстрактная величина от 0.0 до 3.0 ++
-- 0 - нижний уровень (ямы)
-- 1 - уровень земля по умолчанию
-- 2 - верхний уровень (подьем, гора)
-- Юниты имеют высоту 0.5 (если юнит стреляет с земли 1 то пуля имеет высоту 1.5)
-- у обьекто и щитов тоже есть высота (если щит стоит на верхнем уровне (2) и имеет высоту 2 то высота щита 4, все снаряды ниже будут в него врезатся)
CREATE TABLE map_level
(
    id_map int  not null default 0,
    x      int  not null default 0,
    y      int  not null default 0,
    lvl    real not null default 1
);

-- содержит только статичные обьекты
CREATE TABLE map_constructor
(
    id                 SERIAL PRIMARY KEY,
    id_map             INT REFERENCES maps (id), /* ид карты к которой принадлежит координата */
    id_type            INT, /* ид типа координаты */
    texture_over_flore VARCHAR(64), /* название текстуры земли */
    /* говорит в какой последовательности отрисовывать текстуры, ид енподходит т.к. координата уже могла быть в бд перед нанесения текстуры */
    texture_priority   INT,
    /* тоже самое но для обьектов */
    object_priority    INT,
    x                  INT,
    y                  INT,
    rotate             INT, /* говорит на сколько повернуться спрайту обьекта в координате если он есть конечно */
    x_shadow_offset    INT, /* смещение тени по Х от центра координаты */
    y_shadow_offset    INT, /* смещение тени по Y от центра координаты */
    scale              INT /* определяет какой размер должен быть тексуры обьекта, анимации на карте 100 - 100%, 10 - 10%, 200 - 200% от оригинала */
);

CREATE TABLE coordinate_type
(
    id                        SERIAL PRIMARY KEY,

    -- если пусто то тупо обьект, (turret, shield, radar, generator, jammer, missile_defense, meteorite_defense)
    type                      VARCHAR(64),

    /* одновременно может быть либо статичный обьект, либо анимация */
    texture_object            VARCHAR(64), /* имя текстуры обьекта (камень, дерево, стена и тд) если он есть на зоне */
    animate_sprite_sheets     VARCHAR(64), /* имя файла анимации */

    animate_loop              BOOLEAN, /* если координата анимирована говорит что анимация будет всегда по кругу иначе анимацию должно что то активировать */
    /* параметр чисто для отображения, говорит перекроет юнит своим телом этот обьект или нет если надетет на него*/
    unit_overlap              BOOLEAN,

    object_name               text    not null default '',
    object_description        text    not null default '',
    object_inventory          BOOLEAN not null default false, -- для этой координаты по ид создается ящик в таблице ящиков на карте
    object_inventory_capacity int     not null default 0,
    object_hp                 int     not null default -1,    -- -2 - бесмертный и даже нет хп, -1 - бесмертный с хп, 0 - мертвый

    /* геодата привязаная к обьекту,
     работает как и глобальная дата но х у это смещение от центра обьекта и пропадает если обьекта больше нет.
     x, y смещается если обьект повернут, имеет размер отличный от 100%
     radius изменяется если обьект имеет размер отличный от 100%
     [ {"x": 1, "y": 1, "radius":90}, {"x": 2, "y": 2, "radius":90} ]
     */
    geo_data                  json    not null default '{}',

    shadow_intensity          INT, /* сила тени от 0 до 1, (val / 100) */
    animate_speed             INT, /* если координата анимация говорит с какой скоростью ее вопспроизводить, кадров в секунду */
    shadow                    BOOLEAN, /* определяет нужна ли обьекту тень */

    x_shadow_offset           INT, /* смещение тени по Х от центра координаты */
    y_shadow_offset           INT, /* смещение тени по Y от центра координаты */

    height                    int     not null default 0      -- высота обьекта при размере 100
);

CREATE TABLE global_geo_data
(
    id     SERIAL PRIMARY KEY,
    id_map INT REFERENCES maps (id), /* где находится непроходимая точка */
    x      INT,
    y      INT,
    radius INT /* размер непроходимой точки */
);