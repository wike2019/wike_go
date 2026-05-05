package model

type NoParam struct{}

func (this *NoParam) Dog() *Dog {
	return &Dog{Name: "我是一只狗"}
}
func (this *NoParam) Cat() *Cat {
	return &Cat{Name: "我是一只猫"}
}

// 反例子
//下面的程序是跑不起来，因为依赖注入发现了2个Cat对象不知道应该选择哪个所以程序异常退出了

/**
func (this *NoParam) Cat2() *Cat {
	return &Cat{Name: "我是一只猫-2"}
}
*/
