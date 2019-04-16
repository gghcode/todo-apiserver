package api

// Router provides method that return routes.
type Router interface {
	Routes() []RouteInfo
}
