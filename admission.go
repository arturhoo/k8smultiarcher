package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/mattbaird/jsonpatch"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

var MultiarchToleration = corev1.Toleration{
	Key:      "k8smultiarcher",
	Value:    "arm64Supported",
	Operator: corev1.TolerationOpEqual,
	Effect:   "NoSchedule",
}

func ProcessAdmissionReview(cache Cache, requestBody []byte) (*admissionv1.AdmissionReview, error) {
	review, err := AdmissionReviewFromRequest(requestBody)
	if err != nil {
		return nil, err
	}

	if review.Request.Kind.Kind != "Pod" {
		err := fmt.Errorf("got a request for a kind different than pod: %s", review.Request.Kind.Kind)
		slog.Error("invalid request kind", "error", err)
		return nil, err
	}

	obj := review.Request.Object
	pod := &corev1.Pod{}
	err = json.Unmarshal(obj.Raw, pod)
	if err != nil {
		slog.Error("failed to unmarshal pod", "error", err)
		return nil, err
	}

	response := admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
	}

	if !DoesPodSupportArm64(cache, pod) {
		review.Response = &response
		return review, nil
	}

	AddMultiarchTolerationToPod(pod)
	modifiedPod, err := json.Marshal(pod)
	if err != nil {
		slog.Error("failed to marshal pod", "error", err)
		return nil, err
	}

	patch, err := jsonpatch.CreatePatch(obj.Raw, modifiedPod)
	if err != nil {
		slog.Error("failed to create patch for pod", "error", err)
		return nil, err
	}

	jsonPatch, err := json.Marshal(patch)
	if err != nil {
		slog.Error("failed to marshal patch", "error", err)
		return nil, err
	}

	pt := admissionv1.PatchTypeJSONPatch
	response.PatchType = &pt
	response.Patch = jsonPatch
	review.Response = &response

	return review, nil
}

func AdmissionReviewFromRequest(body []byte) (*admissionv1.AdmissionReview, error) {
	var review admissionv1.AdmissionReview
	err := json.Unmarshal(body, &review)
	if err != nil {
		slog.Error("failed to unmarshal request body", "error", err)
		return nil, err
	}

	if review.Request == nil {
		err := fmt.Errorf("got an invalid admission request")
		slog.Error("invalid admission request", "error", err)
		return nil, err
	}

	return &review, nil
}

func DoesPodSupportArm64(cache Cache, pod *corev1.Pod) bool {
	var errs []error
	for _, container := range pod.Spec.Containers {
		if !DoesImageSupportArm64(cache, container.Image) {
			errs = append(errs, fmt.Errorf("image %s lacks arm64 support", container.Image))
		}
	}
	if len(errs) > 0 {
		slog.Info("pod has images without arm64 support", "error", errors.Join(errs...))
		return false
	}
	return true
}

func AddMultiarchTolerationToPod(pod *corev1.Pod) {
	if slices.Contains(pod.Spec.Tolerations, MultiarchToleration) {
		return
	}
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, MultiarchToleration)
}
