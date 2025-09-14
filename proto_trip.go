package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func tripToProto(input *ttoapi.Trip) *trusttrackv1.Trip {
	var output trusttrackv1.Trip
	if input.ObjectID != nil {
		output.SetObjectId(*input.ObjectID)
	}
	if input.TripType != nil {
		tripType := tripTypeToProto(*input.TripType)
		output.SetType(tripType)
		if tripType == trusttrackv1.TripType_TRIP_TYPE_UNKNOWN {
			output.SetUnknownType(string(*input.TripType))
		}
	}
	if input.DriverIds != nil {
		output.SetDriverIds(input.DriverIds)
	}
	if input.TripDuration != nil {
		output.SetDurationS(float64(*input.TripDuration))
	}
	if input.Mileage != nil {
		output.SetMileageKm(*input.Mileage)
	}
	if input.TripStart != nil {
		output.SetStart(tripMetricsToProto(input.TripStart))
	}
	if input.TripEnd != nil {
		output.SetEnd(tripMetricsToProto(input.TripEnd))
	}
	return &output
}

func tripMetricsToProto(input *ttoapi.TripMetrics) *trusttrackv1.Trip_Metrics {
	var output trusttrackv1.Trip_Metrics
	if input.Datetime != nil {
		output.SetTime(timestamppb.New(*input.Datetime))
	}
	if input.Latitude != nil {
		output.SetLatitude(*input.Latitude)
	}
	if input.Longitude != nil {
		output.SetLongitude(*input.Longitude)
	}
	if input.Address != nil {
		output.SetAddress(addressToProto(input.Address))
	}
	return &output
}

func tripTypeToProto(input ttoapi.TripTripType) trusttrackv1.TripType {
	switch input {
	case ttoapi.TripTripTypeNONE:
		return trusttrackv1.TripType_TRIP_TYPE_NONE
	case ttoapi.TripTripTypePRIVATE:
		return trusttrackv1.TripType_PRIVATE
	case ttoapi.TripTripTypeBUSINESS:
		return trusttrackv1.TripType_BUSINESS
	case ttoapi.TripTripTypeWORK:
		return trusttrackv1.TripType_WORK
	case ttoapi.TripTripTypeUNKNOWN:
		return trusttrackv1.TripType_TRIP_TYPE_UNKNOWN
	default:
		return trusttrackv1.TripType_TRIP_TYPE_UNSPECIFIED
	}
}
