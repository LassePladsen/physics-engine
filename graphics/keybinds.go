package graphics

import (
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// keyboardState isolates Ebiten's frame input APIs so keybind behavior can be
// tested without a running Ebiten game.
type keyboardState interface {
	IsKeyPressed(ebiten.Key) bool
	IsKeyJustPressed(ebiten.Key) bool
}

type ebitenKeyboardState struct{}

func (ebitenKeyboardState) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

func (ebitenKeyboardState) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

var keyboard keyboardState = ebitenKeyboardState{}

type registeredKeybind struct {
	keys      []ebiten.Key
	callbacks []keybindCallback
}

// keybindCallbacks is indexed by the canonical form of a combination. The
// registered value retains the normalized keys needed when dispatching.
var keybindCallbacks = make(map[string]*registeredKeybind)

type keybindCallback struct {
	fn          func()
	onJustPress bool // if callback only on the frame the key is pressed
}

// AddKeybind registers a callback for a combination of any number
// of keys. Key order and duplicate keys do not affect the combination.
func AddKeybind(callback func(), onJustPress bool, keys ...ebiten.Key) {
	normalizedKeys := normalizeKeys(keys)
	if len(normalizedKeys) == 0 {
		return
	}

	id := canonicalKey(normalizedKeys)
	keybind, ok := keybindCallbacks[id]
	if !ok {
		keybind = &registeredKeybind{keys: normalizedKeys}
		keybindCallbacks[id] = keybind
	}
	keybind.callbacks = append(keybind.callbacks, keybindCallback{callback, onJustPress})
}

func RunKeybinds() {
	for _, keybind := range keybindCallbacks {
		for _, callback := range keybind.callbacks {
			if areKeysPressed(keybind.keys, callback.onJustPress) {
				callback.fn()
			}
		}
	}
}

func areKeysPressed(keys []ebiten.Key, justPressed bool) bool {
	if len(keys) == 0 {
		return false
	}
	anyJustPressed := false
	for _, key := range keys {
		if !keyboard.IsKeyPressed(key) {
			return false
		}
		if keyboard.IsKeyJustPressed(key) {
			anyJustPressed = true
		}
	}

	return !justPressed || anyJustPressed
}

func normalizeKeys(keys []ebiten.Key) []ebiten.Key {
	normalized := append([]ebiten.Key(nil), keys...)
	sort.Slice(normalized, func(i, j int) bool { return normalized[i] < normalized[j] })

	if len(normalized) == 0 {
		return normalized
	}

	unique := normalized[:1]
	for _, key := range normalized[1:] {
		if key != unique[len(unique)-1] {
			unique = append(unique, key)
		}
	}
	return unique
}

func canonicalKey(keys []ebiten.Key) string {
	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(strconv.Itoa(int(key)))
		builder.WriteByte(',')
	}
	return builder.String()
}
