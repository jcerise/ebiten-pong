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

	BALLSPEED = 3
)

var (
	// Player variables
	playerImage *ebiten.Image
	playerHeight int
	playerPos Position

	// Opponent variables
	opponentImage *ebiten.Image
	opponentWidth int
	opponentHeight int
	opponentPos Position

	// Ball variables
	ballImage *ebiten.Image
	ballWidth int
	ballHeight int
	ballPos Position
	ballVelocity Velocity

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
	playerPos = Position{0, 0}
	_, playerHeight = playerImage.Size()

	opponentWidth, opponentHeight  = opponentImage.Size()
	opponentPos = Position{float64(SCREENWIDTH - opponentWidth), 0}

	ballWidth, ballHeight = ballImage.Size()
	ballPos = Position{float64((SCREENWIDTH / 2) - ballWidth), float64((SCREENHEIGHT / 2) - ballHeight)}
	ballVelocity = Velocity{
		2,
		4,
	}
}

type Game struct{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if playerPos.Y >= float64(SCREENHEIGHT - playerHeight) {
			playerPos.Y = float64(SCREENHEIGHT - playerHeight)
		} else {
			playerPos.Y += 4
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if playerPos.Y <= 0 {
			playerPos.Y = 0
		} else {
			playerPos.Y -= 4
		}
	}

	// If the ball has reached an edge, calculate a new velocity for it based on where it hit
	//var bounceAngle float64
	if ballPos.Y <= 0 || (ballPos.Y + float64(ballHeight)) >= SCREENHEIGHT {
		//var relativeIntersectX = (SCREENWIDTH / 2) - ballPos.X
		//var normalizedIntesectX = relativeIntersectX / (SCREENWIDTH / 2)
		//bounceAngle = normalizedIntesectX * 75
		//
		//ballVelocity.vx = BALLSPEED * math.Cos(bounceAngle)
		//ballVelocity.vy = BALLSPEED * -math.Sin(bounceAngle)
		ballVelocity.vy = -ballVelocity.vy
	}

	if ballPos.X <= 0 || ballPos.X + float64(ballWidth) >= SCREENWIDTH {
		//var relativeIntersectY = (SCREENHEIGHT / 2) - ballPos.Y
		//var normalizedIntesectY = relativeIntersectY / (SCREENHEIGHT / 2)
		//bounceAngle = normalizedIntesectY * 75
		//
		//ballVelocity.vx = BALLSPEED * -math.Cos(bounceAngle)
		//ballVelocity.vy = BALLSPEED * math.Sin(bounceAngle)
		ballVelocity.vx = -ballVelocity.vx
	}

	ballPos.X += ballVelocity.vx
	ballPos.Y += ballVelocity.vy

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	playerOptions := &ebiten.DrawImageOptions{}
	playerOptions.GeoM.Translate(playerPos.X, playerPos.Y)
	screen.DrawImage(playerImage, playerOptions)

	opponentOptions := &ebiten.DrawImageOptions{}
	opponentOptions.GeoM.Translate(opponentPos.X, opponentPos.Y)
	screen.DrawImage(opponentImage, opponentOptions)

	ballOptions := &ebiten.DrawImageOptions{}
	ballOptions.GeoM.Translate(ballPos.X, ballPos.Y)
	screen.DrawImage(ballImage, ballOptions)
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
