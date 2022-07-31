package main

// TODO 出口出不去

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"time"
)

//go:embed assets/FZSTK.TTF
var ttfFile []byte
var fontAlpha font.Face
var fontNum font.Face
var emptyImage = ebiten.NewImage(1, 1)
var imgSlipFloor *ebiten.Image

func init() {
	imgSlipFloor = ebiten.NewImage(gridLen, gridLen)
	imgSlipFloor.Fill(colornames.Darkred)
	emptyImage.Fill(color.Black)
	tt, err := opentype.Parse(ttfFile)
	if err != nil {
		log.Fatal(err)
	}
	fontAlpha, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	fontNum, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
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

type dir point

var (
	up    = dir{0, -1}
	left  = dir{-1, 0}
	down  = dir{0, 1}
	right = dir{1, 0}
)

type floorShapeType uint8

const (
	floorShapeTypeEmpty floorShapeType = iota
	floorShapeTypeSlipFloor
	floorShapeTypeTransferUp
	floorShapeTypeTransferLeft
	floorShapeTypeTransferDown
	floorShapeTypeTransferRight
)

type board struct {
	items                 [][]itemInterface
	itemsCache            [][]itemInterface
	floorShape            [][]floorShapeType
	monster               *monster
	player                []*player
	random                *rand.Rand
	pickedPlayerItem      *playerItem
	pickedPlayerItemCache point
	curPlayer             int
	firstPlayer           int
	smallTurn             int
	bigTurn               int
	alreadyMoveCount      int
}

func newBoard(playerNum int) *board {
	b := &board{
		items:      make([][]itemInterface, height),
		itemsCache: make([][]itemInterface, height),
		floorShape: make([][]floorShapeType, height),
		player:     make([]*player, playerNum),
		random:     rand.New(rand.NewSource(time.Now().UnixMilli())),
		monster:    newMonster(),
	}
	for i := 0; i < height; i++ {
		b.items[i] = make([]itemInterface, width)
		b.itemsCache[i] = make([]itemInterface, width)
		b.floorShape[i] = make([]floorShapeType, width)
	}
	for i := 0; i < 11; i++ {
		(&stoneRegular{}).init(b)
	}
	b.initSlipFloor()
	switch len(b.player) {
	case 4:
		b.player[3] = newPlayer(colornames.Blue, "蓝方")
		fallthrough
	case 3:
		b.player[2] = newPlayer(colornames.Yellow, "黄方")
		fallthrough
	case 2:
		b.player[1] = newPlayer(colornames.Green, "绿方")
		fallthrough
	case 1:
		b.player[0] = newPlayer(colornames.Red, "红方")
		b.player[0].display(b)
	default:
		logger.Fatal("invalid player number")
	}
	return b
}

func (b *board) saveCache() {
	b.pickedPlayerItemCache = b.pickedPlayerItem.pos
	for i := range b.items {
		for j := range b.items[i] {
			b.itemsCache[i][j] = b.items[i][j]
		}
	}
}

func (b *board) loadCache() {
	b.pickedPlayerItem.pos = b.pickedPlayerItemCache
	for i := range b.items {
		for j := range b.items[i] {
			b.items[i][j] = b.itemsCache[i][j]
			if b.items[i][j] != nil {
				b.items[i][j].(*stoneRegular).setPos(point{j, i})
			}
		}
	}
}

func (b *board) Update() error {
	if b.monster.isMoving {
		return nil
	}
	if b.pickedPlayerItem == nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
			item := b.player[b.curPlayer].willMove(1)
			if item != nil {
				b.pickedPlayerItem = item
				b.saveCache()
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
			item := b.player[b.curPlayer].willMove(2)
			if item != nil {
				b.pickedPlayerItem = item
				b.saveCache()
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
			item := b.player[b.curPlayer].willMove(3)
			if item != nil {
				b.pickedPlayerItem = item
				b.saveCache()
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit4) {
			item := b.player[b.curPlayer].willMove(4)
			if item != nil {
				b.pickedPlayerItem = item
				b.saveCache()
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit5) {
			item := b.player[b.curPlayer].willMove(5)
			if item != nil {
				b.pickedPlayerItem = item
				b.saveCache()
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit6) {
			item := b.player[b.curPlayer].willMove(6)
			if item != nil {
				b.pickedPlayerItem = item
				b.saveCache()
			}
		}
	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			b.loadCache()
			b.alreadyMoveCount = 0
			b.pickedPlayerItem = nil
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if b.pickedPlayerItem.checkLegal(b) {
				b.pickedPlayerItem.alreadyMove = true
				b.alreadyMoveCount = 0
				b.pickedPlayerItem = nil
				for {
					b.curPlayer = (b.curPlayer + 1) % len(b.player)
					if b.curPlayer == b.firstPlayer {
						b.smallTurn++
						if b.bigTurn == 0 && b.smallTurn >= 2 || b.smallTurn >= len(b.player[0].items) {
							b.monster.move(b)
							for _, player := range b.player {
								player.nextTurn()
							}
							b.bigTurn++
							b.firstPlayer = (b.firstPlayer + 1) % len(b.player)
							b.curPlayer = b.firstPlayer
							b.smallTurn = 0
						}
					}
					b.player[b.curPlayer].display(b)
					if b.player[b.curPlayer].hasItemToMove() {
						break
					}
				}
			}
		} else {
			if b.alreadyMoveCount < b.pickedPlayerItem.step {
				if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) && b.pickedPlayerItem.tryMove(b, down) ||
					inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) && b.pickedPlayerItem.tryMove(b, left) ||
					inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && b.pickedPlayerItem.tryMove(b, up) ||
					inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) && b.pickedPlayerItem.tryMove(b, right) {
					b.alreadyMoveCount++
				}
			}
		}
	}
	return nil
}

