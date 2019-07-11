package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
)

// 自定义结构体
type DairyFarm struct {
}
// 必须要实现init和invoke方法
func (t * DairyFarm) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("init OK"))
}

func (t * DairyFarm)  Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var opttype = args[0]
	var assetname = args[1]
	var content = args[2]

	// 判断用户要执行的操作
	if opttype == "putvalue" {
		stub.PutState(assetname, []byte(content))
		return shim.Success([]byte("success put " + content))
	}else if opttype == "getvalue" {
		// 查询数据
		value, err := stub.GetState(assetname)
		if err != nil {
			return shim.Error("getstate error")
		}
		return shim.Success(value)
	}else if opttype == "gethistory" {
		keyiter, err := stub.GetHistoryForKey(assetname)
		defer keyiter.Close()
		// 出错退出
		if err != nil {
			return shim.Error("gethistory error")
		}
		// 遍历
		var mylist []string
		for keyiter.HasNext() {
			// 取出当前节点
			res, err := keyiter.Next()
			if err != nil {
				return shim.Error("gethistory call next error")
			}
			// 将当前节点中存储的值取出
			// txid := res.TxId
			// statud := res.IsDelete
			txvalue := res.Value
			txtime := res.Timestamp
			// 时间戳格式化, 得到总秒数
			tm := time.Unix(txtime.Seconds, 0)
			// 转换为常用时间格式
			datastr := tm.Format("2018-10-1 00:23:11")
			// 将数据拼接到一起
			all := fmt.Sprintf("%s:%s", txvalue, datastr)
			// 存储该字符串
			mylist = append(mylist, all)
		}
		// 数据打包
		jsontext, err := json.Marshal(mylist)
		if err != nil {
			return shim.Error("json marshal error")
		}
		return shim.Success(jsontext)
	}else{
		return shim.Success([]byte("目前还没有这个操作..."))
	}
}

// main函数
func main() {
	err := shim.Start(new (DairyFarm))
	if err != nil {
		fmt.Printf("启动失败...")
	}
}