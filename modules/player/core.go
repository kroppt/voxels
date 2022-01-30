package player

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	worldMod     world.Interface
	settingsMod  settings.Interface
	graphicsMod  graphics.Interface
	lastChunkPos chunk.Position
	posAssigned  bool
	position     PositionEvent
	dirAssigned  bool
	direction    DirectionEvent
	firstLoad    bool
}

// chunkRange is the range of chunks between Min and Max.
type chunkRange struct {
	Min chunk.Position
	Max chunk.Position
}

func toVoxelPos(playerPos PositionEvent) chunk.VoxelCoordinate {
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
	return chunk.VoxelCoordinate{
		X: int32(x),
		Y: int32(y),
		Z: int32(z),
	}
}

// forEach executes the given function on every position in the this ChunkRange.
// The return of fn indices whether to stop iterating
func (rng chunkRange) forEach(fn func(pos chunk.Position) bool) {
	for x := rng.Min.X; x <= rng.Max.X; x++ {
		for y := rng.Min.Y; y <= rng.Max.Y; y++ {
			for z := rng.Min.Z; z <= rng.Max.Z; z++ {
				stop := fn(chunk.Position{X: x, Y: y, Z: z})
				if stop {
					return
				}
			}
		}
	}
}

// contains returns whether this ChunkRange contains the given pos.
func (rng chunkRange) contains(pos chunk.Position) bool {
	if pos.X < rng.Min.X || pos.X > rng.Max.X {
		return false
	}
	if pos.Y < rng.Min.Y || pos.Y > rng.Max.Y {
		return false
	}
	if pos.Z < rng.Min.Z || pos.Z > rng.Max.Z {
		return false
	}
	return true
}

func (c *core) updatePosition(posEvent PositionEvent) {
	c.posAssigned = true
	c.position = posEvent
	if c.dirAssigned {
		c.graphicsMod.UpdateView(c.getFrustumCulledChunks(), c.getUpdatedViewMatrix(), c.getUpdatedProjMatrix())
	}
	newChunkPos := chunk.VoxelCoordToChunkCoord(toVoxelPos(posEvent), c.settingsMod.GetChunkSize())
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	old := chunkRange{
		Min: chunk.Position{
			X: c.lastChunkPos.X - renderDistance,
			Y: c.lastChunkPos.Y - renderDistance,
			Z: c.lastChunkPos.Z - renderDistance,
		},
		Max: chunk.Position{
			X: c.lastChunkPos.X + renderDistance,
			Y: c.lastChunkPos.Y + renderDistance,
			Z: c.lastChunkPos.Z + renderDistance,
		},
	}
	new := chunkRange{
		Min: chunk.Position{
			X: newChunkPos.X - renderDistance,
			Y: newChunkPos.Y - renderDistance,
			Z: newChunkPos.Z - renderDistance,
		},
		Max: chunk.Position{
			X: newChunkPos.X + renderDistance,
			Y: newChunkPos.Y + renderDistance,
			Z: newChunkPos.Z + renderDistance,
		},
	}
	new.forEach(func(pos chunk.Position) bool {
		if !old.contains(pos) || c.firstLoad {
			c.worldMod.LoadChunk(pos)
		}
		return false
	})
	old.forEach(func(pos chunk.Position) bool {
		if !new.contains(pos) && !c.firstLoad {
			c.worldMod.UnloadChunk(pos)
		}
		return false
	})
	if c.firstLoad {
		c.firstLoad = false
	}
	c.lastChunkPos = newChunkPos
}

func (c *core) updateDirection(dirEvent DirectionEvent) {
	c.dirAssigned = true
	c.direction = dirEvent
	if c.posAssigned {
		c.graphicsMod.UpdateView(c.getFrustumCulledChunks(), c.getUpdatedViewMatrix(), c.getUpdatedProjMatrix())
	}
}

func (c *core) getUpdatedViewMatrix() mgl.Mat4 {
	if !c.dirAssigned || !c.posAssigned {
		panic("attempted to calc view matrix with unassigned direction or position")
	}
	view := mgl.Ident4()
	cur := c.direction.Rotation.Inverse().Mat4()
	view = view.Mul4(cur)
	pos := mgl.Translate3D(-c.position.X, -c.position.Y, -c.position.Z)
	view = view.Mul4(pos)
	return view
}

func (c *core) getUpdatedProjMatrix() mgl.Mat4 {
	fovRad := mgl.DegToRad(c.settingsMod.GetFOV())
	near := c.settingsMod.GetNear()
	far := c.settingsMod.GetFar()
	width, height := c.settingsMod.GetResolution()
	aspect := float64(width) / float64(height)
	return mgl.Perspective(fovRad, aspect, near, far)
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

func (c *core) isWithinFrustum(cam *camera, pos chunk.Position, chunkSize uint32) bool {
	corner := mgl.Vec3{
		float64(chunkSize) * float64(pos.X),
		float64(chunkSize) * float64(pos.Y),
		float64(chunkSize) * float64(pos.Z),
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
	newChunkPos := chunk.VoxelCoordToChunkCoord(toVoxelPos(c.position), c.settingsMod.GetChunkSize())
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	rng := chunkRange{
		Min: chunk.Position{
			X: newChunkPos.X - renderDistance,
			Y: newChunkPos.Y - renderDistance,
			Z: newChunkPos.Z - renderDistance,
		},
		Max: chunk.Position{
			X: newChunkPos.X + renderDistance,
			Y: newChunkPos.Y + renderDistance,
			Z: newChunkPos.Z + renderDistance,
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

	rng.forEach(func(pos chunk.Position) bool {
		if c.isWithinFrustum(cam, pos, c.settingsMod.GetChunkSize()) {
			viewChunks[pos] = struct{}{}
		}
		return false
	})

	return viewChunks
}
