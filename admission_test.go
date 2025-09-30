package main

import (
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestAddMultiarchTolerationToPod(t *testing.T) {
	pod := &corev1.Pod{
		Spec: corev1.PodSpec{
			Tolerations: []corev1.Toleration{
				{
					Key:      "key1",
					Operator: corev1.TolerationOpEqual,
					Value:    "value1",
				},
				{
					Key:      "key2",
					Operator: corev1.TolerationOpEqual,
					Value:    "value2",
				},
			},
		},
	}
	AddMultiarchTolerationToPod(pod)
	expectedTolerations := []corev1.Toleration{
		{
			Key:      "key1",
			Operator: corev1.TolerationOpEqual,
			Value:    "value1",
		},
		{
			Key:      "key2",
			Operator: corev1.TolerationOpEqual,
			Value:    "value2",
		},
		MultiarchToleration,
	}

	if !slices.Equal(pod.Spec.Tolerations, expectedTolerations) {
		t.Errorf("Unexpected tolerations. Expected: %v, Got: %v", expectedTolerations, pod.Spec.Tolerations)
	}
}
