package DxValue

import (
	"testing"
	"fmt"
)

func TestDxArray_JsonParserFromByte(t *testing.T) {
	arr := NewArray()
	_,err := arr.JsonParserFromByte([]byte(`[  32  ,  "2342"  ,[ 2 , true , false  ,{ "Name" : "DxSoft" , "Age"  :  32 } ] ]`),false,false)
	if err == nil {
		fmt.Println(arr.ToString())
	}else{
		fmt.Println("Paser Error")
	}
}

func TestDxArray_LoadJsonFile(t *testing.T) {
	var v DxValue
	v.LoadJsonFile("DataProxy.config.json",true,false)
	fmt.Println(v.AsString())
}

func TestDxArray_SaveJsonFile(t *testing.T) {
	var v DxValue
	v.LoadJsonFile("DataProxy.config.json",true,false)
	if rec,_ := v.AsRecord();rec != nil{
		if arr := rec.AsArray("list");arr!=nil{
			arr.SaveJsonFile("d:\\1.json",true)
			fmt.Println("SaveJsonOK")
		}
	}
}

func TestDxArray_Append(t *testing.T) {
	arr := NewArray()
	arr.Append(2,"@3423",23,"asdfasdf")
	fmt.Println(arr.String())
}
