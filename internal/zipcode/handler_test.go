package zipcode

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipefrizzo/brazilian-zipcode-api/internal/address"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/address/mocks"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/mock/gomock"
)

func TestHandler_FetchAddressByZipcode(t *testing.T) {
	tests := []struct {
		name             string
		zipcode          string
		mockAddressRepo  func(ctrl *gomock.Controller) address.AddressRepository
		expectedStatus   int
		expectedResponse *address.Address
		expectedError    string
	}{
		{
			name:    "success",
			zipcode: "12345678",
			mockAddressRepo: func(ctrl *gomock.Controller) address.AddressRepository {
				mockRepo := mocks.NewMockAddressRepository(ctrl)
				mockRepo.EXPECT().
					Get(gomock.Any(), "12345678").
					Return(&address.Address{
						FederativeUnit: "SP",
						City:           "São Paulo",
						Neighborhood:   "Centro",
						AddressName:    "Avenida Paulista",
						Complement:     "",
						Zipcode:        "12345678",
					}, nil)
				return mockRepo
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &address.Address{
				FederativeUnit: "SP",
				City:           "São Paulo",
				Neighborhood:   "Centro",
				AddressName:    "Avenida Paulista",
				Complement:     "",
				Zipcode:        "12345678",
			},
		},
		{
			name:    "missing zipcode",
			zipcode: "",
			mockAddressRepo: func(ctrl *gomock.Controller) address.AddressRepository {
				return mocks.NewMockAddressRepository(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "missing zipcode\n",
		},
		{
			name:    "address not found",
			zipcode: "12345678",
			mockAddressRepo: func(ctrl *gomock.Controller) address.AddressRepository {
				mockRepo := mocks.NewMockAddressRepository(ctrl)
				mockRepo.EXPECT().
					Get(gomock.Any(), "12345678").
					Return(nil, errors.New("address not found"))
				return mockRepo
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "address not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := tt.mockAddressRepo(ctrl)
			handler := New(mockRepo)

			var req *http.Request
			if tt.zipcode == "" {
				req = httptest.NewRequest(http.MethodGet, "/zipcode", nil)
			} else {
				req = httptest.NewRequest(http.MethodGet, "/zipcode/"+tt.zipcode, nil)
			}
			rr := httptest.NewRecorder()

			router := httprouter.New()
			if tt.zipcode == "" {
				router.GET("/zipcode", handler.FetchAddressByZipcode())
			} else {
				router.GET("/zipcode/:zipcode", handler.FetchAddressByZipcode())
			}

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedError != "" {
				if rr.Body.String() != tt.expectedError {
					t.Errorf("handler returned unexpected error: got %v want %v", rr.Body.String(), tt.expectedError)
				}
			} else {
				var gotAddress address.Address
				if err := json.Unmarshal(rr.Body.Bytes(), &gotAddress); err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}
				if gotAddress.Zipcode != tt.expectedResponse.Zipcode {
					t.Errorf("handler returned unexpected zipcode: got %v want %v", gotAddress.Zipcode, tt.expectedResponse.Zipcode)
				}
				if gotAddress.City != tt.expectedResponse.City {
					t.Errorf("handler returned unexpected city: got %v want %v", gotAddress.City, tt.expectedResponse.City)
				}
				if gotAddress.FederativeUnit != tt.expectedResponse.FederativeUnit {
					t.Errorf("handler returned unexpected federative unit: got %v want %v", gotAddress.FederativeUnit, tt.expectedResponse.FederativeUnit)
				}
			}
		})
	}
}
