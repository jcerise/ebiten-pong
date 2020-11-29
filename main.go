package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"time"
)

// Direction represents a direction along the Y-axis, Up, Down, or no movement
type Direction int

// Screen is a state in the game, representing different displays of information, depening on the context
type Screen int

const (
	FontSize      = 32
	SmallFontSize = FontSize / 3

	ScreenWidth  = 640
	ScreenHeight = 480
)

// Directions along the Y-Axis
const (
	Up Direction = iota
	Down
	Neutral
)

// Screen states present in the game
const (
	ScreenTitle Screen = iota
	ScreenCredits
	ScreenGame
	ScreenPlayerWin
	ScreenPlayerLose
)

var (
	// Player Paddle variables
	playerImage *ebiten.Image
	player      Paddle

	// Opponent Paddle variables
	opponentImage *ebiten.Image
	opponent      Paddle

	// Ball variables
	ballImage *ebiten.Image
	ball      Ball

	// Font variables
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func init() {
	var err error

	rand.Seed(time.Now().UnixNano())

	// Load both the player, opponent, and ball images. These will only be loaded once, and we'll just re-use them
	playerImage, _, err = ebitenutil.NewImageFromFile("resources/player.png")
	if err != nil {
		log.Fatal(err)
	}

	opponentImage, _, err = ebitenutil.NewImageFromFile("resources/opponent.png")
	if err != nil {
		log.Fatal(err)
	}

	ballImage, _, err = ebitenutil.NewImageFromFile("resources/ball.png")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72

	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    FontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    SmallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set the starting positions for the player (Left side of screen), the opponent (right side of screen), and the
	// ball (center of screen)
	playerPos := Position{0, 0}
	playerWidth, playerHeight := playerImage.Size()
	player = Paddle{playerImage, playerPos, playerWidth, playerHeight, Neutral}

	opponentWidth, opponentHeight := opponentImage.Size()
	opponentPos := Position{float64(ScreenWidth - opponentWidth), 0}
	opponent = Paddle{opponentImage, opponentPos, opponentWidth, opponentHeight, Neutral}

	ballWidth, ballHeight := ballImage.Size()
	ballPos := Position{float64((ScreenWidth / 2) - ballWidth), float64((ScreenHeight / 2) - ballHeight)}
	ballVelocity := Velocity{
		float64(rangeNegative(-4, 4)),
		float64(rangeNegative(-4, 4)),
	}
	ball = Ball{ballImage, ballPos, ballWidth, ballHeight, ballVelocity}
}

// rangeNegative returns a random value, in a given range, which can include negative numbers
func rangeNegative(min, max int) int {
	if min == max {
		return min
	}

	return rand.Intn(max-min+1) + min
}

// Game represents the main game object. All our logic and update loops will occur in this object
type Game struct {
	screen Screen

	// The player and opponent scores
	playerScore   int
	opponentScore int
}

// Update handles each frames logic and state changes
func (g *Game) Update() error {
	switch g.screen {
	case ScreenTitle:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.screen = ScreenGame
		}

		if ebiten.IsKeyPressed(ebiten.KeyC) {
			g.screen = ScreenCredits
		}
	case ScreenCredits:
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.screen = ScreenTitle
		}
	case ScreenGame:
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.movePaddle(&player, Down)
		}

		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.movePaddle(&player, Up)
		}

		// Set a new velocity for the ball. If the ball has not collided with a hard surface, the velocity will remain
		// unchanged.
		var ballVelocity Velocity = g.calculateBallVelocity(ball)

		ball.velocity = ballVelocity

		ball.position.X += ball.velocity.vx
		ball.position.Y += ball.velocity.vy

		// Move the opponent. The opponent always just tracks where the ball is, and moves towards along the Y axis. Not
		// very smart, but hey, it's Pong.
		if ball.position.Y > opponent.position.Y+(float64(opponent.height)/2) {
			// The ball is below the current center of the opponent paddle, move it down to try and intercept
			g.movePaddle(&opponent, Down)
		} else if ball.position.Y < opponent.position.Y+(float64(opponent.height)/2) {
			// The ball is above the current center of the opponent paddle, move it up to try and intercept
			g.movePaddle(&opponent, Up)
		}

		// Check if the ball has gone off the screen. If its gone off the screen on the left (negative X value), the
		// opponent has scored. If its gone off the screen on the right (X larger than screen size), the player has scored
		if ball.position.X <= 0 || ball.position.X >= ScreenWidth {
			if ball.position.X <= 0 {
				// The opponent has scored. Increment the opponents score, and reset the ball position ot the center of the
				// screen with a random velocity towards the player
				g.opponentScore ++
			} else if ball.position.X >= ScreenWidth {
				// The player has scored, Increment the players score, and reset the ball position to the center of the screen
				// with a random velocity towards the opponent
				g.playerScore ++
			}

			// Either way, reset the ball and start a new round
			ball.position = Position{float64((ScreenWidth / 2) - ball.width), float64((ScreenHeight / 2) - ball.height)}
			ball.velocity = Velocity{float64(rangeNegative(-4, 4)), float64(rangeNegative(-4, 4))}
		}

		if g.playerScore >= 3 {
			g.screen = ScreenPlayerWin
		} else if g.opponentScore >= 3 {
			g.screen = ScreenPlayerLose
		}
	case ScreenPlayerWin:
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.screen = ScreenTitle
			g.opponentScore = 0
			g.playerScore = 0
		}
	case ScreenPlayerLose:
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.screen = ScreenTitle
			g.opponentScore = 0
			g.playerScore = 0
		}
	}

	return nil
}

