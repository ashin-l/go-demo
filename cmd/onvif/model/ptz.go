/*
 * @Author: lqc
 * @Date: 2021-09-24 16:44:56
 * @Description: ptz model
 */

package model

type Vector2D struct {
	X     float64 `xml:"x,attr"`
	Y     float64 `xml:"y,attr"`
	Space string  `xml:"space,attr"`
}

type Vector1D struct {
	X     float64 `xml:"x,attr"`
	Space string  `xml:"space,attr"`
}

type PTZVector struct {
	PanTilt Vector2D `xml:"PanTilt"`
	Zoom    Vector1D `xml:"Zoom"`
}

type MoveStatus struct {
	Status string
}

type PTZMoveStatus struct {
	PanTilt MoveStatus
	Zoom    MoveStatus
}

type PTZStatus struct {
	Position   PTZVector
	MoveStatus PTZMoveStatus
	Error      string
	UtcTime    string
}

type GetStatusResponse struct {
	PTZStatus PTZStatus
}
