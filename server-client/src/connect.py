#!/usr/bin/env python3
from kubernetes import client, config

config.kube_config.load_kube_config(config_file="kubeconfig/config.yaml")


class Kubernetes:
    def __init__(self):
        self.corev1 = client.CoreV1Api()
        self.custom_objects = client.CustomObjectsApi()

    def ListNameSpace(self):
        data = []
        for ns in self.corev1.list_namespace().items:
            data.append(ns)
        return data

    def CreateNameSpace(self, name):
        body = client.V1Namespace()
        body.metadata = client.V1ObjectMeta(name=name)
        return self.corev1.create_namespace(body=body)

    def ListDevice(self):
        # /apis/{group}/{version}/namespaces/{namespace}/{plural}
        return self.custom_objects.list_cluster_custom_object("devices.kubeedge.io", "v1alpha1", "devices")

    def GetDevice(self, name, namespace="default"):
        # /apis/devices.kubeedge.io/v1alpha1/namespaces/default/devices/temperature-and-humidity
        # /apis/{group}/{version}/namespaces/{namespace}/{plural}/{name}
        return self.custom_objects.get_namespaced_custom_object(group="devices.kubeedge.io", version="v1alpha1", namespace=namespace, plural="devices", name=name)


k = Kubernetes()
print(k.GetDevice(name="temperature-and-humidity"))
