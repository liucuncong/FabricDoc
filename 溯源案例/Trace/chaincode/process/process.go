package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"time"
)

type Process struct{}


func (t* Process) Init(stub shim.ChaincodeStubInterface) peer.Response{
	return shim.Success([]byte("初始化完成..."))
}

func(t *Process) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	opttype := args[0]
	assetname := args[1]
	content := args[2]

	// 判读用户要进行的操作
	if opttype == "putvalue" {
		stub.PutState(assetname, []byte(content))
		shim.Success([]byte("put value sucess:" + content))
	} else if opttype == "getvalue" {
		value, err := stub.GetState(assetname)
		if err != nil {
			return shim.Error("get value fail")
		}
		return shim.Success(value)
	} else if opttype == "gethistory" {
		keyiter, err := stub.GetHistoryForKey(assetname)
		if err != nil {
			return shim.Error("get history fail...")
		}
		// 遍历
		var list []string
		for keyiter.HasNext() {
			// 获取当前值
			res, err := keyiter.Next()
			if err != nil {
				return shim.Error("get next value fail")
			}
			// 将数据取出
			// txID := res.TxId
			txvalue := res.Value
			// txstatus := res.IsDelete
			txTime := res.Timestamp
			tm := time.Unix(txTime.Seconds, 0)
			datestr := tm.Format("2013-10-11 11:23:45 am")
			// 数据合并
			all := fmt.Sprintf("%s:%s", txvalue, datestr)
			list = append(list, all)
		}
		// 数据打包为json
		jsontext, err := json.Marshal(list)
		if err != nil {
			return shim.Error("json marshal fail")
		}
		return shim.Success(jsontext)
	}
	return shim.Error("error")
}

func main(){
	err := shim.Start(new(Process))
	if err != nil {
		fmt.Println("启动失败...")
	}
}