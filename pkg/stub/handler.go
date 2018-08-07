package stub

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/sjones-sot/uwsgi-operator/pkg/apis/sourceoftruth/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHandler returns a new handler instance
func NewHandler() sdk.Handler {
	return &Handler{}
}

// Handler type - not figured a use for this yet
type Handler struct {
	// Fill me
}

// Handle deals with expected event types on the CR deployment
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.UwsgiApp:
		uwsgiapp := o

		if event.Deleted {
			return nil
		}

		dep := deploymentForUwsgiApp(uwsgiapp)
		err := sdk.Create(dep)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create deployment: %v", err)
			return err
		}

		err = sdk.Get(dep)
		if err != nil {
			logrus.Errorf("failed to get deployment: %v", err)
		}

		replicas := uwsgiapp.Spec.Replicas
		if *dep.Spec.Replicas != replicas {
			dep.Spec.Replicas = &replicas
			err = sdk.Update(dep)
			if err != nil {
				logrus.Errorf("Failed to update deployment: %v", err)
				return err
			}
		}

		podList := podList()
		labelSelector := labels.SelectorFromSet(labelsForUwsgiApp(uwsgiapp.Name)).String()
		listOps := &metav1.ListOptions{LabelSelector: labelSelector}
		err = sdk.List(uwsgiapp.Namespace, podList, sdk.WithListOptions(listOps))
		if err != nil {
			logrus.Errorf("failed to list pods %v: ", err)
			return err
		}
		podNames := getPodNames(podList.Items)
		if !reflect.DeepEqual(podNames, uwsgiapp.Status.Nodes) {
			uwsgiapp.Status.Nodes = podNames
			err := sdk.Update(uwsgiapp)
			if err != nil {
				logrus.Errorf("failed to update UwsgiApp status: %v", err)
				return err
			}
		}

	}
	return nil
}

// deploymentForUwsgiApp returns a UwsgiApp deployment object
func deploymentForUwsgiApp(u *v1alpha1.UwsgiApp) *appsv1.Deployment {
	ls := labelsForUwsgiApp(u.Name)
	replicas := u.Spec.Replicas
	var command []string
	switch {
	case u.Spec.Command != nil:
		command = u.Spec.Command
	default:
		command = append(command, "uwsgi", "--ini", "/etc/uwsgi.ini")
	}

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      u.Name,
			Namespace: u.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   u.Spec.Image,
						Name:    u.Spec.ApplicationName,
						Command: command,
						Ports:   u.Spec.Ports,
					}},
				},
			},
		},
	}
	addOwnerRefToObject(dep, asOwner(u))
	return dep
}

func labelsForUwsgiApp(name string) map[string]string {
	return map[string]string{"k8s_app": "UwsgiApp", "UwsgiApp_cr": name}
}

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

func asOwner(u *v1alpha1.UwsgiApp) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: u.APIVersion,
		Kind:       u.Kind,
		Name:       u.Name,
		UID:        u.UID,
		Controller: &trueVar,
	}
}

// podList returns a v1.PodList object
func podList() *corev1.PodList {
	return &corev1.PodList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
