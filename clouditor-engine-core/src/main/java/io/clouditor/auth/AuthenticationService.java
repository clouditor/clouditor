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

import static com.kosprov.jargon2.api.Jargon2.jargon2Hasher;
import static com.kosprov.jargon2.api.Jargon2.jargon2Verifier;

import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTVerifier;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.exceptions.JWTVerificationException;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.kosprov.jargon2.api.Jargon2.Type;
import io.clouditor.Engine;
import io.clouditor.data_access_layer.HibernatePersistence;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.Date;
import java.util.List;
import javax.inject.Inject;
import javax.ws.rs.NotAuthorizedException;
import org.jvnet.hk2.annotations.Service;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Service
public class AuthenticationService {

  public static final String ISSUER = "clouditor";

  public static final String ROLE_GUEST = "guest";
  public static final String ROLE_USER = "user";
  public static final String ROLE_ADMIN = "admin";

  private static final Logger LOGGER = LoggerFactory.getLogger(AuthenticationService.class);
  public static final String ERROR_MESSAGE_USER_NOT_FOUND = "User does not exist";

  private Engine engine;

  @Inject
  public AuthenticationService(Engine engine) {
    this.engine = engine;
  }

  public void init() {
    // check, if users exist, otherwise create the default user
    var users = new HibernatePersistence().count(User.class);

    if (users == 0) {
      createDefaultUser();
    }
  }

  /**
   * Creates the default admin user according to the credentials configured in the {@link Engine}.
   */
  private void createDefaultUser() {
    var user = new User();
    user.setUsername(engine.getDefaultApiUsername());
    user.setFullName(engine.getDefaultApiUsername());

    user.setRoles(List.of(ROLE_ADMIN, ROLE_USER));
    user.setPassword(hashPassword(engine.getDefaultApiPassword()));

    new HibernatePersistence().saveOrUpdate(user);

    LOGGER.info("Created default user {}.", user.getUsername());
  }

  private String hashPassword(String password) {
    var hasher =
        jargon2Hasher()
            .type(Type.ARGON2id)
            .memoryCost(65536)
            .timeCost(3)
            .parallelism(6)
            .saltLength(16)
            .hashLength(16);

    return hasher.password(password.getBytes()).encodedHash();
  }

  public String createToken(User user) {
    Algorithm algorithm = Algorithm.HMAC256(this.engine.getApiSecret());

    return JWT.create()
        .withIssuer(ISSUER)
        .withSubject(user.getUsername())
        .withClaim("full_name", user.getFullName())
        .withClaim("email", user.getEmail())
        .withExpiresAt(Date.from(Instant.now().plus(1, ChronoUnit.DAYS)))
        .sign(algorithm);
  }

  public User verifyToken(String token) {
    try {
      Algorithm algorithm = Algorithm.HMAC256(this.engine.getApiSecret());

      JWTVerifier verifier =
          JWT.require(algorithm).withIssuer(ISSUER).build(); // Reusable verifier instance
      DecodedJWT jwt = verifier.verify(token);

      return new HibernatePersistence()
          .get(User.class, jwt.getSubject())
          .orElseThrow(() -> new NotAuthorizedException(ERROR_MESSAGE_USER_NOT_FOUND));

    } catch (JWTVerificationException ex) {
      throw new NotAuthorizedException("Invalid token", ex);
    }
  }

  public boolean verifyLogin(LoginRequest request) {
    // fetch user from database
    var referenceOptional = new HibernatePersistence().get(User.class, request.getUsername());

    if (referenceOptional.isEmpty()) {
      return false;
    }

    var reference = referenceOptional.get();

    if (reference.getPassword() == null) {
      return false;
    }

    return jargon2Verifier()
        .hash(reference.getPassword())
        .password(request.getPassword().getBytes())
        .verifyEncoded();
  }

  public List<User> getUsers() {
    return new HibernatePersistence().listAll(User.class);
  }

  /**
   * Creates a new user in the database
   *
   * @param user the {@link User} to be created.
   * @return false, if the user already exists
   */
  public boolean createUser(User user) {
    // check, if user already exists
    var ref = new HibernatePersistence().get(User.class, user.getId());

    if (ref.isPresent()) {
      return false;
    }

    // create the new user
    new HibernatePersistence().saveOrUpdate(user);

    LOGGER.info("Created user {}.", user.getId());

    return true;
  }

  public User getUser(String id) {
    return new HibernatePersistence().get(User.class, id).orElse(null);
  }

  public void updateUser(String id, User user) {
    // fetch existing
    var ref = new HibernatePersistence().get(User.class, id);

    if (ref.isEmpty()) {
      return;
    }

    // make sure, identifiers match
    user.setUsername(id);

    // if password is empty, it means that we do not update it
    if (user.getPassword() == null) {
      user.setPassword(ref.get().getPassword());
    } else {
      // encode hash
      user.setPassword(hashPassword(user.getPassword()));
    }

    // store it
    new HibernatePersistence().saveOrUpdate(user);
  }

  public void deleteUser(String id) {
    // delete it from database
    new HibernatePersistence().delete(User.class, id);
  }
}
