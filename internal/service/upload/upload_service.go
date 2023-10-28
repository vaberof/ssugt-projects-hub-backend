package upload

type UploadService interface {
	SaveProjectFiles(files []*ProjectFile) error
}

type uploadServiceImpl struct {
}

func NewUploadService() UploadService {
	return &uploadServiceImpl{}
}

func (service *uploadServiceImpl) SaveProjectFiles(files []*ProjectFile) error {
	//TODO implement me
	panic("implement me")
}
