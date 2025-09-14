package trusttrack

import (
	"strconv"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

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
		output.SetBatteryCurrentMa(float64(*input.BatteryCurrent))
	}
	if input.BatteryVoltage != nil {
		output.SetBatteryVoltageV(float64(*input.BatteryVoltage))
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
		output.SetCanbusOdometerKm(float64(*input.CanbusDistance))
	}
	if input.CanbusEngineCoolantTemperature != nil {
		output.SetCanbusEngineCoolantTemperatureC(float64(*input.CanbusEngineCoolantTemperature))
	}
	if input.CanbusFuelRate != nil {
		output.SetCanbusFuelRateLPerH(float64(*input.CanbusFuelRate))
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
		output.SetEngineHoursLifetimeH(float64(*input.EngineHours))
	}
	if input.EngineRpm != nil {
		output.SetEngineRpm(float64(*input.EngineRpm))
	}
	if input.FirstDriverID != nil {
		output.SetFirstDriverId(*input.FirstDriverID)
	}
	if input.FuelLevelCan != nil {
		output.SetFuelLevelCanPercent(float64(*input.FuelLevelCan))
	}
	if input.FuelUsed != nil {
		output.SetFuelUsedLifetimeL(float64(*input.FuelUsed))
	}
	if input.GpsAltitude != nil {
		output.SetGpsAltitudeM(float64(*input.GpsAltitude))
	}
	if input.GpsSpeed != nil {
		output.SetGpsSpeedKmh(float64(*input.GpsSpeed))
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
		output.SetPedalPositionPercent(float64(*input.PedalPos))
	}
	if input.PowerSupplyVoltage != nil {
		output.SetPowerSupplyVoltageV(float64(*input.PowerSupplyVoltage))
	}
	if input.SecondDriverID != nil {
		output.SetSecondDriverId(*input.SecondDriverID)
	}
	if input.ServiceDist != nil {
		output.SetServiceDistanceRemainingKm(float64(*input.ServiceDist))
	}
	if input.SpeedTacho != nil {
		output.SetTachoSpeedKmh(float64(*input.SpeedTacho))
	}
	if input.SpeedWheel != nil {
		output.SetWheelSpeedKmh(float64(*input.SpeedWheel))
	}
	if input.VehicleID != nil {
		output.SetVehicleId(*input.VehicleID)
	}
	if input.PcbTemperature != nil {
		output.SetPcbTemperatureC(float64(*input.PcbTemperature))
	}
	if input.VirtualOdometer != nil {
		output.SetVirtualOdometerKm(float64(*input.VirtualOdometer))
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
		output.SetDin1WorkingTimeDiffS(float64(*input.Din1WorkingTimeDiff))
	}
	if input.Din2WorkingTimeDiff != nil {
		output.SetDin2WorkingTimeDiffS(float64(*input.Din2WorkingTimeDiff))
	}
	if input.Din3WorkingTimeDiff != nil {
		output.SetDin3WorkingTimeDiffS(float64(*input.Din3WorkingTimeDiff))
	}
	if input.Din4WorkingTimeDiff != nil {
		output.SetDin4WorkingTimeDiffS(float64(*input.Din4WorkingTimeDiff))
	}
	if input.VirtualOdometerDiff != nil {
		output.SetVirtualOdometerDiff(float64(*input.VirtualOdometerDiff))
	}
	if input.EcodriveFuelUsedInHighestGear != nil {
		output.SetEcodriveFuelUsedInHighestGear(float64(*input.EcodriveFuelUsedInHighestGear))
	}
	if input.CanbusHoursToService != nil {
		// Parse string to float64 for canbus_hours_to_service
		if hoursToService, err := strconv.ParseFloat(*input.CanbusHoursToService, 64); err == nil {
			output.SetCanbusHoursToService(hoursToService)
		}
	}
	// TODO: Split into separate functions and parse all fields.
	return &output
}
