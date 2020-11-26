package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"log"
)

const (
	SCREENWIDTH = 640
	SCREENHEIGHT = 480

	BALLSPEED = 4
)

var (
	// Player variables
	playerImage *ebiten.Image
	player Paddle

	// Opponent variables
	opponentImage *ebiten.Image
	opponentWidth int
	opponentHeight int
	opponentPos Position

	// Ball variables
	ballImage *ebiten.Image
	ball Ball

)

func init() {
	var err error

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

	// Set the starting positions for the player (Left side of screen), the opponent (right side of screen), and the
	// ball (center of screen)
	playerPos := Position{0, 0}
	playerWidth, playerHeight := playerImage.Size()
	player = Paddle{playerImage, playerPos, playerWidth, playerHeight}

	opponentWidth, opponentHeight  = opponentImage.Size()
	opponentPos = Position{float64(SCREENWIDTH - opponentWidth), 0}

	ballWidth, ballHeight := ballImage.Size()
	ballPos := Position{float64((SCREENWIDTH / 2) - ballWidth), float64((SCREENHEIGHT / 2) - ballHeight)}
	ballVelocity := Velocity{
		4,
		2,
	}
	ball = Ball{ballImage, ballPos, ballWidth, ballHeight, ballVelocity}
}

type Game struct{}

func (g *Game) Update() error {

	var movingUp = false
	var movingDown = false

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if player.position.Y >= float64(SCREENHEIGHT - player.height) {
			player.position.Y = float64(SCREENHEIGHT - player.height)
		} else {
			player.position.Y += 4
			movingDown = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if player.position.Y <= 0 {
			player.position.Y = 0
		} else {
			player.position.Y -= 4
			movingUp = true
		}
	}

	// If the ball has reached an edge, calculate a new velocity for it based on where it hit
	//var bounceAngle float64
	var newVelocity Velocity
	if ball.position.Y <= 0 || (ball.position.Y + float64(ball.height)) >= SCREENHEIGHT {
		newVelocity = Velocity{ball.velocity.vx, -ball.velocity.vy}
		ball.velocity = newVelocity
	} else if (ball.position.X + float64(ball.height)) >= SCREENWIDTH {
		newVelocity = Velocity{-ball.velocity.vx, ball.velocity.vy}
		ball.velocity = newVelocity
	} else if checkBallCollision(player, ball) {
		newVelocity = Velocity{-ball.velocity.vx, ball.velocity.vy}

		if movingUp {
			newVelocity.vy -= 2
		} else if movingDown {
			newVelocity.vy += 2
		}

		ball.velocity = newVelocity
	}

	ball.position.X += ball.velocity.vx
	ball.position.Y += ball.velocity.vy

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	playerOptions := &ebiten.DrawImageOptions{}
	playerOptions.GeoM.Translate(player.position.X, player.position.Y)
	screen.DrawImage(playerImage, playerOptions)

	opponentOptions := &ebiten.DrawImageOptions{}
	opponentOptions.GeoM.Translate(opponentPos.X, opponentPos.Y)
	screen.DrawImage(opponentImage, opponentOptions)

	ballOptions := &ebiten.DrawImageOptions{}
	ballOptions.GeoM.Translate(ball.position.X, ball.position.Y)
	screen.DrawImage(ballImage, ballOptions)
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREENWIDTH, SCREENHEIGHT
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(SCREENWIDTH, SCREENHEIGHT)
	ebiten.SetWindowTitle("Ebiten Pong")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
