package xcmd

import (
	"errors"
	"flag"
	"strings"
)

var ErrNoNeedBind = errors.New("no need bind")

func IsErrHelp(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), flag.ErrHelp.Error())
}
