package vars

import (
	"strconv"
	"time"
)

// StringValueDelim 默认的数据分割符
var StringValueDelim = ","

func parseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return int(i), err
}
func parseInt8(s string) (int8, error) {
	i, err := strconv.ParseInt(s, 10, 8)
	return int8(i), err
}
func parseInt16(s string) (int16, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	return int16(i), err
}
func parseInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}
func parseInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return int64(i), err
}

func parseUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	return uint(i), err
}
func parseUint8(s string) (uint8, error) {
	i, err := strconv.ParseUint(s, 10, 8)
	return uint8(i), err
}
func parseUint16(s string) (uint16, error) {
	i, err := strconv.ParseUint(s, 10, 16)
	return uint16(i), err
}
func parseUint32(s string) (uint32, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	return uint32(i), err
}
func parseUint64(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	return uint64(i), err
}

func parseFloat64(s string) (float64, error) {
	i, err := strconv.ParseFloat(s, 64)
	return i, err
}

func parseFloat32(s string) (float32, error) {
	i, err := strconv.ParseFloat(s, 32)
	return float32(i), err
}
func parseString(s string) (string, error) { return s, nil }
func parseTimeDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}
func parseBool(s string) (bool, error) { return strconv.ParseBool(s) }

//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceStrig(string,parseString,SetProviderByFieldType,StringValueDelim)"

//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceInt(int,parseInt,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceInt8(int8,parseInt8,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceInt16(int16,parseInt16,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceInt32(int32,parseInt32,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceInt64(int64,parseInt64,SetProviderByFieldType,StringValueDelim)"

//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceUint(uint,parseUint,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceUint8(uint8,parseUint8,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceUint16(uint16,parseUint16,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceUint32(uint32,parseUint32,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceUint64(uint64,parseUint64,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceFloat32(float32,parseFloat32,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceFloat64(float64,parseFloat64,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xslice" "SliceTimeDuration(time.Duration,parseTimeDuration,SetProviderByFieldType,StringValueDelim)"

//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Bool(bool,parseBool)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "String(string,parseString)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Duration(time.Duration,parseTimeDuration)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Int(int,parseInt)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Int8(int8,parseInt8)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Int16(int16,parseInt16)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Int32(int32,parseInt32)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Int64(int64,parseInt64)"

//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Uint(uint,parseUint)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Uint8(uint8,parseUint8)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Uint16(uint16,parseUint16)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Uint32(uint32,parseUint32)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Uint64(uint64,parseUint64)"

//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Float32(float32,parseFloat32)"
//go:generate gotemplate -outfmt gen_%v "../templates/xvar" "Float64(float64,parseFloat64)"

//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapStringString(string,string,parseString,parseString,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapStringInt(string,int,parseString,parseInt,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapIntString(int,string,parseInt,parseString,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapIntInt(int,int,parseInt,parseInt,SetProviderByFieldType,StringValueDelim)"

//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapStringInt64(string,int64,parseString,parseInt64,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapInt64String(int64,string,parseInt64,parseString,SetProviderByFieldType,StringValueDelim)"
//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapInt64Int64(int64,int64,parseInt64,parseInt64,SetProviderByFieldType,StringValueDelim)"

//go:generate gotemplate -outfmt gen_%v "../templates/xmap" "MapStringTimeDuration(string,time.Duration,parseString,parseTimeDuration,SetProviderByFieldType,StringValueDelim)"
