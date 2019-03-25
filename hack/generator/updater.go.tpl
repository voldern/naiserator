// generated by your friendly code generator. DO NOT EDIT.
// to refresh this file, run `go generate` in your shell.

package updater

import (
    "fmt"

    log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    autoscalingv1 "k8s.io/api/autoscaling/v1"
    extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
    networkingv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    typed_core_v1 "k8s.io/client-go/kubernetes/typed/core/v1"
    typed_apps_v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
    typed_autoscaling_v1 "k8s.io/client-go/kubernetes/typed/autoscaling/v1"
    typed_extensions_v1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
    typed_networking_v1 "k8s.io/client-go/kubernetes/typed/networking/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/kubernetes"
)

{{range .}}

func {{.Name}}(client {{.Interface}}, old, new {{.Type}}) func() error {
    log.Infof("creating or updating {{ .Type }} for %s", new.Name)
	if old == nil {
		return func() error {
			_, err := client.Create(new)
			return err
		}
	}

	CopyMeta(old, new)
	{{if .TransformFunc}}
	    {{.TransformFunc}}(old, new)
    {{end}}

	return func() error {
		_, err := client.Update(new)
		return err
	}
}

{{end}}

func Updater(clientSet kubernetes.Interface, resource runtime.Object) func() error {
	switch new := resource.(type) {
	{{range .}}
		case {{.Type}}:
		c := clientSet.{{.ClientType}}(new.Namespace)
		old, err := c.Get(new.Name, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return func() error { return err }
			}
			return {{.Name}}(c, nil, new)
		}
		return {{.Name}}(c, old, new)
	{{end}}
	default:
		panic(fmt.Errorf("BUG! You didn't specify a case for type '%T' in the file hack/generator/updater.go", new))
	}
}
