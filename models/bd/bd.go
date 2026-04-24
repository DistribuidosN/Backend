package bd

// PaginatedImages representa una galería con información de paginación.
type PaginatedImages struct {
	Images      []Image `json:"images" xml:"images"`
	CurrentPage int     `json:"current_page" xml:"currentPage"`
	Limit       int     `json:"limit" xml:"limit"`
	TotalCount  int     `json:"total_count" xml:"totalCount"`
	HasMore     bool    `json:"has_more" xml:"hasMore"`
}

// Image contiene los datos básicos de una imagen.
type Image struct {
	ImageUuid      string `json:"image_uuid" xml:"imageUuid"`
	BatchUuid      string `json:"batch_uuid" xml:"batchUuid"`
	OriginalName   string `json:"original_name" xml:"originalName"`
	ResultPath     string `json:"result_path" xml:"resultPath"`
	Status         string `json:"status" xml:"status"`
	NodeId         string `json:"node_id" xml:"nodeId"`
	ReceptionTime  string `json:"reception_time" xml:"receptionTime"`
	ConversionTime string `json:"conversion_time" xml:"conversionTime"`
}

// BatchWithCover representa un lote con su imagen de portada.
type BatchWithCover struct {
	Batch          Batch  `json:"batch" xml:"batch"`
	CoverImageUrl  string `json:"cover_image_url" xml:"coverImageUrl"`
	CoverImageUuid string `json:"cover_image_uuid" xml:"coverImageUuid"`
}

// Batch representa un lote de imágenes.
type Batch struct {
	BatchUuid   string `json:"batch_uuid" xml:"batchUuid"`
	UserUuid    string `json:"user_uuid" xml:"userUuid"`
	RequestTime string `json:"request_time" xml:"requestTime"`
	Status      string `json:"status" xml:"status"`
}
