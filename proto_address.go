package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

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
