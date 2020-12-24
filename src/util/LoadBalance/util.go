package LoadBalance


func(this *LoadBalance) allDown() bool{
	for _,server:=range this.nodes{
		if server.GetStatus()=="ok" {
			return  false
		}
	}
	return  true
}