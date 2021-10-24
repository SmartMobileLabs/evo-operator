 export VERSION=0.0.0.48
 microk8s kubectl delete -f config/samples/app.dac.nokia.com_v1alpha1_smlevo.yaml -n sml-evo
 make docker-build  docker-push
 make deploy
 sleep 20
 microk8s kubectl apply -f config/samples/app.dac.nokia.com_v1alpha1_smlevo.yaml  -n sml-evo


