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
	"image/color"

	"github.com/hajimehoshi/ebiten"

	"github.com/hajimehoshi/gopherwalk/internal/scene"
)

func New(id int) *GameScene {
	f := strToField(testFields[id])

	x, y := f.StartPosition()
	p := NewPlayer(x, y)

	return &GameScene{
		id:     id,
		player: p,
		field:  f,
	}
}

type GameScene struct {
	id     int
	player *Player
	field  *Field
}

func (s *GameScene) Update(context scene.Context) error {
	s.field.Update(context)
	s.player.Update(context, s.field)

	if s.player.AtGoal() {
		context.GoToGameScene(s.id + 1)
	}

	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0x99, 0xcc, 0xff, 0xff})
	s.field.Draw(screen)
	s.player.Draw(screen)
}
