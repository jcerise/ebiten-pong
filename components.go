package main

import "github.com/hajimehoshi/ebiten/v2"

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	vx float64
	vy float64
}

type Paddle struct {
	image *ebiten.Image
	position Position
	width int
	height int
	lastState Direction
}

type Ball struct {
	image *ebiten.Image
	position Position
	width int
	height int
	velocity Velocity
}
