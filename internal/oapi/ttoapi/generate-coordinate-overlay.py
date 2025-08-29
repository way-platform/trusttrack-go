#!/usr/bin/env python3

import argparse
import re
import sys

import yaml


def parse_field_from_row(row_parts):
    """Parse a single table row into field information."""
    if len(row_parts) < 4:
        return None

    name = row_parts[0].strip()
    field_type = row_parts[1].strip()
    array_context = row_parts[2].strip() if len(row_parts) > 2 else ""
    description = row_parts[3].strip() if len(row_parts) > 3 else ""
    units = row_parts[4].strip() if len(row_parts) > 4 else ""

    # Skip empty rows or separator rows
    if not name or name == "---" or not field_type or field_type == "---":
        return None

    return {
        "name": name,
        "type": field_type,
        "array_context": array_context,
        "description": description,
        "units": units,
    }


def extract_enum_values(description):
    """Extract enum values from description text."""
    if "Possible values:" not in description:
        return []

    enum_part = description.split("Possible values:")[1].strip()

    # Extract values separated by <br> tags
    br_matches = re.findall(r"<br\s*/?>\s*([^<]*?)(?=\s*<br|$)", enum_part)
    if br_matches:
        enum_values = []
        for match in br_matches:
            value = match.strip()
            if value and not value.isspace():
                enum_values.append(value)
        return enum_values

    # Fallback: look for space-separated uppercase words
    clean_part = re.sub(r"<[^>]+>", " ", enum_part)
    words = clean_part.split()
    enum_values = []
    for word in words:
        word = word.strip(".,;()<>")
        if ((word.isupper() and len(word) > 1) or "_" in word) and word not in [
            "BR",
            "AND",
            "OR",
        ]:
            enum_values.append(word)

    return enum_values


def map_field_to_openapi(field):
    """Convert field information to OpenAPI property definition."""
    field_type = field["type"].lower().strip()
    description = field["description"]
    units = field["units"]

    # Clean up description
    clean_description = re.sub(r"<br\s*/?>", " ", description)
    clean_description = re.sub(r"<[^>]+>", "", clean_description)
    clean_description = re.sub(r"\s+", " ", clean_description).strip()

    # Add units to description if present
    if units:
        clean_description += f" (units: {units})"

    # Map types to OpenAPI - preserve original types by not overriding
    # This function is used to define properties that don't exist in the original spec
    # For fields that do exist, we'll let the original spec take precedence
    if field_type == "number":
        return {"type": "number", "description": clean_description}
    elif field_type == "string":
        return {"type": "string", "description": clean_description}
    elif field_type == "boolean":
        return {"type": "boolean", "description": clean_description}
    elif field_type == "date" or field_type == "datetime":
        return {
            "type": "string",
            "format": "date-time",
            "description": clean_description,
        }
    elif field_type == "array":
        return {
            "type": "array",
            "items": {"type": "string"},
            "description": clean_description,
        }
    elif field_type == "enum":
        enum_values = extract_enum_values(description)
        if enum_values:
            return {
                "type": "string",
                "enum": enum_values,
                "description": clean_description,
            }
        else:
            return {"type": "string", "description": clean_description}
    elif field_type == "bitmap":
        return {"type": "string", "description": clean_description}
    else:
        return {"type": "string", "description": clean_description}


def organize_fields_by_context(fields):
    """Organize fields by their array context."""
    organized = {
        "root": [],
        "other": [],
        "calculated_inputs": [],
        "device_inputs": [],
        "position": [],
        "tires": [],
    }

    for field in fields:
        context = field["array_context"].lower().strip()
        name = field["name"]

        # Skip container fields
        if name in [
            "items",
            "other",
            "calculated_inputs",
            "device_inputs",
            "inputs",
            "tires",
        ]:
            continue

        # Determine the context
        if context == "other":
            organized["other"].append(field)
        elif context == "calculated_inputs":
            organized["calculated_inputs"].append(field)
        elif context == "device_inputs":
            organized["device_inputs"].append(field)
        elif context == "position":
            organized["position"].append(field)
        elif context == "tirexx":
            organized["tires"].append(field)
        elif context == "" or context == "inputs":
            if name not in ["geozone_ids"]:  # These will be handled specially
                organized["root"].append(field)

    return organized


