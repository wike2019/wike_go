package ants_service

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

type AntsCtl struct {
	*ants.Pool
	*sync.WaitGroup
	Total int
	Ok    int
	err   error
	Fail  int
}

func NewPool(size int) (*AntsCtl, error) {
	pool, err := ants.NewPool(size)
	return &AntsCtl{
		Pool:      pool,
		WaitGroup: &sync.WaitGroup{},
	}, err
}
func (this *AntsCtl) SetTotal(total int) {
	this.Total = total
	this.Add(this.Total)
}
func (this *AntsCtl) Submit(task func() error) error {
	return this.Pool.Submit(func() {
		defer this.Done()
		err := task()
		if err != nil {
			this.Fail++
			if this.err == nil {
				this.err = err
			}
			return
		}
		this.Ok++
	})
}
func (this *AntsCtl) Error() error {
	return this.err
}
