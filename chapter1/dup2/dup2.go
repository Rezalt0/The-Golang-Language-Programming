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

var (
	F bool
)

func instructionMessage(b bool) {
	if b {
		fmt.Println("Запуск цикла чтения данных из stdin")
		fmt.Println("Нажмите Ctrl+C, чтобы прервать цикл чтения, и перейти к циклу подсчёта.")
		fmt.Println("После нажатия Ctrl+C, нажмите Ввод")
	}
}
func chanWaitSIGINT() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	instructionMessage(F)
	time.Sleep(2 * time.Second)
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if F {
		fmt.Println("введите данные:")
	} else {
		fmt.Println("чтение данных из файлов")
	}
	if err != nil {
		fmt.Printf("Error run command %v", err)
	}
	return sigChan
}

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		F = true
		countLines(os.Stdin, counts)
	} else {
		F = false
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	sigChan := chanWaitSIGINT()
	input := bufio.NewScanner(f)
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
}
