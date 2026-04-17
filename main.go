package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
)

const (
	BotToken   = "YOUR_TELEGRAM_BOT_TOKEN" // Замените на ваш токен бота
	ChatID     = 000000000000              // Замените на ваш Chat ID
	SteamURL   = "https://store.steampowered.com/search/?maxprice=free&specials=1&ndl=1"
	CheckDelay = 15 * time.Minute
)

type Game struct {
	Title string
	Link  string
}

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Printf("Бот %s запущен. Режим тишины при первом запуске включен.\n", bot.Self.UserName)

	foundGames := make(map[string]bool)

	// Флаг первого запуска
	isFirstRun := true

	for {
		fmt.Printf("[%s] Проверка Steam...\n", time.Now().Format("15:04:05"))
		games := checkSteam()

		newFoundCount := 0
		for _, game := range games {
			if !foundGames[game.Title] {
				// Если это не первый запуск — отправляем в ТГ
				if !isFirstRun {
					msgText := fmt.Sprintf("🎁 **НОВАЯ РАЗДАЧА!**\n\n🎮 *%s*\n\n🔗 [Открыть в Steam](%s)", game.Title, game.Link)
					msg := tgbotapi.NewMessage(ChatID, msgText)
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					fmt.Printf("НОВАЯ ИГРА: %s\n", game.Title)
				}

				// Добавляем в базу виденных
				foundGames[game.Title] = true
				newFoundCount++
			}
		}

		if isFirstRun {
			fmt.Printf("Первое сканирование завершено. Запомнил игр: %d. Жду обновлений...\n", newFoundCount)
			isFirstRun = false
		}

		time.Sleep(CheckDelay)
	}
}

func checkSteam() []Game {
	var games []Game
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "birthtime=946684801; lastagecheckage=1-0-2000")
	})

	c.OnHTML("a.search_result_row", func(e *colly.HTMLElement) {
		title := e.ChildText(".title")
		link := e.Attr("href")
		cleanLink := strings.Split(link, "?")[0]

		if title != "" {
			games = append(games, Game{Title: title, Link: cleanLink})
		}
	})

	c.Visit(SteamURL)
	return games
}
