package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)
func chanWaitSIGINT() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	fmt.Println("Запуск цикла чтения данных")
	fmt.Println("Для остановки цикла нажмите Ctrl+C, чтобы прервать цикл чтения, и перейти к циклу подсчёта.")
	fmt.Println("После нажатия Ctrl+C, нажмите Ввод")
	time.Sleep(2 * time.Second)
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error run command %v", err)
	}
	return sigChan
}
func main() {
	sigChan := chanWaitSIGINT()
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("введите данные:")
SCAN:
	for input.Scan() {
		select {
		case <-sigChan:
			fmt.Print("Получен сигнал прерывания цикла чтения данных\nЗапуск цикла подсчёта повторяющихся строк\n")
			break SCAN
		default:
			counts[input.Text()]++
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