def create_properties_dict(fields):
    """Create OpenAPI properties dictionary from field list."""
    properties = {}
    for field in fields:
        properties[field["name"]] = map_field_to_openapi(field)
    return properties


def generate_component_schemas(organized_fields):
    """Generate separate component schemas for nested types."""
    schemas = {}

    # Position schema
    if organized_fields["position"]:
        schemas["Position"] = {
            "type": "object",
            "description": "Container for all record GPS parameters",
            "properties": create_properties_dict(organized_fields["position"]),
        }

    # CalculatedInputs schema
    if organized_fields["calculated_inputs"]:
        schemas["CalculatedInputs"] = {
            "type": "object",
            "description": "Container for parameters calculated in the system from other parameters according to the configuration",
            "properties": create_properties_dict(organized_fields["calculated_inputs"]),
        }

    # DeviceInputs schema
    if organized_fields["device_inputs"]:
        schemas["DeviceInputs"] = {
            "type": "object",
            "description": "Container for parameters received from hardware",
            "properties": create_properties_dict(organized_fields["device_inputs"]),
        }

    # OtherInputs schema
    if organized_fields["other"]:
        schemas["OtherInputs"] = {
            "type": "object",
            "description": "Container for other system parameters",
            "properties": create_properties_dict(organized_fields["other"]),
        }

    # TireData schema
    if organized_fields["tires"]:
        schemas["TireData"] = {
            "type": "object",
            "description": "TPMS tire data",
            "properties": create_properties_dict(organized_fields["tires"]),
        }

    # CoordinateInputs schema
    inputs_properties = {}
    if "CalculatedInputs" in schemas:
        inputs_properties["calculated_inputs"] = {
            "$ref": "#/components/schemas/CalculatedInputs"
        }
    if "DeviceInputs" in schemas:
        inputs_properties["device_inputs"] = {
            "$ref": "#/components/schemas/DeviceInputs"
        }
    if "OtherInputs" in schemas:
        inputs_properties["other"] = {"$ref": "#/components/schemas/OtherInputs"}
    if "TireData" in schemas:
        inputs_properties["tires"] = {
            "type": "object",
            "description": "Container for TPMS parameters received from device",
            "patternProperties": {
                "^tire(0[1-9]|[12][0-9]|3[0-6])$": {
                    "$ref": "#/components/schemas/TireData"
                }
            },
        }

    if inputs_properties:
        schemas["CoordinateInputs"] = {
            "type": "object",
            "description": "Container for all coordinate input parameters",
            "properties": inputs_properties,
        }

    return schemas


def generate_coordinate_schema(organized_fields):
    """Generate the main Coordinate schema that references component schemas."""
    schema = {"type": "object", "properties": {}}

    # Add well-known root fields with proper types
    schema["properties"]["object_id"] = {
        "type": "string",
        "description": "Object identifier",
    }
    schema["properties"]["datetime"] = {
        "type": "string",
        "format": "date-time",
        "description": "Date and time point of coordinate generated in hardware",
    }
    schema["properties"]["ignition_status"] = {
        "type": "string",
        "enum": ["ON", "OFF", "UNKNOWN"],
        "description": "Indicates if the object's ignition is on",
    }
    schema["properties"]["trip_type"] = {
        "type": "string",
        "enum": ["UNKNOWN", "NONE", "PRIVATE", "BUSINESS", "WORK"],
        "description": "Trip type",
    }

    # Reference Position schema
    if organized_fields["position"]:
        schema["properties"]["position"] = {"$ref": "#/components/schemas/Position"}

    # Add geozone_ids
    schema["properties"]["geozone_ids"] = {
        "type": "array",
        "items": {"type": "string"},
        "description": "Container for all geozones IDs",
    }

    # Reference CoordinateInputs schema
    has_inputs = any(
        organized_fields[key]
        for key in ["calculated_inputs", "device_inputs", "other", "tires"]
    )
    if has_inputs:
        schema["properties"]["inputs"] = {
            "$ref": "#/components/schemas/CoordinateInputs"
        }

    return schema


