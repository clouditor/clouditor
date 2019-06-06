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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

import io.clouditor.AbstractEngineUnitTest;
import org.junit.jupiter.api.Test;

class UserTest extends AbstractEngineUnitTest {

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
    var service = this.engine.getService(UserService.class);

    var user = new User("user", "mypass");

    var token = service.createToken(user);

    assertNotNull(token);

    var decodedUser = service.verifyToken(token);

    assertEquals(user, decodedUser);

    var ctx = new UserContext(decodedUser, false);

    assertEquals(user, ctx.getUserPrincipal());

    // no roles implemented yet

    assertFalse(ctx.isUserInRole("fake_role"));
  }
}
