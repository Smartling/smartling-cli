package files

import (
	"errors"

	sdkerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

func returnError(err error) bool {
	if errors.Is(err, sdkerror.NotAuthorizedError{}) {
		return true
	}

	for {
		smartlingAPIError, isSmartlingAPIError := err.(sdkerror.APIError)
		if isSmartlingAPIError {
			reasons := map[string]struct{}{
				"AUTHENTICATION_ERROR":   {},
				"AUTHORIZATION_ERROR":    {},
				"MAINTENANCE_MODE_ERROR": {},
			}

			_, stopExecution := reasons[smartlingAPIError.Code]
			return stopExecution
		}
		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
}
