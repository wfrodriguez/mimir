package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wfrodriguez/mimir/eeris"
)

// ReadJSON lee un archivo de la ruta `path` y lo asigna en la variable `dst`
func ReadJSON(path string, dst any) error { // {{{
	body, err := os.ReadFile(path)
	if err != nil {
		return eeris.NewWrapError(err, "ReadJSON", fmt.Sprintf("lectura del archivo `%s`", path))
	}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.DisallowUnknownFields()
	err = dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return eeris.NewWrapError(
				syntaxError,
				"ReadJSON",
				fmt.Sprintf("El texto contiene JSON mal formado (en el carácter %d)", syntaxError.Offset),
			)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return eeris.NewWrapError(err, "ReadJSON", "El texto contiene JSON mal formado")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return eeris.NewWrapError(
					unmarshalTypeError,
					"ReadJSON",
					fmt.Sprintf(
						"El texto contiene un tipo JSON incorrecto para el campo `%q`",
						unmarshalTypeError.Field,
					),
				)
			}
			return eeris.NewWrapError(
				unmarshalTypeError,
				"ReadJSON",
				fmt.Sprintf(
					"El texto contiene un tipo JSON incorrecto (en el carácter `%d`)",
					unmarshalTypeError.Offset,
				),
			)
		case errors.Is(err, io.EOF):
			return eeris.NewWrapError(err, "ReadJSON", "El texto no debe estar vacío")
		// Si el JSON contiene un campo que no puede ser mapeado al destino entonces Decode() devolverá ahora un mensaje
		// de error con el formato "json: unknown field "<name>". Comprobamos esto, extraemos el nombre del campo del
		// error y lo interpolamos en nuestro mensaje de error personalizado. Tenga en cuenta que hay un tema abierto en
		// <https://github.com/golang/go/issues/29035> con respecto a convertir esto en un tipo de error distinto en el
		// futuro.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return eeris.NewWrapError(
				err,
				"ReadJSON",
				fmt.Sprintf("El texto contiene una clave desconocida `%s`", fieldName),
			)
		case errors.As(err, &invalidUnmarshalError):
			return eeris.NewWrapError(err, "ReadJSON", "")
		default:
			return eeris.NewWrapError(err, "ReadJSON", "Error general")
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return eeris.NewWrapError(err, "ReadJSON", "El texto sólo debe contener un único valor JSON")
	}

	return nil
} // }}}
