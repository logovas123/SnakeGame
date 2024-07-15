package main

import "go.mod/src/snake"

func main() {
	s := snake.NewSnake()
	s.Reset()
	g := snake.NewGame()
	g.SetSnake(s)
	g.Run()
}
