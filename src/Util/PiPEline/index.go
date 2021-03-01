package PiPEline

import "sync"
type InChan chan interface{}
type OutChan chan interface{}
type CmdFunc func(args ...interface{}) InChan
type PipeCmdFunc func(in InChan) OutChan
type Pipe struct{
	Cmd CmdFunc
	PipeCmd PipeCmdFunc
	Count int
}

func New()  *Pipe {
	return &Pipe{Count:1}
}
func(this *Pipe) SetCmd(c CmdFunc)  {
	this.Cmd=c
}
func(this *Pipe) SetPipeCmd(c PipeCmdFunc,count int )  {
	this.PipeCmd=c
	this.Count=count
}

func(this *Pipe) Exec(args ...interface{}) OutChan  {
	in:=this.Cmd(args)
	out:=make(OutChan)
	wg:=sync.WaitGroup{}
	for i:=0;i<this.Count;i++{
		getChan:=this.PipeCmd(in)
		wg.Add(1)
		go func(input OutChan) {
			defer wg.Done()
			for v:=range input{
				out<-v
			}
		}(getChan)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}