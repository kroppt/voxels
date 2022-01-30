package input_test

import (
	"fmt"
	"testing"

	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/input"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/veandco/go-sdl2/sdl"
)

const radPerPixel60Fov1080p = 0.00106916675777135701865376418730309896037268813747497991241552305271399806859756004655613930392969573839250030872981391876499815549171599432302996328453614659015666711180688934233446017999583370005037294423634816183735651831579120613095493436003836285184792
const radPerPixel90Fov1080p = 0.00185184973497023151985751879482804751407962318981660498547356045285533064137507471543730185961738618138333648949172377718744841394086817553062678181719377138118952338460917919221716142495657779862348491431188709119747144501408936194412833493296089717961251
const radPerPixel60Fov720p = 0.00160374937279330816829202322728175610492195086447145154086328041257942052266469076005433338466058794268480374657841493528031900554291374010332983185292802664704121007621386567449140497010184846281883017096154046972947380855904302476309350952412298983907694

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := input.New(nil, nil, nil)
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func TestModuleRouteEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns false on quit event", func(t *testing.T) {
		t.Parallel()

		graphicsMod := graphics.FnModule{
			FnPollEvent: func() (sdl.Event, bool) {
				return &sdl.QuitEvent{
					Type:      sdl.QUIT,
					Timestamp: 0,
				}, true
			},
			FnDestroyWindow: func() error {
				return nil
			},
		}
		mod := input.New(graphicsMod, nil, nil)

		quitResult := mod.RouteEvents()
		expected := false
		if quitResult != expected {
			t.Fatal("calling quit did not return false")
		}
	})

	t.Run("returns true after consuming all events", func(t *testing.T) {
		t.Parallel()

		graphicsMod := graphics.FnModule{
			FnPollEvent: func() (sdl.Event, bool) {
				return nil, false
			},
			FnDestroyWindow: func() error {
				return nil
			},
		}
		mod := input.New(graphicsMod, nil, nil)

		exhaustEventsResult := mod.RouteEvents()
		expected := true
		if exhaustEventsResult != expected {
			t.Fatal("exhausting events did not return true")
		}
	})

	t.Run("calls DestroyWindow on quit event", func(t *testing.T) {
		t.Parallel()

		destroyed := false
		graphicsMod := graphics.FnModule{
			FnPollEvent: func() (sdl.Event, bool) {
				return &sdl.QuitEvent{
					Type:      sdl.QUIT,
					Timestamp: 0,
				}, true
			},
			FnDestroyWindow: func() error {
				destroyed = true
				return nil
			},
		}
		mod := input.New(graphicsMod, nil, nil)

		mod.RouteEvents()

		if !destroyed {
			t.Fatal("expected tests to fail")
		}
	})

	testCases := []struct {
		message   string
		evtType   uint32
		scancode  sdl.Scancode
		sym       sdl.Keycode
		direction camera.MoveDirection
		pressed   bool
	}{
		{
			message:   "forward",
			evtType:   sdl.KEYDOWN,
			scancode:  sdl.SCANCODE_W,
			sym:       sdl.K_w,
			direction: camera.MoveForwards,
			pressed:   true,
		},
		{
			message:   "backward",
			evtType:   sdl.KEYDOWN,
			scancode:  sdl.SCANCODE_S,
			sym:       sdl.K_s,
			direction: camera.MoveBackwards,
			pressed:   true,
		},
		{
			message:   "left",
			evtType:   sdl.KEYDOWN,
			scancode:  sdl.SCANCODE_A,
			sym:       sdl.K_a,
			direction: camera.MoveLeft,
			pressed:   true,
		},
		{
			message:   "right",
			evtType:   sdl.KEYDOWN,
			scancode:  sdl.SCANCODE_D,
			sym:       sdl.K_d,
			direction: camera.MoveRight,
			pressed:   true,
		},
		{
			message:   "forward",
			evtType:   sdl.KEYUP,
			scancode:  sdl.SCANCODE_W,
			sym:       sdl.K_w,
			direction: camera.MoveForwards,
			pressed:   false,
		},
		{
			message:   "backward",
			evtType:   sdl.KEYUP,
			scancode:  sdl.SCANCODE_S,
			sym:       sdl.K_s,
			direction: camera.MoveBackwards,
			pressed:   false,
		},
		{
			message:   "left",
			evtType:   sdl.KEYUP,
			scancode:  sdl.SCANCODE_A,
			sym:       sdl.K_a,
			direction: camera.MoveLeft,
			pressed:   false,
		},
		{
			message:   "right",
			evtType:   sdl.KEYUP,
			scancode:  sdl.SCANCODE_D,
			sym:       sdl.K_d,
			direction: camera.MoveRight,
			pressed:   false,
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run("creates a "+tC.message+" movement event from a keyboard event", func(t *testing.T) {
			t.Parallel()

			first := true
			moveKeyboardEvent := sdl.KeyboardEvent{
				Type:      tC.evtType,
				Timestamp: 0,
				WindowID:  0,
				State:     0,
				Repeat:    0,
				Keysym: sdl.Keysym{
					Scancode: tC.scancode,
					Sym:      tC.sym,
					Mod:      sdl.KMOD_NONE,
				},
			}
			quitKeyboardEvent := sdl.QuitEvent{
				Type:      sdl.QUIT,
				Timestamp: 0,
			}
			graphicsMod := graphics.FnModule{
				FnPollEvent: func() (sdl.Event, bool) {
					if first {
						first = false
						return &moveKeyboardEvent, true
					}
					return &quitKeyboardEvent, true
				},
				FnDestroyWindow: func() error {
					return nil
				},
			}

			expectEvent := camera.MovementEvent{
				Direction: tC.direction,
				Pressed:   tC.pressed,
			}
			var actualEvent *camera.MovementEvent
			cameraMod := &camera.FnModule{
				FnHandleMovementEvent: func(evt camera.MovementEvent) {
					actualEvent = &evt
				},
			}
			mod := input.New(graphicsMod, cameraMod, nil)

			mod.RouteEvents()

			if actualEvent == nil {
				t.Fatalf("expected %v but got %v", expectEvent, nil)
			}

			if expectEvent != *actualEvent {
				t.Fatalf("expected %v but got %v", expectEvent, *actualEvent)
			}
		})
	}

	testLookCases := []struct {
		xRel int32
		yRel int32
	}{
		{
			xRel: 1,
			yRel: 2,
		},
		{
			xRel: 0,
			yRel: 2,
		},
		{
			xRel: 1,
			yRel: 0,
		},
		{
			xRel: -1,
			yRel: 2,
		},
		{
			xRel: 1,
			yRel: -2,
		},
		{
			xRel: -1,
			yRel: -2,
		},
	}
	for _, tC := range testLookCases {
		tC := tC
		t.Run(fmt.Sprintf("convert sdl.MouseMotionEvent (%v, %v) to LookEvent", tC.xRel, tC.yRel), func(t *testing.T) {
			motionEvent := sdl.MouseMotionEvent{
				Type:      sdl.MOUSEMOTION,
				Timestamp: 0,
				WindowID:  0,
				Which:     0,
				State:     sdl.BUTTON_LEFT,
				X:         0,
				Y:         0,
				XRel:      tC.xRel,
				YRel:      tC.yRel,
			}

			quitEvent := sdl.QuitEvent{
				Type:      sdl.QUIT,
				Timestamp: 0,
			}

			first := true
			graphicsMod := graphics.FnModule{
				FnPollEvent: func() (sdl.Event, bool) {
					if first {
						first = false
						return &motionEvent, true
					}
					return &quitEvent, true
				},
				FnDestroyWindow: func() error {
					return nil
				},
			}
			var evtHandle *camera.LookEvent
			cameraMod := &camera.FnModule{
				FnHandleLookEvent: func(evt camera.LookEvent) {
					evtHandle = &evt
				},
			}
			settingsRepo := settings.New()
			settingsRepo.SetFOV(60)
			settingsRepo.SetResolution(1920, 1080)
			mod := input.New(graphicsMod, cameraMod, settingsRepo)

			mod.RouteEvents()

			xRad := radPerPixel60Fov1080p * float64(tC.xRel)
			yRad := radPerPixel60Fov1080p * float64(tC.yRel)
			expectLookEvent := camera.LookEvent{
				Right: xRad,
				Down:  yRad,
			}

			if evtHandle == nil {
				t.Fatalf("expected %v but got %v", expectLookEvent, nil)
			}

			if expectLookEvent != *evtHandle {
				t.Fatalf("expected %v but got %v", expectLookEvent, *evtHandle)
			}

		})
	}

}

