package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

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
