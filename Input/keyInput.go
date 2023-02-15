package Input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type KeyHandler struct {
	State         map[glfw.Key]bool
	PreviousState map[glfw.Key]bool

	keyEventList *glfwKeyEventList
}

func NewKeyHandler() (*KeyHandler, glfw.KeyCallback) {
	h := &KeyHandler{
		State:         make(map[glfw.Key]bool),
		PreviousState: make(map[glfw.Key]bool),
		keyEventList:  newGlfwKeyEventList(),
	}
	return h, h.getCallback()
}

func (h *KeyHandler) process(events []glfwKeyEvent) {
	for _, event := range events {
		h.setState(event.key, event.action)
	}
}

func (h *KeyHandler) setState(key glfw.Key, action glfw.Action) {
	switch action {
	case glfw.Press:
		h.State[key] = true
		// fmt.Println("Key: ", key, " pressed")
	case glfw.Release:
		h.State[key] = false
		// fmt.Println("Key: ", key, " released")
	}

}

func (h *KeyHandler) getCallback() glfw.KeyCallback {
	if h.keyEventList == nil {
		return nil
	}
	return h.keyEventList.Callback
}

func (h *KeyHandler) Update() {
	h.PreviousState, h.State = h.State, make(map[glfw.Key]bool)
	keyEvents := h.keyEventList.freeze()
	h.process(keyEvents)
}

// ===== Helper Function =====

func (h *KeyHandler) IsKeyDown(key glfw.Key) bool {
	//	fmt.Printf("%t, Key: %s\n", h.State[key], glfw.GetKeyName(key, glfw.GetKeyScancode(key)))
	return h.State[key]
}
func (h *KeyHandler) LeftPressed() bool {
	return h.State[glfw.KeyLeft]
}
func (h *KeyHandler) RightPressed() bool {
	return h.State[glfw.KeyRight]
}
func (h *KeyHandler) UpPressed() bool {
	return h.State[glfw.KeyUp]
}
func (h *KeyHandler) DownPressed() bool {
	return h.State[glfw.KeyDown]
}
func (h *KeyHandler) SpacePressed() bool {
	//fmt.Println(h.State[glfw.KeySpace])
	return h.State[glfw.KeySpace]
}

func (h *KeyHandler) WasKeyDown(key glfw.Key) bool {
	return h.PreviousState[key]
}
func (h *KeyHandler) WasLeftPressed() bool {
	return h.PreviousState[glfw.KeyLeft]
}
func (h *KeyHandler) WasRightPressed() bool {
	return h.PreviousState[glfw.KeyRight]
}
func (h *KeyHandler) WasUpPressed() bool {
	return h.PreviousState[glfw.KeyUp]
}
func (h *KeyHandler) WasDownPressed() bool {
	return h.PreviousState[glfw.KeyDown]
}
func (h *KeyHandler) WasSpacePressed() bool {
	return h.PreviousState[glfw.KeySpace]
}
