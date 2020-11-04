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

package io.clouditor.assurance;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.data_access_layer.PersistentObject;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;
import javax.persistence.*;
import javax.validation.Valid;
import javax.validation.constraints.Size;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.hibernate.annotations.LazyCollection;
import org.hibernate.annotations.LazyCollectionOption;

@Entity(name = "certification")
@Table(name = "certification")
public class Certification implements PersistentObject<String>, Serializable {

  private static final long serialVersionUID = 5983205960445678160L;

  /** A unique identifier for each certification, such as CSA CCM or Azure CIS. */
  @Column(name = "certification_id", nullable = false)
  @Id
  private String id;

  /** A list of controls in the certificate */
  @Size(min = 1)
  @Valid
  @JsonProperty
  @ManyToMany(targetEntity = Control.class)
  @LazyCollection(LazyCollectionOption.FALSE)
  @JoinTable(
      name = "control_to_certification",
      joinColumns =
          @JoinColumn(
              name = "certification_id",
              referencedColumnName = "certification_id",
              nullable = false),
      inverseJoinColumns = @JoinColumn(name = "control_id", referencedColumnName = "control_id"))
  private List<Control> controls = new ArrayList<>();

  @Column(name = "certification_description")
  private String description;

  @Column(name = "publisher")
  private String publisher;

  @Column(name = "website")
  private String website;

  // TODO: startDate, endDate

  public List<Control> getControls() {
    return controls;
  }

  public void setControls(@Size(min = 1) @Valid List<Control> controls) {
    this.controls = controls;
  }

  @Override
  public String getId() {
    return id;
  }

  public void setId(String id) {
    this.id = id;
  }

  public String getDescription() {
    return description;
  }

  public void setDescription(String description) {
    this.description = description;
  }

  public String getPublisher() {
    return publisher;
  }

  public void setPublisher(String publisher) {
    this.publisher = publisher;
  }

  public String getWebsite() {
    return website;
  }

  public void setWebsite(String website) {
    this.website = website;
  }

  @Override
  public String toString() {
    return "Certification{"
        + "id='"
        + id
        + '\''
        + ", controls="
        + controls
        + ", description='"
        + description
        + '\''
        + ", publisher='"
        + publisher
        + '\''
        + ", website='"
        + website
        + '\''
        + '}';
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    Certification that = (Certification) o;

    return new EqualsBuilder()
        .append(getId(), that.getId())
        .append(new ArrayList<>(getControls()), new ArrayList<>(that.getControls()))
        .append(getDescription(), that.getDescription())
        .append(getPublisher(), that.getPublisher())
        .append(getWebsite(), that.getWebsite())
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(getId())
        .append(getControls())
        .append(getDescription())
        .append(getPublisher())
        .append(getWebsite())
        .toHashCode();
  }
}
