package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
)

type itemInterface interface {
	Draw() (*ebiten.Image, *ebiten.DrawImageOptions)
	init(b *board)
	tryMove(b *board, d dir) bool
	forceMove(b *board, d dir)
	setPos(pos point)
}

var (
	//go:embed assets/stone.png
	fileStone []byte
)

var (
	imgStone *ebiten.Image
)

type stoneRegular struct {
	pos point
}

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
		if b.floorShape[y][x] != floorShapeTypeEmpty {
			continue
		}
		b.items[y][x] = i
		i.pos.x, i.pos.y = x, y
		break
	}
}

func (i *stoneRegular) tryMove(b *board, d dir) bool {
	pos := i.pos
	pos.x += d.x
	pos.y += d.y
	if pos.outOfRange() || b.monster.pos == pos || b.floorShape[pos.y][pos.x] >= floorShapeTypeTransferUp || b.items[pos.y][pos.x] != nil {
		return false
	}
	for _, player := range b.player {
		for _, item := range player.items {
			if item.pos == pos {
				return false
			}
		}
	}
	b.items[i.pos.y][i.pos.x] = nil
	b.items[pos.y][pos.x] = i
	i.pos = pos
	if b.floorShape[pos.y][pos.x] == floorShapeTypeSlipFloor {
		i.tryMove(b, d)
	}
	if i.pos.x == width-1 && i.pos.y == height-1 {
		b.items[i.pos.y][i.pos.x] = nil
	}
	return true
}

func (i *stoneRegular) forceMove(b *board, d dir) {
	pos := i.pos
	pos.x += d.x
	pos.y += d.y
	if pos.outOfRange() || b.monster.pos == pos || b.floorShape[pos.y][pos.x] >= floorShapeTypeTransferUp {
		b.items[i.pos.y][i.pos.x] = nil
		return
	}
	if b.items[pos.y][pos.x] != nil {
		b.items[pos.y][pos.x].forceMove(b, d)
	}
	func() {
		for _, player := range b.player {
			for _, item := range player.items {
				if item.pos == pos {
					item.forceMove(b, d)
					return
				}
			}
		}
	}()
	b.items[i.pos.y][i.pos.x] = nil
	b.items[pos.y][pos.x] = i
	i.pos = pos
	if b.floorShape[pos.y][pos.x] == floorShapeTypeSlipFloor {
		i.tryMove(b, d)
	}
}

func (i *stoneRegular) setPos(pos point) {
	i.pos = pos
}
