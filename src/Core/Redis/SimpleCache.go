package Redis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/wike2019/wike_go/src/Result"
	"time"

)

const (
	Serilizer_JSON="json"
	Serilizer_GOB="gob"
)
//简单缓存
type DBGetterFunc func() interface{}
type SimpleCache struct {
	Operation *RedisStringOperation  //操作类
	Expire time.Duration  // 过期时间
	DBGetter  DBGetterFunc // 一旦缓存中没有  DB获取的方法
	Serilizer string  // 序列化方式
	Policy CachePolicy //验证策略
}
func NewSimpleCache(operation *RedisStringOperation, expire time.Duration, serilizer string) *SimpleCache {
	return &SimpleCache{Operation: operation, Expire: expire, Serilizer: serilizer}
}
func (this *SimpleCache)SetCahPolicy(policy CachePolicy)  {
	policy.SetOperation(this.Operation)
	this.Policy=policy
}
// 设置缓存
func(this *SimpleCache) SetCache(key string ,value interface{}){
	this.Operation.Set(key,value,WithExpire(this.Expire)).Unwrap()
}
func(this *SimpleCache) GetCacheForObject(key string,obj interface{})   (ret Result.Any) {
	ret = this.GetCache(key)
	if ret==nil {
		return nil
	}
	if this.Serilizer==Serilizer_JSON{
		err:=json.Unmarshal([]byte(ret.(string)),obj)
		if err!=nil{
			return nil
		}
	}else if   this.Serilizer==Serilizer_GOB{
		var buf =&bytes.Buffer{}
		buf.WriteString(ret.(string))
		dec:=gob.NewDecoder(buf)
		if dec.Decode(obj)!=nil{
			return nil
		}
	}
	return nil
}

func(this *SimpleCache) GetCache(key string) (ret Result.Any){
	if this.Policy!=nil{ //检查策略
		this.Policy.Before(key)
	}
	if this.Serilizer==Serilizer_JSON{
		f:= func()   (ret []Result.Any){
		    obj:= this.DBGetter()
		    if obj==nil{
		    	return []Result.Any{""}
			}
			b,err:=json.Marshal(obj)
			if err!=nil{
				return  []Result.Any{""}
			}
			return []Result.Any{string(b)}
		}
		ret=this.Operation.Get(key).Unwrap_Or_Else(f)[0]
	}else if this.Serilizer==Serilizer_GOB {
		f := func()  (ret []Result.Any){
			obj:= this.DBGetter()
			if obj==nil{
				return  []Result.Any{""}
			}
			var buf= &bytes.Buffer{}
			enc := gob.NewEncoder(buf)
			if err := enc.Encode(obj); err != nil {
				return  []Result.Any{""}
			}
			return []Result.Any{buf.String()}
		}
		ret = this.Operation.Get(key).Unwrap_Or_Else(f)[0]
	}

	   if ret.(string)=="" && this.Policy!=nil {
	   	  this.Policy.IfNil(key,"")

	   }else{
		   this.SetCache(key, ret)
	   }
		return  ret
}