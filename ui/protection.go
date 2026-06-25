package ui

import (
	"reflect"

	"fyne.io/fyne/v2"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// getGLFWWindow reaches into Fyne's internal glfw window wrapper via reflection
// so we can access the native OS window handle.
func getGLFWWindow(w fyne.Window) *glfw.Window {
	v := reflect.ValueOf(w)
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Ptr {
		return nil
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}
	viewport := v.FieldByName("viewport")
	if !viewport.IsValid() || viewport.IsNil() {
		return nil
	}
	return (*glfw.Window)(viewport.UnsafePointer())
}
