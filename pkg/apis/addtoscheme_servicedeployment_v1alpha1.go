package apis

import (
	"github.com/angudadevops/service-deployment-operator/pkg/apis/servicedeployment/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
