package formbinding

import "fmt"

const (
	require_error        = ` cannot be empty.`
	alpha_dash_error     = ` should contain only alphanumeric, dash ('-') and underscore ('_') characters.`
	alpha_dash_dot_error = ` should contain only alphanumeric, dash ('-'), underscore ('_') and dot ('.') characters.`
	size_error           = ` must be size %s.`
	min_size_error       = ` must contain at least %s characters.`
	max_size_error       = ` must contain at most %s characters.`
	email_error          = ` is not a valid email address.`
	url_error            = `"%s" is not a valid URL.`
	include_error        = ` must contain substring "%s".`
	username_error       = ` can only contain alphanumeric chars ('0-9','a-z','A-Z'), dash ('-'), underscore ('_') and dot ('.'). It cannot begin or end with non-alphanumeric chars, and consecutive non-alphanumeric chars are also forbidden.`
	unknown_error        = `Unknown error:`
)

func formatErr(formatter string, specifiers any) string {
	return fmt.Sprintf(formatter, specifiers)
}
