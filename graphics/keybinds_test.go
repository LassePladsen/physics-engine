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

func (k simulatedKeyboard) IsKeyPressed(key ebiten.Key) bool {
	return k.held[key]
}

func (k simulatedKeyboard) IsKeyJustPressed(key ebiten.Key) bool {
	return k.justPressed[key]
}

func useKeyboard(t *testing.T, state keyboardState) {
	t.Helper()
	previousKeyboard, previousKeybinds := keyboard, keybindCallbacks
	keyboard = state
	keybindCallbacks = make(map[string]*registeredKeybind)
	t.Cleanup(func() {
		keyboard = previousKeyboard
		keybindCallbacks = previousKeybinds
	})
}

func TestAddKeybindNormalizesKeysAndAppendsCallbacks(t *testing.T) {
	useKeyboard(t, simulatedKeyboard{})
	first, second := func() {}, func() {}

	RegisterKeybind(first, true, ebiten.KeyC, ebiten.KeyControl, ebiten.KeyC)
	RegisterKeybind(second, true, ebiten.KeyControl, ebiten.KeyC)

	if got := len(keybindCallbacks); got != 1 {
		t.Fatalf("registered bindings = %d, want 1", got)
	}
	for _, binding := range keybindCallbacks {
		if got, want := binding.keys, []ebiten.Key{ebiten.KeyC, ebiten.KeyControl}; !reflect.DeepEqual(got, want) {
			t.Fatalf("keys = %v, want %v", got, want)
		}
		if !binding.onJustPress {
			t.Fatal("onJustPress = false, want true")
		}
		if got, want := len(binding.callbacks), 2; got != want {
			t.Fatalf("callbacks = %d, want %d", got, want)
		}
	}
}

func TestAddKeybindPreservesJustPressedMode(t *testing.T) {
	useKeyboard(t, simulatedKeyboard{})

	RegisterKeybind(func() {}, true, ebiten.KeyP)

	binding := keybindCallbacks[canonicalKey([]ebiten.Key{ebiten.KeyP})]
	if binding == nil {
		t.Fatal("binding was not registered")
	}
	if !binding.onJustPress {
		t.Fatal("onJustPress = false, want true")
	}
}

func TestAddKeybindIgnoresEmptyCombinations(t *testing.T) {
	useKeyboard(t, simulatedKeyboard{})

	RegisterKeybind(func() {}, false)

	if got := len(keybindCallbacks); got != 0 {
		t.Fatalf("registered bindings = %d, want 0", got)
	}
}

func TestRunKeybindsRunsHeldBindingsEveryFrame(t *testing.T) {
	state := &simulatedKeyboard{held: map[ebiten.Key]bool{ebiten.KeyA: true}}
	useKeyboard(t, state)
	called := 0
	RegisterKeybind(func() { called++ }, false, ebiten.KeyA)

	RunKeybinds()
	RunKeybinds()

	if got := called; got != 2 {
		t.Fatalf("callback calls = %d, want 2", got)
	}
}

func TestRunKeybindsRunsJustPressedBindingsOnlyOnPressFrame(t *testing.T) {
	state := &simulatedKeyboard{
		held:        map[ebiten.Key]bool{ebiten.KeyA: true},
		justPressed: map[ebiten.Key]bool{ebiten.KeyA: true},
	}
	useKeyboard(t, state)
	called := 0
	RegisterKeybind(func() { called++ }, true, ebiten.KeyA)

	RunKeybinds()
	state.justPressed[ebiten.KeyA] = false
	RunKeybinds()

	if got := called; got != 1 {
		t.Fatalf("callback calls = %d, want 1", got)
	}
}

func TestRunKeybindsInvokesCombinationOnlyWhenAllKeysArePressed(t *testing.T) {
	state := &simulatedKeyboard{held: map[ebiten.Key]bool{ebiten.KeyControl: true}}
	useKeyboard(t, state)
	called := 0
	RegisterKeybind(func() { called++ }, false, ebiten.KeyC, ebiten.KeyControl)

	RunKeybinds()
	state.held[ebiten.KeyC] = true
	RunKeybinds()

	if got := called; got != 1 {
		t.Fatalf("callback calls = %d, want 1", got)
	}
}

func TestCanonicalKeyDistinguishesKeyBoundaries(t *testing.T) {
	if canonicalKey([]ebiten.Key{1, 23}) == canonicalKey([]ebiten.Key{12, 3}) {
		t.Fatal("different key sequences produced the same canonical key")
	}
}
