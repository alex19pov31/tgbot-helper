package tgbothelper

import (
	"regexp"
	"strings"
	"time"
)

type routeCallback func(message *MessageData)

type route struct {
	comand       string
	buttonText   string
	containText  string
	pregTemplate string
	callback     routeCallback
}

// RouteGroup - группа роутов
type RouteGroup struct {
	routes   []route
	callback chan routeCallback
}

func (r *route) check(message *MessageData, callback chan routeCallback) {
	if r.buttonText != "" &&
		message.Buttons.GetButtonByText(r.buttonText) != nil {
		callback <- r.callback
		return
	}

	context := message.Message
	if context == "" {
		return
	}

	if r.comand != "" && r.comand == context {
		callback <- r.callback
		return
	}

	if r.containText != "" && strings.Contains(context, r.containText) {
		callback <- r.callback
		return
	}

	if r.pregTemplate != "" {
		if checked, _ := regexp.MatchString(r.pregTemplate, context); checked {
			callback <- r.callback
			return
		}
	}
}

// AddContainButtonRoute - проверка наличия кнопки с указанным текстом
func (rg *RouteGroup) AddContainButtonRoute(buttonText string, callback routeCallback) {
	rg.routes = append(rg.routes, route{buttonText: buttonText, callback: callback})
}

// AddEqualTextRoute - проверка наличия эквиваленого текста
func (rg *RouteGroup) AddEqualTextRoute(comand string, callback routeCallback) {
	rg.routes = append(rg.routes, route{comand: comand, callback: callback})
}

// AddContainTextRoute - проверка на вхождение текста
func (rg *RouteGroup) AddContainTextRoute(containText string, callback routeCallback) {
	rg.routes = append(rg.routes, route{containText: containText, callback: callback})
}

// AddPregTextRoute - проверка текста сообщения по регуляному выражению
func (rg *RouteGroup) AddPregTextRoute(pregTemplate string, callback routeCallback) {
	rg.routes = append(rg.routes, route{pregTemplate: pregTemplate, callback: callback})
}

// Run - запуск роутера
func (rg *RouteGroup) Run(message *MessageData) {
	rg.callback = make(chan routeCallback)
	for _, rt := range rg.routes {
		go func(message *MessageData, rt route, callback chan routeCallback) {
			rt.check(message, callback)
		}(message, rt, rg.callback)
	}

	go func(rg *RouteGroup) {
		select {
		case fn := <-rg.callback:
			fn(message)
		case <-time.NewTicker(time.Second).C:
		}
	}(rg)
}

// NewRouteGroup - создает новую группу роутов
func NewRouteGroup() *RouteGroup {
	return &RouteGroup{}
}
