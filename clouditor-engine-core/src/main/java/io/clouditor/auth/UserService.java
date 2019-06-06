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
import java.util.Date;
import javax.inject.Inject;
import javax.ws.rs.NotAuthorizedException;
import org.jvnet.hk2.annotations.Service;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Service
public class UserService {

  private static final Logger LOGGER = LoggerFactory.getLogger(UserService.class);

  private Engine engine;

  @Inject
  public UserService(Engine engine) {
    this.engine = engine;
  }

  public void init() {
    // check, if users exist, otherwise create the default user
    var users = PersistenceManager.getInstance().count(User.class);

    if (users == 0) {
      createDefaultUser();
    }
  }

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

    user.setPassword(hasher.password(engine.getDefaultApiPassword().getBytes()).encodedHash());

    PersistenceManager.getInstance().persist(user);

    LOGGER.info("Created default user {}.", user.getUsername());
  }

  public String createToken(User user) {
    Algorithm algorithm = Algorithm.HMAC256(this.engine.getApiSecret());

    return JWT.create()
        .withIssuer(UserContext.ISSUER)
        .withSubject(user.getName())
        .withExpiresAt(Date.from(Instant.now().plus(1, ChronoUnit.DAYS)))
        .sign(algorithm);
  }

  public User verifyToken(String token) {
    try {
      // for now we use the api password as JWT secret. in the future we need to see if this makes
      // sense
      Algorithm algorithm = Algorithm.HMAC256(this.engine.getApiSecret());

      JWTVerifier verifier =
          JWT.require(algorithm)
              .withIssuer(UserContext.ISSUER)
              .build(); // Reusable verifier instance
      DecodedJWT jwt = verifier.verify(token);

      return new User(jwt.getSubject());
    } catch (JWTVerificationException ex) {
      throw new NotAuthorizedException("Invalid token", ex);
    }
  }

  public boolean verifyUser(User user) {
    // fetch user from database
    var reference = PersistenceManager.getInstance().getById(User.class, user.getId());

    if (reference == null) {
      return false;
    }

    return jargon2Verifier()
        .hash(reference.getPassword())
        .password(user.getPassword().getBytes())
        .verifyEncoded();
  }
}
