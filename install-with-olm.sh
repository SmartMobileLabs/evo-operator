#!/bin/bash
FOO=$0
if [ "$(basename $FOO)" == *"uninstall"* ]; then 
    microk8s kubectl delete -f og.yaml -n sml-evo-system
    microk8s kubectl delete serviceaccount sml-evo-controller-manager -n sml-evo-system
    microk8s kubectl delete -f additionals/og.yaml -n sml-evo-system
    microk8s kubectl delete -f additionals/role.yaml -n sml-evo-system 
    microk8s kubectl delete -f additionals/role_binding.yaml -n sml-evo-system
    microk8s kubectl delete -f bundle/manifests/app.dac.nokia.com.example.com_smlevoes.yaml -n sml-evo-system
    microk8s kubectl delete -f bundle/manifests/sml-evo-controller-manager-metrics-service_v1_service.yaml -n sml-evo-system
    microk8s kubectl delete -f bundle/manifests/sml-evo-manager-config_v1_configmap.yaml -n sml-evo-system
    microk8s kubectl delete -f bundle/manifests/sml-evo-metrics-reader_rbac.authorization.k8s.io_v1_clusterrole.yaml -n sml-evo-system
    microk8s kubectl delete -f bundle/manifests/sml-evo.clusterserviceversion.yaml -n sml-evo-system
    microk8s kubectl delete ns sml-evo-system 
    microk8s kubectl delete ns sml-evo
elif [[ "$(basename $FOO)" == *"install"* ]]; then
    microk8s kubectl create ns sml-evo-system 
    microk8s kubectl create ns sml-evo
    microk8s kubectl create -f additionals/og.yaml -n sml-evo-system
    microk8s kubectl create -f additionals/role.yaml -n sml-evo-system
    microk8s kubectl create -f additionals/role_binding.yaml -n sml-evo-system
    microk8s kubectl create serviceaccount sml-evo-controller-manager -n sml-evo-system
    microk8s kubectl create -f bundle/manifests/app.dac.nokia.com.example.com_smlevoes.yaml -n sml-evo-system
    microk8s kubectl create -f bundle/manifests/sml-evo-controller-manager-metrics-service_v1_service.yaml -n sml-evo-system
    microk8s kubectl create -f bundle/manifests/sml-evo-manager-config_v1_configmap.yaml -n sml-evo-system
    microk8s kubectl create -f bundle/manifests/sml-evo-metrics-reader_rbac.authorization.k8s.io_v1_clusterrole.yaml -n sml-evo-system
    microk8s kubectl create -f bundle/manifests/sml-evo.clusterserviceversion.yaml -n sml-evo-system
fi
