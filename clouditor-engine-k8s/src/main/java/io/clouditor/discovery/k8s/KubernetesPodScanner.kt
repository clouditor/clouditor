package io.clouditor.discovery.k8s

import io.clouditor.discovery.ScanException
import io.clouditor.discovery.ScannerInfo
import io.kubernetes.client.ApiException
import io.kubernetes.client.models.V1Pod
import java.util.function.Function

@ScannerInfo(assetType = "Pod", group = "Kubernetes", service = "Compute")
class KubernetesPodScanner :
    KubernetesScanner<V1Pod>(Function { it.metadata.uid }, Function { it.metadata.name }) {
    @Throws(ScanException::class)

    override fun list(): List<V1Pod> {
        return try {
            val list = api!!.listPodForAllNamespaces(null, null, null, null, null, null, null, null)
            list.items
        } catch (e: ApiException) {
            throw ScanException(e)
        }
    }
}