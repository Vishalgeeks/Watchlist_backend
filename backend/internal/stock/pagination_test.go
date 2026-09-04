package stock

import "testing"

func TestNormalizePagination(t *testing.T) {
	t.Run("valid values", func(t *testing.T) {
		page, limit, err := normalizePagination("2", "25")
		if err != nil {
			t.Fatalf("normalizePagination returned unexpected error: %v", err)
		}
		if page != 2 || limit != 25 {
			t.Fatalf("normalizePagination returned (%d, %d), want (2, 25)", page, limit)
		}
	})

	t.Run("fallback values", func(t *testing.T) {
		page, limit, err := normalizePagination("0", "999")
		if err != nil {
			t.Fatalf("normalizePagination returned unexpected error: %v", err)
		}
		if page != 1 || limit != 50 {
			t.Fatalf("normalizePagination returned (%d, %d), want (1, 50)", page, limit)
		}
	})
}
