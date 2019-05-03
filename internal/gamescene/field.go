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
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	tileWidth  = 16
	tileHeight = 16
)

type Tile int

const (
	TileNone Tile = iota
	TilePart
	TileOutside
	TileWallBig
	TileWallSmall
	TileFFBig
	TileFFSmall
	TileElevator
	TileStart
)

func (t Tile) IsBig() bool {
	switch t {
	case TileWallBig, TileFFBig:
		return true
	}
	return false
}

type Field struct {
	tiles [][]Tile
}

func (f *Field) StartPosition() (x, y int) {
	for j, row := range f.tiles {
		for i, tile := range row {
			if tile == TileStart {
				return i, j
			}
		}
	}
	panic("gamescene: start position is not found")
}

func (f *Field) touchingTiles(rect image.Rectangle) []Tile {
	var tiles []Tile
	minX := rect.Min.X / tileWidth
	maxX := rect.Max.X / tileWidth
	minY := (rect.Min.Y - 1) / tileHeight
	maxY := (rect.Max.Y - 1) / tileHeight
	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			tiles = append(tiles, f.tileAt(i, j))
		}
	}
	return tiles
}

func (f *Field) Conflicts(rect image.Rectangle) bool {
	for _, t := range f.touchingTiles(rect) {
		switch t {
		case TileOutside, TileWallBig, TileWallSmall:
			return true
		}
	}
	return false
}

func (f *Field) ConflictsWithFoot(rect image.Rectangle) bool {
	for _, t := range f.touchingTiles(rect) {
		switch t {
		case TileOutside, TileWallBig, TileWallSmall, TileElevator:
			return true
		}
	}
	return false
}

func (f *Field) InElevator(rect image.Rectangle) bool {
	for _, t := range f.touchingTiles(rect) {
		if t == TileElevator {
			return true
		}
	}
	return false
}

func (f *Field) tileAt(x, y int) Tile {
	if y < 0 || len(f.tiles) <= y {
		return TileOutside
	}
	if x < 0 || len(f.tiles[y]) <= x {
		return TileOutside
	}
	t := f.tiles[y][x]
	if t != TilePart {
		return t
	}
	if y > 0 {
		if x > 0 {
			if t := f.tiles[y-1][x-1]; t.IsBig() {
				return t
			}
		}
		if t := f.tiles[y-1][x]; t.IsBig() {
			return t
		}
	}
	if x > 0 {
		if t := f.tiles[y][x-1]; t.IsBig() {
			return t
		}
	}
	return t
}

func (f *Field) Draw(screen *ebiten.Image) {
	for j, row := range f.tiles {
		for i, tile := range row {
			x := i * tileWidth
			y := j * tileHeight
			w := 0
			h := 0
			if tile.IsBig() {
				w = tileWidth*2 - 1
				h = tileHeight*2 - 1
			} else {
				w = tileWidth - 1
				h = tileHeight - 1
			}
			switch tile {
			case TileWallSmall:
				ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{0x66, 0x66, 0x66, 0xff})
			case TileWallBig:
				ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{0x66, 0x66, 0x66, 0xff})
			case TileFFSmall:
				ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{0xff, 0x00, 0x00, 0x40})
			case TileFFBig:
				ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{0xff, 0x00, 0x00, 0x40})
			case TileElevator:
				ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{0xff, 0xff, 0x00, 0xff})
			}
		}
	}
}

const testField = `
w              w
w              w
w              w
w              w
w              w
w              w
w              w
w              w
w              w
weW.F.F.F.F.W. w
we............ w
we             w
we             w
we            sw
wwwwwwwwwwwwwwww
`

func strToField(str string) *Field {
	var tiles [][]Tile
	for _, line := range strings.Split(testField, "\n") {
		if line == "" {
			continue
		}
		var row []Tile
		for _, c := range line {
			switch c {
			case 'W':
				row = append(row, TileWallBig)
			case 'w':
				row = append(row, TileWallSmall)
			case 'F':
				row = append(row, TileFFBig)
			case 'f':
				row = append(row, TileFFSmall)
			case 'e':
				row = append(row, TileElevator)
			case 's':
				row = append(row, TileStart)
			case '.':
				row = append(row, TilePart)
			default:
				row = append(row, TileNone)
			}
		}
		tiles = append(tiles, row)
	}
	return &Field{
		tiles: tiles,
	}
}
