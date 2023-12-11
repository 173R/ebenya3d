package window

import "C"
import (
	"ebenya3d/src/consts"
	"ebenya3d/src/input"
	"fmt"
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"math"
	"runtime"
)

/*
	func Init(title string) *glfw.Window {
		if err := glfw.Init(); err != nil {
			panic(err)
		}

		glfw.WindowHint(glfw.Resizable, glfw.True)
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		window, err := glfw.CreateWindow(int(consts.Width), int(consts.Height), title, nil, nil)
		if err != nil {
			panic(err)
		}
		window.MakeContextCurrent()
		window.SetFramebufferSizeCallback(framebufferSizeCallback)
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		window.SetCursorPos(float64(consts.Width)/2, float64(consts.Height)/2)

		return window
	}
*/
var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	input.MouseButton1: glfw.MouseButton1,
	input.MouseButton2: glfw.MouseButton2,
	input.MouseButton3: glfw.MouseButton3,
}

var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: input.MouseButton1,
	glfw.MouseButton2: input.MouseButton2,
	glfw.MouseButton3: input.MouseButton3,
}

type GLFW struct {
	io                  imgui.IO
	window              *glfw.Window
	keyMap              map[glfw.Key]imgui.Key
	pressedMouseButtons [3]bool
	time                float64
}

func New(io imgui.IO) (*GLFW, error) {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to init glfw: %w", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(int(consts.Width), int(consts.Height), consts.Title, nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	window.MakeContextCurrent()

	// TODO: vsync
	// glfw.SwapInterval(1)

	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	//window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPos(float64(consts.Width)/2, float64(consts.Height)/2)

	backend := &GLFW{
		io:     io,
		window: window,
	}
	backend.setKeyMapping()
	backend.setCallbacks()

	return backend, nil
}

func (b *GLFW) Terminate() {
	b.window.Destroy()
	glfw.Terminate()
}

func (b *GLFW) ShouldClose() bool {
	return b.window.ShouldClose()
}

func (b *GLFW) PollEvents() {
	glfw.PollEvents()
}

func (b *GLFW) NewFrame() {
	displaySize := b.DisplaySize()
	b.io.SetDisplaySize(imgui.Vec2{X: displaySize[0], Y: displaySize[1]})

	currentTime := glfw.GetTime()
	if b.time > 0 {
		b.io.SetDeltaTime(float32(currentTime - b.time))
	}
	b.time = currentTime

	// Setup inputs
	if b.window.GetAttrib(glfw.Focused) != 0 {
		x, y := b.window.GetCursorPos()
		b.io.SetMousePos(imgui.Vec2{X: float32(x), Y: float32(y)})
	} else {
		b.io.SetMousePos(imgui.Vec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}

	for i := 0; i < len(b.pressedMouseButtons); i++ {
		down := b.pressedMouseButtons[i] || (b.window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		b.io.SetMouseButtonDown(i, down)
		b.pressedMouseButtons[i] = false
	}
}

func (b *GLFW) SwapBuffers() {
	b.window.SwapBuffers()
}

func (b *GLFW) GetKey(key glfw.Key) glfw.Action {
	return b.window.GetKey(key)
}

func (b *GLFW) SetCursorCallback(callback func(xPos float64, yPos float64)) {
	b.window.SetCursorPosCallback(func(w *glfw.Window, xPos float64, yPos float64) {
		callback(xPos, yPos)
	})
}
func (b *GLFW) SetShouldClose(shouldClose bool) {
	b.window.SetShouldClose(shouldClose)
}

// DisplaySize returns the dimension of the display.
func (b *GLFW) DisplaySize() [2]float32 {
	w, h := b.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

// FramebufferSize returns the dimension of the framebuffer.
func (b *GLFW) FramebufferSize() [2]float32 {
	w, h := b.window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

func (b *GLFW) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.

	b.keyMap = map[glfw.Key]imgui.Key{
		glfw.KeyTab:       imgui.KeyTab,
		glfw.KeyLeft:      imgui.KeyLeftArrow,
		glfw.KeyRight:     imgui.KeyRightArrow,
		glfw.KeyUp:        imgui.KeyUpArrow,
		glfw.KeyDown:      imgui.KeyDownArrow,
		glfw.KeyPageUp:    imgui.KeyPageUp,
		glfw.KeyPageDown:  imgui.KeyPageDown,
		glfw.KeyHome:      imgui.KeyHome,
		glfw.KeyEnd:       imgui.KeyEnd,
		glfw.KeyInsert:    imgui.KeyInsert,
		glfw.KeyDelete:    imgui.KeyDelete,
		glfw.KeyBackspace: imgui.KeyBackspace,
		glfw.KeySpace:     imgui.KeySpace,
		glfw.KeyEnter:     imgui.KeyEnter,
		glfw.KeyEscape:    imgui.KeyEscape,
		glfw.KeyA:         imgui.KeyA,
		glfw.KeyC:         imgui.KeyC,
		glfw.KeyV:         imgui.KeyV,
		glfw.KeyX:         imgui.KeyX,
		glfw.KeyY:         imgui.KeyY,
		glfw.KeyZ:         imgui.KeyZ,

		glfw.KeyLeftControl:  imgui.ModCtrl,
		glfw.KeyRightControl: imgui.ModCtrl,
		glfw.KeyLeftAlt:      imgui.ModAlt,
		glfw.KeyRightAlt:     imgui.ModAlt,
		glfw.KeyLeftSuper:    imgui.ModSuper,
		glfw.KeyRightSuper:   imgui.ModSuper,
	}
}

func (b *GLFW) setCallbacks() {
	b.window.SetMouseButtonCallback(b.mouseButtonChange)
	b.window.SetScrollCallback(b.mouseScrollChange)
	b.window.SetKeyCallback(b.keyChange)
	b.window.SetCharCallback(b.charChange)
}

func framebufferSizeCallback(_ *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (b *GLFW) mouseButtonChange(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonIndex, known := glfwButtonIndexByID[button]

	if known && (action == glfw.Press) {
		b.pressedMouseButtons[buttonIndex] = true
	}
}

func (b *GLFW) mouseScrollChange(_ *glfw.Window, x, y float64) {
	b.io.AddMouseWheelDelta(float32(x), float32(y))
}

func (b *GLFW) keyChange(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	imKey := imgui.Key(key)
	if mapped, ok := b.keyMap[key]; ok {
		imKey = mapped
	}

	b.io.AddKeyEvent(imKey, action == glfw.Press)
}

func (b *GLFW) charChange(_ *glfw.Window, char rune) {
	b.io.AddInputCharactersUTF8(string(char))
	//platform.imguiIO.AddInputCharacters(string(char))
}

// ClipboardText returns the current clipboard text, if available.
func (b *GLFW) ClipboardText() (string, error) {
	return b.window.GetClipboardString(), nil
}

// SetClipboardText sets the text as the current clipboard text.
func (b *GLFW) SetClipboardText(text string) {
	b.window.SetClipboardString(text)
}
