package DxMsgPack

import (
	"github.com/suiyunonghen/DxValue"
	"time"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
	"fmt"
	"bytes"
	"errors"
	"reflect"
)

var(
	ErrUnSetOnStartArray  = errors.New("Customer Decode has no ArrayEvent ")
	ErrCannotSet		= errors.New("Value can't set")
    interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	stringType = reflect.TypeOf((*string)(nil)).Elem()
	sliceStringPtrType = reflect.TypeOf((*[]string)(nil))
)

func (coder *MsgPackDecoder)decodeStrMapFunc(mp *map[string]interface{})(error)  {
	var(
		code MsgPackCode
		err  error
	)
	if code,err = coder.readCode();err!=nil{
		return err
	}
	maplen := 0
	switch code {
	case CodeMap16:
		if v16,err := coder.readBigEnd16();err!=nil{
			return err
		}else{
			maplen = int(v16)
		}
	case CodeMap32:
		if v32,err := coder.readBigEnd32();err!=nil{
			return err
		}else{
			maplen = int(v32)
		}
	default:
		if code >= CodeFixedMapLow && code<= CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	for i := 0;i<maplen;i++{
		if k,v,err  := coder.decodeStrMapKvRecord(CodeUnkonw);err!=nil{
			return err
		}else{
			(*mp)[k] = v
		}
	}
	return nil
}


func (coder *MsgPackDecoder)decodeStrValueMapFunc(mp *map[string]string)(error)  {
	var(
		code MsgPackCode
		err  error
	)
	if code,err = coder.readCode();err!=nil{
		return err
	}
	maplen := 0
	switch code {
	case CodeMap16:
		if v16,err := coder.readBigEnd16();err!=nil{
			return err
		}else{
			maplen = int(v16)
		}
	case CodeMap32:
		if v32,err := coder.readBigEnd32();err!=nil{
			return err
		}else{
			maplen = int(v32)
		}
	default:
		if code >= CodeFixedMapLow && code<= CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	for i := 0;i<maplen;i++{
		if k,v,err  := coder.decodeStrMapKvRecord(CodeUnkonw);err!=nil{
			return err
		}else{
			(*mp)[k] = fmt.Sprintf("%v",v)
		}
	}
	return nil
}

func (coder *MsgPackDecoder)decodeIntKeyMapFunc(mp *map[int]interface{})(error)  {
	var(
		code MsgPackCode
		err  error
	)
	if code,err = coder.readCode();err!=nil{
		return err
	}
	maplen := 0
	switch code {
	case CodeMap16:
		if v16,err := coder.readBigEnd16();err!=nil{
			return err
		}else{
			maplen = int(v16)
		}
	case CodeMap32:
		if v32,err := coder.readBigEnd32();err!=nil{
			return err
		}else{
			maplen = int(v32)
		}
	default:
		if code >= CodeFixedMapLow && code<= CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	for i := 0;i<maplen;i++{
		if k,v,err  := coder.decodeIntKeyMapKvRecord(CodeUnkonw);err!=nil{
			return err
		}else{
			(*mp)[int(k)] = v
		}
	}
	return nil
}

func (coder *MsgPackDecoder)decodeIntKeyMapFunc64(mp *map[int64]interface{})(error)  {
	var(
		code MsgPackCode
		err  error
	)
	if code,err = coder.readCode();err!=nil{
		return err
	}
	maplen := 0
	switch code {
	case CodeMap16:
		if v16,err := coder.readBigEnd16();err!=nil{
			return err
		}else{
			maplen = int(v16)
		}
	case CodeMap32:
		if v32,err := coder.readBigEnd32();err!=nil{
			return err
		}else{
			maplen = int(v32)
		}
	default:
		if code >= CodeFixedMapLow && code<= CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	for i := 0;i<maplen;i++{
		if k,v,err  := coder.decodeIntKeyMapKvRecord(CodeUnkonw);err!=nil{
			return err
		}else{
			(*mp)[k] = v
		}
	}
	return nil
}

func (coder *MsgPackDecoder)DecodeUnknownMapStd(strcode MsgPackCode)(interface{},error)  {
	if maplen,err := coder.DecodeMapLen(strcode);err!=nil{
		return nil, err
	}else{
		//判断键值，是Int还是str
		if strcode,err = coder.readCode();err!=nil{
			return nil,err
		}
		if strcode.IsInt(){
			if k,v,err := coder.decodeIntKeyMapKvRecord(strcode);err!=nil{
				return nil,err
			}else{
				intMap := make(map[int64]interface{},maplen)
				intMap[k] = v
				for j := 1;j<maplen;j++{
					if k,v,err = coder.decodeIntKeyMapKvRecord(CodeUnkonw);err!=nil{
						return nil,err
					}
					intMap[k] = v
				}
				return intMap,nil
			}
		}else if strcode.IsStr(){
			if k,v,err := coder.decodeStrMapKvRecord(strcode);err!=nil{
				return nil,err
			}else{
				strMap := make(map[string]interface{},maplen)
				strMap[k] = v
				for j := 1;j<maplen;j++{
					if k,v,err = coder.decodeStrMapKvRecord(CodeUnkonw);err!=nil{
						return nil,err
					}
					strMap[k] = v
				}
				return strMap,nil
			}
		}
		return nil,ErrInvalidateMapKey
	}
}


func (coder *MsgPackDecoder)Decode2Interface()(interface{},error)  {
	code,err := coder.readCode()
	if err!=nil{
		return nil,err
	}
	if code.IsStr(){
		if stbt,err := coder.DecodeString(code);err!=nil{
			return nil,err
		}else{
			return DxCommonLib.FastByte2String(stbt),nil
		}
	}else if code.IsFixedNum(){
		return int8(code),nil
	}else if code.IsInt(){
		if i64,err := coder.DecodeInt(code);err!=nil{
			return nil,err
		}else{
			return i64,nil
		}
	}else if code.IsMap(){
		return coder.DecodeUnknownMapStd(code)
	}else if code.IsArray(){
		return coder.DecodeArrayStd(code)
	}else if code.IsBin(){
		if bin,err := coder.DecodeBinary(code);err!=nil{
			return nil,err
		}else{
			return bin,nil
		}
	}else if code.IsExt(){
		if bin,err := coder.DecodeExtValue(code);err!=nil{
			return nil,err
		}else {
			return bin,nil
		}
	} else{
		switch code {
		case CodeTrue:	return true,nil
		case CodeFalse: return false,nil
		case CodeNil:	return nil,nil
		case CodeFloat:
			if v32,err := coder.readBigEnd32();err!=nil{
				return nil,err
			}else{
				return *(*float32)(unsafe.Pointer(&v32)),nil
			}
		case CodeDouble:
			if v64,err := coder.readBigEnd64();err!=nil{
				return nil,err
			}else{
				return *(*float64)(unsafe.Pointer(&v64)),nil
			}

		case CodeFixExt4:
			if code,err = coder.readCode();err!=nil{
				return nil,err
			}
			if int8(code) == -1{
				if ms,err := coder.readBigEnd32();err!=nil{
					return nil,err
				}else{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					return ntime,nil
				}
			}else{
				var mb [5]byte
				if _,err = coder.r.Read(mb[1:]);err!=nil{
					return nil,err
				}
				mb[0] = byte(code)
				return mb[:],nil
			}
		}
	}
	return nil,nil
}


func (coder *MsgPackDecoder)DecodeArrayStd(code MsgPackCode)([]interface{},error)  {
	var (
		err error
		arrlen int
	)
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return nil,err
		}
	}
	if arrlen,err = coder.DecodeArrayLen(code);err!=nil{
		return nil,err
	}
	arr := make([]interface{},arrlen)
	for i := 0;i<arrlen;i++{
		if v,err := coder.Decode2Interface();err != nil{
			return nil,err
		}else{
			arr[i] = v
		}
	}
	return arr,nil
}

func (coder *MsgPackDecoder)DecodeArrayCustomer(code MsgPackCode)(error)  {
	if coder.OnStartArray==nil{
		return ErrUnSetOnStartArray
	}
	var (
		err error
		arrlen int
	)
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return err
		}
	}
	if arrlen,err = coder.DecodeArrayLen(code);err!=nil{
		return err
	}
	arr := coder.OnStartArray(arrlen)
	if arr!=nil && coder.OnParserArrElement != nil{
		for i := 0;i<arrlen;i++{
			if v,err := coder.Decode2Interface();err != nil{
				return err
			}else{
				coder.OnParserArrElement(arr,i,v)
			}
		}
	}
	return nil
}

