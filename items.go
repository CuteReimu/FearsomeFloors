package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
)

type itemInterface interface {
	drawable
	init(b *board)
}

var (
	//go:embed assets/stone.png
	fileStone []byte
)

var (
	imgStone     *ebiten.Image
	imgSlipFloor *ebiten.Image
)

func init() {
	imgSlipFloor = ebiten.NewImage(gridLen, gridLen)
	imgSlipFloor.Fill(colornames.Red)
}

type stoneRegular struct{}

func init() {
	imageStone, _, err := image.Decode(bytes.NewReader(fileStone))
	if err != nil {
		panic(err)
	}
	imgStone = ebiten.NewImageFromImage(imageStone)
}

func (i *stoneRegular) Draw() (*ebiten.Image, *ebiten.DrawImageOptions) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(0.2, 0.25)
	opt.GeoM.Translate(8, 6)
	return imgStone, opt
}

func (i *stoneRegular) init(b *board) {
	for {
		x, y := b.random.Intn(width), b.random.Intn(height)
		if x < 3 && y < 3 || x >= width-1 && y >= height-1 || y-x >= height-3 || x-y >= width-3 {
			continue
		}
		if b.items[y][x] != nil {
			continue
		}
		b.items[y][x] = i
		break
	}
}

type slipFloor1 struct{}

func (i *slipFloor1) Draw() (*ebiten.Image, *ebiten.DrawImageOptions) {
	return imgSlipFloor, &ebiten.DrawImageOptions{}
}

func (i *slipFloor1) init(b *board) {
	for {
		x, y := b.random.Intn(width), b.random.Intn(height)
		if x < 3 && y < 3 || x >= width-2 && y >= height-2 || y-x >= height-4 || x-y >= width-4 {
			continue
		}
		if b.items[y][x] != nil || b.items[y+1][x] != nil || b.items[y][x+1] != nil || b.items[y+1][x+1] != nil {
			continue
		}
		b.items[y][x] = i
		b.items[y+1][x] = i
		b.items[y][x+1] = i
		b.items[y+1][x+1] = i
		break
	}
}

type slipFloor2 struct{}

func (i *slipFloor2) Draw() (*ebiten.Image, *ebiten.DrawImageOptions) {
	return imgSlipFloor, &ebiten.DrawImageOptions{}
}

func (i *slipFloor2) init(b *board) {
	for {
		if b.random.Intn(2) == 0 {
			x, y := b.random.Intn(width), b.random.Intn(height)
			if x < 3 && y < 3 || x >= width-4 && y >= height-1 || y-x >= height-3 || x-y >= width-6 {
				continue
			}
			if b.items[y][x] != nil || b.items[y][x+1] != nil || b.items[y][x+2] != nil || b.items[y][x+3] != nil {
				continue
			}
			b.items[y][x] = i
			b.items[y][x+1] = i
			b.items[y][x+2] = i
			b.items[y][x+3] = i
		} else {
			x, y := b.random.Intn(width), b.random.Intn(height)
			if x < 3 && y < 3 || x >= width-1 && y >= height-4 || y-x >= height-6 || x-y >= width-3 {
				continue
			}
			if b.items[y][x] != nil || b.items[y+1][x] != nil || b.items[y+2][x] != nil || b.items[y+3][x] != nil {
				continue
			}
			b.items[y][x] = i
			b.items[y+1][x] = i
			b.items[y+2][x] = i
			b.items[y+3][x] = i
		}
		break
	}
}
