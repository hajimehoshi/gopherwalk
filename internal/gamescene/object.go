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

const (
	tileWidth  = 16
	tileHeight = 16
)

type Object interface {
	Touches(rect image.Rectangle) bool
	Conflicts(rect image.Rectangle) bool
	ConflictsWithFoot(rect image.Rectangle) bool

	Update(context scene.Context)
	Draw(screen *ebiten.Image)
}

type ObjectWall struct {
	big bool
	x   int
	y   int
}

func (o *ObjectWall) Touches(rect image.Rectangle) bool {
	w := tileWidth
	h := tileHeight
	if o.big {
		w *= 2
		h *= 2
	}
	return image.Rect(o.x*tileWidth, o.y*tileHeight, o.x*tileWidth+w, o.y*tileHeight+h).Overlaps(rect)
}

func (o *ObjectWall) Conflicts(rect image.Rectangle) bool {
	return o.Touches(rect)
}

func (o *ObjectWall) ConflictsWithFoot(rect image.Rectangle) bool {
	return o.Touches(rect)
}

func (o *ObjectWall) Update(context scene.Context) {
}

func (o *ObjectWall) Draw(screen *ebiten.Image) {
	x := o.x * tileWidth
	y := o.y * tileHeight
	if o.big {
		w := tileWidth*2 - 1
		h := tileWidth*2 - 1
		ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.NRGBA{0x66, 0x66, 0x66, 0xff})
	} else {
		w := tileWidth - 1
		h := tileWidth - 1
		ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.NRGBA{0x66, 0x66, 0x66, 0xff})
	}
}

type ObjectFF struct {
	big bool
	x   int
	y   int

	on bool
}

func (o *ObjectFF) area() image.Rectangle {
	w := tileWidth
	h := tileHeight
	if o.big {
		w *= 2
		h *= 2
	}
	return image.Rect(o.x*tileWidth, o.y*tileHeight, o.x*tileWidth+w, o.y*tileHeight+h)
}

func (o *ObjectFF) Touches(rect image.Rectangle) bool {
	if !o.on {
		return false
	}
	w := tileWidth
	h := tileHeight
	if o.big {
		w *= 2
		h *= 2
	}
	return o.area().Overlaps(rect)
}

func (o *ObjectFF) Conflicts(rect image.Rectangle) bool {
	return o.Touches(rect)
}

func (o *ObjectFF) ConflictsWithFoot(rect image.Rectangle) bool {
	return o.Touches(rect)
}

func (o *ObjectFF) Update(context scene.Context) {
	if !context.Input().IsJustTapped() {
		return
	}
	x, y := context.Input().CursorPosition()
	if !image.Pt(x, y).In(o.area()) {
		return
	}
	o.on = !o.on
}

func (o *ObjectFF) Draw(screen *ebiten.Image) {
	c := color.NRGBA{0xff, 0x00, 0x00, 0x40}
	if o.on {
		c = color.NRGBA{0xff, 0x00, 0x00, 0xff}
	}
	x := o.x * tileWidth
	y := o.y * tileHeight
	if o.big {
		w := tileWidth*2 - 1
		h := tileWidth*2 - 1
		ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), c)
	} else {
		w := tileWidth - 1
		h := tileWidth - 1
		ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), c)
	}
}

type ObjectElevator struct {
	x int
	y int
}

func (o *ObjectElevator) Touches(rect image.Rectangle) bool {
	w := tileWidth
	h := tileHeight
	return image.Rect(o.x*tileWidth, o.y*tileHeight, o.x*tileWidth+w, o.y*tileHeight+h).Overlaps(rect)
}

func (o *ObjectElevator) Conflicts(rect image.Rectangle) bool {
	return false
}

func (o *ObjectElevator) ConflictsWithFoot(rect image.Rectangle) bool {
	return o.Touches(rect)
}

func (o *ObjectElevator) Update(context scene.Context) {
}

func (o *ObjectElevator) Draw(screen *ebiten.Image) {
	x := o.x * tileWidth
	y := o.y * tileHeight
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(tileWidth-1), float64(tileHeight-1), color.NRGBA{0xff, 0xff, 0x00, 0xff})
}
