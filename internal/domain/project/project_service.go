package project

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const projectUploadsRelativePath = "public\\uploads\\projects"

type ProjectService interface {
	Create(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error)
	Get(id domain.ProjectId) (*Project, error)
	ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error)
}

type ProjectConfig struct {
	pathToUploadsDirectory string
}

type projectServiceImpl struct {
	projectStorage ProjectStorage
	config         ProjectConfig
}

func NewProjectService(projectStorage ProjectStorage) ProjectService {
	return &projectServiceImpl{projectStorage: projectStorage}
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
		log.Printf("ERROR!!!")
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
	// TODO: use absolute path from service config
	rootDirectory, _ := os.Getwd()
	projectUploadsDirectory := rootDirectory + "\\" + projectUploadsRelativePath + "\\" + projectId.String()
	_, err := os.Stat(projectUploadsDirectory)
	if err != nil {
		// TODO: use service logs.Logger
		//log.Fatal(err)
		return nil, err
	}

	files, err := os.ReadDir(projectUploadsDirectory)
	if err != nil {
		// TODO: use service logs.Logger
		//log.Fatal(err)
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

			projectFile.Name = fileName

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