func TestMouseMotionOnlyPassedToCameraIfM1Held(t *testing.T) {
	motionEvent := sdl.MouseMotionEvent{
		Type:      sdl.MOUSEMOTION,
		Timestamp: 0,
		WindowID:  0,
		Which:     0,
		State:     0,
		X:         0,
		Y:         0,
		XRel:      1,
		YRel:      1,
	}

	quitEvent := sdl.QuitEvent{
		Type:      sdl.QUIT,
		Timestamp: 0,
	}

	first := true
	graphicsMod := graphics.FnModule{
		FnPollEvent: func() (sdl.Event, bool) {
			if first {
				first = false
				return &motionEvent, true
			}
			return &quitEvent, true
		},
	}
	expected := false
	actual := false
	cameraMod := &camera.FnModule{
		FnHandleLookEvent: func(evt camera.LookEvent) {
			actual = true
		},
	}
	mod := input.New(graphicsMod, cameraMod, settings.FnRepository{})
	mod.RouteEvents()
	if actual != expected {
		t.Fatal("expected no handle look event, but there was one")
	}
}

func withinError(x, y float64, diff float64) bool {
	if x+diff > y && x-diff < y {
		return true
	}
	return false
}

const errMargin = 0.000001

func TestModulePixelsToRadians(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		fov     float64
		resX    uint32
		resY    uint32
		xRel    int32
		yRel    int32
		expectX float64
		expectY float64
	}{
		{
			fov:     60.0,
			resX:    1920,
			resY:    1080,
			xRel:    1,
			yRel:    0,
			expectX: radPerPixel60Fov1080p,
			expectY: 0.0,
		},
		{
			fov:     60.0,
			resX:    1920,
			resY:    1080,
			xRel:    0,
			yRel:    1,
			expectX: 0.0,
			expectY: radPerPixel60Fov1080p,
		},
		{
			fov:     90.0,
			resX:    1920,
			resY:    1080,
			xRel:    1,
			yRel:    0,
			expectX: radPerPixel90Fov1080p,
			expectY: 0,
		},
		{
			fov:     90.0,
			resX:    1920,
			resY:    1080,
			xRel:    0,
			yRel:    1,
			expectX: 0.0,
			expectY: radPerPixel90Fov1080p,
		},
		{
			fov:     60.0,
			resX:    1280,
			resY:    720,
			xRel:    1,
			yRel:    0,
			expectX: radPerPixel60Fov720p,
			expectY: 0,
		},
		{
			fov:     60.0,
			resX:    1280,
			resY:    720,
			xRel:    0,
			yRel:    1,
			expectX: 0,
			expectY: radPerPixel60Fov720p,
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("%v at %vx%v", tC.fov, tC.resX, tC.resY), func(t *testing.T) {
			t.Parallel()

			settingsRepo := settings.New()
			settingsRepo.SetFOV(tC.fov)
			settingsRepo.SetResolution(tC.resX, tC.resY)
			mod := input.New(nil, nil, settingsRepo)

			xRad, yRad := mod.PixelsToRadians(tC.xRel, tC.yRel)

			if !withinError(xRad, tC.expectX, errMargin) {
				t.Fatalf("expected %v but got %v", tC.expectX, xRad)
			}
			if !withinError(yRad, tC.expectY, errMargin) {
				t.Fatalf("expected %v but got %v", tC.expectY, yRad)
			}
		})
	}
}
