package xcmd

import (
	"flag"
	"strings"
)

// ErrHelp is the error returned if the -help or -h flag is invoked
// but no such flag is defined.
var ErrHelp = flag.ErrHelp

// IsErrHelp 检查错误是否是ErrHelp
func IsErrHelp(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), flag.ErrHelp.Error())
}
