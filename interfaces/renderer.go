package interfaces

import (
	"render2go/core"
	gmMath "render2go/math"

	"github.com/fogleman/gg"
)

// Renderer 渲染器接口
type Renderer interface {
	Clear(r, g, b float64)
	Render(object core.Mobject)
	Present()
	SaveFrame(filename string) error
	GetContext() *gg.Context
	GetCoordinateSystem() *gmMath.CoordinateSystem
	SetAutoSaveProjectName(projectName string)
	SetupCoordinateSystem(objects []core.Mobject)
}
