package trusttrack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// ListObjectsRequest is the request for the [Client.ListObjects] method.
type ListObjectsRequest struct{}

// ListObjectsResponse is the response for the [Client.ListObjects] method.
type ListObjectsResponse struct {
	Objects []*trusttrackv1.Object `json:"objects"`
}

// ListObjects lists all objects.
func (c *Client) ListObjects(ctx context.Context, request *ListObjectsRequest) (*ListObjectsResponse, error) {
	httpRequest, err := c.newRequest(ctx, "GET", "/objects?version=1", nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.do(ctx, httpRequest)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var responseBody []*ttoapi.Object
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := ListObjectsResponse{
		Objects: make([]*trusttrackv1.Object, 0, len(responseBody)),
	}
	for _, object := range responseBody {
		response.Objects = append(response.Objects, objectToProto(object))
	}
	return &response, nil
}

func objectToProto(input *ttoapi.Object) *trusttrackv1.Object {
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
