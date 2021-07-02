package world

import (
	"context"
	"fmt"
	"sync"
)

type opType int

const (
	opTypeAdd opType = iota
	opTypeDel
)

type loadOp struct {
	op opType
	ch *Chunk
}

type ChunkManager struct {
	worker    *Worker
	loadOps   chan<- loadOp
	cache     *Cache
	cacheLock sync.Mutex
	wgSave    sync.WaitGroup
	gen       Generator
	ctx       context.Context
	quitFn    context.CancelFunc
}

func NewChunkManager(loadOps chan<- loadOp, gen Generator) *ChunkManager {
	cache, err := NewCache("world_meta", "world_data", cacheThreshold)
	if err != nil {
		panic(fmt.Sprintf("failed to make cache: %v", err))
	}
	ctx, quitFn := context.WithCancel(context.Background())
	worker := NewWorker(100)
	return &ChunkManager{
		loadOps:   loadOps,
		cache:     cache,
		cacheLock: sync.Mutex{},
		wgSave:    sync.WaitGroup{},
		gen:       gen,
		worker:    worker,
		ctx:       ctx,
		quitFn:    quitFn,
	}
}

func (cm *ChunkManager) LoadChunks(new ChunkRange, old ChunkRange) {
	cm.worker.AddJob(func() {
		new.ForEachSub(old, func(key ChunkPos) bool {
			// TODO cancellation of job?
			cm.cacheLock.Lock()
			chunk, loaded := cm.cache.Load(key)
			cm.cacheLock.Unlock()
			if !loaded {
				chunk = NewChunk(ChunkSize, key, cm.gen)
			}
			// TODO make loadOps buffered aswell?
			select {
			case <-cm.ctx.Done():
				return true
			case cm.loadOps <- loadOp{
				op: opTypeAdd,
				ch: chunk,
			}:
				return false
			}
		})
	})
}

// func (cm *ChunkManager) SaveChunks(new ChunkRange, old ChunkRange, loadedChunks map[ChunkPos]*LoadedChunk) {
// 	ranges := old.Sub(new)
// 	count := 0
// 	for _, rng := range ranges {
// 		count += rng.Count()
// 	}
// 	cm.wgSave.Add(count)
// 	cm.worker.AddJob(func() {
// 		old.ForEachSub(new, func(key ChunkPos) bool {
// 			// TODO cancellation of job?
// 			cm.cacheLock.Lock()
// 			cm.cache.Save(loadedChunks[key].chunk)
// 			cm.cacheLock.Unlock()
// 			cm.loadOps <- loadOp{
// 				op: opTypeDel,
// 				ch: loadedChunks[key].chunk,
// 			}
// 			cm.wgSave.Done()
// 			return false
// 		})
// 	})
// }

// func (cm *ChunkManager) SaveChunk(ch *Chunk) {
// 	cm.wgSave.Add(1)
// 	cm.worker.AddJob(func() {
// 		cm.cacheLock.Lock()
// 		cm.cache.Save(ch)
// 		cm.cacheLock.Unlock()
// 		cm.loadOps <- loadOp{
// 			op: opTypeDel,
// 			ch: ch,
// 		}
// 		cm.wgSave.Done()
// 	})
// }

func (cm *ChunkManager) Destroy() {
	// cm.wgSave.Wait() // wait for all pending saves
	cm.quitFn()
	cm.worker.Quit()
	cm.cacheLock.Lock()
	cm.cache.Sync()
	cm.cache.Destroy()
	cm.cacheLock.Unlock()
}
