package graphics

import (
	"reflect"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type simulatedKeyboard struct {
	held        map[ebiten.Key]bool
	justPressed map[ebiten.Key]bool
}

func (keyboard simulatedKeyboard) IsKeyPressed(key ebiten.Key) bool {
	return keyboard.held[key]
}

func (keyboard simulatedKeyboard) IsKeyJustPressed(key ebiten.Key) bool {
	return keyboard.justPressed[key]
}

func resetKeybinds(t *testing.T, state keyboardState) {
	t.Helper()
	previousCallbacks, previousKeyboard := keybindCallbacks, keyboard
	keybindCallbacks = make(map[string]*registeredKeybind)
	keyboard = state
	t.Cleanup(func() {
		keybindCallbacks = previousCallbacks
		keyboard = previousKeyboard
	})
}

func TestAddKeybindRegistersSingleKeyCallback(t *testing.T) {
	resetKeybinds(t, simulatedKeyboard{})
	callback := func() {}
	AddKeybind(callback, true, ebiten.KeyA)

	if len(keybindCallbacks) != 1 {
		t.Fatalf("registered combinations = %d, want 1", len(keybindCallbacks))
	}
	for _, binding := range keybindCallbacks {
		if !reflect.DeepEqual(binding.keys, []ebiten.Key{ebiten.KeyA}) {
			t.Fatalf("keys = %v, want [%v]", binding.keys, ebiten.KeyA)
		}
		if len(binding.callbacks) != 1 || !binding.callbacks[0].onJustPress {
			t.Fatalf("callback registration = %+v, want one just-press callback", binding.callbacks)
		}
	}
}

func TestAddKeybindCombinationNormalizesAndAccumulatesCallbacks(t *testing.T) {
	resetKeybinds(t, simulatedKeyboard{})
	AddKeybind(func() {}, false, ebiten.KeyC, ebiten.KeyControl, ebiten.KeyC)
	AddKeybind(func() {}, true, ebiten.KeyControl, ebiten.KeyC)

	if len(keybindCallbacks) != 1 {
		t.Fatalf("registered combinations = %d, want 1", len(keybindCallbacks))
	}
	for _, binding := range keybindCallbacks {
		if !reflect.DeepEqual(binding.keys, []ebiten.Key{ebiten.KeyC, ebiten.KeyControl}) {
			t.Fatalf("keys = %v, want normalized keys", binding.keys)
		}
		if len(binding.callbacks) != 2 {
			t.Fatalf("callbacks = %d, want 2", len(binding.callbacks))
		}
	}
}

func TestAddKeybindCombinationSupportsMoreThanFourKeys(t *testing.T) {
	state := simulatedKeyboard{held: map[ebiten.Key]bool{
		ebiten.KeyA: true, ebiten.KeyB: true, ebiten.KeyC: true, ebiten.KeyD: true, ebiten.KeyE: true,
	}}
	resetKeybinds(t, state)

	called := 0
	AddKeybind(func() { called++ }, false, ebiten.KeyE, ebiten.KeyD, ebiten.KeyC, ebiten.KeyB, ebiten.KeyA)
	RunKeybinds()

	if called != 1 {
		t.Fatalf("callback calls = %d, want 1", called)
	}
}

func TestAddKeybindCombinationIgnoresEmptyCombination(t *testing.T) {
	resetKeybinds(t, simulatedKeyboard{})
	AddKeybind(func() {}, false)

	if len(keybindCallbacks) != 0 {
		t.Fatalf("registered combinations = %d, want 0", len(keybindCallbacks))
	}
}

func TestRunKeybindsRequiresAllKeysAndHonorsJustPress(t *testing.T) {
	state := &simulatedKeyboard{
		held:        map[ebiten.Key]bool{ebiten.KeyControl: true},
		justPressed: map[ebiten.Key]bool{ebiten.KeyControl: true},
	}
	resetKeybinds(t, state)

	normalCalls, justPressCalls := 0, 0
	AddKeybind(func() { normalCalls++ }, false, ebiten.KeyControl, ebiten.KeyC)
	AddKeybind(func() { justPressCalls++ }, true, ebiten.KeyC, ebiten.KeyControl)
	RunKeybinds()
	if normalCalls != 0 || justPressCalls != 0 {
		t.Fatalf("incomplete combination calls = normal %d, just press %d; want 0, 0", normalCalls, justPressCalls)
	}

	state.held[ebiten.KeyC] = true
	RunKeybinds()
	if normalCalls != 1 || justPressCalls != 1 {
		t.Fatalf("held combination calls = normal %d, just press %d; want 1, 1", normalCalls, justPressCalls)
	}

	state.justPressed = map[ebiten.Key]bool{}
	RunKeybinds()
	if normalCalls != 2 || justPressCalls != 1 {
		t.Fatalf("held combination without new press = normal %d, just press %d; want 2, 1", normalCalls, justPressCalls)
	}
}