func (coder *MsgPackDecoder)DecodeArray2StdSlice(code MsgPackCode,arr *[]interface{})(error)  {
	var (
		err error
		arrlen int
	)
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return err
		}
	}
	if arrlen,err = coder.DecodeArrayLen(code);err!=nil{
		return err
	}
	for i := 0;i<arrlen;i++{
		if v,err := coder.Decode2Interface();err != nil{
			return err
		}else{
			(*arr)[i] = v
		}
	}
	return nil
}

func (coder *MsgPackDecoder)decodeStrMapKvRecord(strcode MsgPackCode)(string,interface{}, error)  {
	keybt,err := coder.DecodeString(strcode)
	if err != nil{
		return "",nil,err
	}
	if strcode,err = coder.readCode();err!=nil{
		return "",nil,err
	}
	keyName := DxCommonLib.FastByte2String(keybt)
	if strcode.IsStr(){
		if keybt,err = coder.DecodeString(strcode);err!=nil{
			return "",nil,err
		}
		return keyName,DxCommonLib.FastByte2String(keybt),nil
	}else if strcode.IsFixedNum(){
		return keyName,int8(strcode),nil
	}else if strcode.IsInt(){
		if i64,err := coder.DecodeInt(strcode);err!=nil{
			return "",nil,err
		}else{
			return keyName,int(i64),nil
		}
	}else if strcode.IsMap(){
		if baseV,err := coder.DecodeUnknownMapStd(strcode);err!=nil{
			return "",nil,err
		}else{
			return keyName,baseV,nil
		}
	}else if strcode.IsArray(){
		if arr,err := coder.DecodeArrayStd(strcode);err!=nil{
			return "",nil,err
		}else{
			return keyName,arr,nil
		}
	}else if strcode.IsBin(){
		if bin,err := coder.DecodeBinary(strcode);err!=nil{
			return "",nil,err
		} else{
			return keyName,bin,nil
		}
	}else if strcode.IsExt(){
		if bin ,err := coder.DecodeExtValue(strcode);err!=nil{
			return "",nil,err
		}else{
			return keyName,bin,nil
		}
	}else{
		switch strcode {
		case CodeTrue:	return keyName,true,nil
		case CodeFalse: return keyName,false,nil
		case CodeNil:   return keyName,nil,nil
		case CodeFloat:
			if v32,err := coder.readBigEnd32();err!=nil{
				return "",nil,err
			}else{
				return keyName,*(*float32)(unsafe.Pointer(&v32)),nil
			}
		case CodeDouble:
			if v64,err := coder.readBigEnd64();err!=nil{
				return "",nil,err
			}else{
				return keyName,*(*float64)(unsafe.Pointer(&v64)),nil
			}
		case CodeFixExt4:
			if strcode,err = coder.readCode();err!=nil{
				return "",nil,err
			}
			if int8(strcode) == -1{
				if ms,err := coder.readBigEnd32();err!=nil{
					return "",nil,err
				}else{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					return keyName, ntime,nil
				}
			}else{
				var mb [5]byte
				if _,err = coder.r.Read(mb[1:]);err!=nil{
					return "",nil,err
				}
				mb[0] = byte(strcode)
				return keyName,mb[:],nil
			}
		}
	}
	return "",nil,err
}

