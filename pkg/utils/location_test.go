package utils

import (
	"fmt"
	"testing"
)

func TestSetLocation(t *testing.T) {
	expectedLat := 52.5170365
	expectedLon := 13.3888599
	err := SetLocation(expectedLat, expectedLon)
	if err != nil {
		t.Error("SetLocation() in TestGetLocation() returned an unexpected error: ", err)
	}
	lat, lon, err := GetLocation()
	if err != nil {
		t.Error("GetLocation() returned an unexpected error: ", err)
	}
	if lat != expectedLat {
		t.Errorf("GetLoation() returned wrong latitude: got \"%v\" want \"%v\"", lat, expectedLat)
	}
	if lon != expectedLon {
		t.Errorf("GetLoation() returned wrong longitude: got \"%v\" want \"%v\"", lon, expectedLon)
	}
	expectedLat = 48.1371079
	expectedLon = 11.5753822
	err = SetLocation(expectedLat, expectedLon)
	if err != nil {
		t.Error("SetLocation() in TestGetLocation() returned an unexpected error: ", err)
	}
	lat, lon, err = GetLocation()
	if err != nil {
		t.Error("GetLocation() returned an unexpected error: ", err)
	}
	if lat != expectedLat {
		t.Errorf("GetLoation() returned wrong latitude: got \"%v\" want \"%v\"", lat, expectedLat)
	}
	if lon != expectedLon {
		t.Errorf("GetLoation() returned wrong longitude: got \"%v\" want \"%v\"", lon, expectedLon)
	}
	t.Run("checking for unexpected errors", func(t *testing.T) {
		expectedLat = 19
		expectedLon = 63
		expectedError := fmt.Errorf("location (Lat: \"%v\" Lon: \"%v\") is not in range!", expectedLat, expectedLon).Error()
		err = SetLocation(expectedLat, expectedLon)
		if err == nil || err.Error() != expectedError {
			t.Error("SetLocation() returned an unexpected error: ", err)
		}
		expectedLat = 49
		expectedLon = 29
		expectedError = fmt.Errorf("location (Lat: \"%v\" Lon: \"%v\") is not in range!", expectedLat, expectedLon).Error()
		err = SetLocation(expectedLat, expectedLon)
		if err == nil || err.Error() != expectedError {
			t.Error("SetLocation() returned an unexpected error: ", err)
		}
	})
}
