package project

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
)

type ProjectService interface {
	Create(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error)
	Get(id domain.ProjectId) (*Project, error)
	ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error)
	SaveFiles(projectId domain.ProjectId, files []*multipart.FileHeader) ([]Filename, error)
	IsFileExists(filename string) error
	GetUploadsPath() string
}

type ProjectServiceConfig struct {
	ProjectUploadsPath string `yaml:"project-uploads-path"`
}

type projectServiceImpl struct {
	projectStorage ProjectStorage
	config         ProjectServiceConfig

	// Usage: appending to output slice of successfully saved files
	mu sync.Mutex
}

func NewProjectService(projectStorage ProjectStorage, config ProjectServiceConfig) ProjectService {
	return &projectServiceImpl{
		projectStorage: projectStorage,
		config:         config,
	}
}

func (service *projectServiceImpl) Create(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error) {
	switch projectType {
	case ProjectTypeStudentScientificConference:
		return service.createSSCProject(userId, projectType, authors, organization, director, projectTemplate, tags)
	case ProjectTypeLaboratory:
		return service.createLaboratoryProject(userId, projectType, authors, organization, director, projectTemplate, tags)
	default:
		return "", errors.New(fmt.Sprintf("unknown project type '%s'", projectType))
	}
}

func (service *projectServiceImpl) Get(id domain.ProjectId) (*Project, error) {
	domainProject, err := service.projectStorage.Get(id)
	if err != nil {
		return nil, err
	}
	domainProjectFiles, err := service.getProjectFiles(id)
	if err != nil {
		//TODO: use logs logger
		log.Printf("projectId=%s, err=%v\n", domainProject.Id, err)
		return nil, err
	}
	domainProject.Files = domainProjectFiles

	fmt.Printf("project files: %+v\n", domainProjectFiles)
	return domainProject, nil
}

func (service *projectServiceImpl) ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error) {
	projects, err := service.projectStorage.ListByFilters(userId, projectType, organizationName, tags)
	if err != nil {
		return nil, err
	}

	type projectToProcessWrapper struct {
		id    domain.ProjectId
		index int
	}

	// The limit of simultaneously goroutines that getting files
	// for listed projects returned from storage`s method ListByFilters()
	workersLimit := 5
	projectsToProcess := make(chan projectToProcessWrapper, workersLimit)

	var wg sync.WaitGroup

	wg.Add(workersLimit)

	for worker := 1; worker <= workersLimit; worker++ {
		go func() {
			defer wg.Done()
			for project := range projectsToProcess {
				domainProjectFiles, err := service.getProjectFiles(project.id)
				if err != nil {
					//TODO: use logs logger
					log.Printf("cant getProjectFiles: project=%s, err=%v\n", projects[project.index].Id, err)
					projects[project.index].Files = nil
				} else {
					projects[project.index].Files = domainProjectFiles
				}
			}
		}()
	}

	for i := 0; i < len(projects); i++ {
		projectsToProcess <- projectToProcessWrapper{projects[i].Id, i}
	}

	close(projectsToProcess)
	wg.Wait()

	return projects, nil
}

func (service *projectServiceImpl) SaveFiles(projectId domain.ProjectId, files []*multipart.FileHeader) ([]Filename, error) {
	_, err := service.Get(projectId)
	if err != nil {
		return nil, err
	}

	pathToProjectsUploads := filepath.Join(service.GetUploadsPath(), projectId.String())

	if _, err := os.Stat(pathToProjectsUploads); err != nil {
		err = os.Mkdir(pathToProjectsUploads, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	var successfullyUploadedFiles []Filename
	var wg sync.WaitGroup

	wg.Add(len(files))

	for _, file := range files {
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			filePath := pathToProjectsUploads + "\\" + file.Filename
			createdFile, err := os.Create(filePath)
			if err != nil {
				// TODO: logs
				return
			}
			defer createdFile.Close()

			openedFile, err := file.Open()
			if err != nil {
				// TODO: logs
				return
			}
			defer openedFile.Close()

			_, err = io.Copy(createdFile, openedFile)
			if err != nil {
				// TODO: logs
				return
			}

			service.mu.Lock()
			successfullyUploadedFiles = append(successfullyUploadedFiles, Filename(file.Filename))
			service.mu.Unlock()
		}(file)
	}

	wg.Wait()

	return successfullyUploadedFiles, nil
}

func (service *projectServiceImpl) IsFileExists(filepath string) error {
	if _, err := os.Stat(filepath); err != nil {
		return err
	}
	return nil
}

func (service *projectServiceImpl) GetUploadsPath() string {
	return service.config.ProjectUploadsPath
}

func (service *projectServiceImpl) createSSCProject(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error) {
	var sscProjectTemplate StudentScientificConferenceProjectTemplate
	err := json.Unmarshal(projectTemplate, &sscProjectTemplate)
	if err != nil {
		return "", err
	}
	return service.projectStorage.CreateSSCProject(userId, projectType, authors, organization, director, &sscProjectTemplate, tags)
}

func (service *projectServiceImpl) createLaboratoryProject(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error) {
	var laboratoryProjectTemplate LaboratoryProjectTemplate
	err := json.Unmarshal(projectTemplate, &laboratoryProjectTemplate)
	if err != nil {
		return "", err
	}
	return service.projectStorage.CreateLaboratoryProject(userId, projectType, authors, organization, director, &laboratoryProjectTemplate, tags)
}

func (service *projectServiceImpl) getProjectFiles(projectId domain.ProjectId) ([]*ProjectFile, error) {
	projectUploadsDirectory := filepath.Join(service.GetUploadsPath(), projectId.String())
	_, err := os.Stat(projectUploadsDirectory)
	if err != nil {
		// TODO: use service logs.Logger
		return nil, err
	}

	files, err := os.ReadDir(projectUploadsDirectory)
	if err != nil {
		// TODO: use service logs.Logger
		return nil, err
	}

	projectFiles := make([]*ProjectFile, len(files))

	var wg sync.WaitGroup

	wg.Add(len(files))

	for i := 0; i < len(files); i++ {
		go func(i int) {
			defer wg.Done()

			projectFile := &ProjectFile{
				ProjectId: projectId,
			}

			fileName := files[i].Name()
			fileExtension := filepath.Ext(fileName)

			projectFile.Name = Filename(fileName)

			if service.isImage(fileExtension) {
				projectFile.Type = FileTypeImage
				fileContentBytes, _ := os.ReadFile(projectUploadsDirectory + "\\" + fileName)
				projectFile.Content = FileContentBase64(service.toBase64(fileContentBytes))
			} else {
				projectFile.Type = FileTypeOther
			}
			projectFiles[i] = projectFile
		}(i)
	}

	wg.Wait()

	return projectFiles, nil
}

func (service *projectServiceImpl) toBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func (service *projectServiceImpl) isImage(fileExtension string) bool {
	return fileExtension == ".png" || fileExtension == ".jpg" || fileExtension == ".jpeg"
}
