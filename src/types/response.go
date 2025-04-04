package types

type Response[M, E any] struct {
	Code  int  `json:"code"`
	Extra bool `json:"extra"`
	Data  M    `json:"data"`
	Other E    `json:"other"`
}

type SimpleResponse[Model any] struct {
	Response[Model, any]
}
