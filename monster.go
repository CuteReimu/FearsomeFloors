package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"time"
)

//go:embed assets/monster.png
var fileMonster []byte
var imgMonster *ebiten.Image

func init() {
	imageMonster, _, err := image.Decode(bytes.NewReader(fileMonster))
	if err != nil {
		panic(err)
	}
	imgMonster = ebiten.NewImageFromImage(imageMonster)
}

type monster struct {
	isMoving bool
	faceTo   dir
	pos      point
	deck     []*card
	lastCard *card
}

func (m *monster) Draw() (*ebiten.Image, *ebiten.DrawImageOptions) {
	bounds := imgMonster.Bounds()
	dx, dy := bounds.Dx(), bounds.Dy()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(gridLen)/float64(dx), float64(gridLen)/float64(dy))
	if m.faceTo != left {
		opt.GeoM.Translate(-float64(gridLen)/2, -float64(gridLen)/2)
		switch m.faceTo {
		case right:
			opt.GeoM.Scale(-1, 1)
		case up:
			opt.GeoM.Rotate(math.Pi / 2)
		case down:
			opt.GeoM.Rotate(-math.Pi / 2)
		}
		opt.GeoM.Translate(float64(gridLen)/2, float64(gridLen)/2)
	}
	opt.GeoM.Translate(edgeX+1+float64(m.pos.x)*gridLen, edgeY+1+float64(m.pos.y)*gridLen)
	return imgMonster, opt
}

func (m *monster) moveOne(b *board, leftStep int, curKillCount, maxKillCount int) {
	if leftStep == 0 {
		m.isMoving = false
		return
	}
	pos := m.pos
	pos.x += m.faceTo.x
	pos.y += m.faceTo.y
	if pos.outOfRange() {
		pos.x = width - 1 - m.pos.x
		pos.y = height - 1 - m.pos.y
	}
	if b.items[pos.y][pos.x] != nil {
		b.items[pos.y][pos.x].forceMove(b, m.faceTo)
	}
	for _, player := range b.player {
		for _, item := range player.items {
			if item.pos == pos {
				item.die(b)
				curKillCount++
				if curKillCount == maxKillCount {
					m.pos = pos
					m.chooseDir(b)
					m.isMoving = false
					return
				}
			}
		}
	}
	m.pos = pos
	if b.floorShape[m.pos.y][m.pos.x] == floorShapeTypeSlipFloor {
		m.moveOne(b, leftStep, curKillCount, maxKillCount)
		return
	}
	m.chooseDir(b)
	time.AfterFunc(time.Second/2, func() { m.moveOne(b, leftStep-1, curKillCount, maxKillCount) })
}

func (m *monster) move(b *board) {
	m.isMoving = true
	idx := b.random.Intn(len(m.deck))
	if b.bigTurn == 0 {
		for m.deck[idx].step >= 20 {
			idx = b.random.Intn(len(m.deck))
		}
	}
	m.chooseDir(b)
	m.moveOne(b, m.deck[idx].step, 0, m.deck[idx].kills)
	m.deck = append(m.deck[:idx], m.deck[idx+1:]...)
	if len(m.deck) <= 1 {
		m.deck = newDeck()
	}
}

func (m *monster) chooseDir(b *board) {
	leftDistance := m.findPlayer(b, left)
	upDistance := m.findPlayer(b, up)
	rightDistance := m.findPlayer(b, right)
	downDistance := m.findPlayer(b, down)
	min := leftDistance
	if upDistance < min {
		min = upDistance
	}
	if rightDistance < min {
		min = rightDistance
	}
	if downDistance < min {
		min = downDistance
	}
	minCount := 0
	var minDir dir
	if leftDistance == min {
		minCount++
		minDir = left
	}
	if rightDistance == min {
		if minCount > 0 {
			return
		}
		minCount++
		minDir = right
	}
	if upDistance == min {
		if minCount > 0 {
			return
		}
		minCount++
		minDir = up
	}
	if downDistance == min {
		if minCount > 0 {
			return
		}
		minDir = down
	}
	m.faceTo = minDir
}

func (m *monster) findPlayer(b *board, d dir) int {
	if m.faceTo.x+d.x == 0 && m.faceTo.y+d.y == 0 {
		return 99
	}
	pos := m.pos
	for i := 1; i < 99; i++ {
		pos.x += d.x
		pos.y += d.y
		if pos.outOfRange() || b.items[pos.y][pos.x] != nil {
			return 99
		}
		for _, player := range b.player {
			for _, item := range player.items {
				if item.pos == pos {
					return i
				}
			}
		}
	}
	logger.Fatal("unreachable code")
	return 99
}

func newMonster() *monster {
	return &monster{
		faceTo: left,
		pos:    point{width - 1, height - 1},
		deck:   newDeck(),
	}
}

type card struct {
	text  string
	step  int
	kills int
}

func newDeck() []*card {
	return []*card{
		{"5", 5, 99},
		{"7", 7, 99},
		{"7", 7, 99},
		{"8", 8, 99},
		{"8", 8, 99},
		{"10", 10, 99},
		{"X", 20, 1},
		{"XX", 20, 2},
	}
}
