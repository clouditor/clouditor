package io.clouditor.discovery.k8s;

import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import io.kubernetes.client.ApiException;
import io.kubernetes.client.apis.NetworkingV1beta1Api;
import io.kubernetes.client.models.NetworkingV1beta1Ingress;
import java.util.List;

@ScannerInfo(assetType = "Ingress", group = "Kubernetes", service = "Networking")
public class KubernetesIngressScanner extends KubernetesScanner<NetworkingV1beta1Ingress> {

  public KubernetesIngressScanner() {
    super(x -> x.getMetadata().getUid(), x -> x.getMetadata().getName());
  }

  @Override
  protected List<NetworkingV1beta1Ingress> list() throws ScanException {
    NetworkingV1beta1Api nw = new NetworkingV1beta1Api(this.api.getApiClient());
    try {
      return nw.listIngressForAllNamespaces(null, null, null, null, null, null, null, null)
          .getItems();
    } catch (ApiException e) {
      throw new ScanException(e);
    }
  }
}
