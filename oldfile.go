package main

import (
  "database/sql"
  "github.com/mattn/go-sqlite3"
  "log"
  "time"
  "strconv"
  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
  chatid int64
  userid int64
  times []int64
}

func (user *User) sum() int64{
  var s int64 = 0
  var t int64 = 0
  for i, time := range user.times {
    if i % 2 == 0 {
      t = time
    } else {
      s += time - t
    }
  }
  return s
}

func main() {
  db, err := sql.Open("sqlite3", "mydatabase.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()
  log.Println("db Open")

  bot, err := tgbotapi.NewBotAPI("1763199303:AAHm07HCRcXvqeo_BYd_G7MSxeZN74doZFg")
  if err != nil {
    log.Panic(err)
  }
  
  users := []User{}

  bot.Debug = true

  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates := bot.GetUpdatesChan(u)

  for update := range updates {
    if update.Message != nil {
      log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)

      found := false
      var iduser int
      for i, user := range users {
        if user.userid == update.Message.Chat.ID && user.chatid == update.Message.From.ID {
          iduser = i
          found = true
        }
      }
      
      if !found {
        users = append(users, User{update.Message.Chat.ID, update.Message.From.ID, []int64{}})
      } else {
        switch update.Message.Command() {
        case "add": 
          users[iduser].times = append(users[iduser].times, time.Now().Unix())
        case "sum":
          msg := tgbotapi.NewMessage(users[iduser].chatid, strconv.Itoa(int(users[iduser].sum())))
          bot.Send(msg)
        }
      }

      msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
      bot.Send(msg)
    }
  }
}
