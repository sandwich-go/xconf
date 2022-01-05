package xflag

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

// PrintDefaults prints, to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set. See the
// documentation for the global function PrintDefaults for more information.
func PrintDefaults(f *flag.FlagSet) {
	f.VisitAll(func(ff *flag.Flag) {
		var b strings.Builder
		fmt.Fprintf(&b, "  -%s", ff.Name) // Two spaces before -; see next two comments.
		varname, usage := flag.UnquoteUsage(ff)
		if varname == "" || varname == "value" {
			if t, ok := ff.Value.(interface{ TypeName() string }); ok {
				varname = t.TypeName()
			}
			if t, ok := ff.Value.(interface{ IsBoolFlag() bool }); ok && t.IsBoolFlag() {
				varname = "bool"
			}

		}
		if len(varname) > 0 {
			b.WriteString(" ")
			b.WriteString(varname)
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if b.Len() <= 4 { // space, space, '-', 'x'.
			b.WriteString("\t")
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			b.WriteString("\n    \t")
		}
		b.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))

		if !isZeroValue(ff, ff.DefValue) {
			if varname == "string" {
				// put quotes on the value
				fmt.Fprintf(&b, " (default %q)", ff.DefValue)
			} else {
				fmt.Fprintf(&b, " (default %v)", ff.DefValue)
			}
		}

		fmt.Fprint(f.Output(), b.String(), "\n")
	})
}

// isZeroValue determines whether the string represents the zero
// value for a flag.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(flag.Value).String()
}
