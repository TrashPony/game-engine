package dbConnect

import (
	"database/sql"
	_const "github.com/TrashPony/game_engine/src/const"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"time"
)

var DB *sql.DB

func init() {

	//Генератор случайных чисел обычно нужно рандомизировать перед использованием, иначе, он, действительно,
	// будет выдавать одну и ту же последовательность.
	rand.Seed(time.Now().UnixNano())

	command := "postgres://" +
		_const.Config.GetParams("dbLogin") + ":" +
		_const.Config.GetParams("dbPass") + "@" +
		_const.Config.GetParams("dbAddress") + "/infinity_game?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", command)
	if err != nil {
		log.Fatal(err)
	}

	DB.SetMaxOpenConns(600) // TODO сделать очередь запросов

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}

func GetDBConnect() *sql.DB {
	return DB
}