func (coder *MsgPackDecoder)decodeIntKeyMapKvRecord(intkeyCode MsgPackCode)(int64,interface{}, error)  {
	intKey,err := coder.DecodeInt(intkeyCode)
	if err != nil{
		return -1,nil,err
	}
	if intkeyCode,err = coder.readCode();err!=nil{
		return -1,nil,err
	}

	if intkeyCode.IsStr(){
		if keybt,err := coder.DecodeString(intkeyCode);err!=nil{
			return -1,nil,err
		}else{
			return intKey,DxCommonLib.FastByte2String(keybt),nil
		}
	}else if intkeyCode.IsFixedNum(){
		return intKey,int8(intkeyCode),nil
	}else if intkeyCode.IsInt(){
		if i64,err := coder.DecodeInt(intkeyCode);err!=nil{
			return -1,nil,err
		}else{
			return intKey,int(i64),nil
		}
	}else if intkeyCode.IsMap(){
		if baseV,err := coder.DecodeUnknownMapStd(intkeyCode);err!=nil{
			return -1,nil,err
		}else{
			return intKey,baseV,nil
		}
	}else if intkeyCode.IsArray(){
		if arr,err := coder.DecodeArrayStd(intkeyCode);err!=nil{
			return -1,nil,err
		}else{
			return intKey,arr,nil
		}
	}else if intkeyCode.IsBin(){
		if bin,err := coder.DecodeBinary(intkeyCode);err!=nil{
			return -1,nil,err
		} else{
			return intKey,bin,nil
		}
	}else if intkeyCode.IsExt(){
		if bin ,err := coder.DecodeExtValue(intkeyCode);err!=nil{
			return -1,nil,err
		}else{
			return intKey,bin,nil
		}
	}else{
		switch intkeyCode {
		case CodeTrue:	return intKey,true,nil
		case CodeFalse: return intKey,false,nil
		case CodeNil:   return intKey,nil,nil
		case CodeFloat:
			if v32,err := coder.readBigEnd32();err!=nil{
				return -1,nil,err
			}else{
				return intKey,*(*float32)(unsafe.Pointer(&v32)),nil
			}
		case CodeDouble:
			if v64,err := coder.readBigEnd64();err!=nil{
				return -1,nil,err
			}else{
				return intKey,*(*float64)(unsafe.Pointer(&v64)),nil
			}

		case CodeFixExt4:
			if intkeyCode,err = coder.readCode();err!=nil{
				return -1,nil,err
			}
			if int8(intkeyCode) == -1{
				if ms,err := coder.readBigEnd32();err!=nil{
					return -1,nil,err
				}else{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					return intKey, ntime,nil
				}
			}else{
				var mb [5]byte
				if _,err = coder.r.Read(mb[1:]);err!=nil{
					return -1,nil,err
				}
				mb[0] = byte(intkeyCode)
				return intKey,mb[:],nil
			}
		}
	}
	return -1,nil,err
}


