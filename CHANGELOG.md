# Changelog

All notable changes to this project are documented in this file.

## v1.1.0 (2026-08-25)

Brings the Go SDK to parity with the PHP, Python and Node.js 1.1.0 releases.

### Bug Fixes

* `Emails.Send` no longer discards the plain-text body when both `Text` and
  `Html` are supplied. The payload used an `if`/`else` that could only ever emit
  one of the two fields. Both `HTMLContent` and `PlainContent` are now sent, and
  `ContentType` is reported as `HTML` whenever an HTML body is present.
* API error messages are no longer raw JSON. The 400/404/422/429/500 branches
  used `string(respBody)` as the message, so callers saw the entire response
  body instead of the API's own text. SendLayer's
  `{"Errors": [{"Code": ..., "Message": ...}]}` envelope is now parsed, with
  multiple messages joined by `; `, falling back to the singular `Error` key and
  then to a per-status default. The 401 branch did not read the body at all.
* Timeouts are now reported as such. A timed-out request returned the generic
  `"HTTP request failed: ..."` with the underlying `net/http` text; it now
  returns `Request timed out after <duration>`.
* `encoding/json` errors no longer escape the SDK. `Emails.Send`,
  `Webhooks.Create/Get` and `Events.Get` returned the raw `json.Unmarshal`
  error, which is not a SendLayer error type, so callers matching on the SDK's
  types missed it. Decode failures now return `SendLayerError`.
* Successful responses with an empty body (e.g. `204 No Content`) no longer fail
  the caller's decode with `unexpected end of JSON input`.
* Query parameters are now escaped. They were concatenated into the URL raw, so
  any value containing a reserved character — a `MessageID` filter, for
  instance — produced a malformed request.
* CI pinned Go 1.21 while `go.mod` requires 1.23, so the release build could not
  compile the module. CI now uses 1.23.

### Features

* Every error now carries `StatusCode`, `Response` and `Errors`, along with a
  `Codes()` helper for branching on SendLayer's numeric error codes. These
  previously existed only on `SendLayerAPIError`.
* Added `AsError`, which reports whether an error came from this SDK and returns
  its shared fields — useful for reading `StatusCode` or `Codes()` without
  switching on the concrete type. Use `errors.As` when the specific type matters.
* Added `ErrorEntry`, the exported shape of a SendLayer `Errors` entry.
* Added `WithHeaders` for extra request headers. They are applied before
  `Authorization`, so a caller cannot accidentally drop authentication.
* Added `WithHTTPClient` for supplying a custom `*http.Client`.
* Added the exported `Version` constant and `DefaultTimeout`. `Version` is
  reported in the `User-Agent` header, which previously hardcoded `1.0.0`, and a
  test now keeps it in step with `.github/VERSION`.

### Documentation

* Documented the error fields, the full error-type table, the
  `SendLayerAPIError` message prefix, and the 30-second default timeout.
* Added a Configuration section and an explicit HTML-with-plain-text-fallback
  example, in both the README and `examples/send_email`.

### Breaking Changes

These affect code written against v1.0.0. The field *reads* most callers rely on
keep working; only composite literals and 5xx type assertions need attention.

* `SendLayerAPIError` now embeds `SendLayerError` rather than declaring its own
  `Message`, `StatusCode` and `Response` fields. Field *reads* are unchanged
  thanks to promotion (`err.StatusCode` still works); composite literals that
  set those fields directly need updating to
  `&SendLayerAPIError{SendLayerError{...}}`.
* `SendLayerError` gained fields, so positional composite literals
  (`SendLayerError{"msg"}`) no longer compile. Use
  `SendLayerError{Message: "msg"}`.
* `500` maps to `SendLayerInternalServerError`, but other 5xx statuses
  (502, 503, 504) now return `SendLayerAPIError` rather than
  `SendLayerInternalServerError`. This matches the Python and Node.js SDKs.

## v1.0.0 (2026-03-16)

- doc: update README to include the updated request format (f488393)
- update: get event logic (eed5c9d)
- update: webhook creation logic, tests, and example (d0405d4)
- update: Get event endpoint logic (1bf824d)
- update: Email example and test case (90f87c4)
- updated: email sending logic to use the SendEmailRequest struct (190d3bf)
- added: new SendEmailRequest struct (513d167)
