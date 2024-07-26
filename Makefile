include build/makelib/common.mk
include build/makelib/plugin.mk

# Image URL to use all building/pushing image targets
IMG ?= quay.io/validator-labs/validator-plugin-network:latest

# Helm vars
CHART_NAME=validator-plugin-network

dev:
	devspace dev -n validator-plugin-network-system

# Static Analysis / CI

chartCrds = chart/validator-plugin-network/crds/validation.spectrocloud.labs_networkvalidators.yaml

reviewable-ext:
	rm $(chartCrds)
	cp config/crd/bases/validation.spectrocloud.labs_networkvalidators.yaml $(chartCrds)