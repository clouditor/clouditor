/*
 * Copyright (c) 2016-2019, Fraunhofer AISEC. All rights reserved.
 *
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
 *
 * Clouditor Community Edition is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Clouditor Community Edition is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * long with Clouditor Community Edition.  If not, see <https://www.gnu.org/licenses/>
 */

package io.clouditor.auth;

import static io.clouditor.auth.AuthenticationService.ROLE_GUEST;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonView;
import io.clouditor.rest.ObjectMapperResolver.DatabaseOnly;
import io.clouditor.util.PersistentObject;
import java.security.Principal;
import java.util.List;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

public class User implements Principal, PersistentObject<String> {

  private String username;

  @JsonView(DatabaseOnly.class)
  private String password;

  @JsonProperty private boolean shadow = false;

  /** The roles of this users. Defaults to {@link AuthenticationService#ROLE_GUEST}. */
  @JsonProperty private List<String> roles = List.of(ROLE_GUEST);

  public User() {}

  public User(String username) {
    this.username = username;
  }

  public User(String username, String password) {
    this.username = username;
    this.password = password;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }

    if (o == null || getClass() != o.getClass()) {
      return false;
    }

    User user = (User) o;

    // the comparison of a user is just done by the name and attributes, not the password!
    return new EqualsBuilder().append(username, user.username).isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37).append(username).append(password).toHashCode();
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
}
