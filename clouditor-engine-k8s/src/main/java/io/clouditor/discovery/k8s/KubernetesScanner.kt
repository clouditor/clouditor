package io.clouditor.discovery.k8s

import io.kubernetes.client.apis.CoreV1Api
import kotlin.Throws
import java.io.IOException
import io.kubernetes.client.models.V1Probe
import io.clouditor.discovery.MixInIgnore
import io.clouditor.credentials.KubernetesAccount
import io.clouditor.data_access_layer.HibernatePersistence
import io.clouditor.discovery.ScannerInfo
import io.clouditor.discovery.k8s.KubernetesScanner
import io.kubernetes.client.models.V1Pod
import io.clouditor.discovery.ScanException
import io.kubernetes.client.models.V1PodList
import io.kubernetes.client.ApiException
import io.kubernetes.client.models.NetworkingV1beta1Ingress
import io.kubernetes.client.apis.NetworkingV1beta1Api
import com.fasterxml.jackson.annotation.JsonTypeName
import io.clouditor.credentials.CloudAccount
import io.kubernetes.client.models.V1NamespaceList
import io.kubernetes.client.util.credentials.AccessTokenAuthentication
import com.auth0.jwt.interfaces.DecodedJWT
import com.auth0.jwt.interfaces.Claim
import io.clouditor.discovery.Scanner
import io.kubernetes.client.util.credentials.KubeconfigAuthentication
import io.kubernetes.client.ApiClient
import io.kubernetes.client.custom.IntOrString
import java.util.function.Function

abstract class KubernetesScanner<T>(idGenerator: Function<T, String>?, nameGenerator: Function<T, String>?) :
    Scanner<CoreV1Api?, T>(null, idGenerator, nameGenerator) {

    @Throws(IOException::class)
    override fun init() {
        super.init()

        MAPPER.addMixIn(V1Probe::class.java, MixInIgnore::class.java)
        // TODO: provide a custom serializer for this instead of ignoring it
        MAPPER.addMixIn(IntOrString::class.java, MixInIgnore::class.java)

        val account = HibernatePersistence()[KubernetesAccount::class.java, "Kubernetes"]
            .orElseThrow { IOException("Kubernetes account not configured") }
            ?: throw IOException("Kubernetes account not configured")

        val builder = account.resolveCredentials()

        api = CoreV1Api(builder.build())
    }
}