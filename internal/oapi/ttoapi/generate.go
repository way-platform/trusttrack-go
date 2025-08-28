package ttoapi

//go:generate echo [ttoapi] converting from Swagger to OpenAPI...
//go:generate npx swagger2openapi@v7.0.8 01-original.json --outfile 02-converted.yaml --yaml

//go:generate echo [ttoapi] applying overlay...
//go:generate sh -c "go tool -modfile ../../../tools/go.mod openapi-overlay apply overlay.yaml 02-converted.yaml > 03-overlayed.yaml"

//go:generate echo [ttoapi] generating code...
//go:generate go tool -modfile ../../../tools/go.mod oapi-codegen -config config.yaml 03-overlayed.yaml
