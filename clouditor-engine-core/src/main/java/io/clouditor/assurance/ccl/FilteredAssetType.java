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

import java.util.Map;
import java.util.Objects;
import javax.persistence.*;

@Entity(name = "filtered_asset_type")
@Table(name = "filtered_asset_type")
@PrimaryKeyJoinColumn(name = "type_value")
public class FilteredAssetType extends AssetType {

  private static final long serialVersionUID = -2355408351894740425L;

  public FilteredAssetType() {
    super();
  }

  @Transient private Expression assetExpression;

  public void setAssetExpression(Expression assetExpression) {
    this.assetExpression = assetExpression;
  }

  public boolean evaluate(Map properties) {
    return this.getAssetExpression().evaluate(properties);
  }

  @Override
  public String toString() {
    return "TYPE_VALUE: " + super.getValue() + ", EXPRESSION: " + getAssetExpression() + ".";
  }

  public Expression getAssetExpression() {
    return assetExpression;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;
    if (o == null || getClass() != o.getClass()) return false;
    if (!super.equals(o)) return false;
    FilteredAssetType that = (FilteredAssetType) o;
    return Objects.equals(getAssetExpression(), that.getAssetExpression());
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getAssetExpression());
  }
}
