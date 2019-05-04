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
	"strings"

	"github.com/hajimehoshi/ebiten"

	"github.com/hajimehoshi/gopherwalk/internal/scene"
)

type Field struct {
	objects []Object
	startX  int
	startY  int
}

func (f *Field) StartPosition() (x, y int) {
	return f.startX, f.startY
}

func (f *Field) Conflicts(rect image.Rectangle) bool {
	for _, t := range f.objects {
		if t.Conflicts(rect) {
			return true
		}
	}
	return false
}

func (f *Field) ConflictsWithFoot(rect image.Rectangle) bool {
	for _, t := range f.objects {
		if t.ConflictsWithFoot(rect) {
			return true
		}
	}
	return false
}

func (f *Field) InElevator(rect image.Rectangle) bool {
	for _, t := range f.objects {
		if _, ok := t.(*ObjectElevator); ok {
			if t.Touches(rect) {
				return true
			}
		}
	}
	return false
}

func (f *Field) Update(context scene.Context) {
	for _, t := range f.objects {
		t.Update(context)
	}
}

func (f *Field) Draw(screen *ebiten.Image) {
	for _, t := range f.objects {
		t.Draw(screen)
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
	f := &Field{}
	for j, line := range strings.Split(strings.TrimSpace(testField), "\n") {
		for i, c := range line {
			switch c {
			case 'W':
				f.objects = append(f.objects, &ObjectWall{big: true, x: i, y: j})
			case 'w':
				f.objects = append(f.objects, &ObjectWall{big: false, x: i, y: j})
			case 'F':
				f.objects = append(f.objects, &ObjectFF{big: true, x: i, y: j})
			case 'f':
				f.objects = append(f.objects, &ObjectFF{big: false, x: i, y: j})
			case 'e':
				f.objects = append(f.objects, &ObjectElevator{x: i, y: j})
			case 's':
				f.startX = i
				f.startY = j
			case '.':
			default:
			}
		}
	}
	return f
}
