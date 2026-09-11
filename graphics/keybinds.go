package graphics

import (
	"slices"
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
	keys        []ebiten.Key
	callbacks   []func()
	onJustPress bool // if callback only on the frame the key is pressed
}

func (r *registeredKeybind) isPressed() bool {
	isKeyPressed := keyboard.IsKeyPressed
	if r.onJustPress {
		isKeyPressed = keyboard.IsKeyJustPressed
	}

	for _, key := range r.keys {
		if isKeyPressed(key) {
			return true
		}
	}

	return false
}

func (r *registeredKeybind) invokeIfPressed() {
	if !r.isPressed() {
		return
	}
	for _, callback := range r.callbacks {
		callback()
	}

}

// keybindCallbacks is indexed by the canonical form of a combination. The
// registered value retains the normalized keys needed when dispatching.
var keybindCallbacks = make(map[string]*registeredKeybind)

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
		keybind = &registeredKeybind{keys: normalizedKeys, onJustPress: onJustPress}
		keybindCallbacks[id] = keybind
	}
	keybind.callbacks = append(keybind.callbacks, callback)
}

// Runs all keybinds if the corresponding keys are pressed
func RunKeybinds() {
	for _, keybind := range keybindCallbacks {
		keybind.invokeIfPressed()
	}
}

func normalizeKeys(keys []ebiten.Key) []ebiten.Key {
	normalized := append([]ebiten.Key(nil), keys...)
	slices.Sort(normalized)
	return slices.Compact(normalized)
}

func canonicalKey(keys []ebiten.Key) string {
	// Delimiters preserve boundaries between key values (for example, 1+23 vs 12+3).
	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(strconv.Itoa(int(key)))
		builder.WriteByte(',')
	}
	return builder.String()
}
