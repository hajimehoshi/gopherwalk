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

package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/gopherwalk/internal/fieldselectorscene"
	"github.com/hajimehoshi/gopherwalk/internal/gamescene"
	"github.com/hajimehoshi/gopherwalk/internal/scene"
	"github.com/hajimehoshi/gopherwalk/internal/titlescene"
)

const (
	screenWidth  = 256
	screenHeight = 240
)

type SceneManager struct {
	current scene.Scene
	next    scene.Scene
	turbo   bool
}

func (s *SceneManager) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		s.turbo = !s.turbo
	}

	if s.next != nil {
		s.current = s.next
		s.next = nil
	}
	if s.current == nil {
		s.current = &titlescene.TitleScene{}
	}

	n := 1
	if s.turbo {
		n = 5
	}
	for i := 0; i < n; i++ {
		if err := s.current.Update(s); err != nil {
			return nil
		}
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	s.current.Draw(screen)
	return nil
}

func (s *SceneManager) GoToTitleScene() {
	s.next = &titlescene.TitleScene{}
}

func (s *SceneManager) GoToFieldSelectorScene() {
	s.next = fieldselectorscene.New()
}

func (s *SceneManager) GoToGameScene(id int) {
	s.next = gamescene.New(id)
}

func (s *SceneManager) Input() scene.Input {
	return s
}

func (s *SceneManager) CursorPosition() (x, y int) {
	return ebiten.CursorPosition()
}

func (s *SceneManager) IsJustTapped() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}
