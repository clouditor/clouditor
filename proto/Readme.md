# Protobuf API
This folder contains all protobuf files that we have defined.
We follow the Google Cloud's [API design guide](https://cloud.google.com/apis/design/), including the following conventions:

* Ordering:
  1. _syntax, package, import, option_ statements
  2. overview documentation
  3. _service_ definition(s)
  4. request and response messages (ordered as respective methods)
  5. resource messages while parents are defined before child resources
* File names: 
  * <lower_case_underscore_separated_names>.proto
  * File name corresponding to (single) service
* proto file only containing resources, consider naming this file simply as resources.proto
* ENUM (Scale): The first value should be named ENUM_TYPE_UNSPECIFIED
* Commenting: Comment services, RPCs and messages


# Compiling
The following compile snippets assume being in the current proto folder.

To compile the _assessment_ protobuf file:

`protoc -I ./ -I ../third_party assessment.proto evidence.proto --go_out=../ --go-grpc_out=../  --openapi_out=../openapi/assessment`

To compile the _auth_ protobuf file:

`protoc -I ./ -I ../third_party auth.proto --go_out=../ --go-grpc_out=../..`

To compile the _discovery_ protobuf file:

`protoc -I ./proto -I ./third_party discovery.proto --go_out=. --go-grpc_out=. --openapi_out=./openapi/discovery`

To compile the _orchestrator_ protobuf file:

`protoc -I ./proto -I ./third_party orchestrator.proto --go_out=. --go-grpc_out=. --openapi_out=./openapi/orchestrator`

To compile the _evidenceStore_ protobuf file:

`protoc -I ./proto -I ./third_party evidenceStore.proto --go_out=. --go-grpc_out=. --openapi_out=./openapi/evidenceStore`