package io.clouditor.credentials;

import com.auth0.jwt.JWT;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.kubernetes.client.ApiException;
import io.kubernetes.client.apis.CoreV1Api;
import io.kubernetes.client.util.ClientBuilder;
import io.kubernetes.client.util.Config;
import io.kubernetes.client.util.credentials.AccessTokenAuthentication;
import io.kubernetes.client.util.credentials.KubeconfigAuthentication;
import java.io.IOException;

@JsonTypeName(value = "Kubernetes")
public class KubernetesAccount extends CloudAccount<ClientBuilder> {

  private String url;
  private String token;
  private String caCertificate;

  @Override
  public void validate() throws IOException {
    var builder = resolveCredentials();

    var api = new CoreV1Api(builder.build());

    try {
      var list = api.listNamespace(null, null, null, null, null, null, null, null);
      list.getItems();

      if (builder.getAuthentication() instanceof AccessTokenAuthentication) {
        // seems to work, lets decode the JWT to set username
        var decoded = JWT.decode(token);
        var claim = decoded.getClaim("kubernetes.io/serviceaccount/service-account.name");

        if (!claim.isNull()) {
          this.setUser(claim.asString());
        }
      } else if (builder.getAuthentication() instanceof KubeconfigAuthentication) {
        var kubeConfig = builder.getAuthentication();

        // TODO: maybe extract user from client cert, but it is hidden in a private field
      }

      this.setAccountId(api.getApiClient().getBasePath());

      LOGGER.info("Account {} validated with user {}.", this.accountId, this.user);
    } catch (ApiException e) {
      throw new IOException(e);
    }
  }

  /**
   * Discovers an Kubernetes account.
   *
   * @return null, if no account was discovered. Otherwise the discovered {@link KubernetesAccount}.
   */
  public static KubernetesAccount discover() {
    try {
      var account = new KubernetesAccount();

      // use the default client config
      var client = Config.defaultClient();
      var api = new CoreV1Api(client);

      var list = api.listNamespace(null, null, null, null, null, null, null, null);
      list.getItems();

      account.setAutoDiscovered(true);
      account.setAccountId(api.getApiClient().getBasePath());
      /*account.setUser(identity.arn());*/

      return account;
    } catch (IOException | ApiException ex) {
      // TODO: log error, etc.
      return null;
    }
  }

  @Override
  public ClientBuilder resolveCredentials() throws IOException {
    if (this.isAutoDiscovered()) {
      return ClientBuilder.standard();
    }

    return new ClientBuilder()
        .setBasePath(url)
        .setAuthentication(new AccessTokenAuthentication(token))
        .setVerifyingSsl(true)
        .setCertificateAuthority(caCertificate.getBytes());
  }
}
