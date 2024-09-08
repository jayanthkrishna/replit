package controllers

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Kconfig struct {
	clientset *kubernetes.Clientset
}

func NewKconfig(c *kubernetes.Clientset) *Kconfig {
	return &Kconfig{
		clientset: c,
	}
}

func KubernetesConfig() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		fmt.Println("Error while configuring the kube config flags")
		return nil
	}

	return clientset

}

func (c *Kconfig) createResources(name string, env string, ns string) (string, error) {

	deployment, err := c.createDeployment(name, env, ns)
	if err != nil {

		fmt.Printf("Failed to create Deployment : %s\n", name)
		return "", err
	}
	fmt.Printf("Created deployment %q.\n", deployment.GetObjectMeta().GetName())

	service, err := c.createService(deployment, ns)
	if err != nil {
		fmt.Printf("Failed to create service : %s\n", name)
		return "", err
	}

	fmt.Printf("Created Service %q.\n", service.GetObjectMeta().GetName())

	ingress, err := c.createIngress(context.TODO(), service)

	if err != nil {
		fmt.Printf("Failed to create ingress : %s\n", name)
		return "", err
	}

	fmt.Printf("Created ingress %q.\n", ingress.GetObjectMeta().GetName())

	url := ingress.Spec.Rules[0].Host + "/" + name

	return url, nil
}
func (c *Kconfig) createDeployment(name string, env string, ns string) (*appsv1.Deployment, error) {

	// deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploymentsClient := c.clientset.AppsV1().Deployments(ns)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: fmt.Sprintf("%s:%s", env, "latest"),
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	return deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

}

func (c *Kconfig) createService(deployment *appsv1.Deployment, ns string) (*apiv1.Service, error) {

	svc := apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: ns,
		},
		Spec: apiv1.ServiceSpec{
			Selector: deplLabels(*deployment),
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       8080,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}

	s, err := c.clientset.CoreV1().Services(ns).Create(context.TODO(), &svc, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf("Error creating service : %s\n", err.Error())
		return nil, err
	}

	fmt.Printf("Created deployment %q.\n", s.GetObjectMeta().GetName())
	return s, nil

}

func (c *Kconfig) createIngress(ctx context.Context, svc *apiv1.Service) (*netv1.Ingress, error) {
	pathType := "Prefix"
	ingress := netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					Host: "example.dev",
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("/%s", svc.Name),
									PathType: (*netv1.PathType)(&pathType),
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: svc.Name,
											Port: netv1.ServiceBackendPort{
												Number: 8080,
											},
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

	return c.clientset.NetworkingV1().Ingresses(svc.Namespace).Create(ctx, &ingress, metav1.CreateOptions{})

}

func deplLabels(depl appsv1.Deployment) map[string]string {
	return depl.Spec.Template.Labels
}

func int32Ptr(i int32) *int32 { return &i }
