package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	nick string
}

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex
)

func main() {
	fmt.Println("Запуск сервера на порту 8080...")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен! Ждем подключений...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	nick, _ := reader.ReadString('\n')
	nick = strings.TrimSpace(nick)

	mu.Lock()
	if _, exists := clients[nick]; exists {
		mu.Unlock()
		conn.Write([]byte("ERROR: Ник уже занят!\n"))
		return
	}

	client := &Client{
		conn: conn,
		nick: nick,
	}
	clients[nick] = client
	mu.Unlock()

	conn.Write([]byte("OK\n"))

	fmt.Printf("Пользователь %s подключился\n", nick)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			mu.Lock()
			delete(clients, nick)
			mu.Unlock()
			fmt.Printf("Пользователь %s отключился\n", nick)
			return
		}

		message = strings.TrimSpace(message)

		parts := strings.SplitN(message, ":", 3)
		if len(parts) != 3 || parts[0] != "TO" {
			conn.Write([]byte("ERROR: Неверный формат сообщения\n"))
			continue
		}

		toNick := parts[1]
		text := parts[2]

		if toNick == nick {
			conn.Write([]byte("ERROR: Нельзя отправить сообщение самому себе\n"))
			continue
		}

		mu.Lock()
		recipient, exists := clients[toNick]
		mu.Unlock()

		if !exists {
			conn.Write([]byte("ERROR: Пользователь не найден\n"))
			continue
		}

		msg := fmt.Sprintf("От %s: %s\n", nick, text)
		recipient.conn.Write([]byte(msg))

		conn.Write([]byte("OK\n"))
	}
}
