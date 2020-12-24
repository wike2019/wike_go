package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
)

func (this * EtcdCtl) Lock(lockName string,fn func(params ...interface{}),params ...interface{}) error {
		// 1. 上锁
		// 1.1 创建租约
		lease := clientv3.NewLease(this.EtcdClient)

		leaseGrantResp, err := lease.Grant(context.Background(), 5);
		if err!=nil {
		return err
		}

		leaseId := leaseGrantResp.ID

		// 1.2 自动续约
		// 创建一个可取消的租约，主要是为了退出的时候能够释放
		ctx, cancelFunc := context.WithCancel(context.TODO())
		// 3. 释放租约
		defer cancelFunc()
		defer lease.Revoke(context.TODO(), leaseId)

		this.KeepAliveWithContext(ctx,leaseId)

		// 1.3 在租约时间内去抢锁（etcd里面的锁就是一个key）
		kv := clientv3.NewKV(this.EtcdClient)

		// 创建事物
		txn := kv.Txn(context.TODO())

		//if 不存在key， then 设置它, else 抢锁失败
		txn.If(clientv3.Compare(clientv3.CreateRevision(lockName), "=", 0)).
		Then(clientv3.OpPut(lockName, "g", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockName))

		// 提交事务
		txnResp, err := txn.Commit();
		if err!=nil {
		return err
		}
		if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return fmt.Errorf("锁被占用:")
		}
		// 2. 抢到锁后执行业务逻辑，没有抢到退出
		fn(params)
		return nil
		// 3. 释放锁，步骤在上面的defer，当defer租约关掉的时候，对应的key被回收了
}

