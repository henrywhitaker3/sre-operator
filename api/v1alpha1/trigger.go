package v1alpha1

type Trigger struct {
	//+kubebuilder:validation:Required
	Name string `json:"name"`

	//+kubebuilder:validation:Required
	Kind string `json:"kind"`
}
