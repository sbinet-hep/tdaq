// Code generated by "stringer -type FrameType -output z_frametype_string.go ."; DO NOT EDIT.

package tdaq

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FrameUnknown-0]
	_ = x[FrameCmd-1]
	_ = x[FrameData-2]
	_ = x[FrameOK-3]
	_ = x[FrameErr-4]
}

const _FrameType_name = "FrameUnknownFrameCmdFrameDataFrameOKFrameErr"

var _FrameType_index = [...]uint8{0, 12, 20, 29, 36, 44}

func (i FrameType) String() string {
	if i >= FrameType(len(_FrameType_index)-1) {
		return "FrameType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FrameType_name[_FrameType_index[i]:_FrameType_index[i+1]]
}