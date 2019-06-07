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
import io.clouditor.util.PersistenceManager;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.function.Consumer;
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
    var users = PersistenceManager.getInstance().count(User.class);

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

    var hasher =
        jargon2Hasher()
            .type(Type.ARGON2id)
            .memoryCost(65536)
            .timeCost(3)
            .parallelism(6)
            .saltLength(16)
            .hashLength(16);

    user.setRoles(List.of(ROLE_ADMIN, ROLE_USER));
    user.setPassword(hasher.password(engine.getDefaultApiPassword().getBytes()).encodedHash());

    PersistenceManager.getInstance().persist(user);

    LOGGER.info("Created default user {}.", user.getUsername());
  }

  public String createToken(String subject) {
    Algorithm algorithm = Algorithm.HMAC256(this.engine.getApiSecret());

    return JWT.create()
        .withIssuer(ISSUER)
        .withSubject(subject)
        .withExpiresAt(Date.from(Instant.now().plus(1, ChronoUnit.DAYS)))
        .sign(algorithm);
  }

  public User verifyToken(String token) {
    try {
      Algorithm algorithm = Algorithm.HMAC256(this.engine.getApiSecret());

      JWTVerifier verifier =
          JWT.require(algorithm).withIssuer(ISSUER).build(); // Reusable verifier instance
      DecodedJWT jwt = verifier.verify(token);

      var user = PersistenceManager.getInstance().getById(User.class, jwt.getSubject());

      if (user == null) {
        throw new NotAuthorizedException(ERROR_MESSAGE_USER_NOT_FOUND);
      }

      return user;
    } catch (JWTVerificationException ex) {
      throw new NotAuthorizedException("Invalid token", ex);
    }
  }

  public boolean verifyLogin(LoginRequest request) {
    // fetch user from database
    var reference = PersistenceManager.getInstance().getById(User.class, request.getUsername());

    if (reference == null) {
      return false;
    }

    return jargon2Verifier()
        .hash(reference.getPassword())
        .password(request.getPassword().getBytes())
        .verifyEncoded();
  }

  public List<User> getUsers() {
    var users = new ArrayList<User>();

    PersistenceManager.getInstance().find(User.class).forEach((Consumer<? super User>) users::add);

    return users;
  }
}
