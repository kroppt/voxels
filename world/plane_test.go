package world_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/kroppt/voxels/world"
)

func TestTablePlaneGet(t *testing.T) {
	t.Parallel()
	r := world.Range{-1, 1}
	p, err := world.NewPlane(nil, r, r, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	testCases := []world.Position{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 0},
		{0, 1, 1},
		{1, 1, 1},
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
		{-1, -1, 0},
		{0, -1, -1},
		{-1, -1, -1},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("at %v", tC), func(t *testing.T) {
			t.Parallel()
			vox, err := p.At(tC)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if vox.Position.X != tC.X || vox.Position.Y != tC.Y {
				t.Fatalf("expected %v but got %v", tC, vox.Position)
			}
		})
	}
}

func TestPlaneAt(t *testing.T) {
	t.Parallel()

	t.Run("default voxel color should be transparent", func(t *testing.T) {
		t.Parallel()
		r := world.Range{0, 0}
		p, err := world.NewPlane(nil, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expect := world.Color{0, 0, 0, 0}

		vox, err := p.At(world.Position{0, 0, 0})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if vox.Color != expect {
			t.Fatalf("expected %v but got %v", expect, vox.Color)
		}
	})

	t.Run("changing voxel fields should persist", func(t *testing.T) {
		t.Parallel()
		r := world.Range{0, 0}
		p, err := world.NewPlane(nil, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		pos := world.Position{0, 0, 0}
		expect := world.Color{1.0, 1.0, 1.0, 1.0}

		vox, err := p.At(pos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		vox.Color = expect
		vox, err = p.At(pos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if vox.Color != expect {
			t.Fatalf("expected %v but got %v", expect, vox.Color)
		}
	})

}

func TestTablePlaneSize(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		x world.Range
		y world.Range
		z world.Range
	}{
		{
			x: world.Range{0, 0},
			y: world.Range{0, 0},
			z: world.Range{0, 0},
		},
		{
			x: world.Range{0, 1},
			y: world.Range{0, 1},
			z: world.Range{0, 1},
		},
		{
			x: world.Range{1, 0},
			y: world.Range{1, 0},
			z: world.Range{1, 0},
		},
		{
			x: world.Range{1, 1},
			y: world.Range{1, 1},
			z: world.Range{1, 1},
		},
		{
			x: world.Range{-1, 1},
			y: world.Range{0, 0},
			z: world.Range{0, 0},
		},
		{
			x: world.Range{0, 0},
			y: world.Range{-2, 2},
			z: world.Range{0, 0},
		},
		{
			x: world.Range{0, 0},
			y: world.Range{0, 0},
			z: world.Range{-3, 3},
		},
		{
			x: world.Range{-1, 1},
			y: world.Range{-2, 2},
			z: world.Range{-3, 3},
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("x: %v, y: %v, z: %v", tC.x, tC.y, tC.z), func(t *testing.T) {
			p, err := world.NewPlane(nil, tC.x, tC.y, tC.z)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			x, y, z := p.Size()

			if x != tC.x {
				t.Fatalf("expected %v but got %v", tC.x, x)
			}
			if y != tC.y {
				t.Fatalf("expected %v but got %v", tC.y, y)
			}
			if z != tC.z {
				t.Fatalf("expected %v but got %v", tC.z, z)
			}
		})
	}
}

func TestTablePlaneRange(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		x world.Range
		y world.Range
		z world.Range
	}{
		{
			x: world.Range{0, 0},
			y: world.Range{0, 0},
			z: world.Range{0, 0},
		},
		{
			x: world.Range{1, 0},
			y: world.Range{1, 0},
			z: world.Range{1, 0},
		},
		{
			x: world.Range{0, 1},
			y: world.Range{0, 1},
			z: world.Range{0, 1},
		},
		{
			x: world.Range{1, 1},
			y: world.Range{1, 1},
			z: world.Range{1, 1},
		},
		{
			x: world.Range{-1, 0},
			y: world.Range{-1, 0},
			z: world.Range{-1, 0},
		},
		{
			x: world.Range{0, -1},
			y: world.Range{0, -1},
			z: world.Range{0, -1},
		},
		{
			x: world.Range{-1, -1},
			y: world.Range{-1, -1},
			z: world.Range{-1, -1},
		},
	}

	for _, tC := range testCases {

		t.Run(fmt.Sprintf("x: %v, y: %v, z: %v", tC.x, tC.y, tC.z), func(t *testing.T) {
			t.Parallel()
			p, err := world.NewPlane(nil, tC.x, tC.y, tC.z)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Run("below minimum x", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Min - 1, tC.y.Min, tC.z.Min})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("below minimum y", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Min, tC.y.Min - 1, tC.z.Min})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("below minimum z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Min, tC.y.Min, tC.z.Min - 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("below minimum x, y", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Min - 1, tC.y.Min - 1, tC.z.Min})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("below minimum y, z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Min, tC.y.Min - 1, tC.z.Min - 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("below minimum x, z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Min - 1, tC.y.Min, tC.z.Min - 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("below minimum x, y, z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Min - 1, tC.y.Min - 1, tC.z.Min - 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("above maximum x", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Max + 1, tC.y.Max, tC.z.Max})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("above maximum y", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Max, tC.y.Max + 1, tC.z.Max})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("above maximum z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Max, tC.y.Max, tC.z.Max + 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("above maximum x, y", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Max + 1, tC.y.Max + 1, tC.z.Max})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("above maximum y, z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Max, tC.y.Max + 1, tC.z.Max + 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("above maximum x, z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Max + 1, tC.y.Max, tC.z.Max + 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			t.Run("above maximum x, y, z", func(t *testing.T) {
				t.Parallel()
				_, err := p.At(world.Position{tC.x.Max + 1, tC.y.Max + 1, tC.z.Max + 1})
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if !errors.Is(err, world.ErrOutOfBounds) {
					t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrOutOfBounds, err)
				}
			})

			for x := tC.x.Min; x <= tC.x.Max; x++ {
				for y := tC.y.Min; y <= tC.y.Max; y++ {
					for z := tC.z.Min; z <= tC.z.Max; z++ {

						x, y, z := x, y, z
						t.Run("value within bounds is successfully gotten", func(t *testing.T) {
							t.Parallel()
							pos := world.Position{x, y, z}
							vox, err := p.At(pos)
							if err != nil {
								t.Fatalf("unexpected error: %v", err)
							}
							if !reflect.DeepEqual(vox.Position, pos) {
								t.Fatalf("expected %v but got %v", pos, vox.Position)
							}
						})

					}
				}
			}
		})

	}
}

