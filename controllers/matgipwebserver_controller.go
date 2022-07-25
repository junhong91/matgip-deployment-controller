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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

var (
	matgipWebServerReplicas int32 = 3
	matgipWebServerImage string = "junhong1991/matgip:1.1.0"
)

// constructDeployment is a method which construct matgip web server deployment
func (r *MatgipRedisServerReconciler) constructDeployment(matgipWebServer *matgipv1.MatgipWebServer) (*appsv1.Deployment, error)  {
	deployment := &appsv1.Deployment {
		ObjectMeta: metav1.ObjectMeta {
			Name: matgipWebServer.Name,
			Namespace: matgipWebServer.Namespace,
		},
		Spec: appsv1.DeploymentSpec {
			Replicas: &matgipWebServerReplicas,
			Selector: &metav1.LabelSelector {
				MatchLabels: map[string]string {
					"app": matgipWebServer.Name,
				},
			},
			Template: corev1.PodTemplateSpec {
				ObjectMeta: metav1.ObjectMeta {
					Name: matgipWebServer.Name,
					Labels: map[string]string {
						"app": matgipWebServer.Name,
						"appNamespace": matgipWebServer.Namespace,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container {
						{
							Name: matgipWebServer.Name,
							Image: matgipWebServerImage,
							Ports: []corev1.ContainerPort {
								{
									Name: matgipWebServer.Name,
									ContainerPort: *matgipWebServer.Spec.PortNumber,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "REDIS_HOST",
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

	return deployment, nil
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

	secret, err := r.constructSecret(&matgipWebServer)
	if err != nil {
		log.Error(err, "unable to construct matgip secret from CRD...")
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, secret); err != nil {
		log.Error(err, "unable to create secret for MatgipWebServer", "secret", secret)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MatgipWebServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&matgipv1.MatgipWebServer{}).
		Complete(r)
}
