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
import java.util.ArrayList;
import java.util.List;
import javax.persistence.*;
import javax.validation.Valid;
import javax.validation.constraints.Size;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.hibernate.annotations.LazyCollection;
import org.hibernate.annotations.LazyCollectionOption;

@Entity(name = "certification")
@Table(name = "certification")
public class Certification implements PersistentObject<String> {

  private static final long serialVersionUID = 5983205960445678160L;

  /** A unique identifier for each certification, such as CSA CCM or Azure CIS. */
  @JsonProperty
  @Column(name = "certification_id", nullable = false)
  @Id
  private String id;

  /** A list of controls in the certificate */
  @Size(min = 1)
  @Valid
  @JsonProperty
  @ManyToMany(targetEntity = Control.class, cascade = CascadeType.ALL)
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

  @JsonProperty
  @Column(name = "certification_description", length = 65535)
  private String description;

  @JsonProperty
  @Column(name = "publisher")
  private String publisher;

  @JsonProperty
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
    return new ToStringBuilder(this)
        .append("id", id)
        .append("controls", controls)
        .append("description", description)
        .append("publisher", publisher)
        .append("website", website)
        .toString();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    Certification that = (Certification) o;

    return new EqualsBuilder()
        .append(id, that.id)
        .append(new ArrayList<>(controls), new ArrayList<>(that.controls))
        .append(description, that.description)
        .append(publisher, that.publisher)
        .append(website, that.website)
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(id)
        .append(controls)
        .append(description)
        .append(publisher)
        .append(website)
        .toHashCode();
  }
}
