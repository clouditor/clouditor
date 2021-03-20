package io.clouditor.credentials

import com.auth0.jwt.JWT
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
import io.kubernetes.client.util.credentials.KubeconfigAuthentication
import io.kubernetes.client.ApiClient
import io.kubernetes.client.util.ClientBuilder
import io.kubernetes.client.util.Config

@JsonTypeName(value = "Kubernetes")
class KubernetesAccount : CloudAccount<ClientBuilder?>() {
    private val url: String? = null
    private val token: String? = null
    private val caCertificate: String? = null
    @Throws(IOException::class)
    override fun validate() {
        val builder = resolveCredentials()
        val api = CoreV1Api(builder.build())
        try {
            val list = api.listNamespace(null, null, null, null, null, null, null, null)
            list.items
            if (builder.authentication is AccessTokenAuthentication) {
                // seems to work, lets decode the JWT to set username
                val decoded = JWT.decode(token)
                val claim = decoded.getClaim("kubernetes.io/serviceaccount/service-account.name")
                if (!claim.isNull) {
                    setUser(claim.asString())
                }
            } else if (builder.authentication is KubeconfigAuthentication) {
                val kubeConfig = builder.authentication

                // TODO: maybe extract user from client cert, but it is hidden in a private field
            }
            setAccountId(api.apiClient.basePath)
            LOGGER.info("Account {} validated with user {}.", accountId, user)
        } catch (e: ApiException) {
            throw IOException(e)
        }
    }

    @Throws(IOException::class)
    override fun resolveCredentials(): ClientBuilder {
        return if (this.isAutoDiscovered) {
            ClientBuilder.standard()
        } else ClientBuilder()
            .setBasePath(url)
            .setAuthentication(AccessTokenAuthentication(token))
            .setVerifyingSsl(true)
            .setCertificateAuthority(caCertificate!!.toByteArray())
    }

    companion object {
        /**
         * Discovers an Kubernetes account.
         *
         * @return null, if no account was discovered. Otherwise the discovered [KubernetesAccount].
         */
        fun discover(): KubernetesAccount? {
            return try {
                val account = KubernetesAccount()

                // use the default client config
                val client = Config.defaultClient()
                val api = CoreV1Api(client)
                val list = api.listNamespace(null, null, null, null, null, null, null, null)
                list.items
                account.isAutoDiscovered = true
                account.setAccountId(api.apiClient.basePath)
                /*account.setUser(identity.arn());*/account
            } catch (ex: IOException) {
                // TODO: log error, etc.
                null
            } catch (ex: ApiException) {
                null
            }
        }
    }
}