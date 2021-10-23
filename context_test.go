package learngocontext

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextE, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextA.Value("b")) // nil
	fmt.Println(contextE.Value("e")) // E
	fmt.Println(contextE.Value("d")) // nil
	fmt.Println(contextE.Value("b")) // B
	fmt.Println(contextE.Value("c")) // nil
}

func CreateCounter(ctx context.Context, group *sync.WaitGroup) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				group.Add(1)
				destination <- counter
				counter++
				time.Sleep(1 * time.Second)
				group.Done()
			}
		}
	}()

	fmt.Println("init")
	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	group := &sync.WaitGroup{}
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx, group)
	for n := range destination {
		fmt.Println("counter", n)
		if n == 10 {
			break
		}
	}

	cancel()
	group.Wait()
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	group := &sync.WaitGroup{}
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounter(ctx, group)
	for n := range destination {
		fmt.Println("counter", n)
	}

	group.Wait()
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	group := &sync.WaitGroup{}
	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel()

	destination := CreateCounter(ctx, group)
	for n := range destination {
		fmt.Println("counter", n)
	}

	group.Wait()
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
