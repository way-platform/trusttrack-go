package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

func otherInputsToProto(input *ttoapi.OtherInputs) *trusttrackv1.OtherInputs {
	var output trusttrackv1.OtherInputs
	if input.CountryCodeGeonames != nil {
		output.SetCountryCodeGeonames(float64(*input.CountryCodeGeonames))
	}
	if input.VirtualGpsOdometer != nil {
		output.SetVirtualGpsOdometerKm(float64(*input.VirtualGpsOdometer))
	}
	return &output
}
