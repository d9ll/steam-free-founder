# 🎮 Steam Freebies Notifier (Go)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Language: Go](https://img.shields.io/badge/Language-Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)

Сверхлегкий монитор бесплатных раздач в Steam. Бот сканирует магазин на наличие временных 100% скидок и мгновенно присылает уведомление в Telegram, чтобы вы не пропустили ни одной игры.

---

## 🔥 Почему это круто?

Большинство ботов просто спамят списком всех бесплатных игр (Free-to-Play). Этот скрипт работает иначе:

* **Smart Filter:** Ищет только платные игры, которые стали бесплатными на время.
* **First Run Silence:** При первом запуске бот просто «запоминает» текущие акции и не спамит старыми уведомлениями.
* **Go Power:** Потребляет меньше 10МБ ОЗУ. Идеально для запуска на Raspberry Pi или слабом VPS.
* **Age Bypass:** Автоматически обходит проверку возраста для игр 18+.

---

## 🛠 Установка и запуск


2. Установите зависимости
go get [github.com/go-telegram-bot-api/telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
go get [github.com/gocolly/colly](https://github.com/gocolly/colly)

или(так надежнее)
```bash
go mod tidy
```
2. Настройка
Откройте main.go и заполните ваши данные:

const (

    BotToken = "ВАШ_ТГ_ТОКЕН"
   
    ChatID   = 000000000 // Ваш ID
   
)

3.Запуск
```bash
go run .
```
📦 Как это работает

Бот использует библиотеку Colly для парсинга страницы поиска Steam по специальному фильтру maxprice=free&specials=1.

Инициализация: При запуске создается карта foundGames.

Мониторинг: Каждые 15 минут (настраивается) бот проверяет обновления.

Уведомление: Если появляется игра, которой нет в карте — вы получаете пуш в Telegram с прямой ссылкой.

🤝 Контрибьютинг

Если у вас есть идеи, как улучшить бота (например, добавить поддержку Docker или базу данных SQLite), смело создавайте Pull Request(я буду рад)
