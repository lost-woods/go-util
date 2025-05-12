package k8sutil

import "k8s.io/client-go/tools/cache"

func (watcher *GenericWatcher) IsRunning() bool {
	return watcher.InformerChannel != nil
}

func (watcher *GenericWatcher) Start(informer cache.SharedIndexInformer) *GenericWatcher {
	if watcher.IsRunning() {
		log.Warnf("Watcher '%s' is already running, skipping", watcher.Name)
		return watcher
	}

	log.Infof("Starting watcher '%s'...", watcher.Name)
	watcher.InformerChannel = make(chan struct{})
	go informer.Run(watcher.InformerChannel)

	if !cache.WaitForCacheSync(watcher.InformerChannel, informer.HasSynced) {
		log.Fatalf("Failed to sync informer cache for watcher '%s'", watcher.Name)
	}

	return watcher
}

func (watcher *GenericWatcher) Wait() {
	if !watcher.IsRunning() {
		log.Warnf("Watcher '%s' is not running, nothing to wait for", watcher.Name)
		return
	}

	<-watcher.InformerChannel
}

func (watcher *GenericWatcher) Stop() {
	if !watcher.IsRunning() {
		log.Warnf("Watcher '%s' is not running, nothing to stop", watcher.Name)
		return
	}

	log.Infof("Stopping watcher '%s'...", watcher.Name)
	close(watcher.InformerChannel)
	watcher.InformerChannel = nil
}