def generate_overlay(coordinate_schema, component_schemas):
    """Generate the complete overlay structure."""
    # Combine coordinate schema with component schemas
    all_schemas = {"Coordinate": coordinate_schema}
    all_schemas.update(component_schemas)

    return {
        "overlay": "1.0.0",
        "info": {"title": "Overlay for TrustTrack API", "version": "1.0.0"},
        "actions": [
            {
                "target": "$.paths",
                "description": "Skip code generation for requests",
                "remove": True,
            },
            {
                "target": "$.components.responses",
                "description": "Skip code generation for responses",
                "remove": True,
            },
            {
                "target": "$.components.securitySchemes",
                "description": "Skip code generation for security schemes",
                "remove": True,
            },
            {"target": '$..[?(@.format == "uuid")].format', "remove": True},
            {
                "target": '$..[?(@.type=="array")]',
                "description": "Skip pointer to Go slices",
                "update": {"x-go-type-skip-optional-pointer": True},
            },
            {
                "target": '$..[?(@.format == "email")].format',
                "description": "Skip email format",
                "remove": True,
            },
            {
                "target": '$..[?(@.format == "int32")].format',
                "description": "Remove int32 format to avoid generator conflicts",
                "remove": True,
            },
            {"target": "$.components.schemas", "update": all_schemas},
        ],
    }


def parse_markdown_table(md_file):
    """Parse the response parameters table from markdown file."""
    with open(md_file, "r", encoding="utf-8") as f:
        content = f.read()

    lines = content.split("\n")

    # Find the "Response parameters" table
    table_start = None
    for i, line in enumerate(lines):
        if "#### Response parameters" in line:
            # Look for the table header after this section
            for j in range(i, min(i + 10, len(lines))):
                if (
                    lines[j].strip().startswith("| Name ")
                    and "Type" in lines[j]
                    and "Description" in lines[j]
                ):
                    table_start = j
                    break
            break

    if table_start is None:
        print("Error: Could not find the response parameters table", file=sys.stderr)
        return []

    # Parse table rows
    fields = []
    for i in range(table_start + 2, len(lines)):  # Skip header and separator
        line = lines[i].strip()

        # Stop if we reach the end of the table
        if not line.startswith("|") or line == "":
            if not line.startswith("|"):
                break
            continue

        # Skip separator rows
        if "---" in line and all(c in "-| " for c in line):
            continue

        # Parse table row
        parts = [part.strip() for part in line.split("|")[1:-1]]
        field = parse_field_from_row(parts)
        if field:
            fields.append(field)

    return fields


def main():
    parser = argparse.ArgumentParser(
        description="Generate OpenAPI overlay from markdown documentation"
    )
    parser.add_argument(
        "markdown_file", help="Path to markdown file containing the field table"
    )
    parser.add_argument("output_file", help="Path to output OpenAPI overlay YAML file")
    parser.add_argument("--verbose", "-v", action="store_true", help="Verbose output")

    args = parser.parse_args()

    # Parse the markdown table
    if args.verbose:
        print(f"Parsing markdown file: {args.markdown_file}")

    fields = parse_markdown_table(args.markdown_file)
    if not fields:
        print("Error: No fields found in markdown table", file=sys.stderr)
        return 1

    if args.verbose:
        print(f"Found {len(fields)} fields")

    # Organize fields by context
    organized = organize_fields_by_context(fields)

    if args.verbose:
        print("Field organization:")
        for context, field_list in organized.items():
            print(f"  {context}: {len(field_list)} fields")

    # Generate component schemas
    component_schemas = generate_component_schemas(organized)

    # Generate the coordinate schema
    coordinate_schema = generate_coordinate_schema(organized)

    # Generate the complete overlay
    overlay = generate_overlay(coordinate_schema, component_schemas)

    # Write the overlay file
    with open(args.output_file, "w", encoding="utf-8") as f:
        yaml.dump(overlay, f, default_flow_style=False, sort_keys=False, indent=2)

    if args.verbose:
        print(f"Generated OpenAPI overlay: {args.output_file}")
        print("Schema includes:")
        print(f"  - Coordinate properties: {len(coordinate_schema['properties'])}")
        print(f"  - Component schemas: {len(component_schemas)}")
        for schema_name, schema_data in component_schemas.items():
            if "properties" in schema_data:
                print(
                    f"    - {schema_name}: {len(schema_data['properties'])} properties"
                )

    return 0


if __name__ == "__main__":
    sys.exit(main())
