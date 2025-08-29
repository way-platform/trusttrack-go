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
		output.SetAverageFuelConsumption(float64(*input.AverageFuelConsumption))
	}
	if input.FuelTankCapacity != nil {
		output.SetFuelTankCapacity(float64(*input.FuelTankCapacity))
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
		output.SetAltitude(float64(*input.Altitude))
	}
	if input.Speed != nil {
		output.SetSpeed(float64(*input.Speed))
	}
	if input.Direction != nil {
		output.SetDirection(float64(*input.Direction))
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
