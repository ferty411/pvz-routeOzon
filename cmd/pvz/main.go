package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"pvz/internal/service"
	"pvz/internal/storage"
	"strings"
	"time"
)

func main() {
	repo, err := storage.NewFileStorage("data.json")
	if err != nil {
		fmt.Printf("Критическая ошибка при инициализации хранилища: %v\n", err)
		os.Exit(1)
	}
	svc := service.NewPvzService(repo)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("=== Система управления ПВЗ ===")
	fmt.Println("Введите 'help' для получения списка команд.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		command := strings.ToLower(args[0])

		switch command {
		case "help":
			printHelp()

		case "exit":
			fmt.Println("Завершение работы...")
			return

		case "receive-order":
			if len(args) < 4 {
				fmt.Println("Использование: receive-order <order_id> <customer_id> <YYYY-MM-DD>")
				continue
			}
			orderID, customerID, dateStr := args[1], args[2], args[3]

			expiryDate, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				fmt.Println("Ошибка: неверный формат даты. Используйте YYYY-MM-DD (например, 2026-12-31).")
				continue
			}

			err = svc.AcceptOrder(orderID, customerID, expiryDate)
			if err != nil {
				fmt.Printf("Ошибка при приеме заказа: %v\n", err)
			} else {
				fmt.Println("Заказ успешно принят от курьера.")
			}

		case "return-to-courier":
			if len(args) < 2 {
				fmt.Println("Использование: return-to-courier <order_id>")
				continue
			}
			orderID := args[1]

			err := svc.ReturnCorier(orderID)
			if err != nil {
				fmt.Printf("Ошибка при возврате курьеру: %v\n", err)
			} else {
				fmt.Println("Заказ успешно возвращен курьеру (удален).")
			}

		case "client-action":
			if len(args) < 4 {
				fmt.Println("Использование: client-action <user_id> <issue|return> <order_id_1> [order_id_2...]")
				continue
			}
			customerID := args[1]
			action := strings.ToLower(args[2]) // issue или return
			orderIDs := args[3:]               // Все оставшиеся аргументы — это ID заказов

			err := svc.ProccessClientAction(action, customerID, orderIDs)
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Println("Операция с клиентом успешно выполнена.")
			}

		case "list-orders":
			if len(args) < 2 {
				fmt.Println("Использование: list-orders <user_id> [--limit <N>] [--in-pvz]")
				continue
			}
			userID := args[1]

			fs := flag.NewFlagSet("list-orders", flag.ContinueOnError)
			limit := fs.Int("limit", 0, "Количество последних заказов")
			inPVZ := fs.Bool("in-pvz", false, "Показать только заказы в ПВЗ")

			if err := fs.Parse(args[2:]); err != nil {
				fmt.Println("Ошибка парсинга флагов:", err)
				continue
			}

			orders, err := svc.GetClientOrders(userID, *limit, *inPVZ)
			if err != nil {
				fmt.Printf("Ошибка получения заказов: %v\n", err)
				continue
			}

			fmt.Printf("=== Заказы клиента %s ===\n", userID)
			if len(orders) == 0 {
				fmt.Println("Заказов не найдено.")
			}
			for _, o := range orders {
				// Выводим дату в красивом формате
				fmt.Printf("- ID: %s | Статус: %s | Обновлен: %s\n", o.OrderID, o.Status, o.UpdatedAt.Format(time.DateTime))
			}

		case "list-returns":
			fs := flag.NewFlagSet("list-returns", flag.ContinueOnError)
			page := fs.Int("page", 1, "Номер страницы")
			limit := fs.Int("limit", 10, "Элементов на странице")

			if err := fs.Parse(args[1:]); err != nil {
				fmt.Println("Ошибка парсинга флагов:", err)
				continue
			}

			returns, err := svc.GetReturnsList(*page, *limit)
			if err != nil {
				fmt.Printf("Ошибка получения возвратов: %v\n", err)
				continue
			}

			fmt.Printf("=== Возвраты (Страница %d) ===\n", *page)
			if len(returns) == 0 {
				fmt.Println("Возвратов нет.")
			}
			for _, o := range returns {
				fmt.Printf("- Заказ: %s | Клиент: %s | Возвращен: %s\n", o.OrderID, o.CustomerID, o.UpdatedAt.Format(time.DateTime))
			}

		case "history":
			history, err := svc.GetOrderHistory()
			if err != nil {
				fmt.Printf("Ошибка получения истории: %v\n", err)
				continue
			}

			fmt.Println("=== История заказов ===")
			if len(history) == 0 {
				fmt.Println("База пуста.")
			}
			for _, o := range history {
				// time.TimeOnly выводит только часы, минуты и секунды
				fmt.Printf("[%s] Заказ: %s | Клиент: %s | Статус: %s\n",
					o.UpdatedAt.Format(time.TimeOnly), o.OrderID, o.CustomerID, o.Status)
			}

		default:
			fmt.Println("Неизвестная команда. Введите 'help' для списка команд.")
		}
	}
}

func printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("  return-to-courier <order_id>                         - Вернуть заказ курьеру (удалить)")
	fmt.Println("  client-action <user_id> <issue|return> <order_ids>   - Выдать заказы или принять возврат от клиента")
	fmt.Println("  list-orders <user_id> [--limit <N>] [--in-pvz]       - Получить список заказов клиента")
	fmt.Println("  list-returns [--page <N>] [--limit <N>]              - Получить список возвратов")
	fmt.Println("  history                                              - Показать историю всех изменений")
	fmt.Println("  help                                                 - Показать это сообщение")
	fmt.Println("  exit                                                 - Выйти из программы")
}
