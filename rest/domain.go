package rest

import (
	"strings"

	"github.com/akitasoftware/akita-cli/cfg"
	"github.com/akitasoftware/akita-cli/printer"
)

// This global setting identifies which backend to use, and defaults to akita.software.
//
// If Postman credentials are used, then the domain is chosen based on the
// selected Postman environment (which may be the default or set in an
// environment variable.)
//
// If the --domain flag is used, it unconditionally overrides this choice.
//
// The special values "akitasoftware.com" and "staging.akitasoftware.com"
// need to be prefixed with an "api" to turn them into a host name.  We'll
// assume everything else is supposed to be used as-is.
//
// (The initial goal of the --domain flag was to allow per-customer instances
// of the Akita backend, i.e., myuser.akitasoftware.com and
// api.myuser.akitasoftware.com, but this usage is not supported.)

var Domain string

// Return the default domain, given the settings in use
func DefaultDomain() string {
	// Check if Postman API key in use
	key, env := cfg.GetPostmanAPIKeyAndEnvironment()
	if key == "" {
		printer.Debugf("No Postman API key, using Akita backend.\n")
		return "akita.software"
	}

	// Dispatch based on Postman environment.
	switch strings.ToUpper(env) {
	case "":
		// Not specified by user, default to PRODUCTION
		return "api.observability.postman.com"
	case "DEV":
		printer.Debugf("Selecting localhost backend for DEV environment.\n")
		return "localhost:50443"
	case "BETA":
		printer.Debugf("Selecting Postman beta backend for pre-production testing.\n")
		return "api.observability.postman-beta.com"
	case "PREVIEW":
		printer.Debugf("Selecting Postman preview backend for pre-production testing.\n")
		return "api.observability.postman-preview.com"
	case "STAGE":
		printer.Debugf("Selecting Postman staging backend for pre-production testing.\n")
		return "api.observability.postman-stage.com"
	case "PRODUCTION":
		printer.Debugf("Selecting Postman production backend.\n")
		return "api.observability.postman.com"
	default:
		printer.Warningf("Unknown Postman environment %q, using production.\n")
		return "api.observability.postman.com"
	}
}

// Convert domain to the specific host to contact.
func DomainToHost(domain string) string {
	switch domain {
	case "akita.software":
		return "api.akita.software"
	case "staging.akita.software":
		return "api.staging.akita.software"
	default:
		return domain
	}
}
