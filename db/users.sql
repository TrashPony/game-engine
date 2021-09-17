CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    vk_user_id int not null default 0,
    name       VARCHAR(64),
    password   VARCHAR(255),
    mail       VARCHAR(64),
    role       text,
    language   text
);

CREATE TABLE players
(
    id        SERIAL PRIMARY KEY,
    user_id   INT REFERENCES users (id),
    name      VARCHAR(64),
    avatar    text,
    biography text,
    title     text
);

CREATE TABLE player_interface
(
    id_players INT REFERENCES players (id),
    data       json -- тут текущее состояние окон ¯\_(ツ)_/¯
);