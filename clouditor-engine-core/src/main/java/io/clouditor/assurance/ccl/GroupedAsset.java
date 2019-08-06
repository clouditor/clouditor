package io.clouditor.assurance.ccl;

import java.util.Map;

public class GroupedAsset extends AssetType {

  private String field;
  private Expression assetExpression;

  public String getField() {
    return this.field;
  }

  public void setField(String field) {
    this.field = field;
  }

  public void setAssetExpression(Expression assetExpression) {
    this.assetExpression = assetExpression;
  }

  @Override
  public boolean evaluate(Map properties) {
    return this.assetExpression.evaluate(properties);
  }
}
