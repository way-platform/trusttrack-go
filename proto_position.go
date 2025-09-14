package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

func positionToProto(input *ttoapi.Position) *trusttrackv1.Position {
	var output trusttrackv1.Position
	if input.Latitude != nil {
		output.SetLatitude(*input.Latitude)
	}
	if input.Longitude != nil {
		output.SetLongitude(*input.Longitude)
	}
	if input.Altitude != nil {
		output.SetAltitudeM(float64(*input.Altitude))
	}
	if input.Speed != nil {
		output.SetSpeedKmh(float64(*input.Speed))
	}
	if input.Direction != nil {
		output.SetDirectionDeg(float64(*input.Direction))
	}
	if input.SatellitesCount != nil {
		output.SetSatellitesCount(int32(*input.SatellitesCount))
	}
	return &output
}
