//package main
//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"log"
//	"net/http"
//	"strconv"
//	"strings"
//
//	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
//)
//
//type binanceResp struct {
//	Price float64 `json:"price,string"` // указываем теги для парсеера
//	Code  int64   `json:"code"`
//}
//
//type wallet map[string]float64 //кошелек
//
//var db = map[int64]wallet{} // БД
//
//func main() {
//	bot, err := tgbotapi.NewBotAPI("5627299146:AAEgDE7_3jHRV4bday2zSt-MnEaTtR0QLq8")
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//
//	updates := bot.GetUpdatesChan(u)
//
//	for update := range updates {
//		if update.Message != nil { // If we got a message
//			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
//
//			msgArr := strings.Split(update.Message.Text, " ") //принимает сообщение от пользователя и вносит в массив в каждый индекс через пробел
//			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgArr[0])
//			//msg.ReplyToMessageID = update.Message.MessageID //бот ссылается на сообщение
//
//			switch msgArr[0] {
//			case "ADD":
//				summ, err := strconv.ParseFloat(msgArr[2], 64)
//				if err != nil {
//					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ERR NOT INITIALIZATION MAP"))
//					continue
//				}
//				if _, ok := db[update.Message.Chat.ID]; !ok {
//					db[update.Message.Chat.ID] = wallet{}
//				}
//
//				db[update.Message.Chat.ID][msgArr[1]] += summ
//
//				msg := fmt.Sprintf("Баланс увеличился и составляет:: %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])
//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
//			case "SUB":
//
//				summ, err := strconv.ParseFloat(msgArr[2], 64)
//				if err != nil {
//					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ERR NOT INITIALIZATION MAP"))
//					continue
//				}
//				if _, ok := db[update.Message.Chat.ID]; !ok {
//					db[update.Message.Chat.ID] = wallet{}
//				}
//
//				db[update.Message.Chat.ID][msgArr[1]] -= summ
//
//				msg := fmt.Sprintf("Баланс уменьшился и составляет: %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])
//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
//			case "DEL":
//				delete(db[update.Message.Chat.ID], msgArr[1])
//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта далена"))
//			case "SHOW":
//				msg := "Баланс:\n"
//				var usdSumm float64
//				for key, value := range db[update.Message.Chat.ID] {
//					coinPrice, err := getPrice(key)
//					if err != nil {
//						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
//					}
//					usdSumm += value * coinPrice
//					msg += fmt.Sprintf("%s: %f [%.3f$]\n", key, value, value*coinPrice)
//				}
//				msg += fmt.Sprintf("Сумма: %.3f$", usdSumm)
//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
//			default:
//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
//			}
//
//		}
//	}
//}
//
//
//
//func getPrice(coin string) (price float64, err error) {
//
//	resp, err := http.Get(fmt.Sprintf("https://www.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))
//	if err != nil {
//		return
//	}
//
//	defer resp.Body.Close() // Закрытие соединения запроса get
//
//	var jsonResp binanceResp
//	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
//
//	if err != nil {
//		return
//	}
//
//	if jsonResp.Code != 0 {
//		err = errors.New("Некорректная валюта")
//	}
//
//	price = jsonResp.Price
//
//	return
//}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	coin2 "go.mod/internal/coin"
	coin "go.mod/internal/coin/db"
	"go.mod/internal/config"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/julienschmidt/httprouter"
)

type binanceResp struct {
	Price float64 `json:"price,string"` // указываем теги для парсеера
	Code  int64   `json:"code"`
}

type wallet map[string]float64 //кошелек

var db = map[int64]wallet{} // БД

func main() {
	//logger := logging.GetLogger()
	//logger.Info("create router")
	//router := httprouter.New()
	//
	//cfg := config.GetConfig()
	//
	//postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//repository := coin.NewRepository(postgreSQLClient, logger)
	//
	//c := coin2.Coin{
	//	Name: "OK",
	//}
	//err = repository.Create(context.TODO(), &c)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//
	//
	//
	//logger.Info("register author handler")
	//
	//coinHandler := coin2.NewHandler(repository, logger)
	//coinHandler.Register(router)
	//
	//start(router, cfg)

	bot, err := tgbotapi.NewBotAPI("5627299146:AAEgDE7_3jHRV4bday2zSt-MnEaTtR0QLq8")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msgArr := strings.Split(update.Message.Text, " ") //принимает сообщение от пользователя и вносит в массив в каждый индекс через пробел
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgArr[0])
			msg.ReplyToMessageID = update.Message.MessageID //бот ссылается на сообщение

			switch msgArr[0] {
			case "ADD":
				//summ, err := strconv.ParseFloat(msgArr[2], 64)
				//if err != nil {
				//	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ERR NOT INITIALIZATION MAP"))
				//	continue
				//}
				//if _, ok := db[update.Message.Chat.ID]; !ok {
				//	db[update.Message.Chat.ID] = wallet{}
				//}
				//
				//db[update.Message.Chat.ID][msgArr[1]] += summ
				//
				//msg := fmt.Sprintf("Баланс увеличился и составляет:: %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])
				//bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

				logger := logging.GetLogger()
				logger.Info("create router")
				router := httprouter.New()

				cfg := config.GetConfig()

				postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
				if err != nil {
					logger.Fatalf("%v", err)
				}
				repository := coin.NewRepository(postgreSQLClient, logger)

				//c := coin2.Coin{
				//	Name: "ETH",
				//}
				//err = repository.Create(context.TODO(), &c)
				//if err != nil {
				//	logger.Fatal(err)
				//}

				all, err := repository.FindAll(context.TODO())
				if err != nil {
					logger.Fatal("%v", err)
				}
				for _, cin := range all {
					logger.Info("%v", cin)
					msg := fmt.Sprintf("Баланс увеличился и составляет:: %s", cin)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
				}

				logger.Info("register author handler")

				coinHandler := coin2.NewHandler(repository, logger)
				coinHandler.Register(router)

				start(router, cfg)

			case "SUB":

				summ, err := strconv.ParseFloat(msgArr[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ERR NOT INITIALIZATION MAP"))
					continue
				}
				if _, ok := db[update.Message.Chat.ID]; !ok {
					db[update.Message.Chat.ID] = wallet{}
				}

				db[update.Message.Chat.ID][msgArr[1]] -= summ

				msg := fmt.Sprintf("Баланс уменьшился и составляет: %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			case "DEL":
				delete(db[update.Message.Chat.ID], msgArr[1])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта далена"))
			case "SHOW":
				msg := "Баланс:\n"
				var usdSumm float64
				for key, value := range db[update.Message.Chat.ID] {
					coinPrice, err := getPrice(key)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					}
					usdSumm += value * coinPrice
					msg += fmt.Sprintf("%s: %f [%.3f$]\n", key, value, value*coinPrice)
				}
				msg += fmt.Sprintf("Сумма: %.3f$", usdSumm)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			case "":

			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
			}

		}
	}
}

func getPrice(coin string) (price float64, err error) {

	resp, err := http.Get(fmt.Sprintf("https://www.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))
	if err != nil {
		return
	}

	defer resp.Body.Close() // Закрытие соединения запроса get

	var jsonResp binanceResp
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)

	if err != nil {
		return
	}

	if jsonResp.Code != 0 {
		err = errors.New("Некорректная валюта")
	}

	price = jsonResp.Price

	return
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
