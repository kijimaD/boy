package window

import "github.com/kijimaD/goboy/pkg/types"

// Window is
type Window interface {
	Render(imageData types.ImageData)
	Run(run func())
	PollKey()
}
