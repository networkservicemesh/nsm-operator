package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newObjectMeta(name string, namespace string, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
}
