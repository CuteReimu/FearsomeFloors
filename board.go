package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"time"
)

//go:embed assets/FZSTK.TTF
var ttfFile []byte
var fontFace font.Face
var emptyImage = ebiten.NewImage(1, 1)

func init() {
	emptyImage.Fill(color.Black)
	tt, err := opentype.Parse(ttfFile)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

const (
	edgeX   = 50
	edgeY   = 100
	width   = 15
	height  = 10
	gridLen = 60
)

type dir uint8

const (
	up dir = iota
	left
	down
	right
)

type drawable interface {
	Draw() (*ebiten.Image, *ebiten.DrawImageOptions)
}

type board struct {
	items   [][]itemInterface
	monster point
	player  []point
	random  *rand.Rand
}

func newBoard(playerNum int) *board {
	b := &board{
		items:   make([][]itemInterface, height),
		monster: point{width - 1, height - 1},
		player:  make([]point, playerNum),
		random:  rand.New(rand.NewSource(time.Now().UnixMilli())),
	}
	for i := 0; i < height; i++ {
		b.items[i] = make([]itemInterface, width)
	}
	for i := 0; i < 11; i++ {
		(&stoneRegular{}).init(b)
	}
	(&slipFloor1{}).init(b)
	(&slipFloor2{}).init(b)
	return b
}

func (b *board) Update() error {
	return nil
}

func (b *board) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if b.items[j][i] != nil {
				img, opt := b.items[j][i].Draw()
				opt.GeoM.Translate(edgeX+1+float64(i)*gridLen, edgeY+1+float64(j)*gridLen)
				screen.DrawImage(img, opt)
			}
		}
	}
	for i := 0; i <= width; i++ {
		textX := edgeX + gridLen*i + gridLen/2 - 6
		textUp, textDown := string(rune('A'+i)), string(rune('A'+14-i))
		opt := &ebiten.DrawImageOptions{}
		if i < 3 {
			text.Draw(screen, textUp, fontFace, textX, edgeY-3, color.Black)
			text.Draw(screen, textDown, fontFace, textX, edgeY+gridLen*(height-3+i)+17, color.Black)
			opt.GeoM.Scale(2, gridLen*(height-3+float64(i)))
			opt.GeoM.Translate(edgeX+gridLen*float64(i), edgeY)
		} else if i >= width-3 {
			if i < width {
				text.Draw(screen, textUp, fontFace, textX, edgeY+gridLen*(i-width+4)-3, color.Black)
				text.Draw(screen, textDown, fontFace, textX, edgeY+gridLen*height+17, color.Black)
			}
			opt.GeoM.Scale(2, gridLen*(width+height-3-float64(i)))
			opt.GeoM.Translate(edgeX+gridLen*float64(i), edgeY+gridLen*(float64(i)-width+3))
		} else {
			text.Draw(screen, textUp, fontFace, textX, edgeY-3, color.Black)
			text.Draw(screen, textDown, fontFace, textX, edgeY+gridLen*height+17, color.Black)
			opt.GeoM.Scale(2, gridLen*height)
			opt.GeoM.Translate(edgeX+gridLen*float64(i), edgeY)
		}
		screen.DrawImage(emptyImage, opt)
	}
	for j := 0; j <= height; j++ {
		textY := edgeY + gridLen*j + gridLen/2 + 6
		textLeft, textRight := string(rune('Z'-j)), string(rune('Z'-9+j))
		opt := &ebiten.DrawImageOptions{}
		if j < 3 {
			text.Draw(screen, textLeft, fontFace, edgeX-20, textY, color.Black)
			text.Draw(screen, textRight, fontFace, edgeX+gridLen*(width-3+j)+6, textY, color.Black)
			opt.GeoM.Scale(gridLen*(width-3+float64(j)), 2)
			opt.GeoM.Translate(edgeX, edgeY+gridLen*float64(j))
		} else if j >= height-3 {
			if j < height {
				text.Draw(screen, textLeft, fontFace, edgeX+gridLen*(j+4-height)-20, textY, color.Black)
				text.Draw(screen, textRight, fontFace, edgeX+gridLen*width+6, textY, color.Black)
			}
			opt.GeoM.Scale(gridLen*(width+height-3-float64(j)), 2)
			opt.GeoM.Translate(edgeX+gridLen*(float64(j)-height+3), edgeY+gridLen*float64(j))
		} else {
			text.Draw(screen, textLeft, fontFace, edgeX-20, textY, color.Black)
			text.Draw(screen, textRight, fontFace, edgeX+gridLen*width+6, textY, color.Black)
			opt.GeoM.Scale(gridLen*width, 2)
			opt.GeoM.Translate(edgeX, edgeY+gridLen*float64(j))
		}
		screen.DrawImage(emptyImage, opt)
	}
}

func (b *board) Layout(int, int) (screenWidth, screenHeight int) {
	return 1024, 768
}

type point struct {
	x, y int
}

func (p point) move(d dir, l int) point {
	p1 := p
	for i := 0; i < l; i++ {
		p2 := p1.moveOne(d)
		if p2.outOfRange() {
			p1.x = width - 1 - p1.x
			p1.y = height - 1 - p1.y
			p1 = p1.moveOne(d)
		} else {
			p1 = p2
		}
	}
	return p1
}

func (p point) moveOne(d dir) point {
	switch d {
	case up:
		return point{p.x, p.y - 1}
	case left:
		return point{p.x - 1, p.y}
	case down:
		return point{p.x, p.y + 1}
	case right:
		return point{p.x + 1, p.y}
	default:
		panic(d)
	}
}

func (p point) outOfRange() bool {
	return p.x < 0 || p.x >= width || p.y < 0 || p.y >= height || p.x-p.y >= width-3 || p.y-p.x >= height-3
}
