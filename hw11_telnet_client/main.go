package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Парсинг аргументов командной строки
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		log.Fatal("Usage: go-telnet [--timeout=duration] host port")
	}
	address := net.JoinHostPort(args[0], args[1])

	// Создание клиента Telnet
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	// Подключение к серверу
	if err := client.Connect(); err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer client.Close()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск горутин для отправки и получения данных
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := client.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
			sigChan <- syscall.SIGTERM
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := client.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "...Send error:", err)
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	fmt.Fprintln(os.Stderr, "...EOF")

	client.Close()
	wg.Wait()
}
