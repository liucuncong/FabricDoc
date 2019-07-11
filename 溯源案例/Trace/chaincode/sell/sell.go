package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strings"
	"time"
)
//销售厂商的智能合约

type Sell struct {}


func(t *Sell) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("OK"))
}

func (t* Sell) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	opttype := args[0]
	assetname := args[1]
	content := args[2]
	if opttype == "putvalue" {
		err := stub.PutState(assetname, []byte(content))
		if err != nil {
			return shim.Error("putvalue fail...")
		}
		return shim.Success([]byte("putvalue sucess"))
	}else if opttype == "getvalue" {
		value, err := stub.GetState(assetname)
		if err != nil {
			return shim.Error("getvalue fail...")
		}
		return shim.Success(value)
	}else if opttype == "gethistory" {
		//通过销售商id进行查询
		//想要查到销售商从哪些加工厂进过牛奶
		keyiter, err := stub.GetHistoryForKey(assetname)
		if err != nil {
			return shim.Error("gethistorykey fail...")
		}
		defer keyiter.Close()
		// 遍历
		var list []string
		var key []string
		for keyiter.HasNext() {
			res, err := keyiter.Next()
			if err != nil {
				return shim.Error("iter next fail...")
			}
			// 读数据
			txvalue := res.Value
			txtime := res.Timestamp
			tm := time.Unix(txtime.Seconds, 0)
			timetext := tm.Format("1978/09/10 10:10:10")
			all := fmt.Sprintf("%s:%s", txvalue, timetext)
			list = append(list, all)// 将所有的key和时间戳都放置在list切片中
			key = append(key, string(txvalue))	// 销售过的所有商品的key
		}
		// 对数组中的第一件溯源
		var history []string
		name := key[0]	// 根据该name去查生成厂商
		history = append(history, assetname)
		history = append(history, name)
		//调用processcc智能合约，获取返回结果
		args := [][]byte{[]byte("invoke"), []byte("gethistory"), []byte(name), []byte("")}
		response := stub.InvokeChaincode("processcc", args, "tracechannel")
		if response.Status != shim.OK {
			return shim.Error("response fail...")
		}
		// json解码
		var buf []string
		err = json.Unmarshal([]byte(response.Payload), buf)
		if err != nil {
			return shim.Error("json unmarshal fail...")
		}
		// 遍历
		for _, v := range buf {
			history = append(history, v)
		}

		// 找其中的一个奶牛场
		tmp := history[1]
		name = strings.Split(tmp, ":")[0]
		args = [][]byte{[]byte("invoke"), []byte("gethistory"), []byte(name), []byte(" ")}
		response = stub.InvokeChaincode("dairycc", args, "tracechannel")
		if response.Status != shim.OK {
			return shim.Error("response fail...")
		}
		// json解码
		err = json.Unmarshal([]byte(response.Payload), buf)
		if err != nil {
			return shim.Error("json unmarshal fail...")
		}
		// 遍历
		for _, v := range buf {
			history = append(history, v)
		}

		myjson, err := json.Marshal(history)
		if err != nil {
			return shim.Error("json marshal fail...")
		}
		return shim.Success(myjson)
	}
	return shim.Error("error")
}

func main() {
	err := shim.Start(new(Sell))
	if err != nil {
		fmt.Println("启动失败...")
	}
}