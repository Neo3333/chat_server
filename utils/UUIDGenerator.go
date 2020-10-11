package utils

import (
	"fmt"
	"math/rand"
)

/**
version 1.0
*/

const (
	DEFAULT_UUID_CNT_CACHE = 64
)


type UUIDGenerator struct {
	Prefix       string
	internalChan chan int
	done         chan struct{}
}

func NewUUIDGenerator(prefix string) *UUIDGenerator {
	gen := &UUIDGenerator{
		Prefix:       prefix,
		internalChan: make(chan int, DEFAULT_UUID_CNT_CACHE),
		done:		  make(chan struct{}),
	}
	gen.startGen()
	return gen
}

func (this *UUIDGenerator) gen() int{
	num := rand.Intn(90000)+10000
	return num
}

//开启 goroutine, 把生成的数字形式的UUID放入缓冲管道
func (this *UUIDGenerator) startGen() {
	go func() {
		for {
			select {
			case <-this.done:
				break
			default:
				this.internalChan <- this.gen()
			}
		}
	}()
}

//获取带前缀的字符串形式的UUID
func (this *UUIDGenerator) Get() string {
	idgen := <-this.internalChan
	return fmt.Sprintf("%s%d", this.Prefix, idgen)
}


func (this *UUIDGenerator) Close(){
	close(this.done)
}

