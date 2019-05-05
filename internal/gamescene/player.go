// Copyright 2019 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gamescene

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/gopherwalk/internal/scene"
)

const PlayerUnit = 32

type Player struct {
	x32     int
	y32     int
	dir     Dir
	falling bool
}

func NewPlayer(x, y int) *Player {
	return &Player{
		x32: x * PlayerUnit,
		y32: y * PlayerUnit,
	}
}

func (p *Player) Update(context scene.Context, f *Field) {
	climbing := false
	if !p.falling && f.InElevator(p.elevatorArea()) {
		p.y32--
		climbing = true
	} else if !f.OverlapsWithFoot(p.footArea()) {
		if !p.falling {
			switch p.dir {
			case DirLeft:
				p.x32 -= 8
			case DirRight:
				p.x32 += 8
			default:
				panic("not reached")
			}
			p.falling = true
		}
		for i := 0; i < 3 && !f.OverlapsWithFoot(p.footArea()); i++ {
			p.y32++
		}
	} else {
		p.falling = false
	}
	if !p.falling {
		if !climbing {
			a := p.conflictionArea()
			switch p.dir {
			case DirLeft:
				a.Min.X--
				a.Max.X--
				if f.Overlaps(a, p.dir) {
					p.dir = DirRight
				} else {
					p.x32--
				}
			case DirRight:
				a.Min.X++
				a.Max.X++
				if f.Overlaps(a, p.dir) {
					p.dir = DirLeft
				} else {
					p.x32++
				}
			default:
				panic("not reached")
			}
		}
		x, y := context.Input().CursorPosition()
		if image.Pt(x, y).In(p.clickableArea()) && context.Input().IsJustTapped() {
			switch p.dir {
			case DirLeft:
				p.dir = DirRight
			case DirRight:
				p.dir = DirLeft
			default:
				panic("not reached")
			}
		}
	}
}

func (p *Player) conflictionArea() image.Rectangle {
	x := p.x32 * tileWidth / PlayerUnit
	y := p.y32 * tileHeight / PlayerUnit
	return image.Rect(x, y, x+tileWidth, y+tileHeight)
}

func (p *Player) elevatorArea() image.Rectangle {
	x := 0
	switch p.dir {
	case DirLeft:
		x = (p.x32*tileWidth)/PlayerUnit + tileWidth*3/4 - 1
	case DirRight:
		x = (p.x32*tileWidth)/PlayerUnit + tileWidth/4
	default:
		panic("not reached")
	}
	y := (p.y32 * tileHeight) / PlayerUnit
	return image.Rect(x, y, x+1, y+tileHeight)
}

func (p *Player) clickableArea() image.Rectangle {
	x := p.x32*tileWidth/PlayerUnit - tileWidth/2
	y := p.y32*tileHeight/PlayerUnit - tileHeight
	return image.Rect(x, y, x+tileWidth*2, y+tileHeight*2)
}

func (p *Player) footArea() image.Rectangle {
	x := 0
	switch p.dir {
	case DirLeft:
		x = (p.x32 * tileWidth) / PlayerUnit
	case DirRight:
		x = (p.x32*tileWidth)/PlayerUnit + tileWidth/2
	default:
		panic("not reached")
	}
	y := (p.y32*tileHeight)/PlayerUnit + tileHeight
	return image.Rect(x, y, x+tileWidth/2, y+1)
}

func (p *Player) Draw(screen *ebiten.Image) {
	a := p.clickableArea()
	ebitenutil.DrawRect(screen, float64(a.Min.X), float64(a.Min.Y), float64(a.Dx()), float64(a.Dy()), color.NRGBA{0, 0, 0xff, 0x40})
	a2 := p.conflictionArea()
	ebitenutil.DrawRect(screen, float64(a2.Min.X), float64(a2.Min.Y), float64(a2.Dx()), float64(a2.Dy()), color.NRGBA{0, 0, 0xff, 0x40})
	a3 := p.elevatorArea()
	ebitenutil.DrawRect(screen, float64(a3.Min.X), float64(a3.Min.Y), float64(a3.Dx()), float64(a3.Dy()), color.NRGBA{0, 0, 0xff, 0xff})
	a4 := p.footArea()
	ebitenutil.DrawRect(screen, float64(a4.Min.X), float64(a4.Min.Y), float64(a4.Dx()), float64(a4.Dy()), color.NRGBA{0, 0, 0xff, 0x80})
}
