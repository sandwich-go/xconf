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
	"strings"

	"github.com/sandwich-go/xconf/xutil"
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

// MustSaveToFile 将内置解析的数据dump到文件，根据文件后缀选择codec，如发生错误会panic
func (x *XConf) MustSaveToFile(fileName string) { panicIfErr(x.SaveToFile(fileName)) }

// MustSaveToWriter 将内置解析的数据dump到writer，需指定ConfigType，如发生错误会panic
func (x *XConf) MustSaveToWriter(ct ConfigType, writer io.Writer) {
	panicIfErr(x.SaveToWriter(ct, writer))
}

// SaveVarToWriterAsYAML 将内置解析的数据解析到yaml，带comment
func (x *XConf) SaveVarToWriterAsYAML(valPtr interface{}, writer io.Writer) error {
	s, err := xutil.YAMLWithComments(valPtr, 0, x.cc.TagName, "usage", strings.ToLower)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(s))
	return err
}

// MustSaveVarToFile 将外部传入的valPtr,写入到fileName中，根据文件后缀选择codec，如发生错误会panic
func (x *XConf) MustSaveVarToFile(valPtr interface{}, fileName string) {
	panicIfErr(x.SaveVarToFile(valPtr, fileName))
}

// MustSaveVarToWriter 将外部传入的valPtr,写入到writer中，类型为ct，如发生错误会panic
func (x *XConf) MustSaveVarToWriter(valPtr interface{}, ct ConfigType, writer io.Writer) {
	panicIfErr(x.SaveVarToWriter(valPtr, ct, writer))
}

// MustSaveToBytes 将内置解析的数据以字节流返回，需指定ConfigType
func (x *XConf) MustSaveToBytes(ct ConfigType) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	x.MustSaveToWriter(ct, bytesBuffer)
	return bytesBuffer.Bytes()
}

// SaveToFile 将内置解析的数据dump到文件，根据文件后缀选择codec
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
