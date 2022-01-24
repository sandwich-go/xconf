package xcmdtest

import "time"

//go:generate optiongen
func ConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"HttpAddress": ":3001",
		"Timeouts": (map[string]time.Duration)(
			map[string]time.Duration{
				"read":  time.Duration(10) * time.Second,
				"write": time.Duration(20) * time.Second,
			},
		),
	}
}

//go:generate optiongen
func LogOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"LogLevelTest": 1,
	}
}
