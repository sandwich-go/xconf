package secconf

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/sandwich-go/xconf/secconf/xxtea"
	"golang.org/x/crypto/openpgp"
)

func DecoderBase64(data []byte) ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(data))
	return io.ReadAll(decoder)
}

func DecoderGZip(data []byte) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	return ioutil.ReadAll(gzReader)
}

func NewDecoderXXTEA(key []byte) CodecFunc {
	return func(data []byte) ([]byte, error) {
		return xxtea.Decrypt(data, key), nil
	}
}

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
		return io.ReadAll(md.UnverifiedBody)
	}
}
