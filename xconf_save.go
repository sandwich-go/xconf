package xconf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func saveToWriter(v map[string]interface{}, ct ConfigType, writer io.Writer) error {
	marshal := GetEncodeFunc(string(ct))
	data, err := marshal(v)
	if err != nil {
		return fmt.Errorf("got error:%s while marshal data", err.Error())
	}
	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("got error:%w while write to writer", err)
	}
	return nil
}

func (x *XConf) MustSaveToFile(fileName string) { panicIfErr(x.SaveToFile(fileName)) }
func (x *XConf) MustSaveToWriter(ct ConfigType, writer io.Writer) {
	panicIfErr(x.SaveToWriter(ct, writer))
}

// SaveVarToFile 将外部传入的valPtr,写入到fileName中，根据文件后缀选择codec
func (x *XConf) MustSaveVarToFile(valPtr interface{}, fileName string) {
	panicIfErr(x.SaveVarToFile(valPtr, fileName))
}

// SaveVarToWriter 将外部传入的valPtr,写入到writer中，类型为ct
func (x *XConf) MustSaveVarToWriter(valPtr interface{}, ct ConfigType, writer io.Writer) {
	panicIfErr(x.SaveVarToWriter(valPtr, ct, writer))
}

func (x *XConf) MustAsBytes(ct ConfigType) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	x.MustSaveToWriter(ct, bytesBuffer)
	return bytesBuffer.Bytes()
}

// SaveTo 将内置解析的数据dump到文件，根据文件后缀选择codec
func (x *XConf) SaveToFile(fileName string) error {
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ""
	}
	marshal := GetEncodeFunc(ext)
	data, err := marshal(x.dataLatestCached)
	if err != nil {
		return fmt.Errorf("got error:%s while marshal data", err.Error())
	}

	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return fmt.Errorf("got error:%s while MkdirAll :%s", err.Error(), filepath.Dir(fileName))
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("got error:%s while WriteFile :%s", err.Error(), fileName)
	}
	return nil
}

// SaveToWriter 将内置解析的数据dump到writer，类型为ct
func (x *XConf) SaveToWriter(ct ConfigType, writer io.Writer) error {
	return saveToWriter(x.dataLatestCached, ct, writer)
}

// SaveVarToWriter 将外部传入的valPtr,写入到writer中，类型为ct
func (x *XConf) SaveVarToWriter(valPtr interface{}, ct ConfigType, writer io.Writer) error {
	if reflect.ValueOf(valPtr).Kind() != reflect.Ptr {
		return errors.New("unsupported type, pass in as ptr")
	}
	data := x.StructMapStructure(valPtr)
	return saveToWriter(data, ct, writer)
}

// SaveVarToFile 将外部传入的valPtr,写入到fileName中，根据文件后缀选择codec
func (x *XConf) SaveVarToFile(valPtr interface{}, fileName string) error {
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := x.SaveVarToWriter(valPtr, ConfigType(extClean(filepath.Ext(fileName))), bytesBuffer)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return fmt.Errorf("got error:%s while MkdirAll :%s", err.Error(), filepath.Dir(fileName))
	}

	if err := ioutil.WriteFile(fileName, bytesBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("got error:%s while WriteFile :%s", err.Error(), fileName)
	}
	return nil
}
