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

package io.clouditor.oauth;

import io.clouditor.Engine;
import io.clouditor.auth.AuthenticationService;
import io.clouditor.auth.LoginResponse;
import io.clouditor.auth.User;
import io.clouditor.util.PersistenceManager;
import java.net.URI;
import java.util.Base64;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.QueryParam;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.ClientRequestFilter;
import javax.ws.rs.client.ClientResponseFilter;
import javax.ws.rs.core.NewCookie;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.UriBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Path("")
public class OAuthResource {

  private final AuthenticationService service;

  private final Engine engine;

  private static final Logger LOGGER = LoggerFactory.getLogger(OAuthResource.class);

  @Inject
  public OAuthResource(AuthenticationService service, Engine engine) {
    this.service = service;
    this.engine = engine;
  }

  @GET
  @Path("profile")
  public Profile getProfile() {
    var profile = new Profile();
    profile.setEnabled(
        this.engine.getOAuthClientId() != null
            && this.engine.getOAuthClientSecret() != null
            && this.engine.getOAuthAuthUrl() != null
            && this.engine.getOAuthTokenUrl() != null
            && this.engine.getOAuthJwtSecret() != null);

    return profile;
  }

  @GET
  @Path("callback")
  public Response callback(@QueryParam("code") String code) {
    var token = retrieveAccessToken(code);

    var user = token.decode(this.engine.getOAuthJwtSecret(), this.engine.getoAuthJwtIssuer());

    if (user == null) {
      // redirect back to the beginning
      return buildRedirect();
    }

    LOGGER.info("Decoded access token as user {}", user.getUsername());

    // try to find the user
    var ref = PersistenceManager.getInstance().getById(User.class, user.getId());

    if (ref == null) {
      service.createUser(user);
    }

    // issue token for our API
    var payload = new LoginResponse();

    payload.setToken(service.createToken(user));

    // TODO: max age, etc.
    /* angular is particular about the hash! it needs to be included.
    we cannot use UriBuilder, since it will remove the hash */
    var uri = URI.create("/#?token=" + payload.getToken());

    return Response.temporaryRedirect(uri)
        .cookie(new NewCookie("authorization", payload.getToken()))
        .build();
  }

  private TokenResponse retrieveAccessToken(String code) {
    var tokenUrl = this.engine.getOAuthTokenUrl();

    LOGGER.info("Trying to retrieve access token from {}", tokenUrl);

    var uri =
        UriBuilder.fromUri(tokenUrl)
            .queryParam("grant_type", "authorization_code")
            .queryParam("code", code)
            .build();

    var client =
        ClientBuilder.newClient()
            .register(
                (ClientRequestFilter)
                    requestContext -> {
                      var headers = requestContext.getHeaders();
                      headers.add(
                          "Authorization",
                          "Basic "
                              + Base64.getEncoder()
                                  .encodeToString(
                                      (this.engine.getOAuthClientId()
                                              + ":"
                                              + this.engine.getOAuthClientSecret())
                                          .getBytes()));
                    })
            .register(
                (ClientResponseFilter)
                    (requestContext, responseContext) -> {
                      // stupid workaround because some oauth servers incorrectly sends two
                      // Content-Type
                      // headers! fix it!
                      var contentType = responseContext.getHeaders().getFirst("Content-Type");
                      responseContext.getHeaders().putSingle("Content-Type", contentType);
                    });

    return client.target(uri).request().post(null, TokenResponse.class);
  }

  private Response buildRedirect() {
    var nonce = 25;

    var uri =
        UriBuilder.fromUri(this.engine.getOAuthAuthUrl())
            .queryParam("redirect_uri", this.engine.getBaseUrl() + "/oauth2/callback")
            .queryParam("client_id", this.engine.getOAuthClientId())
            .queryParam("response_type", "code")
            .queryParam("scope", "openid email full_name")
            .queryParam("nonce", nonce)
            .build();

    return Response.temporaryRedirect(uri).build();
  }

  @GET
  @Path("login")
  public Response login() {
    return buildRedirect();
  }
}
