package base

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"strings"
	"time"
)

type Resource interface {
	GetName() string
	GetNamespace() string
}

func ResourceEventHandler[T Resource](handler *BaseHandler, additionalEvents ...map[string]func()) cache.ResourceEventHandlerFuncs {
	handleEvent := func(obj any) {
		resource := obj.(T)
		// GetList
		go handler.Event.AddEvent(handler.Kind, handler.processListEvents(resource.GetName()))

		streamName := fmt.Sprintf("%s-%s", resource.GetNamespace(), resource.GetName())
		// GetDetails
		go handler.Event.AddEvent(streamName, handler.processDetailsEvents(resource.GetNamespace(), resource.GetName()))

		// GetYAML
		go handler.Event.AddEvent(streamName+"-yaml", handler.processYAMLEvents(resource.GetNamespace(), resource.GetName()))

		for _, event := range additionalEvents {
			for key, e := range event {
				go handler.Event.AddEvent(key, e)
			}
		}
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc:    handleEvent,
		UpdateFunc: func(oldObj, newObj any) { handleEvent(oldObj) },
		DeleteFunc: handleEvent,
	}
}

func (h *BaseHandler) StartInformer(c echo.Context, cache cache.ResourceEventHandlerFuncs) {
	h.baseInformer(c, cache)
	go h.Container.SharedInformerFactory(h.QueryConfig, h.QueryCluster).Start(context.Background().Done())
}

func (h *BaseHandler) StartExtensionInformer(c echo.Context, cache cache.ResourceEventHandlerFuncs) {
	h.baseInformer(c, cache)
	go h.Container.ExtensionSharedFactoryInformer(h.QueryConfig, h.QueryCluster).Start(context.Background().Done())
}

func (h *BaseHandler) StartDynamicInformer(c echo.Context, cache cache.ResourceEventHandlerFuncs) {
	h.baseInformer(c, cache)
	go h.Container.DynamicSharedInformerFactory(h.QueryConfig, h.QueryCluster).Start(context.Background().Done())
}

func (h *BaseHandler) baseInformer(_ echo.Context, cache cache.ResourceEventHandlerFuncs) {
	cacheKey := fmt.Sprintf("%s-%s-%s", h.QueryConfig, h.QueryCluster, h.InformerCacheKey)
	if h.Container.Cache().Has(cacheKey) {
		return
	}
	h.Container.Cache().Set(cacheKey, "started")
	h.Informer.AddEventHandler(cache)
}

func (h *BaseHandler) WaitForSync(c echo.Context) {
	err := wait.PollUntilContextCancel(c.Request().Context(), 100*time.Millisecond, true, func(context.Context) (done bool, err error) {
		hasSynced := h.Informer.HasSynced()
		if hasSynced {
			h.Event.AddEvent(h.Kind, h.processListEvents(""))
		}
		return hasSynced, nil
	})
	if err != nil {
		log.Error("failed to load informer for sync")
		return
	}
}

func (h *BaseHandler) IsNamespaceResource(r string) bool {
	cacheKey := fmt.Sprintf("%s-%s-nonNamespacedResources", h.QueryConfig, h.QueryCluster)
	c, exists := h.Container.Cache().Get(cacheKey)
	if !exists {
		return false
	}
	nonNamespacedResources := c.([]string)

	for _, resource := range nonNamespacedResources {
		if strings.EqualFold(r, resource) {
			return false
		}
	}

	return true
}

func (h *BaseHandler) processListEvents(resourceName string) func() {
	return func() {
		items := h.Informer.GetStore().List()
		data := h.marshalListData(items, resourceName)
		streamID := fmt.Sprintf("%s-%s-%s", h.QueryConfig, h.QueryCluster, h.Kind)
		h.Container.SSE().Publish(streamID, &sse.Event{
			Data: data,
		})
	}
}

func (h *BaseHandler) processDetailsEvents(namespace, name string) func() {
	return func() {
		streamID, item, exists, _ := h.getStreamIDAndItem(namespace, name)
		data := h.marshalDetailData(item, exists)
		h.Container.SSE().Publish(streamID, &sse.Event{
			Data: data,
		})
	}
}

func (h *BaseHandler) processYAMLEvents(namespace, name string) func() {
	return func() {
		streamID, item, exists, _ := h.getStreamIDAndItem(namespace, name)
		data := h.marshalYAML(item, exists)
		h.Container.SSE().Publish(fmt.Sprintf("%s-yaml", streamID), &sse.Event{
			Data: data,
		})
	}
}
