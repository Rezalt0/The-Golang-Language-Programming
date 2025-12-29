// Задание пока не сделал

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"syscall"
	"time"
)

var (
	F bool
	Input bufio.Scanner
	SigChan chan os.Signal
)

type lineInfo struct {
    count int
    files []string
}

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

func countLines(f *os.File, counts map[string]*lineInfo) {
	if F {
		SigChan = chanWaitSIGINT()
	}
	Input = *bufio.NewScanner(f)
SCAN:
	for Input.Scan() {
		select {
		case <-SigChan:
			fmt.Print("Получен сигнал прерывания цикла чтения данных\nЗапуск цикла подсчёта повторяющихся строк\n")
			break SCAN
		default:
			fileName := filepath.Base(f.Name())
			if counts[Input.Text()] == nil {
				counts[Input.Text()] = &lineInfo{}
			}
			counts[Input.Text()].count++
			found := slices.Contains(counts[Input.Text()].files, fileName)
			if !found {
				counts[Input.Text()].files = append(counts[Input.Text()].files, fileName)
			}
		}
	}
}

func main() {
	counts := make(map[string]*lineInfo)
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
	fmt.Printf("Line\tCountLine\tFilesName\n")
	for line, n := range counts {
		if n.count > 1 {
			fmt.Printf("%s\t%d\t\t%s\n", line, n.count, n.files)
		}
	}
}