func (coder *MsgPackDecoder)DecodeUnknownMapCustom(strcode MsgPackCode)(error)  {
	if maplen,err := coder.DecodeMapLen(strcode);err!=nil{
		return  err
	}else{
		//判断键值，是Int还是str
		if strcode,err = coder.readCode();err!=nil{
			return err
		}
		if strcode.IsInt(){
			if k,v,err := coder.decodeIntKeyMapKvRecord(strcode);err!=nil{
				return err
			}else if coder.OnStartMap!=nil{
				intMap := coder.OnStartMap(maplen,false)
				if intMap!=nil && coder.OnParserIntKeyMapKv!=nil{
					coder.OnParserIntKeyMapKv(intMap,k,v)
				}else{
					return nil
				}
				for j := 1;j<maplen;j++{
					if k,v,err = coder.decodeIntKeyMapKvRecord(CodeUnkonw);err!=nil{
						return err
					}
					coder.OnParserIntKeyMapKv(intMap,k,v)
				}
				return nil
			}
		}else if strcode.IsStr(){
			if k,v,err := coder.decodeStrMapKvRecord(strcode);err!=nil{
				return err
			}else if coder.OnStartMap!=nil{
				strMap := coder.OnStartMap(maplen,true)
				if strMap!=nil && coder.OnParserStrMapKv!=nil{
					coder.OnParserStrMapKv(strMap,k,v)
				}else{
					return nil
				}
				for j := 1;j<maplen;j++{
					if k,v,err = coder.decodeStrMapKvRecord(CodeUnkonw);err!=nil{
						return err
					}
					coder.OnParserStrMapKv(strMap,k,v)
				}
				return nil
			}
		}
		return ErrInvalidateMapKey
	}
}

