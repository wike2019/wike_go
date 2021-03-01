package Result


type  Any interface{}
type ErrorResult struct {
	data []Any
	err error
}
func(this *ErrorResult) Unwrap() []Any {
	if this.err!=nil{
		panic(this.err.Error())
	}
	return this.data
}
func(this *ErrorResult) Unwrap_Or(v []Any) []Any {
	if this.err!=nil{
		return v
	}
	return this.data
}
func(this *ErrorResult) Unwrap_Or_Else(f func() []Any) []Any {
	if this.err!=nil{
		return f()
	}
	return this.data
}
func Result(vs ...Any) *ErrorResult {
	for _,item :=range vs{
		if e,ok:=item.(error);ok{
			return &ErrorResult{nil,e}
		}
	}
	return &ErrorResult{vs,nil}
}