module github.com/sandwich-go/xconf

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.13
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/mitchellh/mapstructure v1.4.3
	github.com/smartystreets/goconvey v1.7.2
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

//replace github.com/mitchellh/mapstructure => ../mapstructure
