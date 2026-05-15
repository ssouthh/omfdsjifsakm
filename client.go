package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Не могу подключиться к серверу:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Введите ваш ник: ")
	nick, _ := reader.ReadString('\n')
	nick = strings.TrimSpace(nick)

	conn.Write([]byte(nick + "\n"))

	serverReader := bufio.NewReader(conn)
	response, _ := serverReader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "ERROR: Ник уже занят!" {
		fmt.Println(response)
		return
	}

	fmt.Println("Подключено к серверу!")
	fmt.Println("Формат отправки: <ник_получателя> <сообщение>")
	fmt.Println("Для выхода введите exit")
	fmt.Println()

	go func() {
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("\nОтключено от сервера")
				os.Exit(0)
			}
			msg = strings.TrimSpace(msg)

			if msg == "OK" || strings.HasPrefix(msg, "OK:") {
				continue
			}
			if strings.HasPrefix(msg, "ERROR") {
				fmt.Printf("\n%s\n> ", msg)
				continue
			}

			fmt.Printf("\n%s\n> ", msg)
		}
	}()

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("До свидания!")
			return
		}

		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 2)
		if len(parts) < 2 {
			fmt.Println("ERROR: Введите ник получателя и сообщение")
			continue
		}

		toNick := parts[0]
		message := parts[1]

		msg := fmt.Sprintf("TO:%s:%s\n", toNick, message)
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Ошибка отправки:", err)
			return
		}
	}
}
