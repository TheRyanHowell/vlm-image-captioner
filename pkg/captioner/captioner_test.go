package captioner

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockClient is a mock implementation of the Client interface for testing purposes.
type mockClient struct {
	resp openai.ChatCompletionResponse
	err  error
}

// CreateChatCompletion is the mock implementation of the CreateChatCompletion method.
func (m *mockClient) CreateChatCompletion(_ context.Context, _ openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return m.resp, m.err
}

// newTestCaptioner creates a new OpenAICaptioner with a mock client for testing.
func newTestCaptioner(resp openai.ChatCompletionResponse, err error) *OpenAICaptioner {
	return &OpenAICaptioner{
		client: &mockClient{
			resp: resp,
			err:  err,
		},
		model: "test-model",
	}
}

func TestOpenAICaptioner_Caption(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-image-*.txt")
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("failed to remove temp file: %v", err)
		}
	})

	_, err = tempFile.Write([]byte("test-image-content"))
	require.NoError(t, err)

	successResp := openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Content: "A test caption for the file",
				},
			},
		},
	}

	tests := []struct {
		name        string
		filePath    string
		mockResp    openai.ChatCompletionResponse
		mockErr     error
		expectedCap string
		checkErr    func(*testing.T, error)
	}{
		{
			name:        "successful caption",
			filePath:    tempFile.Name(),
			mockResp:    successResp,
			expectedCap: "A test caption for the file",
			checkErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:     "file not found",
			filePath: "non-existent-file.txt",
			checkErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, os.ErrNotExist)
			},
		},
		{
			name:     "api error",
			filePath: tempFile.Name(),
			mockErr:  errors.New("api error"),
			checkErr: func(t *testing.T, err error) {
				assert.EqualError(t, err, "failed to create chat completion: api error")
			},
		},
		{
			name:     "no choices",
			filePath: tempFile.Name(),
			mockResp: openai.ChatCompletionResponse{},
			checkErr: func(t *testing.T, err error) {
				assert.EqualError(t, err, "no choices returned from API")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCaptioner(tt.mockResp, tt.mockErr)
			caption, err := c.Caption(context.Background(), tt.filePath)

			tt.checkErr(t, err)
			if err == nil {
				assert.Equal(t, tt.expectedCap, caption)
			}
		})
	}
}
