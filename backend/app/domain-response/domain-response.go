package domainresponse

import objectstorage "backend/app/pkg/object-storage"

type DomainResponse struct {
	objDataConvert objectstorage.BucketObjectDataConverter
}

func New(objDataConvert objectstorage.BucketObjectDataConverter) *DomainResponse {
	return &DomainResponse{
		objDataConvert: objDataConvert,
	}
}
