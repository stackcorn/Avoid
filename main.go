package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// 定数の定義
const (
	screenWidth         = 640 // 画面の幅
	screenHeight        = 480 // 画面の高さ
	charSize            = 20  // キャラクターのサイズ
	logicalScreenWidth  = 320 // 論理的な画面の幅
	logicalScreenHeight = 240 // 論理的な画面の高さ
	obstacleSize        = 20  // 障害物のサイズ
)

// Star構造体
type Star struct {
	x, y float64
}

// Game構造体：ゲームの状態を保持
type Game struct {
	x, y             float64    // キャラクターの位置
	obstacles        []Obstacle // 障害物の配列
	stars            []Star     // 星の配列
	isGameOver       bool       // ゲームオーバーかどうか
	score            int        // スコア
	startTime        int64      // ゲーム開始時刻
	lastSecond       int64      // 最後にスコアが更新された時刻
	lastObstacleTime int64      // 最後に障害物が追加された時刻
	charImage        *ebiten.Image
	obstacleImage    *ebiten.Image
}

// Obstacle構造体
type Obstacle struct {
	x, y  float64
	speed float64
}

// NewGame関数：新しいゲームの状態を初期化
func NewGame() *Game {
	charImg, _, err := ebitenutil.NewImageFromFile("assets/char.png")
	if err != nil {
		log.Fatalf("failed to load character image: %v", err)
	}

	obstacleImg, _, err := ebitenutil.NewImageFromFile("assets/obstacle.png")
	if err != nil {
		log.Fatalf("failed to load obstacle image: %v", err)
	}

	// 星をランダムに生成
	stars := make([]Star, 100)
	for i := range stars {
		stars[i] = Star{
			x: float64(rand.Intn(logicalScreenWidth)),
			y: float64(rand.Intn(logicalScreenHeight)),
		}
	}

	return &Game{
		x:                50,
		y:                logicalScreenHeight / 2,
		startTime:        time.Now().UnixNano(),
		lastObstacleTime: -10,
		charImage:        charImg,
		obstacleImage:    obstacleImg,
		stars:            stars,
	}
}

// Updateメソッド：ゲームの状態を更新
func (g *Game) Update() error {
	// ゲームオーバー時の処理
	if g.isGameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			// ゲームをリセット
			*g = *NewGame()
		}
		return nil
	}

	// キャラクターの移動
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

	// 障害物とスコアの更新
	updateObstaclesAndScore(g)

	return nil
}

// updateObstaclesAndScore関数：障害物の更新とスコアの計算
func updateObstaclesAndScore(g *Game) {
	currentTime := time.Now().UnixNano()
	currentSecond := (currentTime - g.startTime) / int64(time.Second)

	// 10秒ごとに障害物を追加
	if currentSecond-g.lastObstacleTime >= 10 && currentSecond < 100 {
		g.obstacles = append(g.obstacles, Obstacle{
			x:     logicalScreenWidth,
			y:     float64(rand.Intn(logicalScreenHeight - obstacleSize)),
			speed: 2 + rand.Float64(),
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

	// スコアの更新
	if !g.isGameOver && currentSecond > g.lastSecond {
		if g.score < 999 {
			g.score++
		}
		g.lastSecond = currentSecond
	}
}

// Drawメソッド：ゲームの画面を描画
func (g *Game) Draw(screen *ebiten.Image) {
	// 星を描画（背景）
	for _, star := range g.stars {
		ebitenutil.DrawRect(screen, star.x, star.y, 1, 1, color.White) // 白い点
	}

	// キャラクターの描画
	drawCharacter(g, screen)

	// 障害物の描画
	drawObstacles(g, screen)

	// スコアの表示
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, logicalScreenWidth-80, 5)

	// ゲームオーバーメッセージの表示
	if g.isGameOver {
		displayGameOverMessage(screen)
	}
}

// drawCharacter関数：キャラクターを描画
func drawCharacter(g *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2) // キャラクターを2倍にスケーリング
	charW, charH := g.charImage.Size()
	charW *= 2
	charH *= 2
	op.GeoM.Translate(-float64(charW)/2, -float64(charH)/2)
	op.GeoM.Translate(g.x, g.y)
	screen.DrawImage(g.charImage, op)
}

// drawObstacles関数：障害物を描画
func drawObstacles(g *Game, screen *ebiten.Image) {
	for _, obstacle := range g.obstacles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2) // 障害物を2倍にスケーリング
		obW, obH := g.obstacleImage.Size()
		obW *= 2
		obH *= 2
		op.GeoM.Translate(-float64(obW)/2, -float64(obH)/2)
		op.GeoM.Translate(obstacle.x, obstacle.y)
		screen.DrawImage(g.obstacleImage, op)
	}
}

// displayGameOverMessage関数：ゲームオーバーメッセージを表示
func displayGameOverMessage(screen *ebiten.Image) {
	gameOverMsg := "GAME OVER"
	x := (logicalScreenWidth - len(gameOverMsg)*7) / 2
	y := logicalScreenHeight / 2
	ebitenutil.DebugPrintAt(screen, gameOverMsg, x, y)

	retryMsg := "RETRY: PRESS [R]"
	retryX := (logicalScreenWidth - len(retryMsg)*7) / 2
	retryY := y + 20
	ebitenutil.DebugPrintAt(screen, retryMsg, retryX, retryY)
}

// Layoutメソッド：外部ウィンドウのサイズが変更されたときのレイアウトを決定
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logicalScreenWidth, logicalScreenHeight
}

// main関数：プログラムのエントリーポイント
func main() {
	rand.Seed(time.Now().UnixNano()) // 乱数のシードを設定
	game := NewGame()                // 新しいゲームを作成
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Avoid Game")

	// ゲームループを開始
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// min関数：2つのfloat64のうち、小さい方を返す
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// max関数：2つのfloat64のうち、大きい方を返す
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
