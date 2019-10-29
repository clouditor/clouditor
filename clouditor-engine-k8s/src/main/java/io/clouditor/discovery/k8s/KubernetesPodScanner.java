package io.clouditor.discovery.k8s;

import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import io.kubernetes.client.ApiException;
import io.kubernetes.client.models.V1Pod;
import java.util.List;

@ScannerInfo(assetType = "Pod", group = "Kubernetes", service = "Compute")
public class KubernetesPodScanner extends KubernetesScanner<V1Pod> {

  public KubernetesPodScanner() {
    super(x -> x.getMetadata().getUid(), x -> x.getMetadata().getName());
  }

  @Override
  protected List<V1Pod> list() throws ScanException {
    try {
      var list = this.api.listPodForAllNamespaces(null, null, null, null, null, null, null, null);

      return list.getItems();
    } catch (ApiException e) {
      throw new ScanException(e);
    }
  }
}
