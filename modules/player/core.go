package player

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	worldMod     world.Interface
	settingsMod  settings.Interface
	chunkSize    uint32
	lastChunkPos chunkPos
	posAssigned  bool
	position     PositionEvent
	dirAssigned  bool
	direction    DirectionEvent
}

func (c *core) playerToChunkPosition(pos voxelPos) chunkPos {
	x, y, z := pos.x, pos.y, pos.z
	chunkSize := int32(c.chunkSize)
	if pos.x < 0 {
		x++
	}
	if pos.y < 0 {
		y++
	}
	if pos.z < 0 {
		z++
	}
	x /= chunkSize
	y /= chunkSize
	z /= chunkSize
	if pos.x < 0 {
		x--
	}
	if pos.y < 0 {
		y--
	}
	if pos.z < 0 {
		z--
	}
	return chunkPos{x, y, z}
}

type voxelPos struct {
	x int32
	y int32
	z int32
}

type chunkPos struct {
	x int32
	y int32
	z int32
}

// chunkRange is the range of chunks between Min and Max.
type chunkRange struct {
	Min chunkPos
	Max chunkPos
}

func toVoxelPos(playerPos PositionEvent) voxelPos {
	x, y, z := playerPos.X, playerPos.Y, playerPos.Z
	if x < 0 {
		x--
	}
	if y < 0 {
		y--
	}
	if z < 0 {
		z--
	}
	return voxelPos{
		x: int32(x),
		y: int32(y),
		z: int32(z),
	}
}

// forEach executes the given function on every position in the this ChunkRange.
// The return of fn indices whether to stop iterating
func (rng chunkRange) forEach(fn func(pos chunkPos) bool) {
	for x := rng.Min.x; x <= rng.Max.x; x++ {
		for y := rng.Min.y; y <= rng.Max.y; y++ {
			for z := rng.Min.z; z <= rng.Max.z; z++ {
				stop := fn(chunkPos{x: x, y: y, z: z})
				if stop {
					return
				}
			}
		}
	}
}

// contains returns whether this ChunkRange contains the given pos.
func (rng chunkRange) contains(pos chunkPos) bool {
	if pos.x < rng.Min.x || pos.x > rng.Max.x {
		return false
	}
	if pos.y < rng.Min.y || pos.y > rng.Max.y {
		return false
	}
	if pos.z < rng.Min.z || pos.z > rng.Max.z {
		return false
	}
	return true
}

func (c *core) updatePosition(posEvent PositionEvent) {
	c.posAssigned = true
	c.position = posEvent
	if c.dirAssigned {
		c.worldMod.UpdateView(c.getFrustumCulledChunks())
	}
	voxelPos := toVoxelPos(posEvent)
	newChunkPos := c.playerToChunkPosition(voxelPos)
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	old := chunkRange{
		Min: chunkPos{
			x: c.lastChunkPos.x - renderDistance,
			y: c.lastChunkPos.y - renderDistance,
			z: c.lastChunkPos.z - renderDistance,
		},
		Max: chunkPos{
			x: c.lastChunkPos.x + renderDistance,
			y: c.lastChunkPos.y + renderDistance,
			z: c.lastChunkPos.z + renderDistance,
		},
	}
	new := chunkRange{
		Min: chunkPos{
			x: newChunkPos.x - renderDistance,
			y: newChunkPos.y - renderDistance,
			z: newChunkPos.z - renderDistance,
		},
		Max: chunkPos{
			x: newChunkPos.x + renderDistance,
			y: newChunkPos.y + renderDistance,
			z: newChunkPos.z + renderDistance,
		},
	}
	new.forEach(func(pos chunkPos) bool {
		if !old.contains(pos) {
			c.worldMod.LoadChunk(chunk.Position{
				X: pos.x,
				Y: pos.y,
				Z: pos.z,
			})
		}
		return false
	})
	old.forEach(func(pos chunkPos) bool {
		if !new.contains(pos) {
			c.worldMod.UnloadChunk(chunk.Position{
				X: pos.x,
				Y: pos.y,
				Z: pos.z,
			})
		}
		return false
	})
	c.lastChunkPos = newChunkPos
}

func (c *core) updateDirection(dirEvent DirectionEvent) {
	c.dirAssigned = true
	c.direction = dirEvent
	if c.posAssigned {
		c.worldMod.UpdateView(c.getFrustumCulledChunks())
	}
}

type camera struct {
	eye   mgl.Vec3
	dir   mgl.Vec3
	left  mgl.Vec3
	right mgl.Vec3
	up    mgl.Vec3
	down  mgl.Vec3
}

func createCamera(rot mgl.Quat, pos mgl.Vec3) *camera {
	inverse := rot.Inverse()
	return &camera{
		eye:   pos,
		dir:   inverse.Rotate(mgl.Vec3{0.0, 0.0, -1.0}),
		left:  inverse.Rotate(mgl.Vec3{-1.0, 0.0, 0.0}),
		right: inverse.Rotate(mgl.Vec3{1.0, 0.0, 0.0}),
		up:    inverse.Rotate(mgl.Vec3{0.0, 1.0, 0.0}),
		down:  inverse.Rotate(mgl.Vec3{0.0, -1.0, 0.0}),
	}

}

type fRange struct {
	Min   float64
	Max   float64
	delta float64
}

