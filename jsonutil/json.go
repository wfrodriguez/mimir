package jsonutil

import (
	"encoding/json"
	gerror "errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wfrodriguez/mimir/errors"
)

// ReadJSON lee un archivo de la ruta `path` y lo asigna en la variable `dst`
func ReadJSON(path string, dst any) error { // {{{
	body, err := os.ReadFile(path)
	if err != nil {
		return errors.NewWrapError(err, "ReadJSON", fmt.Sprintf("lectura del archivo `%s`", path))
	}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.DisallowUnknownFields()
	err = dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case gerror.As(err, &syntaxError):
			return errors.NewWrapError(
				syntaxError,
				"ReadJSON",
				fmt.Sprintf("El texto contiene JSON mal formado (en el carácter %d)", syntaxError.Offset),
			)
		case gerror.Is(err, io.ErrUnexpectedEOF):
			return errors.NewWrapError(err, "ReadJSON", "El texto contiene JSON mal formado")
		case gerror.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return errors.NewWrapError(
					unmarshalTypeError,
					"ReadJSON",
					fmt.Sprintf(
						"El texto contiene un tipo JSON incorrecto para el campo `%q`",
						unmarshalTypeError.Field,
					),
				)
			}
			return errors.NewWrapError(
				unmarshalTypeError,
				"ReadJSON",
				fmt.Sprintf(
					"El texto contiene un tipo JSON incorrecto (en el carácter `%d`)",
					unmarshalTypeError.Offset,
				),
			)
		case gerror.Is(err, io.EOF):
			return errors.NewWrapError(err, "ReadJSON", "El texto no debe estar vacío")
		// Si el JSON contiene un campo que no puede ser mapeado al destino entonces Decode() devolverá ahora un mensaje
		// de error con el formato "json: unknown field "<name>". Comprobamos esto, extraemos el nombre del campo del
		// error y lo interpolamos en nuestro mensaje de error personalizado. Tenga en cuenta que hay un tema abierto en
		// <https://github.com/golang/go/issues/29035> con respecto a convertir esto en un tipo de error distinto en el
		// futuro.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return errors.NewWrapError(
				err,
				"ReadJSON",
				fmt.Sprintf("El texto contiene una clave desconocida `%s`", fieldName),
			)
		case gerror.As(err, &invalidUnmarshalError):
			return errors.NewWrapError(err, "ReadJSON", "")
		default:
			return errors.NewWrapError(err, "ReadJSON", "Error general")
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.NewWrapError(err, "ReadJSON", "El texto sólo debe contener un único valor JSON")
	}

	return nil
} // }}}
