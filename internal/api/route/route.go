package route

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewRoutes,
		NewTemplate,
        NewDataHandler,
	),
)

type Routes []Route

type Route interface {
	Setup()
}

func NewRoutes(
	template templateRoute,
    dataHandler dataHandlerRoute,
) Routes {
	return Routes{
		template,
        dataHandler,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
