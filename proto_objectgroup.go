package trusttrack

import (
	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

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
