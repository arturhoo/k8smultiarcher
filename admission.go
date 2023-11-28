package main

import (
	"encoding/json"
	"fmt"

	"github.com/arturhoo/k8smultiarcher/image"

	"github.com/mattbaird/jsonpatch"
	"github.com/rs/zerolog/log"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

var MultiarchToleration = corev1.Toleration{
	Key:      "k8smultiarcher",
	Value:    "arm64Supported",
	Operator: corev1.TolerationOpEqual,
	Effect:   "NoSchedule",
}

func ProcessAdmissionReview(requestBody []byte) (*admissionv1.AdmissionReview, error) {
	review, err := AdmissionReviewFromRequest(requestBody)
	if err != nil {
		return nil, err
	}

	if review.Request.Kind.Kind != "Pod" {
		err := fmt.Errorf("got a request for a kind different than pod: %s", review.Request.Kind.Kind)
		log.Print(err)
		return nil, err
	}

	obj := review.Request.Object
	pod := &corev1.Pod{}
	err = json.Unmarshal(obj.Raw, pod)
	if err != nil {
		log.Error().Msgf("failed to unmarshal pod: %s", err)
		return nil, err
	}

	response := admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
	}

	if !DoesPodSupportArm64(pod) {
		review.Response = &response
		return review, nil
	}

	AddMultiarchTolerationToPod(pod)
	modifiedPod, err := json.Marshal(pod)
	if err != nil {
		log.Error().Msgf("failed to marshal pod: %s", err)
		return nil, err
	}

	patch, err := jsonpatch.CreatePatch(obj.Raw, modifiedPod)
	if err != nil {
		log.Error().Msgf("failed to create patch for pod: %s", err)
		return nil, err
	}

	jsonPatch, err := json.Marshal(patch)
	if err != nil {
		log.Error().Msgf("failed to marshal patch: %s", err)
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
		log.Printf("got err while unmarshalling request body: %s", err)
		return nil, err
	}

	if review.Request == nil {
		err := fmt.Errorf("got an invalid admission request")
		log.Print(err)
		return nil, err
	}

	return &review, nil
}

func DoesPodSupportArm64(pod *corev1.Pod) bool {
	supported := true
	for _, container := range pod.Spec.Containers {
		if !image.DoesImageSupportArm64(container.Image) {
			supported = false
		}
	}
	return supported
}

func AddMultiarchTolerationToPod(pod *corev1.Pod) {
	for _, toleration := range pod.Spec.Tolerations {
		if toleration == MultiarchToleration {
			return
		}
	}

	pod.Spec.Tolerations = append(pod.Spec.Tolerations, MultiarchToleration)
}
