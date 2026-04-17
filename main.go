package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
)

const (
	BotToken   = "YOUR_TELEGRAM_BOT_TOKEN"
	ChatID     = 000000000000
	SteamURL   = "https://store.steampowered.com/search/?maxprice=free&specials=1&ndl=1"
	CheckDelay = 15 * time.Minute
	DBFile     = "games_db.json" // Файл для хранения истории
)

type Game struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Printf("Бот %s запущен. База загружается из %s\n", bot.Self.UserName, DBFile)

	// Загружаем уже виденные игры из файла
	foundGames := loadDatabase()

	for {
		fmt.Printf("[%s] Проверка Steam...\n", time.Now().Format("15:04:05"))
		games := checkSteam()

		hasChanges := false
		for _, game := range games {
			if !foundGames[game.Title] {
				// Отправляем в ТГ только новые игры
				msgText := fmt.Sprintf("🎁 **НОВАЯ РАЗДАЧА!**\n\n🎮 *%s*\n\n🔗 [Открыть в Steam](%s)", game.Title, game.Link)
				msg := tgbotapi.NewMessage(ChatID, msgText)
				msg.ParseMode = "Markdown"
				bot.Send(msg)
				
				fmt.Printf("НОВАЯ ИГРА ДОБАВЛЕНА: %s\n", game.Title)
				foundGames[game.Title] = true
				hasChanges = true
			}
		}

		// Если нашли что-то новое — сохраняем обновленную базу
		if hasChanges {
			saveDatabase(foundGames)
		} else {
			fmt.Println("Ничего нового не найдено.")
		}

		time.Sleep(CheckDelay)
	}
}

// Загрузка базы из JSON
func loadDatabase() map[string]bool {
	db := make(map[string]bool)
	file, err := os.ReadFile(DBFile)
	if err != nil {
		if os.IsNotExist(err) {
			return db // Если файла нет, возвращаем пустую базу
		}
		log.Printf("Ошибка чтения базы: %v", err)
		return db
	}

	var savedList []string
	if err := json.Unmarshal(file, &savedList); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		return db
	}

	for _, title := range savedList {
		db[title] = true
	}
	return db
}

// Сохранение базы в JSON
func saveDatabase(db map[string]bool) {
	var list []string
	for title := range db {
		list = append(list, title)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Printf("Ошибка маршалинга: %v", err)
		return
	}

	if err := os.WriteFile(DBFile, data, 0644); err != nil {
		log.Printf("Ошибка записи файла: %v", err)
	} else {
		fmt.Println("База данных успешно сохранена.")
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