func (coder *MsgPackDecoder)DecodeCustom()(error)  {
	code,err := coder.readCode()
	if err!=nil{
		return err
	}
	if code.IsStr(){
		if stbt,err := coder.DecodeString(code);err!=nil{
			return err
		}else if coder.OnParserNormalValue!=nil{
			coder.OnParserNormalValue(DxCommonLib.FastByte2String(stbt))
		}
		return nil
	}else if code.IsFixedNum(){
		if coder.OnParserNormalValue!=nil{
			coder.OnParserNormalValue(int8(code))
		}
		return nil
	}else if code.IsInt(){
		if i64,err := coder.DecodeInt(code);err!=nil{
			return err
		}else if coder.OnParserNormalValue!=nil{
			coder.OnParserNormalValue(i64)
		}
		return nil
	}else if code.IsMap(){
		return coder.DecodeUnknownMapCustom(code)
	}else if code.IsArray(){
		return coder.DecodeArrayCustomer(code)
	}else if code.IsBin(){
		if bin,err := coder.DecodeBinary(code);err!=nil{
			return err
		}else if coder.OnParserNormalValue!=nil{
			coder.OnParserNormalValue(bin)
		}
		return nil
	}else if code.IsExt(){
		if bin,err := coder.DecodeExtValue(code);err!=nil{
			return err
		}else if coder.OnParserNormalValue!=nil{
			coder.OnParserNormalValue(bin)
		}
		return nil
	} else{
		switch code {
		case CodeTrue:
			if coder.OnParserNormalValue!=nil {
				coder.OnParserNormalValue(true)
			}
			return nil
		case CodeFalse:
			if coder.OnParserNormalValue!=nil {
				coder.OnParserNormalValue(false)
			}
			return nil
		case CodeNil:
			if coder.OnParserNormalValue!=nil {
				coder.OnParserNormalValue(nil)
			}
			return nil
		case CodeFloat:
			if v32,err := coder.readBigEnd32();err!=nil{
				return err
			}else if coder.OnParserNormalValue!=nil{
				coder.OnParserNormalValue(*(*float32)(unsafe.Pointer(&v32)))
			}
			return nil
		case CodeDouble:
			if v64,err := coder.readBigEnd64();err!=nil{
				return err
			}else if coder.OnParserNormalValue!=nil{
				coder.OnParserNormalValue(*(*float64)(unsafe.Pointer(&v64)))
			}
			return nil
		case CodeFixExt4:
			if code,err = coder.readCode();err!=nil{
				return err
			}
			if int8(code) == -1{
				if ms,err := coder.readBigEnd32();err!=nil{
					return err
				}else if coder.OnParserNormalValue!=nil{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					coder.OnParserNormalValue(ntime)
				}
				return nil
			}else if coder.OnParserNormalValue!=nil{
				var mb [5]byte
				if _,err = coder.r.Read(mb[1:]);err!=nil{
					return err
				}
				mb[0] = byte(code)
				coder.OnParserNormalValue(mb[:])
				return nil
			}
		}
	}
	return nil
}

