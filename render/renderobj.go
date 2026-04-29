package render

import (
	"html/template"
	"log"
	"reflect"
	"sync"
)

func RenderObj(objName string, data any) (objHTML template.HTML) {
	objVal := reflect.ValueOf(data)
	if objVal.Type().Kind() != reflect.Struct {
		log.Println("renderObj: object data was not passed as a struct")
		return ""
	}

	var wg sync.WaitGroup
	elements := make([]template.HTML, objVal.NumField())
	for i := range objVal.NumField() {
		if objVal.Field(i).Type().String() == reflect.TypeFor[Ignored]().String() {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			renderFunc := selectRenderFuncFor(objVal.Field(i))
			name := objVal.Type().Field(i).Name
			data := objVal.Field(i).Interface()
			if renderFunc == nil {
				log.Println(name)
			}
			element := renderFunc(name, data)
			elements[i] = element
		}(i)
	}
	wg.Wait()

	// Render the html links as part of the object div
	objHTML = RenderType[[]template.HTML]("object.html")(objName, elements)
	return
}
