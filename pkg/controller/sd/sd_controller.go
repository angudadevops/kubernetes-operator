package sd

import (
	"context"

	servicedeploymentv1alpha1 "github.com/angudadevops/service-deployment-operator/pkg/apis/servicedeployment/v1alpha1"
        appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var log = logf.Log.WithName("controller_sd")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SD Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSD{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sd-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SD
	err = c.Watch(&source.Kind{Type: &servicedeploymentv1alpha1.SD{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/* TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner SD
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &servicedeploymentv1alpha1.SD{},
	})
	if err != nil {
		return err
	}
	*/

	watchTypes := []runtime.Object{
                &appsv1.Deployment{},
                &corev1.Service{},
        }

        for i := range watchTypes {
                err = c.Watch(&source.Kind{Type: watchTypes[i]}, &handler.EnqueueRequestForOwner{
                        IsController: true,
                        OwnerType:    &servicedeploymentv1alpha1.SD{},
                })
                if err != nil {
                        return err
                }
        }

	return nil
}

// blank assignment to verify that ReconcileSD implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSD{}

// ReconcileSD reconciles a SD object
type ReconcileSD struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SD object and makes changes based on the state read
// and what is in the SD.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSD) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SD")

	// Fetch the SD instance
	instance := &servicedeploymentv1alpha1.SD{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	labels := map[string]string{
                "app": instance.Name,
        }

        deployment, err := r.buildDeployment(instance, labels)
        if err != nil {
                return reconcile.Result{}, err
        }

        service, err := r.buildService(instance, labels)
        if err != nil {
                return reconcile.Result{}, err
        }

        foundDep := &appsv1.Deployment{}
        foundService := &corev1.Service{}

	/*
	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set SD instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	found := &corev1.Pod{}
	*/

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDep)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err == nil && foundDep.Spec.Replicas != instance.Spec.Replicas {
                reqLogger.Info("updating existing Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
                foundDep.Spec.Replicas = instance.Spec.Replicas
                if err := r.client.Update(context.TODO(), foundDep); err != nil {
                        return reconcile.Result{}, err
                }
        } else if err != nil {
		return reconcile.Result{}, err
	}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
        if err != nil && errors.IsNotFound(err) {
                reqLogger.Info("creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
                if err := r.client.Create(context.TODO(), service); err != nil {
                        return reconcile.Result{}, err
                }
        } else if err == nil {
                reqLogger.Info("updating existing Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)

                foundService.Spec.Ports = service.Spec.Ports

                if err := r.client.Update(context.TODO(), foundService); err != nil {
                        return reconcile.Result{}, err
                }
        } else if err != nil {
                return reconcile.Result{}, err
        }

	return reconcile.Result{}, nil
}

func (r *ReconcileSD) buildDeployment(cr *servicedeploymentv1alpha1.SD, labels map[string]string) (*appsv1.Deployment, error) {

        deployment := appsv1.Deployment{
                ObjectMeta: metav1.ObjectMeta{
                        Name:      cr.Name,
                        Namespace: cr.Namespace,
                        Labels:    labels,
                },
                Spec: appsv1.DeploymentSpec{
                        Replicas: cr.Spec.Replicas,
                        Selector: &metav1.LabelSelector{
                                MatchLabels: labels,
                        },
                        Template: corev1.PodTemplateSpec{
                                ObjectMeta: metav1.ObjectMeta{
                                        Labels: labels,
                                },
                                Spec: corev1.PodSpec{
                                        Containers: []corev1.Container{
                                                corev1.Container{
                                                        Name:  cr.Name,
                                                        Image: cr.Spec.Image,
                                                        Ports: []corev1.ContainerPort{
                                                                corev1.ContainerPort{Name: "http", ContainerPort: cr.Spec.ContainerPort},
                                                        },
                                                },
                                        },
                                },
                        },
                },
        }

        // Set Service Deployment instance as the owner and controller
        if err := controllerutil.SetControllerReference(cr, &deployment, r.scheme); err != nil {
                return nil, err
        }

        return &deployment, nil
}
func (r *ReconcileSD) buildService(cr *servicedeploymentv1alpha1.SD, labels map[string]string) (*corev1.Service, error) {
        service := corev1.Service{
                ObjectMeta: metav1.ObjectMeta{
                        Name:      cr.Name,
                        Namespace: cr.Namespace,
                        Labels:    labels,
                },
                Spec: corev1.ServiceSpec{
                        Ports: []corev1.ServicePort{
                                corev1.ServicePort{Name: "http", TargetPort: intstr.FromString("http"), Port: cr.Spec.ContainerPort, NodePort: cr.Spec.NodePort},
                        },
                        Selector: labels,
                        Type: "NodePort",
                },
        }

        // Set Service Deployment instance as the owner and controller
        if err := controllerutil.SetControllerReference(cr, &service, r.scheme); err != nil {
                return nil, err
        }

        return &service, nil
}
