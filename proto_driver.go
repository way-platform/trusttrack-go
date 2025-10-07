package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

func driverToProto(input *ttoapi.V2Driver) *trusttrackv1.Driver {
	var output trusttrackv1.Driver
	if input.ID != nil {
		output.SetId(*input.ID)
	}
	if input.FirstName != nil {
		output.SetFirstName(*input.FirstName)
	}
	if input.LastName != nil {
		output.SetLastName(*input.LastName)
	}
	if input.Address != nil {
		output.SetAddress(*input.Address)
	}
	if input.Phone != nil {
		output.SetPhone(*input.Phone)
	}
	if len(input.Identifiers) > 0 {
		identifiers := make([]*trusttrackv1.DriverIdentifier, 0, len(input.Identifiers))
		for _, identifier := range input.Identifiers {
			identifiers = append(identifiers, driverIdentifierToProto(&identifier))
		}
		output.SetIdentifiers(identifiers)
	}
	return &output
}

func driverIdentifierToProto(input *ttoapi.V2ExternalIdentifier) *trusttrackv1.DriverIdentifier {
	var output trusttrackv1.DriverIdentifier
	if input.Identifier != nil {
		output.SetIdentifier(*input.Identifier)
	}
	if input.Type != nil {
		output.SetType(driverIdentifierTypeToProto(*input.Type))
		if output.GetType() == trusttrackv1.DriverIdentifier_IDENTIFIER_TYPE_UNKNOWN {
			output.SetUnknownIdentifierType(*input.Type)
		}
	}
	return &output
}

func driverIdentifierTypeToProto(input string) trusttrackv1.DriverIdentifier_IdentifierType {
	switch input {
	case "DLT":
		return trusttrackv1.DriverIdentifier_DLT
	case "TACHOGRAPH":
		return trusttrackv1.DriverIdentifier_TACHOGRAPH
	case "WIRELESS":
		return trusttrackv1.DriverIdentifier_WIRELESS
	case "IBUTTON":
		return trusttrackv1.DriverIdentifier_IBUTTON
	default:
		return trusttrackv1.DriverIdentifier_IDENTIFIER_TYPE_UNKNOWN
	}
}
