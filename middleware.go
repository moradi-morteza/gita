package main

type CommandContext struct {
	Command   string
	Args      []string
	CleanArgs []string
	Flags     map[string]bool
}

type Middleware func(*CommandContext) error

var middlewareRegistry = map[string][]Middleware{}

func registerMiddleware(command string, mw Middleware) {
	middlewareRegistry[command] = append(middlewareRegistry[command], mw)
}

func runMiddlewares(ctx *CommandContext) error {
	mws := middlewareRegistry[ctx.Command]
	for _, mw := range mws {
		if err := mw(ctx); err != nil {
			return err
		}
	}
	return nil
}