type testRenderer struct {
	initCalled    bool
	initedPlane   *world.Plane
	renderCalled  bool
	renderedPlane *world.Plane
}

var errInitAlreadyCalled = errors.New("init was already called")

func (tr *testRenderer) Init(plane *world.Plane) error {
	if tr.initCalled {
		return errInitAlreadyCalled
	}
	tr.initCalled = true
	tr.initedPlane = plane
	return nil
}

var errRenderAlreadyCalled = errors.New("render was already called")

func (tr *testRenderer) Render(plane *world.Plane) error {
	if tr.renderCalled {
		return errRenderAlreadyCalled
	}
	tr.renderCalled = true
	tr.renderedPlane = plane
	return nil
}

func (tr *testRenderer) Destroy() {}

func TestNewPlane(t *testing.T) {

	t.Run("stub renderer gets initialized creating plane", func(t *testing.T) {
		t.Parallel()
		var tr testRenderer
		r := world.Range{0, 0}

		_, _ = world.NewPlane(&tr, r, r, r)

		if !tr.initCalled {
			t.Fatal("expected testRenderer.Init to be called")
		}
	})

	t.Run("nil renderer does not panic creating plane", func(t *testing.T) {
		t.Parallel()
		r := world.Range{0, 0}

		_, _ = world.NewPlane(nil, r, r, r)
	})

	t.Run("using testRenderer in multiple locations should error", func(t *testing.T) {
		t.Parallel()
		var tr testRenderer

		r := world.Range{0, 0}
		_, err := world.NewPlane(&tr, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = world.NewPlane(&tr, r, r, r)
		if err == nil {
			t.Fatal("expected error but got none")
		}
		if !errors.Is(err, errInitAlreadyCalled) {
			t.Fatalf("expected \"%v\" but got \"%v\"", errInitAlreadyCalled, err)
		}
	})

}

func TestPlaneRender(t *testing.T) {

	t.Run("stub renderer gets rendered when rendering plane", func(t *testing.T) {
		t.Parallel()
		var tr testRenderer
		r := world.Range{0, 0}
		p, err := world.NewPlane(&tr, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = p.Render()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !tr.renderCalled {
			t.Fatal("expected testRenderer called to be true")
		}
	})

	t.Run("stub renderer gets rendered with the given plane", func(t *testing.T) {
		t.Parallel()
		var tr testRenderer
		r := world.Range{0, 0}
		p, err := world.NewPlane(&tr, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = p.Render()

		if err != nil {
			t.Fatalf("unexpected error calling render on plane: %v", err)
		}
		if !tr.renderCalled {
			t.Fatal("expected testRenderer called to be true")
		}
		if tr.renderedPlane != p {
			t.Fatalf("expected called plane %v to be the given plane %v", tr.renderedPlane, p)
		}
	})

	t.Run("stub renderer returns error when called twice", func(t *testing.T) {
		t.Parallel()
		var tr testRenderer
		r := world.Range{0, 0}
		p, err := world.NewPlane(&tr, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = p.Render()
		if err != nil {
			t.Fatalf("unexpected error calling render on plane: %v", err)
		}

		err = p.Render()

		if err == nil {
			t.Fatal("expected error but got none")
		}
		if !errors.Is(err, errRenderAlreadyCalled) {
			t.Fatalf("expected \"%v\" but got \"%v\"", errRenderAlreadyCalled, err)
		}
	})

	t.Run("nil renderer does not panic calling render", func(t *testing.T) {
		t.Parallel()
		r := world.Range{0, 0}
		p, err := world.NewPlane(nil, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		p.Render()
	})

	t.Run("nil renderer returns nil renderer error", func(t *testing.T) {
		t.Parallel()
		r := world.Range{0, 0}
		p, err := world.NewPlane(nil, r, r, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = p.Render()

		if err == nil {
			t.Fatal("expected error but got none")
		}
		if !errors.Is(err, world.ErrNilRenderer) {
			t.Fatalf("expected \"%v\" but got \"%v\"", world.ErrNilRenderer, err)
		}
	})

}
