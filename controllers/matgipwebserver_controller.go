/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	matgipv1 "matgip.real-estate.corp/matgip-deployment-controller/api/v1"
)

// MatgipWebServerReconciler reconciles a MatgipWebServer object
type MatgipWebServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	matgipWebServerReplicas int32  = 3
	matgipWebServerImage    string = "junhong1991/matgip:1.1.0"
)

// constructSecret is a method which construct matgip web server secret
func (r *MatgipWebServerReconciler) constructSecret(matgipWebServer *matgipv1.MatgipWebServer) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      matgipWebServer.Name,
			Namespace: matgipWebServer.Namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"news-client-id":     []byte(matgipWebServer.Spec.NewsClientId),
			"news-client-secret": []byte(matgipWebServer.Spec.NewsClientSecret),
			"token-secret":       []byte(matgipWebServer.Spec.AuthTokenSecret),
		},
	}
	if err := ctrl.SetControllerReference(matgipWebServer, secret, r.Scheme); err != nil {
		return nil, err
	}

	return secret, nil
}

// constructDeployment is a method which construct matgip web server deployment
func (r *MatgipWebServerReconciler) constructDeployment(matgipWebServer *matgipv1.MatgipWebServer) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      matgipWebServer.Name,
			Namespace: matgipWebServer.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &matgipWebServerReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": matgipWebServer.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: matgipWebServer.Name,
					Labels: map[string]string{
						"app":          matgipWebServer.Name,
						"appNamespace": matgipWebServer.Namespace,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  matgipWebServer.Name,
							Image: matgipWebServerImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: *matgipWebServer.Spec.PortNumber,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "REDIS_HOST",
									Value: matgipWebServer.Spec.DatabaseName,
								},
								{
									Name: "NAVER_CLIENT_ID",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: matgipWebServer.Name,
											},
											Key: "news-client-id",
										},
									},
								},
								{
									Name: "NAVER_CLIENT_SECRET",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: matgipWebServer.Name,
											},
											Key: "news-client-secret",
										},
									},
								},
								{
									Name: "TOKEN_SECRET",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: matgipWebServer.Name,
											},
											Key: "token-secret",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(matgipWebServer, deployment, r.Scheme); err != nil {
		return nil, err
	}

	return deployment, nil
}

func (r *MatgipWebServerReconciler) constructService(matgipWebServer *matgipv1.MatgipWebServer) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      matgipWebServer.Name,
			Namespace: matgipWebServer.Namespace,
			Labels: map[string]string{
				"app": matgipWebServer.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":          matgipWebServer.Name,
				"appNamespace": matgipWebServer.Namespace,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     matgipWebServer.Name + "-port",
					Protocol: corev1.ProtocolTCP,
					Port:     *matgipWebServer.Spec.PortNumber,
					TargetPort: intstr.IntOrString{
						IntVal: *matgipWebServer.Spec.PortNumber,
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(matgipWebServer, service, r.Scheme); err != nil {
		return nil, err
	}

	return service, nil
}

//+kubebuilder:rbac:groups=matgip.matgip.real-estate.corp,resources=matgipwebservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=matgip.matgip.real-estate.corp,resources=matgipwebservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=matgip.matgip.real-estate.corp,resources=matgipwebservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MatgipWebServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *MatgipWebServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var matgipWebServer matgipv1.MatgipWebServer
	if err := r.Get(ctx, req.NamespacedName, &matgipWebServer); err != nil {
		log.Error(err, "unable to fetch MatgipWebServer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.V(2).Info("Successfully fetch MatgipWebServer CRD...")

	findSecret := types.NamespacedName{
		Name:      matgipWebServer.Name,
		Namespace: matgipWebServer.Namespace,
	}
	var foundSecret corev1.Secret
	if err := r.Get(ctx, findSecret, &foundSecret); err != nil && errors.IsNotFound(err) {
		log.V(2).Info("Create matgip web secret...")

		secret, err := r.constructSecret(&matgipWebServer)
		if err != nil {
			log.Error(err, "unable to construct matgip secret from CRD...")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, secret); err != nil {
			log.Error(err, "unable to create secret for MatgipWebServer", "secret", secret)
			return ctrl.Result{}, err
		}
	}

	findDeployment := types.NamespacedName{
		Name:      matgipWebServer.Name,
		Namespace: matgipWebServer.Namespace,
	}
	var foundDeployment appsv1.Deployment
	if err := r.Get(ctx, findDeployment, &foundDeployment); err != nil && errors.IsNotFound(err) {
		deployment, err := r.constructDeployment(&matgipWebServer)
		if err != nil {
			log.Error(err, "unable to construct matgip deployment from CRD...")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, deployment); err != nil {
			log.Error(err, "unable to create deployment for MatgipWebServer", "deployment", deployment)
			return ctrl.Result{}, err
		}
	}

	findService := types.NamespacedName{
		Name:      matgipWebServer.Name,
		Namespace: matgipWebServer.Namespace,
	}
	var foundService corev1.Service
	if err := r.Get(ctx, findService, &foundService); err != nil && errors.IsNotFound(err) {
		service, err := r.constructService(&matgipWebServer)
		if err != nil {
			log.Error(err, "unable to construct matgip service from CRD...")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, service); err != nil {
			log.Error(err, "unable to create service for MatgipWebServer", "service", service)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MatgipWebServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&matgipv1.MatgipWebServer{}).
		Complete(r)
}
