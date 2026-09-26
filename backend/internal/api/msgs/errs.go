package msgs

import "errors"

const (
	Internal    = "internal server error"
	InvalidBody = "invalid request body"
	Unavailable = "service unavailable"

	ImageNotFound = "image not found"

	InvalidWidthFmt  = "invalid width %d: must be in [%d, %d]"
	InvalidHeightFmt = "invalid height %d: must be in [%d, %d]"
	InvalidScaleFmt  = "invalid scale %g: must be in [%g, %g]"

	InvalidNameLenFmt = "invalid name length %d: must be in [%d, %d]"
)

var (
	ErrInvalidStyle      = errors.New("invalid style")
	ErrInvalidPalette    = errors.New("invalid palette")
	ErrInvalidAdditional = errors.New("invalid additional")
)
