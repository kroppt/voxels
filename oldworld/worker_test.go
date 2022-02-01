package oldworld_test

import (
	"sync"
	"testing"

	oldworld "github.com/kroppt/voxels/oldworld"
)

func TestNewWorker(t *testing.T) {
	worker := oldworld.NewWorker(10)
	worker.Quit()
	if worker == nil {
		t.Fatal("got nil worker")
	}
}

func TestAddJob(t *testing.T) {
	worker := oldworld.NewWorker(2)
	ch := make(chan int)
	worker.AddJob(func() {
		ch <- 4
	})
	result := <-ch
	worker.Quit()
	if result != 4 {
		t.Fatal("expected job to be successfully completed")
	}
}

func TestAddJobAfterQuit(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("addjob to quit worker did not panic")
		}
	}()
	worker := oldworld.NewWorker(0)
	worker.Quit()
	worker.AddJob(func() {})
}

func TestWorkerFinishesJobsInOrder(t *testing.T) {
	worker := oldworld.NewWorker(11)
	var wg sync.WaitGroup
	wg.Add(3)
	x := 0
	worker.AddJob(func() {
		x += 2
		wg.Done()
	})
	worker.AddJob(func() {
		x *= 5
		wg.Done()
	})
	worker.AddJob(func() {
		x -= 1
		wg.Done()
	})
	wg.Wait()
	worker.Quit()
	expected := 9
	if x != expected {
		t.Fatalf("expected %v but got %v", expected, x)
	}
}
