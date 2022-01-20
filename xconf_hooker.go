package xconf

// AtomicSetterProvider Atomic设置配置的方法提供接口，配置层提供改方法供XConf自动设定最新配置
type AtomicSetterProvider interface {
	AtomicSetFunc() func(interface{})
}

// GetOptionUsage 配置层提供改方法供XConf获取配置说明打印Usage
type GetOptionUsage interface {
	GetOptionUsage() string
}

// XConfOptions 配置实现改接口，xconf会自动获取xconf配置覆盖默认设置
type XConfOptions interface {
	XConfOptions() []Option
}
