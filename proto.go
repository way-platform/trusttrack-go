package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func objectToProto(input *ttoapi.ExternalComposedObject) *trusttrackv1.Object {
	var output trusttrackv1.Object
	if input.ID != nil {
		output.SetId(*input.ID)
	}
	if input.Name != nil {
		output.SetName(*input.Name)
	}
	if input.Imei != nil {
		output.SetImei(*input.Imei)
	}
	if input.VehicleParams != nil {
		output.SetVehicleParams(vehicleParamsToProto(input.VehicleParams))
	}
	if input.LastCoordinate != nil {
		output.SetLastCoordinate(coordinateToProto(input.LastCoordinate))
	}
	return &output
}

func vehicleParamsToProto(input *ttoapi.ExternalVehicleParams) *trusttrackv1.VehicleParams {
	var output trusttrackv1.VehicleParams
	if input.VIN != nil {
		output.SetVin(*input.VIN)
	}
	if input.Make != nil {
		output.SetMake(*input.Make)
	}
	if input.Model != nil {
		output.SetModel(*input.Model)
	}
	if input.PlateNumber != nil {
		output.SetPlateNumber(*input.PlateNumber)
	}
	if input.AverageFuelConsumption != nil {
		output.SetAverageFuelConsumptionLPer_100Km(float64(*input.AverageFuelConsumption))
	}
	if input.FuelTankCapacity != nil {
		output.SetFuelTankCapacityL(float64(*input.FuelTankCapacity))
	}
	if input.FuelType != nil {
		output.SetFuelType(fuelTypeToProto(*input.FuelType))
		if output.GetFuelType() == trusttrackv1.VehicleParams_FUEL_TYPE_UNKNOWN {
			output.SetUnknownFuelType(string(*input.FuelType))
		}
	}
	return &output
}

func fuelTypeToProto(input ttoapi.ExternalVehicleParamsFuelType) trusttrackv1.VehicleParams_FuelType {
	switch input {
	case ttoapi.ExternalVehicleParamsFuelTypeExternalFuelTypeDIESEL, "DIESEL":
		return trusttrackv1.VehicleParams_DIESEL
	case ttoapi.ExternalVehicleParamsFuelTypeExternalFuelTypeELECTRICITY, "ELECTRICITY":
		return trusttrackv1.VehicleParams_ELECTRICITY
	case ttoapi.ExternalVehicleParamsFuelTypeExternalFuelTypeLPG, "LPG":
		return trusttrackv1.VehicleParams_LPG
	case ttoapi.ExternalVehicleParamsFuelTypeExternalFuelTypeOTHER, "OTHER":
		return trusttrackv1.VehicleParams_OTHER
	case ttoapi.ExternalVehicleParamsFuelTypeExternalFuelTypePETROL, "PETROL":
		return trusttrackv1.VehicleParams_PETROL
	case ttoapi.ExternalVehicleParamsFuelTypeExternalFuelTypeUNKNOWN, "UNKNOWN":
		return trusttrackv1.VehicleParams_FUEL_TYPE_NOT_AVAILABLE
	default:
		return trusttrackv1.VehicleParams_FUEL_TYPE_UNKNOWN
	}
}

func coordinateToProto(input *ttoapi.ExternalLastCoordinate) *trusttrackv1.Coordinate {
	var output trusttrackv1.Coordinate
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
		output.SetSatellitesCount(*input.SatellitesCount)
	}
	if input.ServerDatetime != nil {
		output.SetServerTime(timestamppb.New(*input.ServerDatetime))
	}
	if input.LastValidGpsDatetime != nil {
		output.SetLastValidGpsTime(timestamppb.New(*input.LastValidGpsDatetime))
	}
	return &output
}

func objectGroupToProto(input *ttoapi.ExternalObjectGroup) *trusttrackv1.ObjectGroup {
	var output trusttrackv1.ObjectGroup
	if input.ID != nil {
		output.SetId(*input.ID)
	}
	if input.Name != nil {
		output.SetName(*input.Name)
	}
	if input.ObjectsIds != nil {
		output.SetObjectIds(input.ObjectsIds)
	}
	return &output
}

func tripToProto(input *ttoapi.Trip) *trusttrackv1.Trip {
	var output trusttrackv1.Trip
	if input.ObjectID != nil {
		output.SetObjectId(*input.ObjectID)
	}
	if input.TripType != nil {
		tripType := tripTypeToProto(*input.TripType)
		output.SetType(tripType)
		if tripType == trusttrackv1.Trip_TYPE_UNKNOWN {
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

func addressToProto(input *ttoapi.Address) *trusttrackv1.Address {
	var output trusttrackv1.Address
	if input.Country != nil {
		output.SetCountry(*input.Country)
	}
	if input.CountryCode != nil {
		output.SetCountryCode(*input.CountryCode)
	}
	if input.County != nil {
		output.SetCounty(*input.County)
	}
	if input.HouseNumber != nil {
		output.SetHouseNumber(*input.HouseNumber)
	}
	if input.Locality != nil {
		output.SetLocality(*input.Locality)
	}
	if input.Region != nil {
		output.SetRegion(*input.Region)
	}
	if input.Street != nil {
		output.SetStreet(*input.Street)
	}
	if input.Zip != nil {
		output.SetZip(*input.Zip)
	}
	return &output
}

func tripTypeToProto(input ttoapi.TripTripType) trusttrackv1.Trip_Type {
	switch input {
	case ttoapi.TripTripTypeNONE:
		return trusttrackv1.Trip_NONE
	case ttoapi.TripTripTypePRIVATE:
		return trusttrackv1.Trip_PRIVATE
	case ttoapi.TripTripTypeBUSINESS:
		return trusttrackv1.Trip_BUSINESS
	case ttoapi.TripTripTypeWORK:
		return trusttrackv1.Trip_WORK
	case ttoapi.TripTripTypeUNKNOWN:
		return trusttrackv1.Trip_TYPE_UNKNOWN
	default:
		return trusttrackv1.Trip_TYPE_UNSPECIFIED
	}
}
