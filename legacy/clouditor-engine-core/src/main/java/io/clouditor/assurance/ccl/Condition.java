/*
 * Copyright 2016-2019 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *            $$\                           $$\ $$\   $$\
 *            $$ |                          $$ |\__|  $$ |
 *   $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 *  $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 *  $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ |  \__|
 *  $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 *  \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *   \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package io.clouditor.assurance.ccl;

import io.clouditor.discovery.AssetProperties;
import java.io.Serializable;
import javax.persistence.*;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

@Embeddable
public class Condition implements Serializable {

  private static final long serialVersionUID = -7530930851665073637L;

  @Column private String source;

  @ManyToOne
  @JoinColumn(name = "type_value")
  private AssetType assetType;

  @Transient private Expression expression;

  public Expression getExpression() {
    return this.expression;
  }

  public void setExpression(Expression expression) {
    this.expression = expression;
  }

  public AssetType getAssetType() {
    return this.assetType;
  }

  public void setAssetType(AssetType assetType) {
    this.assetType = assetType;
  }

  public boolean evaluate(AssetProperties properties) {
    return this.expression.evaluate(properties);
  }

  public void setSource(String source) {
    this.source = source;
  }

  public String getSource() {
    return this.source;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    Condition condition = (Condition) o;

    return new EqualsBuilder()
        .append(source, condition.source)
        .append(assetType, condition.assetType)
        .append(expression, condition.expression)
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(source)
        .append(assetType)
        .append(expression)
        .toHashCode();
  }
}
