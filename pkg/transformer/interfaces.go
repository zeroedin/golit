package transformer

import "github.com/sspriggs/golit/pkg/jsengine"

// ElementRenderer is the interface the transformer uses to load and render
// custom elements. *jsengine.Engine satisfies this interface. Defining it
// here decouples the HTML expansion logic from the concrete QJS runtime,
// allowing mock implementations in tests.
type ElementRenderer interface {
	LoadBundleForTag(tagName string, registry *jsengine.Registry) (bool, error)
	RenderBatch(requests []jsengine.BatchRequest) ([]jsengine.BatchResult, error)
	IsRegistered(tagName string) bool
}
