package models

type Sticker struct {
	CodigoProducto  string `json:"codigo_producto"`
	CodigoUbicacion string `json:"codigo_ubicacion"`
	Sede            string `json:"sede"`
	Area            string `json:"area"`
}
