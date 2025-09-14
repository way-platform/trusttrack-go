package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

func calculatedInputsToProto(input *ttoapi.CalculatedInputs) *trusttrackv1.CalculatedInputs {
	var output trusttrackv1.CalculatedInputs
	if input.FuelConsumption != nil {
		output.SetFuelConsumptionLifetimeL(float64(*input.FuelConsumption))
	}
	if input.FuelLevel != nil {
		output.SetFuelLevelPercent(float64(*input.FuelLevel))
	}
	if input.Mileage != nil {
		output.SetOdometerKm(float64(*input.Mileage))
	}
	if input.Rpm != nil {
		output.SetEngineRpm(float64(*input.Rpm))
	}
	if input.Temperature != nil {
		output.SetTemperatureC(float64(*input.Temperature))
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
		output.SetWeightKg(float64(*input.Weight))
	}
	return &output
}
