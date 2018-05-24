// +build !windows

package draw

import (
	"fmt"
	"image"
)

func moveTo(p image.Point) error {
	return fmt.Errorf("moveTo: TODO for this os")
}
