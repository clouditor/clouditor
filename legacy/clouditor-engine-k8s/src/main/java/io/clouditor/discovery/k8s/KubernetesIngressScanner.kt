package io.clouditor.discovery.k8s

import io.clouditor.discovery.ScanException
import io.clouditor.discovery.ScannerInfo
import io.kubernetes.client.ApiException
import io.kubernetes.client.apis.NetworkingV1beta1Api
import io.kubernetes.client.models.NetworkingV1beta1Ingress
import java.util.function.Function

@ScannerInfo(assetType = "Ingress", group = "Kubernetes", service = "Networking")
class KubernetesIngressScanner :
    KubernetesScanner<NetworkingV1beta1Ingress>(Function { it.metadata.uid },
        Function { it.metadata.name }) {

    @Throws(ScanException::class)
    override fun list(): List<NetworkingV1beta1Ingress> {
        val nw = NetworkingV1beta1Api(api!!.apiClient)

        return try {
            nw.listIngressForAllNamespaces(null, null, null, null, null, null, null, null)
                .items
        } catch (e: ApiException) {
            throw ScanException(e)
        }
    }
}