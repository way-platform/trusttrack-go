package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func externalLastCoordinateToProto(input *ttoapi.ExternalLastCoordinate) *trusttrackv1.Position {
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
	if input.Datetime != nil {
		output.SetTime(timestamppb.New(*input.Datetime))
	}
	if input.SatellitesCount != nil {
		output.SetSatellitesCount(int32(*input.SatellitesCount))
	}
	if input.ServerDatetime != nil {
		output.SetServerTime(timestamppb.New(*input.ServerDatetime))
	}
	if input.LastValidGpsDatetime != nil {
		output.SetLastValidGpsTime(timestamppb.New(*input.LastValidGpsDatetime))
	}
	return &output
}
