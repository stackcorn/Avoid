package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth         = 640 // 画面の幅
	screenHeight        = 480 // 画面の高さ
	charSize            = 20  // キャラクターのサイズ
	logicalScreenWidth  = 320 // 論理的な画面の幅
	logicalScreenHeight = 240 // 論理的な画面の高さ
)

type Game struct {
	x, y float64
}

func (g *Game) Update() error {
	// キーボード入力に応じて位置を更新
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

	// 画面の範囲内にキャラクターが収まるように調整
	g.x = max(0, min(newX, logicalScreenWidth-charSize))
	g.y = max(0, min(newY, logicalScreenHeight-charSize))

	return nil
}

// min は2つの値のうち小さい方を返します。
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// max は2つの値のうち大きい方を返します。
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, g.x, g.y, 20, 20, color.RGBA{255, 0, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logicalScreenWidth, logicalScreenHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Avoid Game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
