package controllers

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	taskv1alpha1 "github.com/digital-plumbers-union/taskwarrior-operator/api/v1alpha1"
)

// TaskwarriorReconciler reconciles a Taskwarrior object
type TaskwarriorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=task.dpu.sh,resources=taskwarriors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=task.dpu.sh,resources=taskwarriors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list

func (r *TaskwarriorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("taskwarrior", req.NamespacedName)

	// Fetch the Taskwarrior instance
	taskwarrior := &taskv1alpha1.Taskwarrior{}
	err := r.Get(ctx, req.NamespacedName, taskwarrior)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Taskwarrior resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}

		log.Error(err, "Failed to get Taskwarrior")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: taskwarrior.Name, Namespace: taskwarrior.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create new deployment
		dep := r.deploymentForTaskwarrior(taskwarrior)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err := r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := taskwarrior.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err := r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the Taskwarrior status with the pod names
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(taskwarrior.Namespace),
		client.MatchingLabels(labelsForTaskwarrior(taskwarrior.Name)),
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Taskwarrior.Namespace", taskwarrior.Namespace, "Taskwarrior.Name", taskwarrior.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, taskwarrior.Status.Nodes) {
		taskwarrior.Status.Nodes = podNames
		err := r.Status().Update(ctx, taskwarrior)
		if err != nil {
			log.Error(err, "Failed to update Taskwarrior status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// deploymentForTaskwarrior returns a Taskwarrior Deployment object
func (r *TaskwarriorReconciler) deploymentForTaskwarrior(t *taskv1alpha1.Taskwarrior) *appsv1.Deployment {
	ls := labelsForTaskwarrior(t.Name)
	replicas := t.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      t.Name,
			Namespace: t.Namespace,
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
					Containers: []corev1.Container{
						{
							Image:   "dpush/taskwarrior-server:edge",
							Name:    "taskd",
							Command: []string{"taskd", "--help"},
							Ports: []corev1.ContainerPort{{
								ContainerPort: 53589,
								Protocol:      corev1.ProtocolTCP,
								Name:          "taskd",
							}},
						},
						// TODO: Container for REST API
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(t, dep, r.Scheme)
	return dep
}

func labelsForTaskwarrior(name string) map[string]string {
	return map[string]string{"app": "taskwarrior", "taskwarrior_cr": name}
}

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func (r *TaskwarriorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1alpha1.Taskwarrior{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
