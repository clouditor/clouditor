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

import static io.clouditor.auth.AuthenticationService.ROLE_ADMIN;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.clouditor.AbstractEngineUnitTest;
import io.clouditor.data_access_layer.HibernatePersistence;
import java.util.List;
import javax.ws.rs.NotAuthorizedException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

class UserTest extends AbstractEngineUnitTest {

  @Override
  @BeforeEach
  protected void setUp() {
    super.setUp();
    this.engine.setDBName("UserTestDB");
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

    new HibernatePersistence().saveOrUpdate(user);

    var token = service.createToken(user);

    assertNotNull(token);

    var decodedUser = service.verifyToken(token);

    assertEquals(user, decodedUser);

    var ctx = new UserContext(decodedUser, false);

    assertEquals(user, ctx.getUserPrincipal());

    assertTrue(ctx.isUserInRole(ROLE_ADMIN));

    // remove the user again, otherwise, other tests that may be based on the existence of
    // a clean user DB might fail

    new HibernatePersistence().delete(user);
  }

  @Test
  void testUserNotFound() {
    var service = this.engine.getService(AuthenticationService.class);

    var token = service.createToken(new User("maybe-existed-before-but-not-anymore"));

    // the token itself is valid, but verification should fail since the user is not in the DB
    assertThrows(NotAuthorizedException.class, () -> service.verifyToken(token));
  }
}
