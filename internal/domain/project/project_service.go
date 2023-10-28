package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"log"
	"os"
	"path/filepath"
)

const projectUploadsRelativePath = "public\\uploads\\projects"

type ProjectService interface {
	Create(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error)
	Get(id domain.ProjectId) (*Project, error)
	ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error)
}

type projectServiceImpl struct {
	projectStorage ProjectStorage
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
		log.Printf("projectId=%s, err=%v\n", domainProject.Id, err)
	}
	domainProject.Files = domainProjectFiles

	return domainProject, nil
}

func (service *projectServiceImpl) ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error) {
	domainProjects, err := service.projectStorage.ListByFilters(userId, projectType, organizationName, tags)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(domainProjects); i++ {
		domainProjectFiles, err := service.getProjectFiles(domainProjects[i].Id)
		if err != nil {
			//TODO: use logs logger
			log.Printf("projectId=%s, err=%v\n", domainProjects[i].Id, err)
		}
		domainProjects[i].Files = domainProjectFiles
	}

	return domainProjects, nil
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
	projectRootDirectory, _ := os.Getwd()
	pathToUploadsProjectDirectory := projectRootDirectory + "\\" + projectUploadsRelativePath + "\\" + projectId.String()
	_, err := os.Stat(pathToUploadsProjectDirectory)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	files, err := os.ReadDir(pathToUploadsProjectDirectory)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	projectFiles := make([]*ProjectFile, len(files))

	for i := 0; i < len(files); i++ {
		projectFile := ProjectFile{
			ProjectId: projectId,
		}
		fileName := files[i].Name()
		fileExtension := filepath.Ext(fileName)
		//originalFileName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))

		fmt.Printf("fileExtension=%s\n", fileExtension)

		/*if fileExtension == "png" || fileExtension == "jpg" || fileExtension == "jpeg" {
			projectFile.Type = FileTypeImage
		}*/

		projectFile.Name = fileName
		fileContentBytes, _ := os.ReadFile(pathToUploadsProjectDirectory + "\\" + fileName)
		projectFile.Content = fileContentBytes

		projectFiles[i] = &projectFile
	}

	return projectFiles, nil
}
