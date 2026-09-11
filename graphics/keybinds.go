package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Keybind struct {
	Keys [4]ebiten.Key
	Len  int
}

var keybindCallbacks = make(map[Keybind][]keybindCallback)

type keybindCallback struct {
	fn          func()
	onJustPress bool // if callback only on the frame the key is pressed
}

// onJustPress is whether the callback should invoke only the frame the key was pressed
func AddKeybind(key ebiten.Key, callback func(), onJustPress bool) {
	AddKeybindCombination([]ebiten.Key{key}, callback, onJustPress)
}

// max 4 keys. onJustPress is whether the callback should invoke only the frame the key was pressed
func AddKeybindCombination(keys []ebiten.Key, callback func(), onJustPress bool) {
	keybind := Keybind{
		Len: len(keys),
	}

	copy(keybind.Keys[:], keys)

	keybindCallbacks[keybind] = append(keybindCallbacks[keybind], keybindCallback{callback, onJustPress})
}

func RunKeybinds() {
	for key, callbacks := range keybindCallbacks {
		for _, callback := range callbacks {
			if IsKeybindPressed(key, callback.onJustPress) {
				callback.fn()
			}
		}
	}
}

func IsKeybindPressed(keybind Keybind, justPressed bool) bool {
	pressed := false
	for i, key := range keybind.Keys {
		if i >= keybind.Len {
			break
		}
		if justPressed {
			pressed = inpututil.IsKeyJustPressed(key)
		} else {
			pressed = ebiten.IsKeyPressed(key)
		}
	}

	return pressed
}
