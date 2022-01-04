package kv

type ContentChange = func(string, []byte)
type WatchError = func(error)

//go:generate optiongen --option_with_struct_name=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"OnContentChange": ContentChange(nil),
		"OnWatchError":    WatchError(nil),
	}
}
