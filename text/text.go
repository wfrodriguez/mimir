package text

import (
	"github.com/gosimple/slug"
)

// Slugify genera una cadena slug a partir de una cadena Unicode
func Slugify(s string) string {
	return slug.Make(s)
}

// Slugify genera una cadena slug a partir de una cadena Unicode, URL-amigable con soporte para múltiples idiomas.
// Si no se encuentra el idioma, pasa a "en" por defecto
func SlugifyLang(s string, lang string) string {
	return slug.MakeLang(s, lang)
}
