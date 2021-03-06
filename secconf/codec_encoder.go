package secconf

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"

	"github.com/sandwich-go/xconf/secconf/xxtea"
	"golang.org/x/crypto/openpgp"
)

// EncoderBase64 base64编码器
func EncoderBase64(data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	_, err := encoder.Write(data)
	if err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncoderGZip gzip编码器
func EncoderGZip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	writer.Close()
	return buf.Bytes(), err
}

// NewEncoderMagic magic编码器，指定magic，编码时会自动将magic添加到字段头
func NewEncoderMagic(magic []byte) CodecFunc {
	return func(data []byte) ([]byte, error) {
		if len(magic) == 0 {
			return data, nil
		}
		if bytes.HasPrefix(data, magic) {
			return data, errors.New("data should not have magic prefix:" + string(magic))
		}
		return append(magic[:], data[:]...), nil
	}
}

// NewEncoderXXTEA xxtea编码器，指定key
func NewEncoderXXTEA(key []byte) CodecFunc {
	return func(data []byte) ([]byte, error) {
		return xxtea.Encrypt(data, key), nil
	}
}

// NewEncoderOpenPGP OpenPGP编码器，指定key
func NewEncoderOpenPGP(secertKeyring []byte) CodecFunc {
	return func(data []byte) ([]byte, error) {
		entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(secertKeyring))
		if err != nil {
			return nil, err
		}
		buffer := new(bytes.Buffer)
		pgpWriter, err := openpgp.Encrypt(buffer, entityList, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		if _, err = pgpWriter.Write(data); err != nil {
			return nil, err
		}
		if err = pgpWriter.Close(); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
}
