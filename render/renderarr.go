package render

import (
	"html/template"
	"log"
	"reflect"
	"runtime/debug"
	"strconv"
	"sync"
)

func RenderArr(arrName string, data any) (arrHTML template.HTML) {
	if reflect.TypeOf(data).Kind() != reflect.Array &&
		reflect.TypeOf(data).Kind() != reflect.Slice {
		log.Println("renderArr: array data was not passed as a slice or array", string(debug.Stack()))
		return ""
	}

	arrVal := reflect.ValueOf(data)
	renderFunc := selectRenderFuncFor(arrVal.Index(0))
	var wg sync.WaitGroup
	elements := make([]template.HTML, arrVal.Len())

	for i := range arrVal.Len() {
		if arrVal.Index(i).Type().String() == reflect.TypeFor[Ignored]().String() {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			index := strconv.Itoa(i)
			element := renderFunc(index, arrVal.Index(i).Interface())
			elements[i] = element
		}(i)
	}
	wg.Wait()

	// Render the html as part of the array list
	arrHTML = RenderType[[]template.HTML]("array.html")(arrName, elements)
	return
}
