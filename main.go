package main

import (
	//"fmt"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	tele "gopkg.in/telebot.v4"
)

var (
	alexCount  int = 0
	alenaCount int = 0
)

func main() {
	// create telegramm connection
	pref := tele.Settings{
		Token:     os.Getenv("TOKEN"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	// create database connection
	db, err := sqlx.Connect("sqlite3", "./love.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS love_counter (
   name TEXT PRIMARY KEY,
   count INTEGER)
`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS love_counter (
 		  name TEXT PRIMARY KEY,
  		  count INTEGER)`)
	var counts []int
	db.Select(&counts, "SELECT count FROM love_counter;")
	if len(counts) > 0 {
		alexCount = counts[0]
		alenaCount = counts[1]
	}
	if err != nil {
		log.Fatal(err)
	}
	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello! Я помогу вам определить, кто кого больше любит)))")
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		msg := c.Message().Text
		if strings.Contains(strings.ToLower(msg), "люблю") && strings.Contains(strings.ToLower(msg), "тебя") {
			if c.Sender().Username == "qb1110" {
				_, err := db.Exec("UPDATE love_counter SET count = count + 1 WHERE name = 'Alex';")
				if err != nil {
					log.Fatal(err)
				}
				alexCount += 1
				return c.Send(fmt.Sprintf(`❤️

<i>Алексеюшка</i>: <b>%v</b>
<i>Веснушка</i>: <b>%v</b>
		`, alexCount, alenaCount))
			} else {
				_, err := db.Exec("UPDATE love_counter SET count = count + 1 WHERE name = 'Alena';")
				if err != nil {
					log.Fatal(err)
				}
				alenaCount += 1
				return c.Send(fmt.Sprintf(`❤️

<i>Алексеюшка</i>: <b>%v</b>
<i>Веснушка</i>: <b>%v</b>
		`, alexCount, alenaCount))
			}
		}
		return nil
	})

	b.Start()
}
