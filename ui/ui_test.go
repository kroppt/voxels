package ui_test

import (
	"encoding/json"
	"testing"

	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/ui"
)

type GfxStub struct {
}

type GfxVAO struct {
}

func (*GfxStub) NewVAO(mode uint32, layout []int32) *gfx.VAO {
	return &gfx.VAO{}
}

func (*GfxStub) VAOLoad(vao *gfx.VAO, data []float32, usage uint32) error {
	return nil
}

func (*GfxStub) NewShader(source string, shaderType uint32) (gfx.Shader, error) {
	return gfx.Shader{}, nil
}

func (*GfxStub) NewProgram(shaders ...gfx.Shader) (gfx.Program, error) {
	return gfx.Program{}, nil
}

func TestUINew(t *testing.T) {
	t.Run("creates a UI wtih a new buffer object", func(t *testing.T) {
		stub := &GfxStub{}
		ui.New(stub)
		expect := &GfxStub{}
		if *stub != *expect {
			ex, err := json.MarshalIndent(*expect, "", "    ")
			if err != nil {
				t.Fatalf(err.Error())
			}
			st, err := json.MarshalIndent(*stub, "", "    ")
			if err != nil {
				t.Fatalf(err.Error())
			}
			t.Fatalf("expected %s but got %s", string(ex), string(st))
		}
	})
}
