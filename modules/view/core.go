package view

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	graphicsMod  graphics.Interface
	settingsRepo settings.Interface
	viewState    ViewState
	trees        map[chunk.ChunkCoordinate]*Octree
}

func (c *core) updateSelection() {
	c.graphicsMod.UpdateSelection(c.getSelection())
}

func (c *core) getSelection() (chunk.VoxelCoordinate, bool) {
	eye := c.viewState.Pos
	dir := c.viewState.Dir.Rotate(mgl.Vec3{0.0, 0.0, -1.0})
	var found bool
	var lowestDist float64
	var closestVox chunk.VoxelCoordinate
	for _, root := range c.trees {
		// chunks out of viewing frustum cannot be intersected
		// TODO optimization here
		// if _, ok := viewableChunks[chPos]; !ok {
		// 	continue
		// }
		vc, dist, ok := root.FindClosestIntersect(eye, dir)
		if ok && (dist < lowestDist || !found) {
			lowestDist = dist
			closestVox = vc
			found = true
		}
	}
	return closestVox, found
}

func (c *core) updateView(vs ViewState) {
	c.viewState = vs
	viewableChunks := c.getViewableChunks()
	c.graphicsMod.UpdateView(viewableChunks, c.getUpdatedViewMatrix())
}

func (c *core) addTree(cc chunk.ChunkCoordinate, root *Octree) {
	if _, ok := c.trees[cc]; ok {
		panic("unintended use: adding tree but one already exists for this chunk")
	}
	c.trees[cc] = root
}

func (c *core) removeTree(cc chunk.ChunkCoordinate) {
	if _, ok := c.trees[cc]; !ok {
		panic("unintended use: removing tree but none existed for this chunk")
	}
	delete(c.trees, cc)
}

func (c *core) addNode(vc chunk.VoxelCoordinate) {
	cc := chunk.VoxelCoordToChunkCoord(vc, c.settingsRepo.GetChunkSize())
	if _, ok := c.trees[cc]; !ok {
		panic("tree not found for add node, unintended use")
	}
	c.trees[cc] = c.trees[cc].AddLeaf(&vc)
}
func (c *core) removeNode(vc chunk.VoxelCoordinate) {
	cc := chunk.VoxelCoordToChunkCoord(vc, c.settingsRepo.GetChunkSize())
	if _, ok := c.trees[cc]; !ok {
		panic("tree not found for add node, unintended use")
	}
	c.trees[cc], _ = c.trees[cc].Remove(vc)
}

func (c *core) getUpdatedViewMatrix() mgl.Mat4 {
	view := mgl.Ident4()
	cur := c.viewState.Dir.Inverse().Mat4()
	view = view.Mul4(cur)
	pos := mgl.Translate3D(-c.viewState.Pos[0], -c.viewState.Pos[1], -c.viewState.Pos[2])
	view = view.Mul4(pos)
	return view
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

type camera struct {
	eye   mgl.Vec3
	dir   mgl.Vec3
	left  mgl.Vec3
	right mgl.Vec3
	up    mgl.Vec3
	down  mgl.Vec3
}

func createCamera(rot mgl.Quat, pos mgl.Vec3) *camera {
	return &camera{
		eye:   pos,
		dir:   rot.Rotate(mgl.Vec3{0.0, 0.0, -1.0}),
		left:  rot.Rotate(mgl.Vec3{-1.0, 0.0, 0.0}),
		right: rot.Rotate(mgl.Vec3{1.0, 0.0, 0.0}),
		up:    rot.Rotate(mgl.Vec3{0.0, 1.0, 0.0}),
		down:  rot.Rotate(mgl.Vec3{0.0, -1.0, 0.0}),
	}

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

// chunkRange is the range of chunks between Min and Max.
type chunkRange struct {
	Min chunk.ChunkCoordinate
	Max chunk.ChunkCoordinate
}

func toVoxelPos(playerPos mgl.Vec3) chunk.VoxelCoordinate {
	x, y, z := playerPos.X(), playerPos.Y(), playerPos.Z()
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
func (rng chunkRange) forEach(fn func(pos chunk.ChunkCoordinate) bool) {
	for x := rng.Min.X; x <= rng.Max.X; x++ {
		for y := rng.Min.Y; y <= rng.Max.Y; y++ {
			for z := rng.Min.Z; z <= rng.Max.Z; z++ {
				stop := fn(chunk.ChunkCoordinate{X: x, Y: y, Z: z})
				if stop {
					return
				}
			}
		}
	}
}

func (c *core) isWithinFrustum(cam *camera, pos chunk.ChunkCoordinate, chunkSize uint32) bool {
	corner := mgl.Vec3{
		float64(chunkSize) * float64(pos.X),
		float64(chunkSize) * float64(pos.Y),
		float64(chunkSize) * float64(pos.Z),
	}
	near := c.settingsRepo.GetNear()
	far := c.settingsRepo.GetFar()
	fovyDeg := c.settingsRepo.GetFOV()
	width, height := c.settingsRepo.GetResolution()
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

func (c *core) getViewableChunks() map[chunk.ChunkCoordinate]struct{} {
	newChunkPos := chunk.VoxelCoordToChunkCoord(toVoxelPos(c.viewState.Pos), c.settingsRepo.GetChunkSize())
	renderDistance := int32(c.settingsRepo.GetRenderDistance())
	rng := chunkRange{
		Min: chunk.ChunkCoordinate{
			X: newChunkPos.X - renderDistance,
			Y: newChunkPos.Y - renderDistance,
			Z: newChunkPos.Z - renderDistance,
		},
		Max: chunk.ChunkCoordinate{
			X: newChunkPos.X + renderDistance,
			Y: newChunkPos.Y + renderDistance,
			Z: newChunkPos.Z + renderDistance,
		},
	}
	viewChunks := map[chunk.ChunkCoordinate]struct{}{}
	cam := createCamera(c.viewState.Dir, c.viewState.Pos)

	rng.forEach(func(pos chunk.ChunkCoordinate) bool {
		if c.isWithinFrustum(cam, pos, c.settingsRepo.GetChunkSize()) {
			viewChunks[pos] = struct{}{}
		}
		return false
	})

	return viewChunks
}
