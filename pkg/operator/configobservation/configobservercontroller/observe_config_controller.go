package configobservercontroller

import (
	"k8s.io/client-go/tools/cache"

	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-kube-controller-manager-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-kube-controller-manager-operator/pkg/operator/configobservation/cloudprovider"
	"github.com/openshift/cluster-kube-controller-manager-operator/pkg/operator/configobservation/network"
)

type ConfigObserver struct {
	*configobserver.ConfigObserver
}

func NewConfigObserver(
	operatorClient v1helpers.OperatorClient,
	operatorConfigInformers operatorv1informers.SharedInformerFactory,
	configinformers configinformers.SharedInformerFactory,
	resourceSyncer resourcesynccontroller.ResourceSyncer,
	eventRecorder events.Recorder,
) *ConfigObserver {
	c := &ConfigObserver{
		ConfigObserver: configobserver.NewConfigObserver(
			operatorClient,
			eventRecorder,
			configobservation.Listers{
				InfrastructureLister: configinformers.Config().V1().Infrastructures().Lister(),
				NetworkLister:        configinformers.Config().V1().Networks().Lister(),
				ResourceSync:         resourceSyncer,
				PreRunCachesSynced: []cache.InformerSynced{
					configinformers.Config().V1().Infrastructures().Informer().HasSynced,
					configinformers.Config().V1().Networks().Informer().HasSynced,
				},
			},
			cloudprovider.ObserveCloudProviderNames,
			network.ObserveClusterCIDRs,
			network.ObserveServiceClusterIPRanges,
		),
	}

	operatorConfigInformers.Operator().V1().KubeControllerManagers().Informer().AddEventHandler(c.EventHandler())
	configinformers.Config().V1().Infrastructures().Informer().AddEventHandler(c.EventHandler())
	configinformers.Config().V1().Networks().Informer().AddEventHandler(c.EventHandler())

	return c
}
