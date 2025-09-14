package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func fuelEventToProto(input *ttoapi.ExternalFuelEvent) *trusttrackv1.FuelEvent {
	var output trusttrackv1.FuelEvent
	if input.ObjectID != nil {
		output.SetObjectId(*input.ObjectID)
	}
	if input.DriverID != nil {
		output.SetDriverId(*input.DriverID)
	}
	if input.EventType != nil {
		eventType := fuelEventTypeToProto(*input.EventType)
		output.SetEventType(eventType)
		if eventType == trusttrackv1.FuelEvent_EVENT_TYPE_UNKNOWN {
			output.SetUnknownEventType(string(*input.EventType))
		}
	}
	if input.Latitude != nil {
		output.SetLatitude(*input.Latitude)
	}
	if input.Longitude != nil {
		output.SetLongitude(*input.Longitude)
	}
	if input.FuelLevelStart != nil {
		output.SetFuelLevelStartPercent(float64(*input.FuelLevelStart))
	}
	if input.FuelLevelEnd != nil {
		output.SetFuelLevelEndPercent(float64(*input.FuelLevelEnd))
	}
	if input.Difference != nil {
		output.SetFuelLevelDifferencePercent(float64(*input.Difference))
	}
	if input.StartDate != nil {
		output.SetStartTime(timestamppb.New(*input.StartDate))
	}
	if input.EndDate != nil {
		output.SetEndTime(timestamppb.New(*input.EndDate))
	}
	return &output
}

func fuelEventTypeToProto(input ttoapi.ExternalFuelEventEventType) trusttrackv1.FuelEvent_EventType {
	switch input {
	case "DRAIN":
		return trusttrackv1.FuelEvent_DRAIN
	case "REFUEL":
		return trusttrackv1.FuelEvent_REFUEL
	default:
		return trusttrackv1.FuelEvent_EVENT_TYPE_UNKNOWN
	}
}
