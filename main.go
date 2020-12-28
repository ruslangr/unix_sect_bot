package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"./botсonfd"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	users     map[string]string // map ник - статус
	chatID    int64
	names     map[string]string // map имя - ник
	listUsers map[string]string // map имя - статус, для отчета
	statusMsg string
)

func main() {

	users = make(map[string]string) // ники в telegram
	users[botсonfd.SDM] = "В отпуске, здоров"
	users[botсonfd.GRR] = "Дома - здоров"
	users[botсonfd.AAA] = "Дома, здоров."
	users[botсonfd.SAY] = "Дома - здоров"
	users[botсonfd.DDP] = "Дома, здоров"
	users[botсonfd.AIY] = "Загородом, здоров"

	names = make(map[string]string) // настоящие имена
	names["Дима"] = botсonfd.SDM
	names["Руслан"] = botсonfd.GRR
	names["Айтуган"] = botсonfd.AAA
	names["Айдар"] = botсonfd.SAY
	names["Денис"] = botсonfd.DDP
	names["Ильнур"] = botсonfd.AIY

	listUsers = make(map[string]string) // map для отчета

	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(botсonfd.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	go sendReport(bot)

	// инициализируем канал, куда будут прилетать обновления от API
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	// читаем обновления из канала
	for update := range updates {

		reply := ""
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// логируем от кого какое сообщение пришло
		//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		userNick := update.Message.From.UserName
		userStatus := update.Message.Text
		//userState := userNick + " - " + userStatus // для отладки

		//sendUserState := tgbotapi.NewMessage(update.Message.Chat.ID, userState)

		if strings.Contains(userStatus, "/") {
			switch update.Message.Command() {
			case "status_sect": // пришлет статус всех в меню бота
				reply = fmt.Sprintln(listUsers)
				replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(replyMsg)
			case "send": // пошлет в чат chatID статус
				go sendNotifications(bot)
			}
		} else {
			if _, ok := users[userNick]; ok {
				users[userNick] = userStatus
				for key1, val1 := range users {
					for key2, val2 := range names {
						if key1 == val2 {
							listUsers[key2] = val1
						}
					}
				}
			}
		}
	}
}

func sendNotifications(bot *tgbotapi.BotAPI) {
	chatID = -474929898
	var stringMap string
	stringMap = mapToList(listUsers)
	//stringMap = fmt.Sprintln(listUsers)
	bot.Send(tgbotapi.NewMessage(chatID, stringMap))
}

func sendReport(bot *tgbotapi.BotAPI) {
	for {
		t := time.Now()
		ts := t.Format(time.UnixDate)
		tsArr := strings.Split(ts, " ")
		tts := tsArr[3]
		ttsArr := strings.Split(tts, ":")
		hour := ttsArr[0]
		minute := ttsArr[1]
		if hour == botсonfd.SendHour1 && minute == botсonfd.SendMinute1 {
			sendNotifications(bot)
			time.Sleep(time.Minute * 1)
		}
		if hour == botсonfd.SendHour2 && minute == botсonfd.SendMinute2 {
			sendNotifications(bot)
			time.Sleep(time.Minute * 1)
		}
	}
}

func mapToList(m map[string]string) string {
	var str string
	var strAll string
	for key, val := range m {
		str = fmt.Sprintf("%s - %s", key, val)
		strAll = strAll + "\n" + str
	}
	return strAll
}
