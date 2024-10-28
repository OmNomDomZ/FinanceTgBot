package main

import (
	"database/sql"
	"fmt"
	"github.com/mymmrac/telego"
	"log"
	"os"
	"strconv"
	"time"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	_ "github.com/lib/pq"
)

type UserState struct {
	awaitingAmount  bool
	transactionType string // "доход", "расход", "удаление"
}

var userStates = make(map[int64]*UserState)

func sendKeyboardButton(bot *telego.Bot, update telego.Update) {
	chatID := tu.ID(update.Message.Chat.ID)

	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Добавить доход"),
			tu.KeyboardButton("Добавить расход"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("Вывести все транзакции"),
			tu.KeyboardButton("Стереть все транзакции"),
		),
	)

	message := tu.Message(
		chatID,
		"Выберите действие:",
	).WithReplyMarkup(keyboard)

	bot.SendMessage(message)
}

func handleAddTransaction(bot *telego.Bot, update telego.Update, db *sql.DB) {
	chatID := update.Message.Chat.ID

	state, exists := userStates[chatID]
	if exists && state.awaitingAmount {

		amount, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			message := tu.Message(
				tu.ID(chatID),
				"Неверный формат суммы. Попробуйте ещё раз.",
			)
			bot.SendMessage(message)
			return
		}

		_, err = db.Exec("INSERT INTO transactions (amount, transaction_type) VALUES ($1, $2)", amount, state.transactionType)
		if err != nil {
			message := tu.Message(
				tu.ID(chatID),
				"Произошла ошибка при записи в базу данных. Попробуйте позже.",
			)
			bot.SendMessage(message)
			return
		}

		message := tu.Message(
			tu.ID(chatID),
			fmt.Sprintf("%s в размере %.2f успешно добавлен!", state.transactionType, amount),
		)
		bot.SendMessage(message)

		delete(userStates, chatID)
	}
}

func handleIncomeCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	userStates[chatID] = &UserState{
		awaitingAmount:  true,
		transactionType: "доход",
	}

	message := tu.Message(
		tu.ID(chatID),
		"Напишите сумму дохода: ",
	)

	bot.SendMessage(message)
}

func handleExpenseCommand(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	userStates[chatID] = &UserState{
		awaitingAmount:  true,
		transactionType: "расход",
	}

	message := tu.Message(
		tu.ID(chatID),
		"Напишите сумму расхода: ",
	)

	bot.SendMessage(message)
}

func handleTransactionsInfo(bot *telego.Bot, update telego.Update, db *sql.DB) {
	chatID := tu.ID(update.Message.Chat.ID)

	message := tu.Message(
		chatID,
		"Все транзакции: ",
	)
	bot.SendMessage(message)

	rows, err := db.Query("SELECT * FROM transactions")
	if err != nil {
		log.Printf("Ошибка при запросе к базе данных: %v", err)
		message := tu.Message(
			chatID,
			"Произошла ошибка при чтении базы данных. Попробуйте позже.",
		)
		bot.SendMessage(message)
		return
	}
	defer rows.Close()

	for rows.Next() {
		index := 1
		var id int
		var amount float64
		var transactionType string
		var createdAt time.Time
		if err := rows.Scan(&id, &amount, &transactionType, &createdAt); err != nil {
			log.Printf("Ошибка при чтении строки: %v", err)
			continue
		}
		message := tu.Message(
			chatID,
			fmt.Sprintf("%d: Amount: %.2f\n Transaction_type: %s\n Created_at: %s", index, amount, transactionType, createdAt.Format("02-01-2006 15:04:05")),
		)
		bot.SendMessage(message)
		index++
	}
}

func handleClearTable(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	userStates[chatID] = &UserState{
		awaitingAmount:  false,
		transactionType: "удаление",
	}

	message := tu.Message(
		tu.ID(chatID),
		"Вы уверены, что хотите стереть все данные?\n[Да / Нет]",
	)
	bot.SendMessage(message)
}

func handleConfirmation(bot *telego.Bot, update telego.Update, db *sql.DB) {
	chatID := update.Message.Chat.ID

	state, exists := userStates[chatID]
	if !exists || state.transactionType != "удаление" {
		return
	}
	agreement := update.Message.Text

	if agreement == "Да" {
		_, err := db.Exec("TRUNCATE TABLE transactions RESTART IDENTITY;\n")
		if err != nil {
			log.Printf("Ошибка при удалении данных: %v", err)
			message := tu.Message(
				tu.ID(chatID),
				"Произошла ошибка при удалении данных. Попробуйте позже.",
			)
			bot.SendMessage(message)
			return
		}

		message := tu.Message(
			tu.ID(chatID),
			"Все транзакции успешно стерты.",
		)
		bot.SendMessage(message)
	} else if agreement == "Нет" {
		message := tu.Message(
			tu.ID(chatID),
			"Удаление данных отменено.",
		)
		bot.SendMessage(message)
	} else {
		message := tu.Message(
			tu.ID(chatID),
			"Пожалуйста, ответьте 'Да' или 'Нет'.",
		)
		bot.SendMessage(message)
		return
	}

	delete(userStates, chatID)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	defer db.Close()

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		return
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(sendKeyboardButton, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		handleIncomeCommand(bot, update)
	}, th.TextEqual("Добавить доход"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		handleExpenseCommand(bot, update)
	}, th.TextEqual("Добавить расход"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		handleTransactionsInfo(bot, update, db)
	}, th.TextEqual("Вывести все транзакции"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		handleClearTable(bot, update)
	}, th.TextEqual("Стереть все транзакции"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := update.Message.Chat.ID

		state, exists := userStates[chatID]
		if !exists {
			return
		}

		if state.transactionType == "удаление" {
			handleConfirmation(bot, update, db)
			return
		}

		if state.awaitingAmount {
			handleAddTransaction(bot, update, db)
			return
		}

		sendKeyboardButton(bot, update)
	}, th.AnyMessage())

	bh.Start()
}
