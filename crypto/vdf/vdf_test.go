package vdf

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

const (
	N         = 1000
	GOROUTINE = 20
)

var wg = &sync.WaitGroup{}

func TestVDFPool(t *testing.T) {
	jobCh := genJob(100)
	recCh := make(chan [516]byte, 100)
	start := time.Now()
	workerPool(GOROUTINE, jobCh, recCh)
	duration := time.Now().Sub(start)
	log.Println(fmt.Sprintf("VDF computation finished, time spent %s", duration.String()))
	//select {}

}

func genJob(n int) chan [32]byte {
	jobCh := make(chan [32]byte, 200)
	go func() {
		for i := 0; i < n; i++ {
			jobCh <- [32]byte{0xde}
		}
		close(jobCh)
	}()

	return jobCh
}

func workerPool(goroutine int, jobCh <-chan [32]byte, retCh chan<- [516]byte) {
	for i := 0; i < goroutine; i++ {
		go work(i, jobCh, retCh)
	}
}

func work(id int, jobCh <-chan [32]byte, retCh chan<- [516]byte) {
	idx := 0
	for job := range jobCh {
		idx++
		fmt.Println("count:", idx)
		vdf := NewVDF(200, job)
		vdf.GetOutputChannel()
		vdf.Execute()
	}

}

func TestGenerateVDFTime(t *testing.T) {
	input := [32]byte{0xde}
	start := time.Now()
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			vdf := NewVDF(2, input)
			//outputChannel := vdf.GetOutputChannel()
			vdf.GetOutputChannel()
			vdf.Execute()
			defer wg.Done()
		}()
	}
	wg.Wait()

	duration := time.Now().Sub(start)
	//output := <-outputChannel

	//log.Println(fmt.Sprintf("VDF computation finished, result is  %s", hex.EncodeToString(output[:])))
	log.Println(fmt.Sprintf("VDF computation finished, time spent %s", duration.String()))

	//inputVDF, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	//
	//var vdfBytes [516]byte
	//copy(vdfBytes[:], inputVDF)
	//
	//start = time.Now()
	//result := vdf.Verify(vdfBytes)
	//duration = time.Now().Sub(start)
	//
	//log.Println(fmt.Sprintf("VDF verification finished, time spent %s", duration.String()))
	//assert.Equal(t, true, result, "failed verifying vdf proof")

}
