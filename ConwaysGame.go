package main

import (
	"ConwaysGame/Input"
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

var (
	triangle = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
	pause = false
)

const (
	width              = 900
	height             = 900
	vertexShaderSource = `
    #version 460
    layout (location = 0) in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 450
    out vec4 frag_color;
    void main() {
        frag_color = vec4(0.51, 0.21, 0.71, 1);
    }
` + "\x00"

	rows      = 50
	columns   = 50
	threshold = 0.0
	fps       = 20
)

type cell struct {
	x int
	y int

	alive     bool
	aliveNext bool

	drawable uint32
}

func main() {
	runtime.LockOSThread()
	window := initGlfw()

	defer glfw.Terminate()
	program := initOpenGL()

	keyboardHandler, callback := Input.NewKeyHandler()
	window.SetKeyCallback(callback)
	mouseHandler, mouseCallback := Input.NewMouseHandler()
	window.SetMouseButtonCallback(mouseCallback)

	cells := makeCells()
	for !window.ShouldClose() {
		keyboardHandler.Update()
		mouseHandler.Update()
		if keyboardHandler.SpacePressed() {
			pause = !pause
			drawtext()
		}
		if mouseHandler.Mouse1Pressed() {
			x, y := window.GetCursorPos()
			c := cells[int(x/(width/columns))][int(rows-y/(height/rows))]
			c.alive = !c.alive
			c.aliveNext = c.alive
			c.draw()
		}

		t := time.Now()
		if !pause {
			for x := range cells {
				for _, c := range cells[x] {
					c.checkState(cells)
				}
			}
		}
		draw(cells, window, program, pause)
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}

}
func drawtext() {

}
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func draw(cells [][]*cell, window *glfw.Window, program uint32, draw bool) {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
	for x := range cells {
		for _, c := range cells[x] {
			c.draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func (c *cell) draw() {
	if !c.alive {
		return
	}
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

}
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func makeCells() [][]*cell {
	cells := make([][]*cell, rows, rows)
	for x := 0; x < columns; x++ {
		for y := 0; y < rows; y++ {
			c := newCell(x, y)
			c.alive = rand.Float64() < threshold
			//c.alive = false
			c.aliveNext = c.alive

			cells[x] = append(cells[x], c)
		}
	}
	return cells
}
func newCell(x, y int) *cell {
	points := make([]float32, len(triangle), len(triangle))
	copy(points, triangle)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(columns)
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(rows)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}

	}
	return &cell{drawable: makeVao(points),
		x: x,
		y: y,
	}
}
func (c *cell) checkState(cells [][]*cell) {
	c.alive = c.aliveNext
	c.aliveNext = c.alive

	liveCount := c.liveNeighbors(cells)

	if c.alive {
		if liveCount < 2 {
			c.aliveNext = false
		}

		if liveCount == 2 || liveCount == 3 {
			c.aliveNext = true
		}

		if liveCount > 3 {
			c.aliveNext = false
		}
	} else {
		if liveCount == 3 {
			c.aliveNext = true
		}
	}
}

func (c *cell) liveNeighbors(cells [][]*cell) int {
	var liveCount int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == len(cells) {
			x = 0
		} else if x == -1 {
			x = len(cells) - 1
		}
		if y == len(cells[x]) {
			y = 0
		} else if y == -1 {
			y = len(cells[x]) - 1
		}

		if cells[x][y].alive {
			liveCount++
		}
	}

	add(c.x-1, c.y)   // To the left
	add(c.x+1, c.y)   // To the right
	add(c.x, c.y+1)   // up
	add(c.x, c.y-1)   // down
	add(c.x-1, c.y+1) // top-left
	add(c.x+1, c.y+1) // top-right
	add(c.x-1, c.y-1) // bottom-left
	add(c.x+1, c.y-1) // bottom-right
	return liveCount
}
