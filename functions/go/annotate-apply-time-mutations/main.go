package main

import (
	"fmt"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-catalog/functions/go/annotate-apply-time-mutations/annotateapplytimemutations"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
)

const (
	annotationKey = "config.kubernetes.io/apply-time-mutation"
)

type AnnotateApplyTimeMutationsProcessor struct{}

func (rp *AnnotateApplyTimeMutationsProcessor) Process(resourceList *framework.ResourceList) error {
	resourceList.Result = &framework.Result{
		Name: "annotate-apply-time-mutations",
	}
	for _, node := range resourceList.Items {
		fileName, _, err := kioutil.GetFileAnnotations(node)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Processing resource in file", fileName)
		ra := annotateapplytimemutations.ResourceAnnotator{}
		results, err := ra.AnnotateResource(node, fileName)
		resourceList.Result.Items = append(resourceList.Result.Items, results...)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	rp := AnnotateApplyTimeMutationsProcessor{}
	cmd := command.Build(&rp, command.StandaloneEnabled, false)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}