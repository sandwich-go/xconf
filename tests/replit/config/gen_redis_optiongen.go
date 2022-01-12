// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package config

import (
	"sync/atomic"
	"unsafe"
)

// Redis should use NewRedis to initialize it
type Redis struct {
	Endpoints      []string `xconf:"endpoints"`
	Cluster        bool     `xconf:"cluster"`
	TimeoutsStruct Timeouts `xconf:"timeouts_struct"`
}

// ApplyOption apply mutiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *Redis) ApplyOption(opts ...RedisOption) []RedisOption {
	var previous []RedisOption
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// RedisOption option func
type RedisOption func(cc *Redis) RedisOption

// WithRedisEndpoints option func for filed Endpoints
func WithRedisEndpoints(v ...string) RedisOption {
	return func(cc *Redis) RedisOption {
		previous := cc.Endpoints
		cc.Endpoints = v
		return WithRedisEndpoints(previous...)
	}
}

// WithRedisCluster option func for filed Cluster
func WithRedisCluster(v bool) RedisOption {
	return func(cc *Redis) RedisOption {
		previous := cc.Cluster
		cc.Cluster = v
		return WithRedisCluster(previous)
	}
}

// WithRedisTimeoutsStruct option func for filed TimeoutsStruct
func WithRedisTimeoutsStruct(v Timeouts) RedisOption {
	return func(cc *Redis) RedisOption {
		previous := cc.TimeoutsStruct
		cc.TimeoutsStruct = v
		return WithRedisTimeoutsStruct(previous)
	}
}

// NewRedis new Redis
func NewRedis(opts ...RedisOption) *Redis {
	cc := newDefaultRedis()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogRedis != nil {
		watchDogRedis(cc)
	}
	return cc
}

// InstallRedisWatchDog the installed func will called when NewRedis  called
func InstallRedisWatchDog(dog func(cc *Redis)) { watchDogRedis = dog }

// watchDogRedis global watch dog
var watchDogRedis func(cc *Redis)

// newDefaultRedis new default Redis
func newDefaultRedis() *Redis {
	cc := &Redis{}

	for _, opt := range [...]RedisOption{
		WithRedisEndpoints([]string{"192.168.0.1", "192.168.0.2"}...),
		WithRedisCluster(true),
		WithRedisTimeoutsStruct(Timeouts{}),
	} {
		opt(cc)
	}

	return cc
}

// AtomicSetFunc used for XConf
func (cc *Redis) AtomicSetFunc() func(interface{}) { return AtomicRedisSet }

// atomicRedis global *Redis holder
var atomicRedis unsafe.Pointer

// onAtomicRedisSet global call back when  AtomicRedisSet called by XConf.
// use RedisInterface.ApplyOption to modify the updated cc
// if passed in cc not valid, then return false, cc will not set to atomicRedis
var onAtomicRedisSet func(cc RedisInterface) bool

// InstallCallbackOnAtomicRedisSet install callback
func InstallCallbackOnAtomicRedisSet(callback func(cc RedisInterface) bool) {
	onAtomicRedisSet = callback
}

// AtomicRedisSet atomic setter for *Redis
func AtomicRedisSet(update interface{}) {
	cc := update.(*Redis)
	if onAtomicRedisSet != nil && !onAtomicRedisSet(cc) {
		return
	}
	atomic.StorePointer(&atomicRedis, (unsafe.Pointer)(cc))
}

// AtomicRedis return atomic *RedisVisitor
func AtomicRedis() RedisVisitor {
	current := (*Redis)(atomic.LoadPointer(&atomicRedis))
	if current == nil {
		defaultOne := newDefaultRedis()
		if watchDogRedis != nil {
			watchDogRedis(defaultOne)
		}
		atomic.CompareAndSwapPointer(&atomicRedis, nil, (unsafe.Pointer)(defaultOne))
		return (*Redis)(atomic.LoadPointer(&atomicRedis))
	}
	return current
}

// all getter func
func (cc *Redis) GetEndpoints() []string      { return cc.Endpoints }
func (cc *Redis) GetCluster() bool            { return cc.Cluster }
func (cc *Redis) GetTimeoutsStruct() Timeouts { return cc.TimeoutsStruct }

// RedisVisitor visitor interface for Redis
type RedisVisitor interface {
	GetEndpoints() []string
	GetCluster() bool
	GetTimeoutsStruct() Timeouts
}

// RedisInterface visitor + ApplyOption interface for Redis
type RedisInterface interface {
	RedisVisitor
	ApplyOption(...RedisOption) []RedisOption
}
