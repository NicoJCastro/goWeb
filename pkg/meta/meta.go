package meta

//Meta de request. Manejamos la informacion del body, ej: cantidad de registros, pagina actual, etc.

type Meta struct {
	TotalCount int `json:"total_count"`
}

func New(totalCount int) (*Meta, error) {
	return &Meta{
		TotalCount: totalCount,
	}, nil
}
