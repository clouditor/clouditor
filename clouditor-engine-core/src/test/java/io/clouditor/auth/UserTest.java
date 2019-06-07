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

import static io.clouditor.auth.AuthenticationService.ROLE_ADMIN;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.clouditor.AbstractEngineUnitTest;
import io.clouditor.util.PersistenceManager;
import java.util.List;
import javax.ws.rs.NotAuthorizedException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

class UserTest extends AbstractEngineUnitTest {

  @Override
  @BeforeEach
  protected void setUp() {
    super.setUp();

    this.engine.initDB();
  }

  @Test
  void testEquals() {
    var user = new User();
    user.setUsername("clouditor");

    // compare with self
    assertEquals(user, user);

    // compare with null
    assertNotEquals(user, null);

    // compare with other
    assertNotEquals(user, new User());

    // compare with wrong class
    assertNotEquals(user, new Object());
  }

  @Test
  void testVerifyAuthentication() {
    var service = this.engine.getService(AuthenticationService.class);

    var user = new User("user", "mypass");
    user.setRoles(List.of(ROLE_ADMIN));

    PersistenceManager.getInstance().persist(user);

    var token = service.createToken(user.getUsername());

    assertNotNull(token);

    var decodedUser = service.verifyToken(token);

    assertEquals(user, decodedUser);

    var ctx = new UserContext(decodedUser, false);

    assertEquals(user, ctx.getUserPrincipal());

    assertTrue(ctx.isUserInRole(ROLE_ADMIN));
  }

  @Test
  void testUserNotFound() {
    var service = this.engine.getService(AuthenticationService.class);

    var token = service.createToken("maybe-existed-before-but-not-anymore");

    // the token itself is valid, but verification should fail since the user is not in the DB
    assertThrows(NotAuthorizedException.class, () -> service.verifyToken(token));
  }
}
