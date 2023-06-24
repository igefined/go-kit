package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	Empty struct{}

	ValidatorError struct {
		Namespace       string      `json:"namespace,omitempty"`
		Field           string      `json:"field"`
		StructNamespace string      `json:"struct_namespace,omitempty"`
		StructField     string      `json:"struct_field,omitempty"`
		Tag             string      `json:"tag"`
		ActualTag       string      `json:"actual_tag"`
		Kind            interface{} `json:"kind,omitempty"`
		Type            interface{} `json:"type,omitempty"`
		Value           interface{} `json:"value"`
		Param           string      `json:"param"`
	}
)

func InternalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError

	if errors.Is(err, sql.ErrNoRows) || strings.HasSuffix(err.Error(), sql.ErrNoRows.Error()) {
		err = errors.New("entity not found")
		code = http.StatusNotFound
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		code = http.StatusBadRequest
		Respond(w, code, map[string]string{"validator": err.Error()})

		return
	}

	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		res := make(map[string][]ValidatorError)

		for _, err := range validatorErrors {
			res["validator"] = append(res["validator"], ValidatorError{
				Field:     err.Field(),
				Tag:       err.Tag(),
				ActualTag: err.ActualTag(),
				Value:     err.Value(),
				Param:     err.Param(),
			})
		}
		Respond(w, code, res)

		return
	}

	Respond(w, code, map[string]string{"error": err.Error()})
}

func Forbidden(w http.ResponseWriter, err error) {
	var errMsg string

	if err != nil {
		errMsg = err.Error()
	}

	Respond(w, http.StatusForbidden, map[string]string{"error": errMsg})
}

func Conflict(w http.ResponseWriter, err error) {
	var errMsg string

	if err != nil {
		errMsg = err.Error()
	}

	Respond(w, http.StatusConflict, map[string]string{"error": errMsg})
}

func BadRequest(w http.ResponseWriter, err error) {
	Respond(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func Respond(w http.ResponseWriter, code int, data interface{}) {
	if code != 0 {
		w.WriteHeader(code)
	}

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\":\"error decode respond\"}"))
		}
	}
}

func NoContent(w http.ResponseWriter) {
	Respond(w, http.StatusNoContent, nil)
}

func Successfully(w http.ResponseWriter, data interface{}) {
	Respond(w, http.StatusOK, data)
}
