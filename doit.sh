 export VERSION=0.0.66
# microk8s kubectl delete -f config/samples/app.dac.nokia.com_v1alpha1_smlevo.yaml -n sml-evo
 make docker-build  docker-push bundle
 microk8s kubectl create ns sml-evo
 microk8s kubectl create ns sml-evo-system
# microk8s kubectl apply -f config/samples/app.dac.nokia.com_v1alpha1_smlevo.yaml  -n sml-evo