func (coder *MsgPackDecoder)DecodeStand(v interface{})(error)  {
	switch value := v.(type) {
	case *string:
		if strbt,err := coder.DecodeString(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = DxCommonLib.FastByte2String(strbt)
		}
	case *[]interface{}:
		return coder.DecodeArray2StdSlice(CodeUnkonw,value)
	case *DxValue.DxBaseValue:
		switch value.ValueType() {
		case DxValue.DVT_Record:
			rec,_ := value.AsRecord()
			return coder.DecodeStrMap(CodeUnkonw,rec)
		case DxValue.DVT_RecordIntKey:
			rec,_ := value.AsIntRecord()
			return coder.DecodeIntKeyMap(CodeUnkonw,rec)
		case DxValue.DVT_Array:
			arr,_ := value.AsArray()
			return coder.Decode2Array(CodeUnkonw,arr)
		}
	case *time.Time:
		if dt,err := coder.DecodeDateTime_Go(CodeUnkonw);err !=nil{
			return err
		}else{
			*value = dt
		}
	case *int8:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = int8(i64)
		}
	case *int16:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = int16(i64)
		}
	case *int32:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = int32(i64)
		}
	case *int64:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = i64
		}
	case *uint8:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint8(i64)
		}
	case *uint16:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint16(i64)
		}
	case *uint32:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint32(i64)
		}
	case *uint64:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint64(i64)
		}
	case *float32:
		if vf,err := coder.DecodeFloat(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = float32(vf)
		}
	case *float64:
		if vf,err := coder.DecodeFloat(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = vf
		}
	case *bool:
		if code,err := coder.readCode();err!=nil{
			return err
		}else if code == CodeFalse{
			*value = false
		}else if code == CodeTrue{
			*value = true
		}
	case *[]byte:
		if bt,err := coder.DecodeBinary(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = bt
		}
	case *map[string]interface{}:
		return coder.decodeStrMapFunc(value)
	case *map[int]interface{}:
		return coder.decodeIntKeyMapFunc(value)
	case *map[int64]interface{}:
		return coder.decodeIntKeyMapFunc64(value)
	case *map[string]string:
		return coder.decodeStrValueMapFunc(value)
	case *time.Duration:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = time.Duration(i64)
		}
	default:
		v := reflect.ValueOf(value)
		if !v.IsValid() {
			return errors.New("msgpack: Decode(nil)")
		}
		if v.Kind() != reflect.Ptr {
			return fmt.Errorf("msgpack: Decode(nonsettable %T)", value)
		}
		v = v.Elem()
		if !v.IsValid() {
			return fmt.Errorf("msgpack: Decode(nonsettable %T)", value)
		}

		return coder.DecodeValue(v)
	}
	return nil
}

func (coder *MsgPackDecoder)decodeStruct(v reflect.Value)(err error)  {
	maplen := 0
	if maplen,err = coder.DecodeMapLen(CodeUnkonw);err!=nil{
		return err
	}
	strcode := CodeUnkonw
	//判断键值，是Int还是str
	if strcode,err = coder.readCode();err!=nil{
		return err
	}
	if !strcode.IsStr(){
		return errors.New("Struct Can only use String Key")
	}
	if k,v,err := coder.decodeStrMapKvRecord(strcode);err!=nil{
		return nil
	}else{
		structs
		//
		for j := 1;j<maplen;j++{
			if k,v,err = coder.decodeStrMapKvRecord(CodeUnkonw);err!=nil{
				return err
			}
			//strMap[k] = v
		}
	}
}

func (coder *MsgPackDecoder)DecodeValue(v reflect.Value)(error)  {
	typ := v.Type()
	if !v.CanSet(){
		return ErrCannotSet
	}
	switch typ.Kind() {
	case reflect.Bool:
		if code,err := coder.readCode();err!=nil{
			return err
		}else if code == CodeFalse{
			v.SetBool(false)
		}else if code == CodeTrue{
			v.SetBool(true)
		}
	case reflect.Struct:
	case reflect.Map:
	case reflect.Slice:
		elem := typ.Elem()
		switch elem.Kind() {
		case reflect.Uint8:
			if bt,err := coder.DecodeBinary(CodeUnkonw);err!=nil{
				return err
			}else{
				v.SetBytes(bt)
				return nil
			}
		}
		switch elem {
		case stringType:
			if arrlen,err := coder.DecodeArrayLen(CodeUnkonw);err!=nil{
				return err
			}else if arrlen ==-1{
				return nil
			}else{
				ptr := v.Addr().Convert(sliceStringPtrType).Interface().(*[]string)
				ss := setStringsCap(*ptr,arrlen)
				for i := 0; i < arrlen; i++ {
					s, err := coder.DecodeString(CodeUnkonw)
					if err != nil {
						return err
					}
					ss = append(ss, DxCommonLib.FastByte2String(s))
				}
				*ptr = ss
				return nil
			}
		}
	default:
		if v.CanInterface(){
			vt := v.Interface()
			coder.DecodeStand(vt)
		}
	}
	return DxValue.ErrValueType
}

func Unmarshal(data []byte, v...interface{}) error {
	coder := NewDecoder(bytes.NewReader(data))
	for _,vdst := range v{
		if err := coder.DecodeStand(vdst);err!=nil{
			return err
		}
	}
	return nil
}