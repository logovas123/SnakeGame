package snake

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const (
	GameW = 720.0
	GameH = 720.0
)

type Game struct {
	cv  *canvas.Canvas
	wnd *sdlcanvas.Window

	snake    Snake
	gameOver bool
	speed    int
	needMove bool

	food []Point
}

func NewGame() *Game {
	wnd, cv, err := sdlcanvas.CreateWindow(1080, 750, "Hello, Snake!")
	if err != nil {
		panic(err)
	}

	g := &Game{
		cv:       cv,
		wnd:      wnd,
		speed:    500,
		gameOver: false,
	}

	return g
}

func (g *Game) SetSnake(s *Snake) {
	g.snake = *s
}

func (g *Game) Run() {
	go g.snakeMovement()
	go g.foodGeneration()
	g.renderLoop()
}

func (g *Game) foodGeneration() {
	var foodTimer *time.Timer
	resetTimer := func() {
		foodTimer = time.NewTimer(3 * time.Second)
	}
	resetTimer()

	for {
		<-foodTimer.C
		if !g.gameOver {
			min := 1
			max := 20 - 1
			randX := rand.Intn(max-min) + min
			randY := rand.Intn(max-min) + min
			newPoint := Point{float64(randX), float64(randY)}

			check := true
			if g.snake.IsSnake(newPoint) {
				check = false
			}
			for _, p := range g.food {
				if p.X == newPoint.X && p.Y == newPoint.Y {
					check = false
					break
				}
			}

			if check {
				g.food = append(g.food, newPoint)
			}
		}
		resetTimer()
	}
}

func (g *Game) snakeMovement() {
	var snakeTimer *time.Timer
	var snakeDir Dir = Right
	// var snakeLock sync.Mutex

	resetTime := func() {
		snakeTimer = time.NewTimer(time.Duration(g.speed) * time.Millisecond)
	}
	resetTime()

	//keyboard

	g.wnd.KeyUp = func(code int, rn rune, name string) {
		if code < 79 && code > 82 || g.needMove {
			return
		}

		newDir := snakeDir
		switch code {
		case 80: //left
			newDir = Left
		case 82: //top
			newDir = Bottom
		case 79: //right
			newDir = Right
		case 81: //bottom
			newDir = Top
		}
		if !snakeDir.CheckParallel(newDir) {
			snakeDir = newDir
			g.needMove = true

		}

	}

	//loop
	for {

		<-snakeTimer.C

		if !g.gameOver {
			newPos := snakeDir.Exec(g.snake.Head())
			if newPos.X <= 0 || newPos.X >= 20-1 ||
				newPos.Y <= 0 || newPos.Y >= 20-1 {
				g.gameOver = true
			}

			g.snake.CutIfSnake(newPos)

			//isFood
			isFood := false
			for i := range g.food {
				if newPos.X == g.food[i].X && newPos.Y == g.food[i].Y {
					g.food = append(g.food[:i], g.food[i+1:]...)
					g.snake.Add(newPos)
					g.speed -= 10
					isFood = true
					break
				}
			}

			if !isFood {
				g.snake.Move(snakeDir)
				g.needMove = false
			}
		}
		resetTime()
	}

}

func (g *Game) renderLoop() {

	gameAreaSP := Point{X: 15, Y: 15}
	gameAreaEP := Point{X: 15 + GameW, Y: 15 + GameH}

	cellW := GameW / 20
	cellH := GameH / 20

	font, err := g.cv.LoadFont("./tahoma.ttf")
	if err != nil {
		panic(err)
	}

	g.wnd.MainLoop(func() {
		// clear
		g.cv.ClearRect(0, 0, 1080, 750)

		// render world
		g.cv.BeginPath()
		g.cv.SetFillStyle("#333")
		g.cv.FillRect(gameAreaSP.X, gameAreaSP.Y, gameAreaEP.X-15, gameAreaEP.Y-15)
		g.cv.Stroke()

		g.cv.BeginPath()
		g.cv.SetStrokeStyle("#FFF001")
		g.cv.SetLineWidth(1)
		for i := 0; i < 20+1; i++ {
			g.cv.MoveTo(gameAreaSP.X+float64(i)*cellH, gameAreaSP.Y)
			g.cv.LineTo(gameAreaSP.X+float64(i)*cellH, gameAreaEP.Y)
		}
		for i := 0; i < 20+1; i++ {
			g.cv.MoveTo(gameAreaSP.X, gameAreaSP.Y+float64(i)*cellW)
			g.cv.LineTo(gameAreaEP.X, gameAreaSP.Y+float64(i)*cellW)
		}
		g.cv.Stroke()

		g.cv.BeginPath()
		g.cv.SetFillStyle("#ccc")
		//top
		for i := 0; i < 20; i++ {
			g.cv.FillRect(
				gameAreaSP.X+float64(i)*cellW+1,
				gameAreaSP.Y,
				cellW-1*2,
				cellH)
		}
		//bottom
		for i := 0; i < 20; i++ {
			g.cv.FillRect(
				gameAreaSP.X+float64(i)*cellW+1,
				gameAreaSP.Y+cellH*(20-1),
				cellW-1*2,
				cellH)
		}
		//left
		for i := 0; i < 20; i++ {
			g.cv.FillRect(
				gameAreaSP.X,
				gameAreaSP.Y+float64(i)*cellH+1,
				cellW-1,
				cellH-1*2)
		}
		//right
		for i := 0; i < 20; i++ {
			g.cv.FillRect(
				gameAreaSP.X+cellW*(20-1),
				gameAreaSP.Y+float64(i)*cellH+1,
				cellW,
				cellH-1*2)
		}
		g.cv.Stroke()

		// render snake
		g.cv.BeginPath()
		g.cv.SetFillStyle("#FFF")
		for _, p := range g.snake.Parts {
			g.cv.FillRect(
				gameAreaSP.X+p.X*cellW+1,
				gameAreaSP.Y+p.Y*cellH+1,
				cellW-1*2,
				cellH-1*2,
			)
		}
		g.cv.Stroke()

		//render food
		g.cv.BeginPath()
		g.cv.SetFillStyle("#F15555")
		for _, p := range g.food {
			g.cv.FillRect(
				gameAreaSP.X+p.X*cellW+1,
				gameAreaSP.Y+p.Y*cellH+1,
				cellW-1*2,
				cellH-1*2)
		}
		g.cv.Stroke()

		//render score
		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text := fmt.Sprintf("Score: %d", g.snake.Len())
		g.cv.FillText(text, 720+50, 50)

		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text = fmt.Sprintf("Food: %d", len(g.food))
		g.cv.FillText(text, 720+50, 85)

		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text = fmt.Sprintf("Speed: %d", 500-g.speed)
		g.cv.FillText(text, 720+50, 120)

		if g.gameOver {
			g.cv.BeginPath()
			g.cv.SetFont(font, 30)
			text = fmt.Sprintf("Game Over :(")
			g.cv.FillText(text, 720+50, 175)
		}
	})

}