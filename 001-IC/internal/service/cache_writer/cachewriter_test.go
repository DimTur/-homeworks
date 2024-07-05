package cachewriter_test

import (
	"os"
	"testing"

	cachewriter "github.com/DimTur/multi_user_rw_sys/internal/service/cache_writer"
	"github.com/DimTur/multi_user_rw_sys/internal/service/validator"
	"github.com/DimTur/multi_user_rw_sys/models"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
)

func TestWriteMsgs2Cache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := validator.NewMockTokenValidator(ctrl)
	mockValidator.EXPECT().Validate("valid_token").Return(true)
	mockValidator.EXPECT().Validate("invalid_token").Return(false)

	cache := cachewriter.NewMainMsgCache(mockValidator)

	validMsg := models.Message{Token: "valid_token", FileID: "file_1", Data: "Valid message"}
	invalidMsg := models.Message{Token: "invalid_token", FileID: "file_2", Data: "Invalid message"}

	cache.WriteMsgs2Cache(validMsg)
	cache.WriteMsgs2Cache(invalidMsg)

	if len(cache.Cache["file_1"]) != 1 {
		t.Errorf("Expected 1 Msg in cache for file_1, got %d", len(cache.Cache["file_1"]))
	}

	if len(cache.Cache["file_2"]) != 0 {
		t.Errorf("Expected 0 Msgs in cache for file_2, got %d", len(cache.Cache["file_2"]))
	}
}

func TestFlushToFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := validator.NewMockTokenValidator(ctrl)
	cache := cachewriter.NewMainMsgCache(mockValidator)

	mockMessages := []models.Message{
		{Token: "valid_token", FileID: "file_1", Data: "Message 1"},
		{Token: "valid_token", FileID: "file_1", Data: "Message 2"},
		{Token: "invalid_token", FileID: "file_2", Data: "Message 2"},
	}

	mockValidator.EXPECT().Validate("valid_token").Return(true).Times(2)
	mockValidator.EXPECT().Validate("invalid_token").Return(false).Times(1)

	cache.WriteMsgs2Cache(mockMessages[0])
	cache.WriteMsgs2Cache(mockMessages[1])
	cache.WriteMsgs2Cache(mockMessages[2])

	mockValidator.EXPECT().Validate(gomock.Any()).Return(true).AnyTimes()

	// Temporary path for test is set
	tmpPath := "../data"
	os.Mkdir(tmpPath, 0755)
	defer os.RemoveAll(tmpPath)

	cache.FlushToFiles()

	// Verifies that file has been created for file_1
	filePath1 := tmpPath + "/file_1.txt"
	_, err := os.Stat(filePath1)
	if os.IsNotExist(err) {
		t.Errorf("Expected file %s to exists, got error: %v", filePath1, err)
	}

	// Verifies that file has not been created for file_2
	filePath2 := tmpPath + "/file_2.txt"
	_, err = os.Stat(filePath2)
	if !os.IsNotExist(err) {
		t.Errorf("Expected file %s not to exist, but it exists", filePath2)
	}

	// Checks that the file contains messages for file_1
	file, err := os.Open(filePath1)
	if err != nil {
		t.Fatalf("Error opening file %s: %v", filePath1, err)
	}
	defer file.Close()

	buf := make([]byte, 100)
	n, err := file.Read(buf)
	if err != nil {
		t.Fatalf("Error reading file %s: %v", filePath1, err)
	}

	fileData := string(buf[:n])
	expectedData := "Message 1\nMessage 2\n"
	if fileData != expectedData {
		t.Errorf("Expected file data '%s', got '%s'", expectedData, fileData)
	}

	assert.Equal(t, expectedData, fileData)
}
