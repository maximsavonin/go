package main

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "time"
  "strconv"
  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
func sectostr(t int) string {
  var text string = ""
  a := t / (60*60*24)
  if a != 0 {
    text = strconv.Itoa(a) + " days "
  }
  text += strconv.Itoa(t / (60 * 60) % 24) + " hours " + strconv.Itoa(t / 60 % 60) + " min " + strconv.Itoa(t % 60) + " sec"
  return text
}

func main() {
  db, err := sql.Open("sqlite3", "./mydatabase.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()
  log.Println("db Open")
  
  statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, user_id INTEGER, chat_id INTEGER)")
  statement.Exec()

  statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS times (id INTEGER PRIMARY KEY, time INTEGER, user_id INTEGER, FOREIGN KEY (user_id) REFERENCES users (id))")
  statement.Exec()

  bot, err := tgbotapi.NewBotAPI("")
  if err != nil {
    log.Panic(err)
  }

  bot.Debug = true

  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates := bot.GetUpdatesChan(u)

  for update := range updates {
    if update.Message != nil {
      log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)
      
      rows, _ := db.Query("SELECT id FROM users WHERE chat_id = (?)", update.Message.Chat.ID)
      
      var id int
      
      if rows.Next() {
        rows.Scan(&id)
      } else {
        statement, _ = db.Prepare("INSERT INTO user (user_id, chat_id) VALUES (?, ?)")
        statement.Exec(update.Message.Chat.ID, update.Message.From.ID)
      }
      rows.Close()
      switch update.Message.Command() {
      case "add":
        statement, _ = db.Prepare("INSERT INTO times (time, user_id) VALUES (?, ?)")
        statement.Exec(time.Now().Unix(), id)
      case "sum":
        rowstimes, _ := db.Query("SELECT time FROM times WHERE user_id = (?)", id)
        defer rowstimes.Close()
          
        var tfir int
        var tsec int
        var s int = 0
        var i int = 0
        
        for rowstimes.Next() {
          if i % 2 == 0 {
            rowstimes.Scan(&tfir)
          } else {
            rowstimes.Scan(&tsec)
            s += tsec - tfir
          }
          i += 1
        }
        
        var text string = ""

        if i % 2 == 1 {
          s += int(time.Now().Unix()) - tfir
          text = "  ->"
        }
        

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, sectostr(s) + text)
        bot.Send(msg)
      default:
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What??")
        bot.Send(msg)
      }
    }
  }
}
