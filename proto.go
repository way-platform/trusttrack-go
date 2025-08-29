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
		output.SetLastPosition(externalLastCoordinateToProto(input.LastCoordinate))
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
		output.SetFuelLevelStart(float64(*input.FuelLevelStart))
	}
	if input.FuelLevelEnd != nil {
		output.SetFuelLevelEnd(float64(*input.FuelLevelEnd))
	}
	if input.Difference != nil {
		output.SetDifference(float64(*input.Difference))
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

func coordinateToProto(input *ttoapi.Coordinate) *trusttrackv1.Coordinate {
	var output trusttrackv1.Coordinate
	if input.ObjectID != nil {
		output.SetObjectId(*input.ObjectID)
	}
	if input.Datetime != nil {
		output.SetVehicleTime(timestamppb.New(*input.Datetime))
	}
	if input.IgnitionStatus != nil {
		ignitionState := ignitionStateToProto(*input.IgnitionStatus)
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

func ignitionStateToProto(input ttoapi.CoordinateIgnitionStatus) trusttrackv1.IgnitionState {
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

func calculatedInputsToProto(input *ttoapi.CalculatedInputs) *trusttrackv1.CalculatedInputs {
	var output trusttrackv1.CalculatedInputs
	if input.FuelConsumption != nil {
		output.SetFuelConsumption(float64(*input.FuelConsumption))
	}
	if input.FuelLevel != nil {
		output.SetFuelLevel(float64(*input.FuelLevel))
	}
	if input.Mileage != nil {
		output.SetMileage(float64(*input.Mileage))
	}
	if input.Rpm != nil {
		output.SetRpm(float64(*input.Rpm))
	}
	if input.Temperature != nil {
		output.SetTemperature(float64(*input.Temperature))
	}
	if input.CustomInput1 != nil {
		output.SetCustomInput_1(float64(*input.CustomInput1))
	}
	if input.CustomInput2 != nil {
		output.SetCustomInput_2(float64(*input.CustomInput2))
	}
	if input.CustomInput3 != nil {
		output.SetCustomInput_3(float64(*input.CustomInput3))
	}
	if input.CustomInput4 != nil {
		output.SetCustomInput_4(float64(*input.CustomInput4))
	}
	if input.CustomInput5 != nil {
		output.SetCustomInput_5(float64(*input.CustomInput5))
	}
	if input.CustomInput6 != nil {
		output.SetCustomInput_6(float64(*input.CustomInput6))
	}
	if input.CustomInput7 != nil {
		output.SetCustomInput_7(float64(*input.CustomInput7))
	}
	if input.CustomInput8 != nil {
		output.SetCustomInput_8(float64(*input.CustomInput8))
	}
	if input.Din1WorkingTime != nil {
		output.SetDin1WorkingTime(float64(*input.Din1WorkingTime))
	}
	if input.Din2WorkingTime != nil {
		output.SetDin2WorkingTime(float64(*input.Din2WorkingTime))
	}
	if input.Din3WorkingTime != nil {
		output.SetDin3WorkingTime(float64(*input.Din3WorkingTime))
	}
	if input.Din4WorkingTime != nil {
		output.SetDin4WorkingTime(float64(*input.Din4WorkingTime))
	}
	if input.Weight != nil {
		output.SetWeight(float64(*input.Weight))
	}
	return &output
}

func otherInputsToProto(input *ttoapi.OtherInputs) *trusttrackv1.OtherInputs {
	var output trusttrackv1.OtherInputs
	if input.CountryCodeGeonames != nil {
		output.SetCountryCodeGeonames(float64(*input.CountryCodeGeonames))
	}
	if input.VirtualGpsOdometer != nil {
		output.SetVirtualGpsOdometer(float64(*input.VirtualGpsOdometer))
	}
	return &output
}

func tireDataToProto(input *ttoapi.TireData) *trusttrackv1.TireData {
	if input == nil {
		return nil
	}
	var output trusttrackv1.TireData
	if input.TirePressureThresholdDetection != nil {
		output.SetTirePressureThresholdDetection(float64(*input.TirePressureThresholdDetection))
	}
	if input.TireSensorElectricalFault != nil {
		output.SetTireSensorElectricalFault(float64(*input.TireSensorElectricalFault))
	}
	if input.TireStatus != nil {
		output.SetTireStatus(float64(*input.TireStatus))
	}
	if input.TireTemperature != nil {
		output.SetTireTemperature(float64(*input.TireTemperature))
	}
	if input.TireAirLeakageRate != nil {
		output.SetTireAirLeakageRate(float64(*input.TireAirLeakageRate))
	}
	if input.TirePressure != nil {
		output.SetTirePressure(float64(*input.TirePressure))
	}
	if input.TireSensorEnableStatus != nil {
		output.SetTireSensorEnableStatus(float64(*input.TireSensorEnableStatus))
	}
	if input.TireLocation != nil {
		output.SetTireLocation(float64(*input.TireLocation))
	}
	if input.TireExtendedTirePressureSupport != nil {
		output.SetTireExtendedTirePressureSupport(float64(*input.TireExtendedTirePressureSupport))
	}
	return &output
}

func deviceInputsToProto(input *ttoapi.DeviceInputs) *trusttrackv1.DeviceInputs {
	var output trusttrackv1.DeviceInputs
	// Map the most important telemetry fields to proto
	if input.AnalogInput1 != nil {
		output.SetAnalogInput_1(float64(*input.AnalogInput1))
	}
	if input.AnalogInput2 != nil {
		output.SetAnalogInput_2(float64(*input.AnalogInput2))
	}
	if input.AxleCount != nil {
		output.SetAxleCount(float64(*input.AxleCount))
	}
	if input.BatteryCurrent != nil {
		output.SetBatteryCurrent(float64(*input.BatteryCurrent))
	}
	if input.BatteryVoltage != nil {
		output.SetBatteryVoltage(float64(*input.BatteryVoltage))
	}
	if input.CanbusBrakeSwitch != nil {
		output.SetCanbusBrakeSwitch(string(*input.CanbusBrakeSwitch))
	}
	if input.CanbusClutchSwitch != nil {
		output.SetCanbusClutchSwitch(string(*input.CanbusClutchSwitch))
	}
	if input.CanbusCruiseControlState != nil {
		output.SetCanbusCruiseControlState(string(*input.CanbusCruiseControlState))
	}
	if input.CanbusDistance != nil {
		output.SetCanbusDistance(float64(*input.CanbusDistance))
	}
	if input.CanbusEngineCoolantTemperature != nil {
		output.SetCanbusEngineCoolantTemperature(float64(*input.CanbusEngineCoolantTemperature))
	}
	if input.CanbusFuelRate != nil {
		output.SetCanbusFuelRate(float64(*input.CanbusFuelRate))
	}
	if input.CanbusRequestSupported != nil {
		output.SetCanbusRequestSupported(string(*input.CanbusRequestSupported))
	}
	if input.CanbusDiagnosticsSupported != nil {
		output.SetCanbusDiagnosticsSupported(string(*input.CanbusDiagnosticsSupported))
	}
	if input.CanbusVehicleMotion != nil {
		output.SetCanbusVehicleMotion(string(*input.CanbusVehicleMotion))
	}
	if input.CanbusDriver1Card != nil {
		output.SetCanbusDriver_1Card(string(*input.CanbusDriver1Card))
	}
	if input.CanbusDriver1Time != nil {
		output.SetCanbusDriver_1Time(string(*input.CanbusDriver1Time))
	}
	if input.CanbusDriver1Time != nil {
		output.SetCanbusDriver_1Time(string(*input.CanbusDriver1Time))
	}
	if input.CanbusDriver2Card != nil {
		output.SetCanbusDriver_2Card(string(*input.CanbusDriver2Card))
	}
	if input.CanbusDriver2Card != nil {
		output.SetCanbusDriver_2Card(string(*input.CanbusDriver2Card))
	}
	if input.CanbusDriver2Time != nil {
		output.SetCanbusDriver_2Time(string(*input.CanbusDriver2Time))
	}
	if input.CanbusDriver2Time != nil {
		output.SetCanbusDriver_2Time(string(*input.CanbusDriver2Time))
	}
	if input.DigitalInput1 != nil {
		output.SetDigitalInput_1(*input.DigitalInput1)
	}
	if input.DigitalInput2 != nil {
		output.SetDigitalInput_2(*input.DigitalInput2)
	}
	if input.DigitalInput2 != nil {
		output.SetDigitalInput_2(*input.DigitalInput2)
	}
	if input.DigitalInput3 != nil {
		output.SetDigitalInput_3(*input.DigitalInput3)
	}
	if input.DigitalInput3 != nil {
		output.SetDigitalInput_3(*input.DigitalInput3)
	}
	if input.DigitalInput4 != nil {
		output.SetDigitalInput_4(*input.DigitalInput4)
	}
	if input.DigitalInput4 != nil {
		output.SetDigitalInput_4(*input.DigitalInput4)
	}
	if input.EngineHours != nil {
		output.SetEngineHours(float64(*input.EngineHours))
	}
	if input.EngineRpm != nil {
		output.SetEngineRpm(float64(*input.EngineRpm))
	}
	if input.FirstDriverID != nil {
		output.SetFirstDriverId(*input.FirstDriverID)
	}
	if input.FuelLevelCan != nil {
		output.SetFuelLevelCan(float64(*input.FuelLevelCan))
	}
	if input.FuelUsed != nil {
		output.SetFuelUsed(float64(*input.FuelUsed))
	}
	if input.GpsAltitude != nil {
		output.SetGpsAltitude(float64(*input.GpsAltitude))
	}
	if input.GpsSpeed != nil {
		output.SetGpsSpeed(float64(*input.GpsSpeed))
	}
	if input.GsmSignalStrength != nil {
		output.SetGsmSignalStrength(float64(*input.GsmSignalStrength))
	}
	if input.Hdop != nil {
		output.SetHdop(*input.Hdop)
	}
	if input.Ibutton != nil {
		output.SetIbutton(*input.Ibutton)
	}
	if input.Movement != nil {
		output.SetMovement(string(*input.Movement))
	}
	if input.Panic != nil {
		output.SetPanic(*input.Panic)
	}
	if input.PedalPos != nil {
		output.SetPedalPos(float64(*input.PedalPos))
	}
	if input.PowerSupplyVoltage != nil {
		output.SetPowerSupplyVoltage(float64(*input.PowerSupplyVoltage))
	}
	if input.SecondDriverID != nil {
		output.SetSecondDriverId(*input.SecondDriverID)
	}
	if input.ServiceDist != nil {
		output.SetServiceDist(float64(*input.ServiceDist))
	}
	if input.SpeedTacho != nil {
		output.SetSpeedTacho(float64(*input.SpeedTacho))
	}
	if input.SpeedWheel != nil {
		output.SetSpeedWheel(float64(*input.SpeedWheel))
	}
	if input.VehicleID != nil {
		output.SetVehicleId(*input.VehicleID)
	}
	if input.PcbTemperature != nil {
		output.SetPcbTemperature(float64(*input.PcbTemperature))
	}
	if input.VirtualOdometer != nil {
		output.SetVirtualOdometer(float64(*input.VirtualOdometer))
	}
	if input.InputTrigger != nil {
		output.SetInputTrigger(float64(*input.InputTrigger))
	}
	if input.Priority != nil {
		output.SetPriority(string(*input.Priority))
	}
	if input.Operator != nil {
		output.SetOperator(float64(*input.Operator))
	}
	if input.Din1WorkingTimeDiff != nil {
		output.SetDin1WorkingTimeDiff(float64(*input.Din1WorkingTimeDiff))
	}
	if input.Din2WorkingTimeDiff != nil {
		output.SetDin2WorkingTimeDiff(float64(*input.Din2WorkingTimeDiff))
	}
	if input.Din3WorkingTimeDiff != nil {
		output.SetDin3WorkingTimeDiff(float64(*input.Din3WorkingTimeDiff))
	}
	if input.Din4WorkingTimeDiff != nil {
		output.SetDin4WorkingTimeDiff(float64(*input.Din4WorkingTimeDiff))
	}
	if input.VirtualOdometerDiff != nil {
		output.SetVirtualOdometerDiff(float64(*input.VirtualOdometerDiff))
	}
	if input.EcodriveFuelUsedInHighestGear != nil {
		output.SetEcodriveFuelUsedInHighestGear(float64(*input.EcodriveFuelUsedInHighestGear))
	}
	// TODO: Split into separate functions and parse all fields.
	return &output
}
