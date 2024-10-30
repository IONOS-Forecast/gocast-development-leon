package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestSetLocation(t *testing.T) {
	tests := []struct {
		expectedLat float64
		expectedLon float64
		wantErr     bool
	}{
		{52.5170365, 13.3888599, false},
		{48.1371079, 11.5753822, false},
		{19, 63, true},
		{49, 29, true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test-%v", i), func(t *testing.T) {
			if tt.wantErr {
				log.Print("Checking for unexpected errors...")
				expectedError := fmt.Errorf("location (Lat: \"%v\" Lon: \"%v\") is not in range!", tt.expectedLat, tt.expectedLon).Error()
				err := SetLocation(tt.expectedLat, tt.expectedLon)
				if err == nil || err.Error() != expectedError {
					t.Error("SetLocation() returned an unexpected error: ", err)
				}
				log.Print("No unexpected errors found")
			} else {
				err := SetLocation(tt.expectedLat, tt.expectedLon)
				if err != nil {
					t.Error("SetLocation() in TestGetLocation() returned an unexpected error: ", err)
				}
				lat, lon, err := GetLocation()
				if err != nil {
					t.Error("GetLocation() returned an unexpected error: ", err)
				}
				if lat != tt.expectedLat {
					t.Errorf("GetLoation() returned wrong latitude: got \"%v\" want \"%v\"", lat, tt.expectedLat)
				}
				if lon != tt.expectedLon {
					t.Errorf("GetLoation() returned wrong longitude: got \"%v\" want \"%v\"", lon, tt.expectedLon)
				}
			}
		})
	}
}
