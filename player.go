package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"sort"
	"strconv"
	"strings"
)

type playerItem struct {
	alreadyMove bool
	step        int
	pos         point
	color       color.Color
}

func (p *playerItem) Draw() (*ebiten.Image, *ebiten.DrawImageOptions) {
	if p.alreadyMove {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(2, 2)
		img := ebiten.NewImage(gridLen/2, gridLen/2)
		img.Fill(color.Alpha{})
		text.Draw(img, "〇", fontNum, -8, 33, p.color)
		return img, opt
	} else {
		opt := &ebiten.DrawImageOptions{}
		img := ebiten.NewImage(gridLen, gridLen)
		img.Fill(color.Alpha{})
		switch p.step {
		case 1:
			text.Draw(img, "①", fontNum, 6, 46, p.color)
		case 2:
			text.Draw(img, "②", fontNum, 6, 46, p.color)
		case 3:
			text.Draw(img, "③", fontNum, 6, 46, p.color)
		case 4:
			text.Draw(img, "④", fontNum, 6, 46, p.color)
		case 5:
			text.Draw(img, "⑤", fontNum, 6, 46, p.color)
		case 6:
			text.Draw(img, "⑥", fontNum, 6, 46, p.color)
		}
		return img, opt
	}
}

func (p *playerItem) die(b *board) {
	if b.bigTurn > 7 {
		p.pos.y = -2
	} else {
		p.pos.y = -1
	}
	p.pos.x = 0
}

func (p *playerItem) isDead() bool {
	return p.pos.x == 0 && p.pos.y == -2
}

func (p *playerItem) isFinished() bool {
	return p.pos.x == width && p.pos.y == height
}

func (p *playerItem) tryMove(b *board, d dir) bool {
	pos := p.pos
	pos.x += d.x
	pos.y += d.y
	if pos.x == width && pos.y == height-1 || pos.x == width-1 && pos.y == height {
		p.pos.x = width
		p.pos.y = height
		return true
	}
	if pos.outOfRange() || b.monster.pos == pos {
		return false
	}
	if b.floorShape[pos.y][pos.x] >= floorShapeTypeTransferUp {
		return false
	}
	if b.items[pos.y][pos.x] != nil && !b.items[pos.y][pos.x].tryMove(b, d) {
		return false
	}
	p.pos = pos
	if b.floorShape[pos.y][pos.x] == floorShapeTypeSlipFloor {
		p.tryMove(b, d)
	}
	return true
}

func (p *playerItem) forceMove(b *board, d dir) {
	pos := p.pos
	pos.x += d.x
	pos.y += d.y
	if pos.outOfRange() || b.floorShape[pos.y][pos.x] >= floorShapeTypeTransferUp {
		p.die(b)
		return
	}
	if b.items[pos.y][pos.x] != nil {
		b.items[pos.y][pos.x].forceMove(b, d)
	}
	for _, player := range b.player {
		for _, item := range player.items {
			if pos == item.pos {
				item.forceMove(b, d)
			}
		}
	}
	p.pos = pos
	if b.floorShape[pos.y][pos.x] == floorShapeTypeSlipFloor {
		p.tryMove(b, d)
	}
}

func (p *playerItem) checkLegal(b *board) bool {
	if p.pos.x == 0 && p.pos.y == 1 || p.pos.x == width && p.pos.y == height {
		return true
	}
	for i := range b.items {
		for j := range b.items[i] {
			if b.floorShape[i][j] != floorShapeTypeEmpty && p.pos.x == j && p.pos.y == i {
				return false
			}
		}
	}
	for _, player := range b.player {
		for _, item := range player.items {
			if item != p && item.pos == p.pos {
				return false
			}
		}
	}
	return true
}

type player struct {
	items []*playerItem
	text  string
}

func (p *player) willMove(num int) *playerItem {
	for _, item := range p.items {
		if item.step == num {
			if item.alreadyMove || item.isFinished() || item.isDead() {
				return nil
			}
			return item
		}
	}
	return nil
}

func (p *player) hasItemToMove() bool {
	for _, item := range p.items {
		if !item.alreadyMove && !item.isFinished() && !item.isDead() {
			return true
		}
	}
	return false
}

func (p *player) nextTurn() {
	for _, item := range p.items {
		item.alreadyMove = false
		item.step = 7 - item.step
	}
}

func newPlayer(c color.Color, text string) *player {
	return &player{
		text: text,
		items: []*playerItem{
			{step: 1, pos: point{0, -1}, color: c},
			{step: 3, pos: point{0, -1}, color: c},
			{step: 4, pos: point{0, -1}, color: c},
			{step: 5, pos: point{0, -1}, color: c},
		},
	}
}

func (p *player) display(b *board) {
	s := fmt.Sprintf("剩余%d张牌，", len(b.monster.deck))
	if b.monster.lastCard != nil {
		s = "怪物的上一张牌是" + b.monster.lastCard.text + s
	}
	var canMoveItems []int
	for _, item := range p.items {
		if !item.alreadyMove && !item.isFinished() && !item.isDead() {
			canMoveItems = append(canMoveItems, item.step)
		}
	}
	sort.Ints(canMoveItems)
	s += "轮到" + p.text
	if canMoveItems != nil {
		var canMoveItemsString []string
		for _, item := range canMoveItems {
			canMoveItemsString = append(canMoveItemsString, strconv.Itoa(item))
		}
		s += "，能移动的棋子有" + strings.Join(canMoveItemsString, "，")
	}
	ebiten.SetWindowTitle(s)
}
