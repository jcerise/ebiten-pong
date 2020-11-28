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

type Direction int
type Screen int

const (
	FontSize = 32
	SmallFontSize = FontSize / 2

	ScreenWidth  = 640
	ScreenHeight = 480
	
	Up Direction = iota
	Down
	Neutral

	ScreenTitle Screen = iota
	ScreenGame
	ScreenGameOver
)

var (
	// Player Paddle variables
	playerImage *ebiten.Image
	player Paddle

	// Opponent Paddle variables
	opponentImage *ebiten.Image
	opponent Paddle

	// Ball variables
	ballImage *ebiten.Image
	ball Ball

	// Font variables
	arcadeFont font.Face
	smallArcadeFont font.Face
)

func init() {
	var err error

	rand.Seed(time.Now().UnixNano())

	// Load both the player, opponent, and ball images. These will only be loaded once, and we'll just re-use them
	playerImage, _, err = ebitenutil.NewImageFromFile("player.png")
	if err != nil {
		log.Print("Something went wrong initing image")
		log.Fatal(err)
	}

	opponentImage, _, err = ebitenutil.NewImageFromFile("opponent.png")
	if err != nil {
		log.Fatal(err)
	}

	ballImage, _, err = ebitenutil.NewImageFromFile("ball.png")
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

	opponentWidth, opponentHeight  := opponentImage.Size()
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

func rangeNegative(min, max int) int {
	if min == max {
		return min
	}

	return rand.Intn(max-min+1) + min
}

type Game struct{
	screen Screen

	playerScore int
	opponentScore int
}

func (g *Game) Update() error {
	switch g.screen {
	case ScreenTitle:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.screen = ScreenGame
		}
	case ScreenGame:
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			movePaddle(&player, Down)
		}

		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			movePaddle(&player, Up)
		}

		var ballVelocity Velocity = calculateBallVelocity(ball)

		// Add some "spin" on the ball if the players paddle is moving up or down on collision. Effectively, this just
		// increases the Y velocity negatively or positively depending on how the players paddle is moving when the ball
		// collides with it


		ball.velocity = ballVelocity

		ball.position.X += ball.velocity.vx
		ball.position.Y += ball.velocity.vy

		// Move the opponent. The opponent always just tracks where the ball is, and moves towards along the Y axis. Not
		// very smart, but hey, it's Pong.
		if ball.position.Y > opponent.position.Y + (float64(opponent.height) / 2) {
			// The ball is below the current center of the opponent paddle, move it down to try and intercept
			movePaddle(&opponent, Down)
		} else if ball.position.Y < opponent.position.Y + (float64(opponent.height) / 2) {
			// The ball is above the current center of the opponent paddle, move it up to try and intercept
			movePaddle(&opponent, Up)
		}

		// Check if the ball has gone off the screen. If its gone off the screen on the left (negative X value), the
		// opponent has scored. If its gone off the screen on the right (X larger than screen size), the player has scored
		if ball.position.X <= 0 {
			// The opponent has scored. Increment the opponents score, and reset the ball position ot the center of the
			// screen with a random velocity towards the player
			g.opponentScore += 1
			ball.position = Position{float64((ScreenWidth / 2) - ball.width), float64((ScreenHeight / 2) - ball.height)}
			ball.velocity = Velocity{float64(rangeNegative(-4, 4)), float64(rangeNegative(-4, 4))}
		} else if ball.position.X >= ScreenWidth {
			// The player has scored, Increment the players score, and reset the ball position to the center of the screen
			// with a random velocity towards the opponent
			g.playerScore += 1
			ball.position = Position{float64((ScreenWidth / 2) - ball.width), float64((ScreenHeight / 2) - ball.height)}
			ball.velocity = Velocity{float64(rangeNegative(-4, 4)), float64(rangeNegative(-4, 4))}
		}
	}


	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.screen {
	case ScreenTitle:
		title := "Ebiten Pong"
		xAlign := (ScreenWidth - len(title)*FontSize) / 2
		text.Draw(screen, title, arcadeFont, xAlign, 4 * FontSize, color.White)
		instructions := "Press Space to Begin"
		xAlign = (ScreenWidth - len(instructions)*SmallFontSize) / 2
		text.Draw(screen, instructions, smallArcadeFont, xAlign, 8 * FontSize, color.White)
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
		text.Draw(screen, playerScore, arcadeFont, (ScreenWidth / 2) - (ScreenWidth / 3) , 50, color.White)
		text.Draw(screen, opponentScore, arcadeFont, (ScreenWidth / 2) + (ScreenWidth / 5) , 50, color.White)
	}
}

// checkBallCollision checks if a ball struct is currently colliding with a paddle struct. This is accomplished by
// taking the rectangles around each struct, and checking if there is a gap between any of the four sides of the
// bounding rectangle created by the paddle and ball. If a gap is found, no collision is occurring.
func checkBallCollision(paddle Paddle, ball Ball) bool {
	if paddle.position.X < (ball.position.X + float64(ball.width)) &&
		(paddle.position.X + float64(paddle.width)) > ball.position.X &&
		paddle.position.Y < (ball.position.Y + float64(ball.height)) &&
		(paddle.position.Y + float64(paddle.height)) > ball.position.Y {
		// A collision has been detected between the paddle and ball
		return true
	}

	return false
}

// movePaddle moves a Paddle in the indicated direction (either up or down). When moving up, a check is made to ensure
// the top of the paddle does not go beyond the edge of the screen, and likewise, when moving down, a check is made to
// ensure the bottom of the paddle does not go beyond the bottom edge of the screen. Sets the lastState property of the
// Paddle based on what action the paddle takes
func movePaddle(paddle *Paddle, direction Direction) {
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
		if paddle.position.Y >= float64(ScreenHeight- paddle.height) {
			paddle.position.Y = float64(ScreenHeight - player.height)
			paddle.lastState = Neutral
		} else {
			paddle.position.Y += 4
			paddle.lastState = Down
		}
	}
}

func calculateBallVelocity(ball Ball) Velocity{
	// If the ball has reached an edge, calculate a new velocity for it based on where it hit
	//var bounceAngle float64
	var newVelocity Velocity
	if ball.position.Y <= 0 || (ball.position.Y + float64(ball.height)) >= ScreenHeight {
		newVelocity = Velocity{ball.velocity.vx, -ball.velocity.vy}
	} else if checkBallCollision(player, ball) {
		newVelocity = Velocity{-ball.velocity.vx, ball.velocity.vy}
		if player.lastState == Up {
			newVelocity.vy -= 2
		} else if player.lastState == Down {
			newVelocity.vy += 2
		}
	} else if checkBallCollision(opponent, ball) {
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
