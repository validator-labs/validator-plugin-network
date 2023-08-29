package validators

import (
	"github.com/go-logr/logr"
	ktypes "k8s.io/apimachinery/pkg/types"

	"github.com/spectrocloud-labs/valid8or-plugin-network/api/v1alpha1"
	"github.com/spectrocloud-labs/valid8or/pkg/types"
)

type NetworkService struct {
	log logr.Logger
}

func NewNetworkService(log logr.Logger) *NetworkService {
	return &NetworkService{
		log: log,
	}
}

func (n *NetworkService) ReconcileDNSRule(nn ktypes.NamespacedName, rule v1alpha1.DNSRule) (*types.ValidationResult, error) {
	return nil, nil
}

func (n *NetworkService) ReconcileICMPRule(nn ktypes.NamespacedName, rule v1alpha1.ICMPRule) (*types.ValidationResult, error) {
	return nil, nil
}

func (n *NetworkService) ReconcileIPRangeRule(nn ktypes.NamespacedName, rule v1alpha1.IPRangeRule) (*types.ValidationResult, error) {
	return nil, nil
}

func (n *NetworkService) ReconcileMTURule(nn ktypes.NamespacedName, rule v1alpha1.MTURule) (*types.ValidationResult, error) {
	return nil, nil
}

func (n *NetworkService) ReconcileNetcatRule(nn ktypes.NamespacedName, rule v1alpha1.NetcatRule) (*types.ValidationResult, error) {
	return nil, nil
}
