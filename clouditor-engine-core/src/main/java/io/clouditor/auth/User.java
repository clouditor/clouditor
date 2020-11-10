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

package io.clouditor.auth;

import static io.clouditor.auth.AuthenticationService.ROLE_GUEST;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.data_access_layer.PersistentObject;
import java.security.Principal;
import java.util.ArrayList;
import java.util.List;
import javax.persistence.*;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.hibernate.annotations.LazyCollection;
import org.hibernate.annotations.LazyCollectionOption;

@Entity
@Table(name = "user")
public class User implements Principal, PersistentObject<String> {

  private static final long serialVersionUID = -1503934816997542987L;

  @JsonProperty
  @Id
  @Column(name = "user_name", nullable = false)
  private String username;

  @JsonProperty
  @Column(name = "password")
  private String password;

  @JsonProperty
  @Column(name = "full_name")
  private String fullName;

  @JsonProperty
  @Column(name = "email")
  private String email;

  @JsonProperty
  @Column(name = "shadow")
  private boolean shadow = false;

  /** The roles of this users. Defaults to {@link AuthenticationService#ROLE_GUEST}. */
  @JsonProperty
  @ElementCollection(targetClass = String.class)
  @CollectionTable(name = "role", joinColumns = @JoinColumn(name = "user_name"))
  @LazyCollection(LazyCollectionOption.FALSE)
  @Column(name = "role_name")
  private List<String> roles = new ArrayList<>(List.of(ROLE_GUEST));

  public User() {}

  public User(String username) {
    this.username = username;
  }

  public User(String username, String password) {
    this.username = username;
    this.password = password;
  }

  public String getUsername() {
    return username;
  }

  public void setUsername(String username) {
    this.username = username;
  }

  public String getPassword() {
    return password;
  }

  public void setPassword(String password) {
    this.password = password;
  }

  @Override
  @JsonIgnore
  public String getName() {
    return this.username;
  }

  @Override
  public String getId() {
    return this.username;
  }

  public List<String> getRoles() {
    return this.roles;
  }

  public void setRoles(List<String> roles) {
    this.roles = roles;
  }

  public boolean hasRole(String role) {
    return this.roles.contains(role);
  }

  public boolean isShadow() {
    return shadow;
  }

  public void setShadow(boolean shadow) {
    this.shadow = shadow;
  }

  public String getFullName() {
    return fullName;
  }

  public void setEmail(String email) {
    this.email = email;
  }

  public void setFullName(String fullName) {
    this.fullName = fullName;
  }

  public String getEmail() {
    return email;
  }

  @Override
  public String toString() {
    return new ToStringBuilder(this)
        .append("username", username)
        .append("fullName", fullName)
        .append("email", email)
        .append("shadow", shadow)
        .append("roles", roles)
        .toString();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    User user = (User) o;

    return new EqualsBuilder()
        .append(shadow, user.shadow)
        .append(username, user.username)
        .append(fullName, user.fullName)
        .append(email, user.email)
        .append(new ArrayList<>(roles), new ArrayList<>(user.roles))
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(username)
        .append(fullName)
        .append(email)
        .append(shadow)
        .append(roles)
        .toHashCode();
  }
}
