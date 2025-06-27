package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// ParamMiddleware -
func (rnr *Runner) ParamMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	// Used to verify path parameters are resolved
	pathParams := extractPathParams(hc.Path)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, _ *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
		l = Logger(l, "ParamMiddleware")

		vals := r.URL.Query()

		var queryParamTypes map[string]jsonschema.JSONType
		if hc.MiddlewareConfig.ValidateParamsConfig != nil {
			for _, pathParam := range pathParams {
				pathParamValue := pp.ByName(pathParam)

				if !hc.MiddlewareConfig.ValidateParamsConfig.ExcludePathParamsFromQueryParams {
					l.Info("(core) adding path param >%s< Value >%s<", pathParam, pathParamValue)
					vals[pathParam] = []string{
						pathParamValue,
					}
				}
			}

			err := validateParams(l, vals, hc.MiddlewareConfig.ValidateParamsConfig)
			if err != nil {
				l.Warn("(core) failed to validate params >%#v< >%v<", vals, err)
				return err
			}

			queryParamTypes = hc.MiddlewareConfig.ValidateParamsConfig.queryParamTypes
		}

		qp, err := queryparam.BuildQueryParams(l, vals, queryParamTypes)
		if err != nil {
			l.Warn("(core) failed to build query params >%#v< >%v<", vals, err)
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

func extractPathParams(p string) []string {
	parts := strings.Split(p, "/")
	params := []string{}
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			params = append(params, strings.TrimPrefix(part, ":"))
		}
	}
	return params
}

// validateParams validates any provided parameters
func validateParams(l logger.Logger, q url.Values, paramCfg *ValidateParamsConfig) error {
	if len(q) == 0 {
		return nil
	}

	if paramCfg.Schema.IsEmpty() {
		for k := range q {
			return coreerror.NewParamError(fmt.Sprintf("Parameter >%s< not allowed.", k))
		}
	}

	qJSON, err := paramsToJSON(q, paramCfg.queryParamTypes)
	if err != nil {
		l.Warn("(core) failed to convert params to JSON >%v<", err)
		return err
	}

	result, err := jsonschema.Validate(paramCfg.Schema, qJSON)
	if err != nil {
		l.Warn("(core) failed validate params with schema >%s< >%v<", paramCfg.Schema, err)
		return err
	}

	if !result.Valid() {
		err := coreerror.NewSchemaValidationError(result.Errors())
		l.Warn("(core) failed validate params >%#v<", err)
		return err
	}

	l.Debug("(core) all parameters okay")

	return nil
}

func paramsToJSON(q url.Values, queryParamTypes map[string]jsonschema.JSONType) (string, error) {
	if len(q) == 0 {
		return "", nil
	}

	jsonBuilder := strings.Builder{}
	jsonBuilder.WriteString("{")

	for k, v := range q {
		k := strings.Split(k, ":")[0]
		jsonBuilder.WriteString(`"` + k + `"`)
		jsonBuilder.WriteString(":")

		// TODO: consider comma or pipe separators for OR instead of [] suffix.
		// Currently, query params with the [] suffix are not being validated.
		// Currently, we can express (OR *) (AND *); but with the above suggested
		// change, it would be possible to express ((OR *) [AND])*
		// However, this would also mean that multiple rounds of jsonschema
		// validation for the query param would be needed, as a string type field
		// cannot be an array. So each item in the string array would require one
		// round of validation.
		if strings.HasSuffix(k, "[]") || queryParamTypes[k].IsArray {
			paramKey := strings.ReplaceAll(k, "[]", "")
			arr := "["

			for _, v := range v {
				value, err := parseValue(paramKey, v, queryParamTypes)
				if err != nil {
					return "", err
				}

				arr += value
				arr += ","
			}

			if len(arr) > 1 {
				arr = arr[:len(arr)-1] // remove extra comma
			}
			arr += "]"
			jsonBuilder.WriteString(arr)
		} else if len(v) == 0 {
			jsonBuilder.WriteString(`""`)
		} else {
			value, err := parseValue(k, v[0], queryParamTypes)
			if err != nil {
				return "", err
			}

			jsonBuilder.WriteString(value)
		}

		jsonBuilder.WriteString(",")
	}

	qpJSON := jsonBuilder.String()
	qpJSON = qpJSON[0 : len(qpJSON)-1] // remove extra comma
	qpJSON += "}"
	return qpJSON, nil
}

func parseValue(k string, v string, queryParamTypes map[string]jsonschema.JSONType) (string, error) {
	switch queryParamTypes[k].ElemType {
	case "number":
		i, err := strconv.Atoi(v)
		if err != nil {
			return "", fmt.Errorf("failed to parse number >%s< >%s< >%v<", k, v, err)
		}
		return fmt.Sprintf("%#v", i), nil
	case "boolean":
		b, err := strconv.ParseBool(v)
		if err != nil {
			return "", fmt.Errorf("failed to parse boolean >%s< >%s< >%v<", k, v, err)
		}

		return fmt.Sprintf("%#v", b), nil
	}

	return fmt.Sprintf("%#v", v), nil
}

func (rnr *Runner) resolveHandlerQueryParamsConfig(hc HandlerConfig) (HandlerConfig, error) {

	if hc.MiddlewareConfig.ValidateParamsConfig == nil {
		return hc, nil
	}
	if hc.MiddlewareConfig.ValidateParamsConfig.QueryParams == nil {
		return hc, fmt.Errorf("handler >%s< has ValidateParamsConfig without QueryParams", hc.Name)
	}
	hc.MiddlewareConfig.ValidateParamsConfig.queryParamTypes = jsonschema.CreateJSONTypeMap(hc.MiddlewareConfig.ValidateParamsConfig.QueryParams)

	return hc, nil
}
