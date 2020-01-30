package apimachinery

import (
	errors "github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func hasEphemeralContainer(ecs []v1.EphemeralContainer, ecName string) bool {
	for _, ec := range ecs {
		if ec.Name == ecName {
			return true
		}
	}
	return false
}

func createDebugContainerObject(name, image string) v1.EphemeralContainer {
	return v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{
			Name:                     name,
			Image:                    image,
			ImagePullPolicy:          v1.PullIfNotPresent,
			TerminationMessagePolicy: v1.TerminationMessageReadFile,
			Stdin:                    true,
			//SecurityContext: &v1.SecurityContext{
			//	Capabilities: &v1.Capabilities{
			//		Add: []v1.Capability{"NET_ADMIN"},
			//	},
			//},
		},
	}
}

func CreateOrUpdateEphemeralContainer(config *rest.Config, namespace, podName, ephemeralContainerName, image string) error {
	debugContainer := createDebugContainerObject(ephemeralContainerName, image)
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	pods := client.CoreV1().Pods(namespace)
	ec, err := pods.GetEphemeralContainers(podName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrap(err, "ephemeral containers are not enabled for this cluster")
		}
		return err
	}

	if !hasEphemeralContainer(ec.EphemeralContainers, debugContainer.Name) {
		ec.EphemeralContainers = append(ec.EphemeralContainers, debugContainer)
		_, err = pods.UpdateEphemeralContainers(podName, ec)
		if err != nil {
			return err
		}
	}
	return nil
}
