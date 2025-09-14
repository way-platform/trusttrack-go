package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func coordinateToProto(input *ttoapi.Coordinate) *trusttrackv1.Coordinate {
	var output trusttrackv1.Coordinate
	if input.ObjectID != nil {
		output.SetObjectId(*input.ObjectID)
	}
	if input.Datetime != nil {
		output.SetVehicleTime(timestamppb.New(*input.Datetime))
	}
	if input.IgnitionStatus != nil {
		ignitionState := coordinateIgnitionStateToProto(*input.IgnitionStatus)
		output.SetIgnitionState(ignitionState)
	}
	if input.TripType != nil {
		tripType := coordinateTripTypeToProto(*input.TripType)
		output.SetTripType(tripType)
	}
	if input.Position != nil {
		output.SetPosition(positionToProto(input.Position))
	}
	if input.GeozoneIds != nil {
		output.SetGeozoneIds(input.GeozoneIds)
	}
	if input.Inputs != nil {
		if input.Inputs.CalculatedInputs != nil {
			output.SetCalculatedInputs(calculatedInputsToProto(input.Inputs.CalculatedInputs))
		}
		if input.Inputs.DeviceInputs != nil {
			output.SetDeviceInputs(deviceInputsToProto(input.Inputs.DeviceInputs))
		}
		if input.Inputs.Other != nil {
			output.SetOther(otherInputsToProto(input.Inputs.Other))
		}
		if input.Inputs.Tires != nil {
			tires := make(map[string]*trusttrackv1.TireData)
			for key, value := range *input.Inputs.Tires {
				// Convert interface{} to TireData through JSON marshaling/unmarshaling
				if valueMap, ok := value.(map[string]interface{}); ok {
					tireData := &ttoapi.TireData{}
					// Map common tire fields manually since we have interface{}
					if v, exists := valueMap["tire_pressure"]; exists {
						if f, ok := v.(float64); ok {
							f32 := float32(f)
							tireData.TirePressure = &f32
						}
					}
					if v, exists := valueMap["tire_temperature"]; exists {
						if f, ok := v.(float64); ok {
							f32 := float32(f)
							tireData.TireTemperature = &f32
						}
					}
					if v, exists := valueMap["tire_location"]; exists {
						if f, ok := v.(float64); ok {
							f32 := float32(f)
							tireData.TireLocation = &f32
						}
					}
					if convertedTireData := tireDataToProto(tireData); convertedTireData != nil {
						tires[key] = convertedTireData
					}
				}
			}
			output.SetTires(tires)
		}
	}
	return &output
}

func coordinateIgnitionStateToProto(input ttoapi.CoordinateIgnitionStatus) trusttrackv1.IgnitionState {
	switch input {
	case "OFF":
		return trusttrackv1.IgnitionState_OFF
	case "ON":
		return trusttrackv1.IgnitionState_ON
	default:
		return trusttrackv1.IgnitionState_IGNITION_STATE_UNKNOWN
	}
}

func coordinateTripTypeToProto(input ttoapi.CoordinateTripType) trusttrackv1.TripType {
	switch input {
	case "NONE":
		return trusttrackv1.TripType_TRIP_TYPE_NONE
	case "PRIVATE":
		return trusttrackv1.TripType_PRIVATE
	case "BUSINESS":
		return trusttrackv1.TripType_BUSINESS
	case "WORK":
		return trusttrackv1.TripType_WORK
	case "UNKNOWN":
		return trusttrackv1.TripType_TRIP_TYPE_NOT_AVAILABLE
	default:
		return trusttrackv1.TripType_TRIP_TYPE_UNKNOWN
	}
}