type worldRange struct {
	X fRange
	Y fRange
	Z fRange
}

func (rng worldRange) ForEach(fn func(mgl.Vec3) bool) {
	for x := rng.X.Min; x <= rng.X.Max; x += rng.X.delta {
		for y := rng.Y.Min; y <= rng.Y.Max; y += rng.Y.delta {
			for z := rng.Z.Min; z <= rng.Z.Max; z += rng.Z.delta {
				stop := fn(mgl.Vec3{x, y, z})
				if stop {
					return
				}
			}
		}
	}
}

func approxZero(a float64) bool {
	epsilon := 0.000001
	if a > 0 {
		return a <= epsilon
	}
	return -a <= epsilon
}

func pointOutsidePlane(p, a, b, c mgl.Vec3) bool {
	ap, ab, ac := p.Sub(a), b.Sub(a), c.Sub(a)
	abac := ab.Cross(ac)
	d := ap.Dot(abac)
	return d > 0 || approxZero(d)
}

func (c *core) isWithinFrustum(cam *camera, pos chunkPos, chunkSize uint32) bool {
	corner := mgl.Vec3{
		float64(chunkSize) * float64(pos.x),
		float64(chunkSize) * float64(pos.y),
		float64(chunkSize) * float64(pos.z),
	}
	near := c.settingsMod.GetNear()
	far := c.settingsMod.GetFar()
	fovyDeg := c.settingsMod.GetFOV()
	width, height := c.settingsMod.GetResolution()
	aspect := float64(width) / float64(height)
	// far plane math
	farDist := cam.dir.Mul(far)
	farCenter := cam.eye.Add(farDist)
	fovyRad := mgl.DegToRad(fovyDeg / 2.0)
	fhh := far * math.Tan(fovyRad)
	fhw := aspect * fhh
	farLeftOff := cam.left.Mul(fhw)
	farRightOff := cam.right.Mul(fhw)
	farUpOff := cam.up.Mul(fhh)
	farDownOff := cam.down.Mul(fhh)
	ftl := farCenter.Add(farLeftOff)
	ftl = ftl.Add(farUpOff)
	fbl := farCenter.Add(farLeftOff)
	fbl = fbl.Add(farDownOff)
	ftr := farCenter.Add(farRightOff)
	ftr = ftr.Add(farUpOff)
	fbr := farCenter.Add(farRightOff)
	fbr = fbr.Add(farDownOff)
	// near plane math
	nearDist := cam.dir.Mul(near)
	nearCenter := cam.eye.Add(nearDist)
	nhh := near * math.Tan(fovyRad/2.0)
	nhw := aspect * nhh
	nearLeftOff := cam.left.Mul(nhw)
	nearUpOff := cam.up.Mul(nhh)
	nleft := nearCenter.Add(nearLeftOff)
	nup := nearCenter.Add(nearUpOff)

	planeTriangles := [6][3]mgl.Vec3{
		{cam.eye, ftl, fbl},      // left
		{cam.eye, ftr, ftl},      // top
		{cam.eye, fbr, ftr},      // right
		{cam.eye, fbl, fbr},      // bottom
		{fbl, ftl, ftr},          // far
		{nearCenter, nup, nleft}, // near
	}
	cubeRange := worldRange{
		X: fRange{corner.X(), corner.X() + float64(chunkSize), float64(chunkSize)},
		Y: fRange{corner.Y(), corner.Y() + float64(chunkSize), float64(chunkSize)},
		Z: fRange{corner.Z(), corner.Z() + float64(chunkSize), float64(chunkSize)},
	}
	for _, tri := range planeTriangles {
		in := 0
		cubeRange.ForEach(func(v mgl.Vec3) bool {
			// every corner of cube
			if !pointOutsidePlane(v, tri[0], tri[1], tri[2]) {
				in++
				return true
			}
			return false
		})
		if in == 0 {
			return false
		}
	}

	return true
}

func (c *core) getFrustumCulledChunks() map[chunk.Position]struct{} {
	if !c.dirAssigned || !c.posAssigned {
		panic("position and direction required for frustum culling calculations")
	}
	voxelPos := toVoxelPos(c.position)
	newChunkPos := c.playerToChunkPosition(voxelPos)
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	rng := chunkRange{
		Min: chunkPos{
			x: newChunkPos.x - renderDistance,
			y: newChunkPos.y - renderDistance,
			z: newChunkPos.z - renderDistance,
		},
		Max: chunkPos{
			x: newChunkPos.x + renderDistance,
			y: newChunkPos.y + renderDistance,
			z: newChunkPos.z + renderDistance,
		},
	}
	viewChunks := map[chunk.Position]struct{}{}
	cam := createCamera(mgl.Quat{
		W: c.direction.Rotation.W,
		V: mgl.Vec3{
			c.direction.Rotation.X(),
			c.direction.Rotation.Y(),
			c.direction.Rotation.Z(),
		},
	}, mgl.Vec3{
		c.position.X,
		c.position.Y,
		c.position.Z,
	})

	rng.forEach(func(pos chunkPos) bool {
		if c.isWithinFrustum(cam, pos, c.chunkSize) {
			key := chunk.Position{
				X: pos.x,
				Y: pos.y,
				Z: pos.z,
			}
			viewChunks[key] = struct{}{}
		}
		return false
	})

	return viewChunks
}