func (b *board) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if b.floorShape[j][i] == floorShapeTypeSlipFloor {
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(edgeX+1+float64(i)*gridLen, edgeY+1+float64(j)*gridLen)
				screen.DrawImage(imgSlipFloor, opt)
			}
		}
	}
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if b.items[j][i] != nil {
				img, opt := b.items[j][i].Draw()
				opt.GeoM.Translate(edgeX+1+float64(i)*gridLen, edgeY+1+float64(j)*gridLen)
				screen.DrawImage(img, opt)
			}
		}
	}
	for _, player := range b.player {
		for _, item := range player.items {
			img, opt := item.Draw()
			opt.GeoM.Translate(edgeX+1+float64(item.pos.x)*gridLen, edgeY+1+float64(item.pos.y)*gridLen)
			screen.DrawImage(img, opt)
		}
	}
	img, opt := b.monster.Draw()
	screen.DrawImage(img, opt)
	for i := 0; i <= width; i++ {
		textX := edgeX + gridLen*i + gridLen/2 - 6
		textUp, textDown := string(rune('A'+i)), string(rune('A'+14-i))
		opt := &ebiten.DrawImageOptions{}
		if i < 3 {
			text.Draw(screen, textUp, fontAlpha, textX, edgeY-3, color.Black)
			text.Draw(screen, textDown, fontAlpha, textX, edgeY+gridLen*(height-3+i)+17, color.Black)
			opt.GeoM.Scale(2, gridLen*(height-3+float64(i))+2)
			opt.GeoM.Translate(edgeX+gridLen*float64(i), edgeY)
		} else if i >= width-3 {
			if i < width {
				text.Draw(screen, textUp, fontAlpha, textX, edgeY+gridLen*(i-width+4)-3, color.Black)
				text.Draw(screen, textDown, fontAlpha, textX, edgeY+gridLen*height+17, color.Black)
			}
			opt.GeoM.Scale(2, gridLen*(width+height-3-float64(i))+2)
			opt.GeoM.Translate(edgeX+gridLen*float64(i), edgeY+gridLen*(float64(i)-width+3))
		} else {
			text.Draw(screen, textUp, fontAlpha, textX, edgeY-3, color.Black)
			text.Draw(screen, textDown, fontAlpha, textX, edgeY+gridLen*height+17, color.Black)
			opt.GeoM.Scale(2, gridLen*height+2)
			opt.GeoM.Translate(edgeX+gridLen*float64(i), edgeY)
		}
		screen.DrawImage(emptyImage, opt)
	}
	for j := 0; j <= height; j++ {
		textY := edgeY + gridLen*j + gridLen/2 + 6
		textLeft, textRight := string(rune('Z'-j)), string(rune('Z'-9+j))
		opt := &ebiten.DrawImageOptions{}
		if j < 3 {
			text.Draw(screen, textLeft, fontAlpha, edgeX-20, textY, color.Black)
			text.Draw(screen, textRight, fontAlpha, edgeX+gridLen*(width-3+j)+6, textY, color.Black)
			opt.GeoM.Scale(gridLen*(width-3+float64(j))+2, 2)
			opt.GeoM.Translate(edgeX, edgeY+gridLen*float64(j))
		} else if j >= height-3 {
			if j < height {
				text.Draw(screen, textLeft, fontAlpha, edgeX+gridLen*(j+4-height)-20, textY, color.Black)
				text.Draw(screen, textRight, fontAlpha, edgeX+gridLen*width+6, textY, color.Black)
			}
			opt.GeoM.Scale(gridLen*(width+height-3-float64(j))+2, 2)
			opt.GeoM.Translate(edgeX+gridLen*(float64(j)-height+3), edgeY+gridLen*float64(j))
		} else {
			text.Draw(screen, textLeft, fontAlpha, edgeX-20, textY, color.Black)
			text.Draw(screen, textRight, fontAlpha, edgeX+gridLen*width+6, textY, color.Black)
			opt.GeoM.Scale(gridLen*width+2, 2)
			opt.GeoM.Translate(edgeX, edgeY+gridLen*float64(j))
		}
		screen.DrawImage(emptyImage, opt)
	}
}

