package heimdall

import (
	"net/http"
)

// Plugin defines the interface that a Heimdall plugin must have
// plugins can be added to a Heimdall client using the `AddPlugin` method
type Plugin interface {
	OnRequestStart(*http.Request)
	OnRequestEnd(*http.Request, *http.Response)
	OnError(*http.Request, error)
}
