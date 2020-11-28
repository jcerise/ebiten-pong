package main

import "github.com/hajimehoshi/ebiten/v2"

// Position represents a single point in 2D space within the game
type Position struct {
	X float64
	Y float64
}

// Velocity represents the movement in the X and Y directions simultaneously
type Velocity struct {
	vx float64
	vy float64
}

// Paddle is a hard object (has collision), which moves along the Y-axis
type Paddle struct {
	image     *ebiten.Image
	position  Position
	width     int
	height    int
	lastState Direction
}

// Ball is a hard object (has collision), has a velocity, and moves along the X and Y-axiss
type Ball struct {
	image    *ebiten.Image
	position Position
	width    int
	height   int
	velocity Velocity
}
