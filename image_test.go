package main

import "testing"

func TestDoesImageSupportArm64(t *testing.T) {
	cache := NewInMemoryCache()
	cache.Set("image_with_arm_support", true)
	cache.Set("image_without_arm_support", false)

	type args struct {
		cache Cache
		name  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "image supports arm64",
			args: args{
				cache: cache,
				name:  "image_with_arm_support",
			},
			want: true,
		},
		{
			name: "image that does not support arm64",
			args: args{
				cache: cache,
				name:  "image_without_arm_support",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DoesImageSupportArm64(tt.args.cache, tt.args.name); got != tt.want {
				t.Errorf("DoesImageSupportArm64() = %v, want %v", got, tt.want)
			}
		})
	}
}
