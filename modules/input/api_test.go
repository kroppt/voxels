package input_test

import (
	"fmt"
	"testing"

	"github.com/kroppt/voxels/modules/input"
	"github.com/kroppt/voxels/modules/player"
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

type fnGraphicsMod struct {
	FnPollEvent     func() (sdl.Event, bool)
	FnDestroyWindow func() error
}

func (fn fnGraphicsMod) PollEvent() (sdl.Event, bool) {
	return fn.FnPollEvent()
}

func (fn fnGraphicsMod) DestroyWindow() error {
	return fn.FnDestroyWindow()
}

type fnPlayerMod struct {
	FnHandleMovementEvent func(player.MovementEvent)
	FnHandleLookEvent     func(player.LookEvent)
}

func (fn fnPlayerMod) HandleMovementEvent(evt player.MovementEvent) {
	fn.FnHandleMovementEvent(evt)
}

func (fn fnPlayerMod) HandleLookEvent(evt player.LookEvent) {
	fn.FnHandleLookEvent(evt)
}

func TestModuleRouteEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns on quit event", func(t *testing.T) {
		t.Parallel()

		graphicsMod := fnGraphicsMod{
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

		mod.RouteEvents()
	})

	t.Run("calls DestroyWindow on quit event", func(t *testing.T) {
		t.Parallel()

		destroyed := false
		graphicsMod := fnGraphicsMod{
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
		scancode  sdl.Scancode
		sym       sdl.Keycode
		direction player.MoveDirection
	}{
		{
			message:   "forward",
			scancode:  sdl.SCANCODE_W,
			sym:       sdl.K_w,
			direction: player.MoveForwards,
		},
		{
			message:   "backward",
			scancode:  sdl.SCANCODE_S,
			sym:       sdl.K_s,
			direction: player.MoveBackwards,
		},
		{
			message:   "left",
			scancode:  sdl.SCANCODE_A,
			sym:       sdl.K_a,
			direction: player.MoveLeft,
		},
		{
			message:   "right",
			scancode:  sdl.SCANCODE_D,
			sym:       sdl.K_d,
			direction: player.MoveRight,
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run("creates a "+tC.message+" movement event from a keyboard event", func(t *testing.T) {
			t.Parallel()

			first := true
			moveKeyboardEvent := sdl.KeyboardEvent{
				Type:      sdl.KEYDOWN,
				Timestamp: 0,
				WindowID:  0,
				State:     sdl.PRESSED,
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
			graphicsMod := fnGraphicsMod{
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

			movementEvent := player.MovementEvent{
				Direction: tC.direction,
			}
			var evtHandle *player.MovementEvent
			playerMod := &fnPlayerMod{
				FnHandleMovementEvent: func(evt player.MovementEvent) {
					evtHandle = &evt
				},
			}
			mod := input.New(graphicsMod, playerMod, nil)

			mod.RouteEvents()

			if evtHandle == nil {
				t.Fatalf("expected %v but got %v", movementEvent, nil)
			}

			if movementEvent != *evtHandle {
				t.Fatalf("expected %v but got %v", movementEvent, *evtHandle)
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
				State:     0,
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
			graphicsMod := fnGraphicsMod{
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
			var evtHandle *player.LookEvent
			playerMod := &fnPlayerMod{
				FnHandleLookEvent: func(evt player.LookEvent) {
					evtHandle = &evt
				},
			}
			settingsRepo := settings.New(nil)
			settingsRepo.SetFOV(60)
			settingsRepo.SetResolution(1920, 1080)
			mod := input.New(graphicsMod, playerMod, settingsRepo)

			mod.RouteEvents()

			xRad := float32(radPerPixel60Fov1080p * float64(tC.xRel))
			yRad := float32(radPerPixel60Fov1080p * float64(tC.yRel))
			expectLookEvent := player.LookEvent{
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

func TestModulePixelsToRadians(t *testing.T) {
	t.Parallel()

	t.Run("x at 60 fov 1080p", func(t *testing.T) {
		settingsRepo := settings.New(nil)
		settingsRepo.SetFOV(60)
		settingsRepo.SetResolution(1920, 1080)
		mod := input.New(nil, nil, settingsRepo)

		xRad, yRad := mod.PixelsToRadians(1, 0)

		expectedX := float32(radPerPixel60Fov1080p)
		expectedY := float32(0)
		if xRad != expectedX {
			t.Fatalf("expected %v but got %v", expectedX, xRad)
		}
		if yRad != expectedY {
			t.Fatalf("expected %v but got %v", expectedY, yRad)
		}
	})

	t.Run("y at 60 fov 1080p", func(t *testing.T) {
		settingsRepo := settings.New(nil)
		settingsRepo.SetFOV(60)
		settingsRepo.SetResolution(1920, 1080)
		mod := input.New(nil, nil, settingsRepo)

		xRad, yRad := mod.PixelsToRadians(0, 1)

		expectedX := float32(0.0)
		expectedY := float32(radPerPixel60Fov1080p)
		if xRad != expectedX {
			t.Fatalf("expected %v but got %v", expectedX, xRad)
		}
		if yRad != expectedY {
			t.Fatalf("expected %v but got %v", expectedY, yRad)
		}
	})

	t.Run("x at 90 fov 1080p", func(t *testing.T) {
		settingsRepo := settings.New(nil)
		settingsRepo.SetFOV(90)
		settingsRepo.SetResolution(1920, 1080)
		mod := input.New(nil, nil, settingsRepo)

		xRad, yRad := mod.PixelsToRadians(1, 0)

		expectedX := float32(radPerPixel90Fov1080p)
		expectedY := float32(0)
		if xRad != expectedX {
			t.Fatalf("expected %v but got %v", expectedX, xRad)
		}
		if yRad != expectedY {
			t.Fatalf("expected %v but got %v", expectedY, yRad)
		}
	})

	t.Run("y at 90 fov 1080p", func(t *testing.T) {
		settingsRepo := settings.New(nil)
		settingsRepo.SetFOV(90)
		settingsRepo.SetResolution(1920, 1080)
		mod := input.New(nil, nil, settingsRepo)

		xRad, yRad := mod.PixelsToRadians(0, 1)

		expectedX := float32(0.0)
		expectedY := float32(radPerPixel90Fov1080p)
		if xRad != expectedX {
			t.Fatalf("expected %v but got %v", expectedX, xRad)
		}
		if yRad != expectedY {
			t.Fatalf("expected %v but got %v", expectedY, yRad)
		}
	})

	t.Run("x at 60 fov 720p", func(t *testing.T) {
		settingsRepo := settings.New(nil)
		settingsRepo.SetFOV(60)
		settingsRepo.SetResolution(1280, 720)
		mod := input.New(nil, nil, settingsRepo)

		xRad, yRad := mod.PixelsToRadians(1, 0)

		expectedX := float32(radPerPixel60Fov720p)
		expectedY := float32(0)
		if xRad != expectedX {
			t.Fatalf("expected %v but got %v", expectedX, xRad)
		}
		if yRad != expectedY {
			t.Fatalf("expected %v but got %v", expectedY, yRad)
		}
	})

	t.Run("y at 60 fov 720p", func(t *testing.T) {
		settingsRepo := settings.New(nil)
		settingsRepo.SetFOV(60)
		settingsRepo.SetResolution(1280, 720)
		mod := input.New(nil, nil, settingsRepo)

		xRad, yRad := mod.PixelsToRadians(0, 1)

		expectedX := float32(0.0)
		expectedY := float32(radPerPixel60Fov720p)
		if xRad != expectedX {
			t.Fatalf("expected %v but got %v", expectedX, xRad)
		}
		if yRad != expectedY {
			t.Fatalf("expected %v but got %v", expectedY, yRad)
		}
	})
}
