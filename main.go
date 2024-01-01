package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth         = 640 // 画面の幅
	screenHeight        = 480 // 画面の高さ
	charSize            = 20  // キャラクターのサイズ
	logicalScreenWidth  = 320 // 論理的な画面の幅
	logicalScreenHeight = 240 // 論理的な画面の高さ
	obstacleSize        = 20  // 障害物のサイズ
)

type Game struct {
	x, y             float64    // キャラクターの位置
	obstacles        []Obstacle // 複数の障害物
	isGameOver       bool       // ゲームオーバーの状態
	score            int        // スコア
	startTime        int64      // ゲーム開始時刻（Unixナノ秒）
	lastSecond       int64      // 最後にスコアが更新された時刻（秒）
	lastObstacleTime int64      // 最後に障害物が追加された時刻（秒）
	charImage        *ebiten.Image
	obstacleImage    *ebiten.Image
}

type Obstacle struct {
	x, y  float64
	speed float64
}

func NewGame() *Game {
	charImg, _, err := ebitenutil.NewImageFromFile("assets/char.png")
	if err != nil {
		log.Fatalf("failed to load character image: %v", err)
	}

	obstacleImg, _, err := ebitenutil.NewImageFromFile("assets/obstacle.png")
	if err != nil {
		log.Fatalf("failed to load obstacle image: %v", err)
	}

	return &Game{
		x:                50,
		y:                logicalScreenHeight / 2,
		startTime:        time.Now().UnixNano(),
		lastObstacleTime: -10,
		charImage:        charImg,
		obstacleImage:    obstacleImg,
	}
}

func (g *Game) Update() error {
	// ゲームオーバー時の処理
	if g.isGameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			// ゲームをリセットする処理
			*g = *NewGame()
		}
		return nil
	}

	// キーボード入力に応じてキャラクターの位置を更新
	newX, newY := g.x, g.y
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newY--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newY++
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newX--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newX++
	}
	g.x = max(0, min(newX, logicalScreenWidth-charSize))
	g.y = max(0, min(newY, logicalScreenHeight-charSize))

	currentTime := time.Now().UnixNano()
	currentSecond := (currentTime - g.startTime) / int64(time.Second)

	// 100秒経過していない、かつ10秒ごとに新しい障害物を追加
	if currentSecond-g.lastObstacleTime >= 10 && currentSecond < 100 {
		g.obstacles = append(g.obstacles, Obstacle{
			x:     logicalScreenWidth,
			y:     float64(rand.Intn(logicalScreenHeight - obstacleSize)),
			speed: 2 + rand.Float64(), // 障害物の速度をランダム化
		})
		g.lastObstacleTime = currentSecond
	}

	for i := range g.obstacles {
		// 障害物を左に移動
		g.obstacles[i].x -= g.obstacles[i].speed

		// 画面の左端に到達したら、右端から再スタート
		if g.obstacles[i].x < -obstacleSize {
			g.obstacles[i].x = logicalScreenWidth
			g.obstacles[i].y = float64(rand.Intn(logicalScreenHeight - obstacleSize))
		}

		// 衝突判定
		if !g.isGameOver && g.x < g.obstacles[i].x+obstacleSize && g.x+charSize > g.obstacles[i].x &&
			g.y < g.obstacles[i].y+obstacleSize && g.y+charSize > g.obstacles[i].y {
			g.isGameOver = true
		}
	}

	if !g.isGameOver && currentSecond > g.lastSecond {
		if g.score < 999 {
			g.score++
		}
		g.lastSecond = currentSecond
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	charW, charH := g.obstacleImage.Size()
	charW *= 2
	charH *= 2
	op.GeoM.Translate(-float64(charW)/2, -float64(charH)/2)
	op.GeoM.Translate(g.x, g.y)
	screen.DrawImage(g.charImage, op)

	for _, obstacle := range g.obstacles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		obW, obH := g.obstacleImage.Size()
		obW *= 2
		obH *= 2
		op.GeoM.Translate(-float64(obW)/2, -float64(obH)/2)
		op.GeoM.Translate(obstacle.x, obstacle.y)
		screen.DrawImage(g.obstacleImage, op)
	}

	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, logicalScreenWidth-80, 5)

	if g.isGameOver {
		gameOverMsg := "GAME OVER"
		x := (logicalScreenWidth - len(gameOverMsg)*7) / 2
		y := logicalScreenHeight / 2
		ebitenutil.DebugPrintAt(screen, gameOverMsg, x, y)

		retryMsg := "RETRY: PRESS [R]"
		retryX := (logicalScreenWidth - len(retryMsg)*7) / 2
		retryY := y + 20
		ebitenutil.DebugPrintAt(screen, retryMsg, retryX, retryY)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logicalScreenWidth, logicalScreenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Avoid Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
