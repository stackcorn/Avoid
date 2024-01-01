package main

import (
	"fmt"
	"image/color"
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
}

type Obstacle struct {
	x, y  float64 // 障害物の位置
	speed float64 // 障害物の速度
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

func NewGame() *Game {
	return &Game{
		x:                50,
		y:                logicalScreenHeight / 2,
		startTime:        time.Now().UnixNano(),
		lastObstacleTime: -10,
		// その他の初期化が必要なフィールド
	}
}

func (g *Game) Update() error {
	// ゲームオーバー時の処理
	if g.isGameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			// ゲームをリセットする処理
			*g = *NewGame() // NewGameはGameの初期状態を生成する関数
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
	if currentSecond < 100 && currentSecond-g.lastObstacleTime >= 10 {
		g.obstacles = append(g.obstacles, Obstacle{
			x:     logicalScreenWidth,
			y:     float64(rand.Intn(logicalScreenHeight-obstacleSize)) + 10,
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
			g.obstacles[i].y = float64(rand.Intn(logicalScreenHeight-obstacleSize)) + 10
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
	ebitenutil.DrawRect(screen, g.x, g.y, charSize, charSize, color.RGBA{255, 0, 0, 255})

	for _, obstacle := range g.obstacles {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(obstacle.x, obstacle.y)
		screen.DrawImage(obstacleImage(screen), opts)
	}

	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, logicalScreenWidth-80, 5)

	if g.isGameOver {
		msg := "GAME OVER"
		x := (logicalScreenWidth - len(msg)*7) / 2
		y := logicalScreenHeight / 2
		ebitenutil.DebugPrintAt(screen, msg, x, y)

		retryMsg := "RETRY: PRESS [R]"
		retryX := (logicalScreenWidth - len(retryMsg)*7) / 2
		retryY := y + 20 // GAME OVERメッセージの下に表示
		ebitenutil.DebugPrintAt(screen, retryMsg, retryX, retryY)
	}
}

func obstacleImage(screen *ebiten.Image) *ebiten.Image {
	radius := 10.0
	obstacleColor := color.RGBA{0, 128, 0, 255}

	obstacleImage := ebiten.NewImage(int(radius*2), int(radius*2))
	for y := -radius; y < radius; y++ {
		for x := -radius; x < radius; x++ {
			if x*x+y*y <= radius*radius {
				obstacleImage.Set(int(x+radius), int(y+radius), obstacleColor)
			}
		}
	}
	return obstacleImage
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logicalScreenWidth, logicalScreenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := &Game{
		x:                50,
		y:                logicalScreenHeight / 2,
		startTime:        time.Now().UnixNano(),
		lastObstacleTime: -10, // 初期値を-10に設定して、ゲーム開始時に障害物が追加されるようにする
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Avoid Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
