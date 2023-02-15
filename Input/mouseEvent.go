package Input

import "github.com/go-gl/glfw/v3.3/glfw"

type glfwMouseEvent struct {
	button glfw.MouseButton
	action glfw.Action
	mods   glfw.ModifierKey
}

type glfwMouseEventList []glfwMouseEvent

func newGlfwMouseEventList() *glfwMouseEventList {
	eventList := glfwMouseEventList(make([]glfwMouseEvent, 0, eventListCap))
	return &eventList
}

func (mouseEventList *glfwMouseEventList) mouseCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	event := glfwMouseEvent{button, action, mods}
	*mouseEventList = append(*mouseEventList, event)
}

func (mouseEventList *glfwMouseEventList) freeze() []glfwMouseEvent {
	frozen := *mouseEventList
	*mouseEventList = make([]glfwMouseEvent, 0, eventListCap)
	return frozen
}
