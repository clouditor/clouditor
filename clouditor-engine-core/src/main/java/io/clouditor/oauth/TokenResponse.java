package io.clouditor.oauth;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.auth.User;

public class TokenResponse {

  @JsonProperty("access_token")
  private String accessToken;

  @JsonProperty("token_type")
  private String tokenType;

  @JsonProperty("expires_in")
  private int expiresIn;

  @JsonProperty("id_token")
  private String idToken;

  User decode(String secret, String issuer) {
    var algorithm = Algorithm.HMAC512(secret);

    var verifier = JWT.require(algorithm).withIssuer(issuer).build();
    var jwt = verifier.verify(this.idToken);

    var user = new User();
    user.setShadow(true);
    user.setUsername(jwt.getSubject());

    return user;
  }
}
