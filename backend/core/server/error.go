package server

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// WriteError errs may contain coreerror.Error, coreerror.Errors, or `error`.
// `error` and coreerror.InternalServerError are considered internal errors.
// errs must contain errors of only one status code.
func WriteError(l logger.Logger, w http.ResponseWriter, errs ...error) {
	if len(errs) == 0 {
		// This is a logic error!
		l.Error("(core) no errors passed to WriteError")
		writeSystemError(l, w)
		return
	}

	coreerrs, err := coreerror.ToErrors(errs...)
	if err != nil {
		l.Error("(core) system error >%#v< >%v<", errs, err)
		writeSystemError(l, w)
		return
	}

	for _, e := range coreerrs {
		if e.HttpStatusCode == http.StatusInternalServerError {
			l.Error("(core) system error >%#v<", errs)
			writeSystemError(l, w)
			return
		}
	}

	status := coreerrs[0].HttpStatusCode
	l.Warn("(core) %s", coreerrs[0].Error())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(status)

	l.Info("Write response status >%d<", status)

	if err := json.NewEncoder(w).Encode(coreerrs); err != nil {
		l.Error("(core) failed writing response >%v<", err)
		writeSystemError(l, w)
		return
	}
}

func WriteMalformedError(l logger.Logger, w http.ResponseWriter, err error) {
	e := coreerror.NewMalformedDataError(err.Error())
	l.Warn("(core) malformed data >%v< >%v", err, e)
	WriteError(l, w, e)
}

func WriteNotFoundError(l logger.Logger, w http.ResponseWriter, entity string, id string) {
	e := coreerror.NewNotFoundError(entity, id)
	l.Warn("(core) not found error >%v<", e)

	WriteError(l, w, e)
}

func WriteUnauthorizedError(l logger.Logger, w http.ResponseWriter, err error) {
	e := coreerror.NewUnauthorizedError()
	l.Error("(core) unauthorized error >%v< >%v<", err, e)

	WriteError(l, w, e)
}

func WriteUnavailableError(l logger.Logger, w http.ResponseWriter, err error) {
	e := coreerror.NewUnavailableError()
	l.Error("(core) unavailable error >%v< >%v<", err, e)

	WriteError(l, w, e)
}

func WriteSystemError(l logger.Logger, w http.ResponseWriter, err error) {
	l.Error("(core) system error >%v<", err)

	writeSystemError(l, w)
}

func writeSystemError(l logger.Logger, w http.ResponseWriter) {
	e := coreerror.GetRegistryError(coreerror.ErrorCodeInternal)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(e.HttpStatusCode)

	l.Info("Write response status >%d<", e.HttpStatusCode)

	if err := json.NewEncoder(w).Encode([]coreerror.Error{e}); err != nil {
		l.Error("(core) failed writing response >%v<", e)
	}
}

// WriteXMLErrorResponse responds with an 200 HTTP Status Code. For Service Cloud to retry message delivery, a nack (false) should be sent instead.
func WriteXMLErrorResponse(l logger.Logger, w http.ResponseWriter, s interface{}, e error) {
	l.Debug("writing error response >%+v<", s)

	if e != nil && !coreerror.IsError(e) {
		l.Error("(core) system error >%v<", e)
	}

	w.Header().Set("Content-Type", HeaderContentTypeXML+"; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		l.Error("(core) failed writing response >%v<", err)
		return
	}

	if err := xml.NewEncoder(w).Encode(s); err != nil {
		l.Error("(core) failed encoding response >%v<", err)
	}
}
