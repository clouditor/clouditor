package io.clouditor.assurance;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.discovery.Asset;
import java.util.HashMap;
import java.util.Map;
import javax.validation.constraints.NotNull;

public class RuleEvaluation {

  /** The rule. */
  @JsonProperty @NotNull private Rule rule;

  @JsonProperty private Map<String, Boolean> compliance = new HashMap<>();

  @JsonCreator
  public RuleEvaluation(@JsonProperty("rule") Rule rule) {
    this.rule = rule;
  }

  void addCompliant(@NotNull Asset asset) {
    this.compliance.put(asset.getId(), true);
  }

  void addNonCompliant(@NotNull Asset asset) {
    this.compliance.put(asset.getId(), false);
  }

  public boolean isOk() {
    return this.compliance.values().stream().allMatch(status -> status);
  }

  @JsonProperty
  public long getNumberOfCompliant() {
    return this.compliance.values().stream().filter(status -> status).count();
  }

  @JsonProperty
  public long getNumberOfNonCompliant() {
    return this.compliance.values().stream().filter(status -> !status).count();
  }
}
