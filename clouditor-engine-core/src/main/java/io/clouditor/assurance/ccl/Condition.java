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

import com.vladmihalcea.hibernate.type.json.JsonBinaryType;
import io.clouditor.discovery.AssetProperties;
import java.io.Serializable;
import javax.persistence.*;
import org.hibernate.annotations.Type;
import org.hibernate.annotations.TypeDef;

@Entity(name = "condition")
@Table(name = "condition")
@TypeDef(name = "jsonb", typeClass = JsonBinaryType.class)
public class Condition {

  @Id @Embedded private final ConditionPK conditionPK = new ConditionPK();

  @Column(name = "expression")
  @Type(type = "jsonb")
  private Expression expression;

  public ConditionPK getConditionPK() {
    return conditionPK;
  }

  public Expression getExpression() {
    return expression;
  }

  public void setExpression(Expression expression) {
    this.expression = expression;
  }

  public AssetType getAssetType() {
    return getConditionPK().getAssetType();
  }

  public void setAssetType(AssetType assetType) {
    getConditionPK().setAssetType(assetType);
  }

  public boolean evaluate(AssetProperties properties) {
    return this.expression.evaluate(properties);
  }

  public void setSource(String source) {
    getConditionPK().setSource(source);
  }

  public String getSource() {
    return getConditionPK().getSource();
  }

  @Embeddable
  static final class ConditionPK implements Serializable {

    private static final long serialVersionUID = -503140484349205605L;

    @Column(name = "source")
    private String source;

    @ManyToOne
    @JoinColumn(name = "type_value")
    private AssetType assetType;

    private String getSource() {
      return source;
    }

    private AssetType getAssetType() {
      return assetType;
    }

    private void setSource(String source) {
      this.source = source;
    }

    private void setAssetType(AssetType assetType) {
      this.assetType = assetType;
    }
  }
}
