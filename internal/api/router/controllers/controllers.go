package controllers

type HTTPSuccess struct {
	Success bool `json:"success"        example:"true"`
	Data    Data `json:"data,omitempty" example:"key:value" swaggertype:"object,string"`
}

type HTTPError struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error"   example:"Not Found"`
}

type Data map[string]interface{}
