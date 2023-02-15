package Input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type MouseHandler struct {
	State         map[glfw.MouseButton]bool
	PreviousState map[glfw.MouseButton]bool

	mouseEventList *glfwMouseEventList
}

func NewMouseHandler() (*MouseHandler, glfw.MouseButtonCallback) {
	h := &MouseHandler{
		State:          make(map[glfw.MouseButton]bool),
		PreviousState:  make(map[glfw.MouseButton]bool),
		mouseEventList: newGlfwMouseEventList(),
	}
	return h, h.getCallback()
}

func (h *MouseHandler) process(events []glfwMouseEvent) {
	for _, event := range events {
		h.setState(event.button, event.action)
	}
}

func (h *MouseHandler) setState(button glfw.MouseButton, action glfw.Action) {
	switch action {
	case glfw.Press:
		h.State[button] = true
		// fmt.Println("Key: ", key, " pressed")
	case glfw.Release:
		h.State[button] = false
		// fmt.Println("Key: ", key, " released")
	}

}

func (h *MouseHandler) getCallback() glfw.MouseButtonCallback {
	if h.mouseEventList == nil {
		return nil
	}
	return h.mouseEventList.mouseCallback
}

func (h *MouseHandler) Update() {
	h.PreviousState, h.State = h.State, make(map[glfw.MouseButton]bool)
	mouseEvent := h.mouseEventList.freeze()
	h.process(mouseEvent)
}

// ======== Helper Function ========

func (h *MouseHandler) IsButtonDown(button glfw.MouseButton) bool {
	return h.State[button]
}

func (h *MouseHandler) Mouse1Pressed() bool {
	//fmt.Println(h.State[glfw.KeySpace])
	return h.State[glfw.MouseButton1]
}
