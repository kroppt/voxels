package world

import (
	"github.com/EngoEngine/math"
	"github.com/engoengine/glm"
)

// Intersect returns whether the given ray intersects the given box and the
// distance if it does.
func Intersect(box AABC, pos, dir glm.Vec3) (dist float32, hit bool) {
	boxPos := box.Pos.AsVec3()
	boxSize := float32(box.Size)
	boxmin := func(d int) float32 {
		return boxPos[d]
	}
	boxmax := func(d int) float32 {
		return boxPos[d] + boxSize
	}

	invx := float32(1.0) / dir[0]
	tx1 := (boxmin(0) - pos[0]) * invx
	tx2 := (boxmax(0) - pos[0]) * invx
	txmin := math.Min(tx1, tx2)
	txmax := math.Max(tx1, tx2)
	min := txmin
	max := txmax

	invy := float32(1.0) / dir[1]
	ty1 := (boxmin(1) - pos[1]) * invy
	ty2 := (boxmax(1) - pos[1]) * invy
	tymin := math.Min(ty1, ty2)
	tymax := math.Max(ty1, ty2)
	min = math.Max(min, tymin)
	max = math.Min(max, tymax)

	invz := float32(1.0) / dir[2]
	tz1 := (boxmin(2) - pos[2]) * invz
	tz2 := (boxmax(2) - pos[2]) * invz
	tzmin := math.Min(tz1, tz2)
	tzmax := math.Max(tz1, tz2)
	min = math.Max(min, tzmin)
	max = math.Min(max, tzmax)

	hit = (max >= min) && max >= 0.0
	dist = min

	return
}

// GetClosest returns the closest voxel in the list to the eye pos
func GetClosest(eye glm.Vec3, positions []*Voxel) (*Voxel, float32) {
	var closestDist float32
	var found bool
	var closestVoxel *Voxel
	for i := 0; i < len(positions); i++ {
		v := positions[i]
		pos := positions[i].Pos.AsVec3()
		diff := pos.Sub(&eye)
		dist := (&diff).Len()
		if !found || dist < closestDist {
			found = true
			closestDist = dist
			closestVoxel = v
		}
	}
	return closestVoxel, closestDist
}
