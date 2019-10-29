package io.clouditor.discovery.k8s;

import io.clouditor.credentials.KubernetesAccount;
import io.clouditor.discovery.MixInIgnore;
import io.clouditor.discovery.Scanner;
import io.clouditor.util.PersistenceManager;
import io.kubernetes.client.apis.CoreV1Api;
import java.io.IOException;
import java.util.function.Function;

public abstract class KubernetesScanner<T> extends Scanner<CoreV1Api, T> {

  public KubernetesScanner(Function<T, String> idGenerator, Function<T, String> nameGenerator) {
    super(null, idGenerator, nameGenerator);
  }

  @Override
  public void init() throws IOException {
    super.init();

    MAPPER.addMixIn(io.kubernetes.client.models.V1Probe.class, MixInIgnore.class);
    // TODO: provide a custom serializer for this instead of ignoring it
    MAPPER.addMixIn(io.kubernetes.client.custom.IntOrString.class, MixInIgnore.class);

    var account = PersistenceManager.getInstance().getById(KubernetesAccount.class, "Kubernetes");

    if (account == null) {
      throw new IOException("Kubernetes account not configured");
    }

    var builder = account.resolveCredentials();

    this.api = new CoreV1Api(builder.build());
  }
}
