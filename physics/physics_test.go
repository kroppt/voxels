package physics_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/physics"
)

func TestVelocityAsPosition(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		v      physics.Vel
		dt     time.Duration
		expect physics.Pos
	}{
		{
			physics.Vel{glm.Vec3{1, 0, 0}},
			time.Second,
			physics.Pos{glm.Vec3{1, 0, 0}},
		},
		{
			physics.Vel{glm.Vec3{1.5, 1.5, 1.5}},
			time.Minute,
			physics.Pos{glm.Vec3{90, 90, 90}},
		},
		{
			physics.Vel{glm.Vec3{0, -1, 0}},
			time.Second * 5,
			physics.Pos{glm.Vec3{0, -5, 0}},
		},
	}

	for _, tC := range testCases {

		tC := tC
		t.Run(fmt.Sprintf("v=%v dt=%v", tC.v, tC.dt), func(t *testing.T) {
			t.Parallel()
			p := tC.v.AsPosition(tC.dt)
			if p != tC.expect {
				t.Fatalf("expected %v but got %v", tC.expect, p)
			}
		})

	}

}

func TestAccelerationAsVelocity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		a      physics.Acc
		dt     time.Duration
		expect physics.Vel
	}{
		{
			physics.Acc{glm.Vec3{1, 0, 0}},
			time.Second,
			physics.Vel{glm.Vec3{1, 0, 0}},
		},
		{
			physics.Acc{glm.Vec3{1.5, 1.5, 1.5}},
			time.Minute,
			physics.Vel{glm.Vec3{90, 90, 90}},
		},
		{
			physics.Acc{glm.Vec3{0, -1, 0}},
			time.Second * 5,
			physics.Vel{glm.Vec3{0, -5, 0}},
		},
	}

	for _, tC := range testCases {

		tC := tC
		t.Run(fmt.Sprintf("a=%v dt=%v", tC.a, tC.dt), func(t *testing.T) {
			t.Parallel()
			p := tC.a.AsVelocity(tC.dt)
			if p != tC.expect {
				t.Fatalf("expected %v but got %v", tC.expect, p)
			}
		})

	}

}
