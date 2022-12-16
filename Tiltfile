allow_k8s_contexts('kind-epg-exporter')
k8s_yaml('deploy/kubernetes/local/bundle.spilo.yaml')
k8s_resource('spilo', port_forwards='8008:8008')
k8s_resource('spilo', port_forwards='5432:5432')