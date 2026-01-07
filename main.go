package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"sync/atomic"
)

/*Synchronization mechanism for all the goroutines and the main() which
  does the directory traversal and assignes work to the goroutines. When
  all the traversal is done, main() sets done = true. */
var done atomic.Bool
var signal sync.WaitGroup


func grep(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	line_no := 1
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), os.Args[2]) {
			fmt.Println("line no is:", line_no)
			fmt.Println(scanner.Text())
			fmt.Println("file path is:", fileName)
		}
		line_no++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func workers(ch <- chan string, workerID int) {
	for {
		select {
		case msg := <- ch:
		//fmt.Println(msg, "received from workerID:", workerID)
			grep(msg)
	
		default:	//fmt.Println("no msg received", len(ch), ": wid", workerID)
					if len(ch) == 0 {
						if done.Load() == true {
							goto end
						} else {
							time.Sleep(10 * time.Microsecond)
							continue
						}
					}	
		}
	}	
end:	
	signal.Done()
}

func createGoroutines(ch chan string) {
	// since it is CPU-intensive workload, create only finite no. of goroutines
	// prefrebely equivalent to the no. of cores in the machine or motherboard
	totalGoroutines := runtime.GOMAXPROCS(runtime.NumCPU())
	totalGoroutines -= 1 // one core must be occupied by the main()

	for count := 0; count < totalGoroutines; count++ {
		signal.Add(1)
		go workers(ch, count)
	}
}

func main() {
	if len(os.Args) > 3 {
		fmt.Println("Supply only two arguments: the directory path and the string to search")
		os.Exit(1)
	}

	ch := make(chan string, 1024)

	done.Store(false)

	createGoroutines(ch)

	startPath := os.Args[1]  // assign the source directory path
	traverseDir(startPath, ch)
	done.Store(true)	// full traversal is over
	signal.Wait()
}

func traverseDir(path string, ch chan <- string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", path, err)
		return
	}
	
	count := 0
	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			//fmt.Printf("<Dir>      %s\n", fullPath)
			traverseDir(fullPath, ch) // Recursive call
		} else {
			//fmt.Printf("<File>     %s\n", fullPath)
			ch <- fullPath	// place the file into the channel so that a goroutine
							// could pick it up from there
			count++
			
			if count == 1024 {
				count = 0
				runtime.Gosched()
			}
		}
	}
}
