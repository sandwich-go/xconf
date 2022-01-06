package xconf

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (x *XConf) loadReaders(readers ...io.Reader) (map[string]interface{}, error) {
	finalData := make(map[string]interface{})
	for i, v := range readers {
		data, err := x.loadReader(v)
		if err != nil {
			return finalData, fmt.Errorf("got error: %w while loader readers at index:%d", err, i)
		}
		err = x.mergeMap(fmt.Sprintf("reader:%d", i), "reader", data, finalData)
		if err != nil {
			return finalData, fmt.Errorf("got error: %w while merge® reader at index:%d", err, i)
		}
	}
	return finalData, nil
}

func (x *XConf) loadReader(in io.Reader) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(in)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	return data, GetDecodeFunc("")(buf.Bytes(), data)
}

type remoteReader struct {
	url      string
	timeout  time.Duration
	reader   io.Reader
	headerkv []string
}

// NewRemoteReader 返回一个远程Reader，指定url及超时时间
func NewRemoteReader(url string, timeout time.Duration, headerkv ...string) io.Reader {
	return &remoteReader{url: url, timeout: timeout}
}

func (r *remoteReader) do() error {
	// create http client
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return fmt.Errorf("RemoteReader NewRequest with error:%w", err)
	}
	err = kvWithFunc(func(k, v string) bool {
		req.Header.Set(k, v)
		return true
	}, r.headerkv...)
	if err != nil {
		return fmt.Errorf("RemoteReader kv2header error:%w", err)
	}
	client := http.Client{Timeout: r.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("RemoteReader got err:%w while get from:%s ", err, r.url)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("RemoteReader got invalid status code:%d", resp.StatusCode)
	}
	// read response content
	bb, err := readAll(resp.Body)
	if err == nil {
		return fmt.Errorf("RemoteReader got err:%w while read body", err)
	}
	r.reader = bytes.NewReader(bb)
	return nil
}

// Read 实现io.Reader接口
func (r *remoteReader) Read(p []byte) (n int, err error) {
	if r.reader == nil {
		err = r.do()
		if err != nil {
			return 0, err
		}
	}
	return r.reader.Read(p)
}
