package xcmd

import (
	"context"
)

// Executer 命令执行方法
type Executer = func(ctx context.Context, cmd *Command) error

var defaultExecuter = func(ctx context.Context, cmd *Command) error {
	cmd.Usage()
	return ErrHelp
}

// MiddlewareFunc 中间件方法
// cmd *Command : 为当前执行的命令对象
// next : 下一步要执行的方法，可能是下一个中间件或者目标Executer方法
type MiddlewareFunc = func(ctx context.Context, cmd *Command, next Executer) error

// ChainMiddleware middleware chain
func ChainMiddleware(middlewares ...MiddlewareFunc) MiddlewareFunc {
	n := len(middlewares)
	return func(ctx context.Context, cmd *Command, next Executer) error {
		chain := func(currMiddleware MiddlewareFunc, currDispatcher Executer) Executer {
			return func(ctx context.Context, cmd *Command) error {
				return currMiddleware(ctx, cmd, currDispatcher)
			}
		}
		chainHandlerFunc := next
		for i := n - 1; i >= 0; i-- {
			chainHandlerFunc = chain(middlewares[i], chainHandlerFunc)
		}
		return chainHandlerFunc(ctx, cmd)
	}
}

// ChainMiddlewareWithExecuter middleware chain + Executer
func ChainMiddlewareWithExecuter(executer Executer, middlewares ...MiddlewareFunc) Executer {
	m := ChainMiddleware(middlewares...)
	return func(ctx context.Context, cmd *Command) error {
		return m(ctx, cmd, executer)
	}
}
