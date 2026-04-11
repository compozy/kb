package convert

import "testing"

func TestImageConverterAcceptsSupportedExtensionsAndMIMETypes(t *testing.T) {
	t.Parallel()

	converter := ImageConverter{}

	cases := []struct {
		ext      string
		mimeType string
		want     bool
	}{
		{ext: ".png", want: true},
		{ext: ".jpg", want: true},
		{ext: ".jpeg", want: true},
		{ext: ".tiff", want: true},
		{ext: ".bmp", want: true},
		{mimeType: "image/png", want: true},
		{mimeType: "image/jpeg", want: true},
		{mimeType: "image/tiff", want: true},
		{mimeType: "image/bmp", want: true},
		{ext: ".gif", want: false},
		{mimeType: "image/gif", want: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.ext+tc.mimeType, func(t *testing.T) {
			t.Parallel()

			if got := converter.Accepts(tc.ext, tc.mimeType); got != tc.want {
				t.Fatalf("Accepts(%q, %q) = %t, want %t", tc.ext, tc.mimeType, got, tc.want)
			}
		})
	}
}
