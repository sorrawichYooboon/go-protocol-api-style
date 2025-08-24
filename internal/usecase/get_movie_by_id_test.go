package usecase

import (
	"testing"

	"github.com/sorrawichYooboon/go-protocol-api-style/internal/domain"
	mockRepo "github.com/sorrawichYooboon/go-protocol-api-style/internal/mock"
	"github.com/stretchr/testify/assert"
)

func Test_movieUsecase_getMovieByID(t *testing.T) {
	mockMovieRepo := mockRepo.NewMovieRepository(t)

	clearAllMock := func() {
		mockMovieRepo.ClearAll()
	}

	tests := []struct {
		name           string
		mockServiceReq int64

		wantServiceOrRepoCallWithAndResponse func()
		wantServiceOrRepoCallTimes           map[string]map[string]int
		wantMainServiceError                 error
		wantMainServiceResponse              interface{}
	}{
		{
			name:           "Test should return error when movie repository GetByID returns error",
			mockServiceReq: 1,
			wantServiceOrRepoCallWithAndResponse: func() {
				mockMovieRepo.On("GetByID", int64(1)).Return(nil, assert.AnError)
			},
			wantServiceOrRepoCallTimes: map[string]map[string]int{
				"movieRepository": {
					"GetByID": 1,
				},
			},
			wantMainServiceError:    assert.AnError,
			wantMainServiceResponse: nil,
		},
		{
			name:           "Test should return nil when movie repository GetByID returns nil",
			mockServiceReq: 2,
			wantServiceOrRepoCallWithAndResponse: func() {
				mockMovieRepo.On("GetByID", int64(2)).Return(nil, nil)
			},
			wantServiceOrRepoCallTimes: map[string]map[string]int{
				"movieRepository": {
					"GetByID": 1,
				},
			},
			wantMainServiceError:    nil,
			wantMainServiceResponse: nil,
		},
		{
			name:           "Test should return movie when movie repository GetByID returns movie",
			mockServiceReq: 3,
			wantServiceOrRepoCallWithAndResponse: func() {
				mockMovieRepo.On("GetByID", int64(3)).Return(&domain.Movie{ID: 3, Title: "Test Movie"}, nil)
			},
			wantServiceOrRepoCallTimes: map[string]map[string]int{
				"movieRepository": {
					"GetByID": 1,
				},
			},
			wantMainServiceError:    nil,
			wantMainServiceResponse: &domain.Movie{ID: 3, Title: "Test Movie"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer clearAllMock()

			if test.wantServiceOrRepoCallWithAndResponse != nil {
				test.wantServiceOrRepoCallWithAndResponse()
			}

			movieUsecase := NewMovieUsecase(mockMovieRepo)
			response, err := movieUsecase.GetMovieByID(test.mockServiceReq)

			if test.wantMainServiceError != nil {
				assert.Equal(t, test.wantMainServiceError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.wantMainServiceResponse != nil {
				assert.Equal(t, test.wantMainServiceResponse, response)
			} else {
				assert.Nil(t, response)
			}

			for serviceName, serviceCallTimes := range test.wantServiceOrRepoCallTimes {
				for methodName, times := range serviceCallTimes {
					switch serviceName {
					case "movieRepository":
						mockMovieRepo.AssertNumberOfCalls(t, methodName, times)
					default:
						t.Errorf("service %s or method %s not found", serviceName, methodName)
					}
				}
			}
		})
	}
}