func (b *board) Layout(int, int) (screenWidth, screenHeight int) {
	return 1024, 768
}

func (b *board) initSlipFloor() {
	for {
		x, y := b.random.Intn(width), b.random.Intn(height)
		if x < 3 && y < 3 || x >= width-2 && y >= height-2 || x > width-2 || y > height-2 || y-x >= height-4 || x-y >= width-4 {
			continue
		}
		if b.items[y][x] != nil || b.items[y+1][x] != nil || b.items[y][x+1] != nil || b.items[y+1][x+1] != nil {
			continue
		}
		if b.floorShape[y][x] != floorShapeTypeEmpty || b.floorShape[y+1][x] != floorShapeTypeEmpty {
			continue
		}
		if b.floorShape[y][x+1] != floorShapeTypeEmpty || b.floorShape[y+1][x+1] != floorShapeTypeEmpty {
			continue
		}
		b.floorShape[y][x] = floorShapeTypeSlipFloor
		b.floorShape[y+1][x] = floorShapeTypeSlipFloor
		b.floorShape[y][x+1] = floorShapeTypeSlipFloor
		b.floorShape[y+1][x+1] = floorShapeTypeSlipFloor
		break
	}
	for {
		if b.random.Intn(2) == 0 {
			x, y := b.random.Intn(width), b.random.Intn(height)
			if x < 3 && y < 3 || x >= width-4 && y >= height-1 || x > width-4 || y-x >= height-3 || x-y >= width-6 {
				continue
			}
			if b.items[y][x] != nil || b.items[y][x+1] != nil || b.items[y][x+2] != nil || b.items[y][x+3] != nil {
				continue
			}
			if b.floorShape[y][x] != floorShapeTypeEmpty || b.floorShape[y][x+1] != floorShapeTypeEmpty {
				continue
			}
			if b.floorShape[y][x+2] != floorShapeTypeEmpty || b.floorShape[y][x+3] != floorShapeTypeEmpty {
				continue
			}
			b.floorShape[y][x] = floorShapeTypeSlipFloor
			b.floorShape[y][x+1] = floorShapeTypeSlipFloor
			b.floorShape[y][x+2] = floorShapeTypeSlipFloor
			b.floorShape[y][x+3] = floorShapeTypeSlipFloor
		} else {
			x, y := b.random.Intn(width), b.random.Intn(height)
			if x < 3 && y < 3 || x >= width-1 && y >= height-4 || y > height-4 || y-x >= height-6 || x-y >= width-3 {
				continue
			}
			if b.items[y][x] != nil || b.items[y+1][x] != nil || b.items[y+2][x] != nil || b.items[y+3][x] != nil {
				continue
			}
			if b.floorShape[y][x] != floorShapeTypeEmpty || b.floorShape[y+1][x] != floorShapeTypeEmpty {
				continue
			}
			if b.floorShape[y+2][x] != floorShapeTypeEmpty || b.floorShape[y+3][x] != floorShapeTypeEmpty {
				continue
			}
			b.floorShape[y][x] = floorShapeTypeSlipFloor
			b.floorShape[y+1][x] = floorShapeTypeSlipFloor
			b.floorShape[y+2][x] = floorShapeTypeSlipFloor
			b.floorShape[y+3][x] = floorShapeTypeSlipFloor
		}
		break
	}
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
	return point{p.x + d.x, p.y + d.y}
}

func (p point) outOfRange() bool {
	return p.x < 0 || p.x >= width || p.y < 0 || p.y >= height || p.x-p.y >= width-3 || p.y-p.x >= height-3
}
