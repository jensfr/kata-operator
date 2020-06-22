# kata-operator

## Installation

Of course you should read deploy.sh before executing it on your machine. 
When you trust it does the right things, perform the following two steps to deploy
the kata operator on your cluster.

  curl https://raw.githubusercontent.com/jensfr/master/deploy/deploy.sh | bash
  oc get kataconfig example-kataconfig

This will create roles, role bindings and a service account as well as deploy the kata-operator.

## Selecting nodes for the installation

To select nodes where kata is to be installed you need to add labels to those node and then specify
it in deploy/crds/kataconfiguration.openshift.io_v1alpha1_kataconfig_cr.yaml

For example:
  spec:
  kataConfigPoolSelector:
    matchLabels:
       custom-kata1: test
       
To start the installation you need to create the Custom Resource by doing

  oc apply -f deploy/crds/kataconfiguration.openshift.io_v1alpha1_kataconfig_cr.yaml
