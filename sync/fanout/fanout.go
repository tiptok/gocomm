package fanout

import "sync"

func Merge(size int,cs ...<-chan interface{})<-chan interface{}{
	var wg sync.WaitGroup
	outChan :=make(chan interface{},size)

	outputFunc :=func(c <-chan interface{}){
		for n:=range c{
			outChan <-n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _,c :=range cs{
		go outputFunc(c)
	}
	go func(){
		wg.Wait()
		close(outChan)
	}()
	return outChan
}
