package io.clouditor.assurance.ccl;

import com.fasterxml.jackson.annotation.JsonTypeInfo;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "@class")
public abstract class AssetType {

  public abstract boolean evaluate(Map properties);

  public abstract String getField();

  Object getValueFromField(Map asset, String fieldName) {
    // first, try to resolve it directly
    if (asset.containsKey(fieldName)) {
      return asset.get(fieldName);
    }

    // split it at .
    var names = fieldName.split("\\.");

    var base = asset;
    var value = (Object) null;

    for (var name : names) {
      value = base.get(name);

      if (value instanceof Map) {
        base = (Map) value;
      }
    }

    return value;
  }
}
