# Contributing Guidelines

## Protobuf API

### Style Conventions

This folder contains all protobuf files that we have defined. We follow the Google Cloud's [API design
guide](https://cloud.google.com/apis/design/), including the following conventions:

- Ordering:
  1. _syntax, package, import, option_ statements
  2. overview documentation
  3. _service_ definition(s)
  4. request and response messages (ordered as respective methods)
  5. resource messages while parents are defined before child resources
- File names:
  - <lower_case_underscore_separated_names>.proto
  - File name corresponding to (single) service
- proto file only containing resources, consider naming this file simply as resources.proto
- ENUM (Scale): The first value should be named ENUM_TYPE_UNSPECIFIED
- Commenting: Comment services, RPCs and messages
- Request and response messages

  - A **custom method** should have a response message even if it is empty (see [Cloud APIs Common design patterns](https://cloud.google.com/apis/design/design_patterns#empty_responses)).

  - **Standard methods** use request and response messages according to the following table (see [Cloud APIs Naming
    conventions](https://cloud.google.com/apis/design/naming_convention#method_names)).

| Method name | Request Body      | Response Body         |
| ----------- | ----------------- | --------------------- |
| ListBooks   | ListBooksRequest  | ListResponse          |
| GetBook     | GetBookRequest    | Book                  |
| CreateBook  | CreateBookRequest | Book                  |
| UpdateBook  | UpdateBookRequest | Book                  |
| RenameBook  | RenameBookRequest | RenameBookResponse    |
| DeleteBook  | DeleteBookRequest | google.brotobuf.Empty |

- Even if we _transfer_ a single resource (e.g. in the case of `StoreEvidence`) we create a corresponding `XxxRequest`
  message (`StoreEvidenceRequest`). `XxxRequest` allows adding new fields to the request (e.g. metadata) while not
  breaking the code.
- In case of an error we return nil instead of the `XxxResponse`.

### Generate Go files

In order to generate all necessary Go files, the command `go generate ./...` can be used.

### Writing Assertions for Protobuf Message

While in theory protobuf message are simple structs, they contain un-exported fields that do some internal protobuf
magic. Therefore, we cannot use `reflect.DeepEqual` because it also compares un-exported fields. This can fail even though
the messages are regarded as equal.

Therefore, the usage of our internal `assert` package is mandatory (since it uses go-cmp under the hood) and the
following code is recommended if you want to compare protobuf message and still want to have a proper diff between
messages:

```go
// Simple comparison
assert.Equal(t, want, got)

// Exclusion of certain fields
assert.Equal(t, want, got, protocmp.IgnoreFields(&ontology.Container{}, "creation_time", "raw"))
```

This makes use of the [`cmp`](https://pkg.go.dev/github.com/google/go-cmp/cmp) package of
https://pkg.go.dev/github.com/google/go-cmp as well as a special package
[`protocmp`](https://pkg.go.dev/google.golang.org/protobuf/testing/protocmp) which is part of Go's protobuf
implementation.

For convenience, we wrapped the above construct in our own function
`prototest.Equal` (and `prototest.EqualSlice` for slices) in the
`internal/testutil/prototest` package.