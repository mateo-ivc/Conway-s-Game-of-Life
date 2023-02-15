package Input

import "github.com/go-gl/glfw/v3.3/glfw"

type glfwKeyEvent struct {
	key      glfw.Key
	scancode int
	action   glfw.Action
	mods     glfw.ModifierKey
}

const eventListCap = 10

type glfwKeyEventList []glfwKeyEvent

func newGlfwKeyEventList() *glfwKeyEventList {
	eventList := glfwKeyEventList(make([]glfwKeyEvent, 0, eventListCap))
	return &eventList
}
func (keyEventList *glfwKeyEventList) Callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	event := glfwKeyEvent{key, scancode, action, mods}
	*keyEventList = append(*keyEventList, event)
}

func (keyEventList *glfwKeyEventList) freeze() []glfwKeyEvent {
	frozen := *keyEventList
	*keyEventList = make([]glfwKeyEvent, 0, eventListCap)
	return frozen
}
