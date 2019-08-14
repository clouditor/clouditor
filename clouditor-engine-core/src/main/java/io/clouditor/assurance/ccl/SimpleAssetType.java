package io.clouditor.assurance.ccl;

import java.util.Map;

public class SimpleAsset extends AssetType {

  private String field;

  public String getField() {
    return this.field;
  }

  public void setField(String field) {
    this.field = field;
  }

  @Override
  public boolean evaluate(Map properties) {
    return true;
  }
}
