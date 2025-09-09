package mappers

import (
	webdavSchemas "burrowfs/api/schemas/webdav"
	"burrowfs/core/db/models"
	"burrowfs/core/webdav"
)

// FilesToWebDAVResponses converts DB file models to WebDAV XML schema responses
func FilesToWebDAVResponses(files []models.File) []webdavSchemas.FileResponse {
	responses := make([]webdavSchemas.FileResponse, len(files))
	for i, f := range files {
		isCollection := false
		if f.FileID == nil {
			isCollection = true
		}
		href := f.Path
		if isCollection && href[len(href)-1] != '/' {
			href += "/"
		}
		responses[i] = webdavSchemas.FileResponse{
			Href: href,
			PropStat: webdavSchemas.PropStat{
				Prop: webdavSchemas.Prop{
					DisplayName:   f.Name,
					ContentLength: f.Size,
					ContentType:   f.ContentType,
					LastModified:  webdav.FormatModifiedDate(f.UpdatedAt),
					ETag:          f.ETag,
					ResourceType:  isCollection,
					CreationDate:  webdav.FormatCreationDate(f.CreatedAt),
				},
				Status: "HTTP/1.1 200 OK",
			},
		}
	}
	return responses
}