// Draw renders text and graphics to the screen each frame, based on the current game state
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.screen {
	case ScreenTitle:
		lines := []string{"EBITEN PONG", "", "", "", "[Space] - Begin", "", "[C] - Credits"}
		g.printText(lines, arcadeFont, FontSize, screen)
	case ScreenCredits:
		lines := []string{"Credits", "", "", "Programmed by Jeremy Cerise, 2020", "", "Images made by Nicolás A. Ortega (Deathsbreed)", "https://opengameart.org/content/pong-graphics", "", "All works licensed CC BY-SA 2.0"}
		g.printText(lines, smallArcadeFont, SmallFontSize, screen)
	case ScreenGame:
		playerOptions := &ebiten.DrawImageOptions{}
		playerOptions.GeoM.Translate(player.position.X, player.position.Y)
		screen.DrawImage(playerImage, playerOptions)

		opponentOptions := &ebiten.DrawImageOptions{}
		opponentOptions.GeoM.Translate(opponent.position.X, opponent.position.Y)
		screen.DrawImage(opponentImage, opponentOptions)

		ballOptions := &ebiten.DrawImageOptions{}
		ballOptions.GeoM.Translate(ball.position.X, ball.position.Y)
		screen.DrawImage(ballImage, ballOptions)

		playerScore := fmt.Sprintf("%02d", g.playerScore)
		opponentScore := fmt.Sprintf("%02d", g.opponentScore)
		text.Draw(screen, playerScore, arcadeFont, (ScreenWidth/2)-(ScreenWidth/3), 50, color.White)
		text.Draw(screen, opponentScore, arcadeFont, (ScreenWidth/2)+(ScreenWidth/5), 50, color.White)
	case ScreenPlayerWin:
		lines := []string{"YOU WIN!", "", "Press 'R' to play again"}
		g.printText(lines, arcadeFont, FontSize, screen)
	case ScreenPlayerLose:
		lines := []string{"YOU LOSE, TRY AGAIN!", "", "Press 'R' to play again"}
		g.printText(lines, arcadeFont, FontSize, screen)
	}
}

// checkBallCollision checks if a ball struct is currently colliding with a paddle struct. This is accomplished by
// taking the rectangles around each struct, and checking if there is a gap between any of the four sides of the
// bounding rectangle created by the paddle and ball. If a gap is found, no collision is occurring.
func (g *Game) checkBallCollision(paddle Paddle, ball Ball) bool {
	if paddle.position.X < (ball.position.X+float64(ball.width)) &&
		(paddle.position.X+float64(paddle.width)) > ball.position.X &&
		paddle.position.Y < (ball.position.Y+float64(ball.height)) &&
		(paddle.position.Y+float64(paddle.height)) > ball.position.Y {
		// A collision has been detected between the paddle and ball
		return true
	}

	return false
}

// movePaddle moves a Paddle in the indicated direction (either up or down). When moving up, a check is made to ensure
// the top of the paddle does not go beyond the edge of the screen, and likewise, when moving down, a check is made to
// ensure the bottom of the paddle does not go beyond the bottom edge of the screen. Sets the lastState property of the
// Paddle based on what action the paddle takes
func (g *Game) movePaddle(paddle *Paddle, direction Direction) {
	switch direction {
	case Up:
		if paddle.position.Y <= 0 {
			paddle.position.Y = 0
			paddle.lastState = Neutral
		} else {
			paddle.position.Y -= 4
			paddle.lastState = Up
		}
	case Down:
		if paddle.position.Y >= float64(ScreenHeight-paddle.height) {
			paddle.position.Y = float64(ScreenHeight - player.height)
			paddle.lastState = Neutral
		} else {
			paddle.position.Y += 4
			paddle.lastState = Down
		}
	}
}

// calculateBallVelocity checks if the ball has reached a hard surface (either the top or bottom of the screen, or a
// paddle), and calculates a new velocity for the ball based on the surface it collides with. If it has collided with
// a paddle, the last direction of the paddle is used to apply an extra amount of momentum in the same direction,
// allowing the player to "spin" the ball somewhat
func (g *Game) calculateBallVelocity(ball Ball) Velocity {
	// If the ball has reached an edge, calculate a new velocity for it based on where it hit
	//var bounceAngle float64
	var newVelocity Velocity
	if ball.position.Y <= 0 || (ball.position.Y+float64(ball.height)) >= ScreenHeight {
		newVelocity = Velocity{ball.velocity.vx, -ball.velocity.vy}
	} else if g.checkBallCollision(player, ball) {
		newVelocity = Velocity{-ball.velocity.vx, ball.velocity.vy}
		if player.lastState == Up {
			newVelocity.vy -= 2
		} else if player.lastState == Down {
			newVelocity.vy += 2
		}
	} else if g.checkBallCollision(opponent, ball) {
		newVelocity = Velocity{-ball.velocity.vx, ball.velocity.vy}
		if opponent.lastState == Up {
			newVelocity.vy -= 2
		} else if opponent.lastState == Down {
			newVelocity.vy += 2
		}
	} else {
		newVelocity.vx = ball.velocity.vx
		newVelocity.vy = ball.velocity.vy
	}

	return newVelocity
}

// printText takes an array of strings, and prints them, one by one, in vertical order to the screen, centered on the
// X axis
func (g *Game) printText(lines []string, textFont font.Face, fontSize int, screen *ebiten.Image) {
	for i, line := range lines {
		xAlign := (ScreenWidth - len(line)*fontSize) / 2
		text.Draw(screen, line, textFont, xAlign, (i+4)*fontSize, color.White)
	}
}

// Layout defines the size of our game window
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	game := &Game{}
	game.screen = ScreenTitle
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Ebiten Pong")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
