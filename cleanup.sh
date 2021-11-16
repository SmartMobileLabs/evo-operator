#!/bin/bash
	kubectl -n sml-evo get smlevoes.app.dac.nokia.com
	kubectl -n sml-evo delete smlevoes.app.dac.nokia.com  sml-evo-1
	kubectl -n sml-evo patch smlevoes.app.dac.nokia.com sml-evo-1 -p '{"metadata": {"finalizers": null}}' --type=merge
	# kubectl -n sml-evo patch resourcerequest.ops.dac.nokia.com/resource-for-sml-evo -p '{"metadata": {"finalizers": null}}' --type=merge
	kubectl -n sml-evo get smlevoes.app.dac.nokia.com
	kubectl delete ns sml-evo
	kubectl api-resources --verbs=list --namespaced -o name | xargs -n 1 microk8s kubectl get --show-kind --ignore-not-found -A | grep sml
	kubectl create ns sml-evo


