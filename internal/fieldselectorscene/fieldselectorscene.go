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

package fieldselectorscene

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/gopherwalk/internal/scene"
)

type FieldSelectorScene struct {
	fieldButtons []*Button
	selected     int
}

func New() *FieldSelectorScene {
	const (
		w = 24
		h = 16
	)

	s := &FieldSelectorScene{}

	var bs []*Button
	for i := 0; i < 30; i++ {
		id := i + 1
		x := (i%10)*w + 8
		y := (i/10)*h + 8
		b := NewButton(image.Rect(x, y, x+w, y+h), fmt.Sprintf("%d", id))
		b.SetOnTap(func() {
			s.selected = id
		})
		bs = append(bs, b)
	}
	s.fieldButtons = bs

	return s
}

func (s *FieldSelectorScene) Update(context scene.Context) error {
	for _, b := range s.fieldButtons {
		b.Update(context.Input())
	}
	if s.selected != 0 {
		context.GoToGameScene(s.selected)
	}
	return nil
}

func (s *FieldSelectorScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for _, b := range s.fieldButtons {
		b.Draw(screen)
	}
}

const (
	lineHeight = 16
)

type Button struct {
	rect  image.Rectangle
	text  string
	ontap func()
	hover bool
}

func NewButton(rect image.Rectangle, text string) *Button {
	return &Button{
		rect: rect,
		text: text,
	}
}

func (b *Button) SetOnTap(f func()) {
	b.ontap = f
}

func (b *Button) Update(input scene.Input) {
	x, y := input.CursorPosition()
	b.hover = image.Pt(x, y).In(b.rect)
	if b.hover && input.IsJustTapped() && b.ontap != nil {
		b.ontap()
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	bound, _ := font.BoundString(bitmapfont.Gothic12r, b.text)
	bw := (bound.Max.X - bound.Min.X).Ceil()
	bh := (bound.Max.Y - bound.Min.Y).Ceil()
	x := b.rect.Min.X + (b.rect.Dx()-bw)/2 + 4
	y := b.rect.Min.Y + (b.rect.Dy()-bh)/2 + 12
	clr := color.NRGBA{0, 0, 0, 0xff}
	if b.hover {
		clr = color.NRGBA{0xff, 0, 0, 0xff}
	}
	text.Draw(screen, b.text, bitmapfont.Gothic12r, x, y, clr)
}
