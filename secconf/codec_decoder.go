package secconf

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"

	"github.com/sandwich-go/xconf/secconf/xxtea"
	"golang.org/x/crypto/openpgp"
)

// NewDecoderMagic 新建一个magic decoder，指定magic字段，解码的时候会检测该字段
func NewDecoderMagic(magic []byte) CodecFunc {
	return func(data []byte) ([]byte, error) {
		if len(magic) == 0 {
			return data, nil
		}
		if !bytes.HasPrefix(data, magic) {
			return data, errors.New("data should have magic prefix:" + string(magic))
		}
		return data[len(magic):], nil
	}
}

// DecoderBase64 base64解码
func DecoderBase64(data []byte) ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(data))
	return readAll(decoder)
}

// DecoderGZip gzip解码
func DecoderGZip(data []byte) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	return readAll(gzReader)
}

// NewDecoderXXTEA 新建xxtea解码器，指定key
func NewDecoderXXTEA(key []byte) CodecFunc {
	return func(data []byte) ([]byte, error) {
		return xxtea.Decrypt(data, key), nil
	}
}

// NewDecoderOpenPGP 新建OpenPGP解码器，指定key
func NewDecoderOpenPGP(secertKeyring []byte) CodecFunc {
	return func(data []byte) ([]byte, error) {
		entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(secertKeyring))
		if err != nil {
			return nil, err
		}
		md, err := openpgp.ReadMessage(bytes.NewBuffer(data), entityList, nil, nil)
		if err != nil {
			return nil, err
		}
		return readAll(md.UnverifiedBody)
	}
}
