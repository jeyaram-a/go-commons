package gocommons

import (
	"fmt"
	"os"
	"testing"
)

func cleanup(folder string) {
	os.RemoveAll(folder)
}

func TestFolderExistsSuccess(t *testing.T) {
	folderPath := "/tmp/test_folder"
	defer cleanup(folderPath)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		t.Error("error in setting up folder", err)
		t.Fail()
		return
	}

	exists := FolderExists(folderPath)

	if !exists {
		t.Error("folder exists")
		t.Fail()
	}
}

func TestFolderExistsFailureItIsAFile(t *testing.T) {
	folderPath := "/tmp/test_folder"
	defer cleanup(folderPath)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		t.Error("error in setting up folder", err)
		t.Fail()
		return

	}
	filePath := fmt.Sprintf("%s/a.txt", folderPath)
	if _, err := os.Create(filePath); err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	exists := FolderExists(filePath)
	if exists {
		t.Error(fmt.Sprintf("%s is a file", filePath))
		t.Fail()
	}
}

func TestFolderExistsFailureNonExistentPath(t *testing.T) {
	nonExistentPath := "tmp/some_non_existent_folder"
	exists := FolderExists(nonExistentPath)
	if exists {
		t.Error(fmt.Sprintf("%s is non existent", nonExistentPath))
		t.Fail()
	}
}

func TestFolderIsEmptySuccess(t *testing.T) {
	folderPath := "/tmp/test_folder"
	defer cleanup(folderPath)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		t.Error("error in setting up folder", err)
		t.Fail()
		return
	}

	empty, err := IsFolderEmpty(folderPath)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if !empty {
		t.Error(fmt.Sprintf("%s is empty", folderPath))
		t.Fail()
		return
	}
}

func TestIsFolderEmptyFailureNotEmpty(t *testing.T) {
	folderPath := "/tmp/test_folder"
	defer cleanup(folderPath)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		t.Error("error in setting up folder", err)
		t.Fail()
		return

	}
	filePath := fmt.Sprintf("%s/a.txt", folderPath)
	if _, err := os.Create(filePath); err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	empty, err := IsFolderEmpty(folderPath)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if empty {
		t.Error(fmt.Sprintf("%s is a file", filePath))
		t.Fail()
		return
	}

	empty, err = IsFolderEmpty(filePath)

	if err == nil {
		t.Error("supposed to return a error")
		t.Fail()
		return
	}
}
