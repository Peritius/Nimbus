package weather

import (
	"math"
	"testing"
	"time"
)

func TestCachedNilBeforeStore(t *testing.T) {
	current, daily := Cached(50.45, 30.52)
	if current != nil || daily != nil {
		t.Errorf("Cached before Store returned (%v, %v), want (nil, nil)", current, daily)
	}
}

func TestStoreThenCachedSameCoords(t *testing.T) {
	now := time.Date(2026, 7, 29, 12, 0, 0, 0, time.UTC)
	w := &WeatherData{
		Temperature:  22.5,
		WeatherCode:  0,
		Location:     "Kyiv",
		FetchedAt:    now,
	}
	d := []DailyForecast{{Date: "2026-07-29", TempMax: 25, TempMin: 18}}
	Store(w, d, 50.45, 30.52)

	cw, cd := Cached(50.45, 30.52)
	if cw == nil {
		t.Fatal("Cached returned nil current after Store")
	}
	if cw.Temperature != 22.5 {
		t.Errorf("current.Temperature = %v, want 22.5", cw.Temperature)
	}
	if cd == nil {
		t.Fatal("Cached returned nil daily after Store")
	}
	if len(cd) != 1 || cd[0].Date != "2026-07-29" {
		t.Errorf("daily = %v, want [2026-07-29]", cd)
	}
}

func TestCachedDifferentCoords(t *testing.T) {
	w := &WeatherData{Temperature: 15}
	Store(w, nil, 50.45, 30.52)

	cw, cd := Cached(40.71, -74.01)
	if cw != nil || cd != nil {
		t.Errorf("Cached with different coords returned (%v, %v), want (nil, nil)", cw, cd)
	}
}

func TestCachedSameCoordsButNil(t *testing.T) {
	Store(nil, nil, 50.45, 30.52)
	cw, cd := Cached(50.45, 30.52)
	if cw != nil || cd != nil {
		t.Errorf("Cached after Store(nil) returned (%v, %v), want (nil, nil)", cw, cd)
	}
}

func TestCachedReplacedByNewStore(t *testing.T) {
	w1 := &WeatherData{Temperature: 10, WeatherCode: 1}
	w2 := &WeatherData{Temperature: 20, WeatherCode: 2}
	d2 := []DailyForecast{{Date: "2026-07-30"}}
	Store(w1, nil, 50.45, 30.52)
	Store(w2, d2, 50.45, 30.52)

	cw, cd := Cached(50.45, 30.52)
	if cw == nil || cw.Temperature != 20 {
		t.Errorf("current.Temperature = %v, want 20 after second Store", cw.Temperature)
	}
	if cd == nil || len(cd) != 1 || cd[0].Date != "2026-07-30" {
		t.Errorf("daily = %v, want [2026-07-30] after second Store", cd)
	}
}

func TestCacheMatchedCoordsToleratesFloatsEpsilon(t *testing.T) {
	w := &WeatherData{Temperature: 18}
	Store(w, nil, 50.4501, 30.5234)

	cases := []struct {
		lat, lon float64
		match    bool
	}{
		{50.4501, 30.5234, true},
		{50.4500, 30.5234, false},
		{50.4501, 30.5235, false},
	}
	for _, c := range cases {
		cw, _ := Cached(c.lat, c.lon)
		if c.match && cw == nil {
			t.Errorf("Cached(%.4f, %.4f) = nil, wanted a match", c.lat, c.lon)
		}
		if !c.match && cw != nil {
			t.Errorf("Cached(%.4f, %.4f) = %v, wanted nil", c.lat, c.lon, cw)
		}
	}
}

func TestCacheIsSafeForConcurrentAccess(t *testing.T) {
	w := &WeatherData{Temperature: 22.5}
	done := make(chan bool, 4)
	for i := 0; i < 4; i++ {
		go func() {
			Store(w, nil, 50.45, 30.52)
			Cached(50.45, 30.52)
			done <- true
		}()
	}
	for i := 0; i < 4; i++ {
		<-done
	}
}

func TestFetchedAtSetOnStore(t *testing.T) {
	before := time.Now()
	w := &WeatherData{Temperature: 10}
	Store(w, nil, 50.45, 30.52)
	after := time.Now()

	cw, _ := Cached(50.45, 30.52)
	if cw == nil {
		t.Fatal("Cached returned nil")
	}
	if cw.FetchedAt.Before(before) || cw.FetchedAt.After(after) {
		t.Errorf("FetchedAt = %v, want between %v and %v", cw.FetchedAt, before, after)
	}
}

func TestStoreUpdatesFetchedAt(t *testing.T) {
	w := &WeatherData{Temperature: 10}
	Store(w, nil, 50.45, 30.52)
	time.Sleep(time.Millisecond)
	prev, _ := Cached(50.45, 30.52)
	prevTime := prev.FetchedAt

	w2 := &WeatherData{Temperature: 15}
	Store(w2, nil, 50.45, 30.52)
	cur, _ := Cached(50.45, 30.52)
	if !cur.FetchedAt.After(prevTime) {
		t.Error("FetchedAt did not advance on second Store")
	}
}

var _ = math.NaN
