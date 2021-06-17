package ui_test

import (
	"encoding/json"
	"testing"

	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/ui"
)

var _ ui.Gfx = (*GfxStub)(nil)

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

func (*GfxStub) VAODraw(vao *gfx.VAO) {
}

func (*GfxStub) NewShader(source string, shaderType uint32) (gfx.Shader, error) {
	return gfx.Shader{}, nil
}

func (*GfxStub) NewProgram(shaders ...gfx.Shader) (gfx.Program, error) {
	return gfx.Program{}, nil
}

func (*GfxStub) ProgramUploadUniform(program *gfx.Program, uniformName string, data ...float32) error {
	return nil
}

func (*GfxStub) ProgramBind(program *gfx.Program) {
}

func (*GfxStub) ProgramUnbind(program *gfx.Program) {
}

func TestNew(t *testing.T) {
	t.Run("creates a UI with a new buffer object", func(t *testing.T) {
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

func TestElements(t *testing.T) {
	t.Run("successfully adds an element", func(t *testing.T) {
		stub := &GfxStub{}
		uiPtr, err := ui.New(stub)
		if err != nil {
			t.Fatal(err.Error())
		}
		bg := ui.NewBackground(stub, 0, 0)
		if bg.GetVAO() == nil {
			t.Fatal("Background VAO was nil")
		}
		err = uiPtr.AddElement(bg)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
