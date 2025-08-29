package trusttrack

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

var update = flag.Bool("update", false, "update golden files")

func TestCoordinateToProtoGolden(t *testing.T) {
	testDataDir := filepath.Join("testdata", "coordinates-history-v2")
	// Discover all JSON files (except golden files) in the directory
	var testFiles []string
	err := filepath.Walk(testDataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories and non-JSON files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		// Skip golden files
		if strings.HasSuffix(info.Name(), ".golden.json") {
			return nil
		}
		testFiles = append(testFiles, path)
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk test data directory: %v", err)
	}
	if len(testFiles) == 0 {
		t.Fatal("No test files found in coordinates-history-v2 directory")
	}
	// Process each test file
	for _, testFilePath := range testFiles {
		t.Run(filepath.Base(testFilePath), func(t *testing.T) {
			// Read the test data
			testData, err := os.ReadFile(testFilePath)
			if err != nil {
				t.Fatalf("Failed to read test data from %s: %v", testFilePath, err)
			}
			// Parse the JSON as a coordinate collection
			var collection ttoapi.CoordinateCollection
			err = json.Unmarshal(testData, &collection)
			if err != nil {
				t.Fatalf("Failed to parse test data from %s: %v", testFilePath, err)
			}
			// Convert each coordinate to proto
			var protoCoordinates []*trusttrackv1.Coordinate
			for _, coordinate := range collection.Items {
				protoCoord := coordinateToProto(&coordinate)
				protoCoordinates = append(protoCoordinates, protoCoord)
			}
			// Convert to JSON using protojson for consistent formatting
			marshaler := protojson.MarshalOptions{
				Multiline:       true,
				Indent:          "  ",
				UseProtoNames:   false,
				UseEnumNumbers:  false,
				EmitUnpopulated: false,
			}
			var jsonResults []json.RawMessage
			for _, coord := range protoCoordinates {
				protoJSON, err := marshaler.Marshal(coord)
				if err != nil {
					t.Fatalf("Failed to marshal proto coordinate: %v", err)
				}
				jsonResults = append(jsonResults, json.RawMessage(protoJSON))
			}
			// Marshal the array of coordinates
			result, err := json.MarshalIndent(jsonResults, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal result: %v", err)
			}
			// Generate golden file path by replacing .json with .golden.json
			goldenPath := strings.TrimSuffix(testFilePath, ".json") + ".golden.json"
			if *update {
				// Update the golden file
				err = os.WriteFile(goldenPath, result, 0o644)
				if err != nil {
					t.Fatalf("Failed to update golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenPath)
				return
			}
			// Read existing golden file
			expected, err := os.ReadFile(goldenPath)
			if os.IsNotExist(err) {
				// Golden file doesn't exist, create it
				err = os.WriteFile(goldenPath, result, 0o644)
				if err != nil {
					t.Fatalf("Failed to create golden file: %v", err)
				}
				t.Logf("Created initial golden file: %s", goldenPath)
				t.Log("Run test again to validate against golden file")
				return
			}
			if err != nil {
				t.Fatalf("Failed to read golden file: %v", err)
			}
			// Compare with golden file
			if string(expected) != string(result) {
				t.Errorf("Output differs from golden file. Run with -update flag to update the golden file if the change is expected.")
				t.Logf("Expected length: %d, Got length: %d", len(expected), len(result))
				// Save actual result for debugging
				actualPath := strings.TrimSuffix(testFilePath, ".json") + ".actual.json"
				_ = os.WriteFile(actualPath, result, 0o644)
				t.Logf("Actual result saved to: %s", actualPath)
			}
		})
	}
}
