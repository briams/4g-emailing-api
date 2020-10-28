package utils

// ResponseData Mensaje de respuesta
type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Href    string `json:"href"`
}

// ResponseMessage contiene los mensajes y errores
type ResponseMessage struct {
	Data     interface{}    `json:"data"`
	Messages []ResponseData `json:"messages"`
	Errors   []ResponseData `json:"errors"`
}

// AddMessage permite agregar un mensaje a ResponseMessage
func (rm *ResponseMessage) AddMessage(code int, message, href string) {
	rd := ResponseData{
		Code:    code,
		Message: message,
		Href:    href,
	}
	rm.Messages = append(rm.Messages, rd)
}

// AddError permite agregar un error a ResponseMessage
func (rm *ResponseMessage) AddError(code int, message, href string) {
	rd := ResponseData{
		Code:    code,
		Message: message,
		Href:    href,
	}
	rm.Errors = append(rm.Errors, rd)
}
