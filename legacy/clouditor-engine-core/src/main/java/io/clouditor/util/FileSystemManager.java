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

package io.clouditor.util;

import io.clouditor.Component;
import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.nio.file.FileSystem;
import java.nio.file.FileSystemNotFoundException;
import java.nio.file.FileSystems;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collections;
import javax.validation.constraints.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * A small utility package for handling different file systems.
 *
 * @author Banse, Christian
 */
public final class FileSystemManager {

  private static final String KEYS_DIR = "keys/";

  private static final Logger LOGGER = LoggerFactory.getLogger(FileSystemManager.class);

  private static FileSystemManager instance;

  private FileSystem jarFileSystem;

  private FileSystemManager() {}

  /**
   * Returns the current instance of {@link FileSystemManager}.
   *
   * @return the current instance
   */
  public static synchronized FileSystemManager getInstance() {
    if (instance == null) {
      instance = new FileSystemManager();
    }
    return instance;
  }

  /** Cleans up created file systems, such as the {@link FileSystem} for JARs. */
  public void cleanup() {
    try {
      if (jarFileSystem != null) {
        jarFileSystem.close();
      }
    } catch (IOException e) {
      LOGGER.error("Error while trying to close the jarFileSystem: {}", e);
    }
  }

  public Path getPathForResource(@NotNull String resource) throws IOException {
    return this.getPathForResource(resource, Component.class);
  }

  /**
   * Returns a {@link Path} for a given path to a resource.
   *
   * @param resource the resource path
   * @param clazz the class to get the resource from
   * @return a {@link Path} to the resource
   * @throws IOException if the resource was not found or another error occured
   */
  public Path getPathForResource(@NotNull String resource, Class<?> clazz) throws IOException {
    URL url = Component.class.getClassLoader().getResource(resource);
    // try directly with class
    if (url == null) {
      // then try with class loader, i.e. for tests
      url = clazz.getResource(resource);
      // if we cannot create an URL, directly try Paths.get()
      if (url == null) {
        return Paths.get(resource);
      }
    }

    URI uri;
    try {
      uri = url.toURI();
    } catch (URISyntaxException ex) {
      throw new IOException(ex);
    }

    /*
     * we need to "register" the specific file system first before we
     * can use Paths.get()
     */
    if ("jar".equals(uri.getScheme())) {
      // this will just register the handler so we can use it for Paths
      try {
        FileSystems.getFileSystem(uri);
      } catch (FileSystemNotFoundException e) {
        jarFileSystem = FileSystems.newFileSystem(uri, Collections.singletonMap("create", true));
        // just to be safe, ignore if we somehow accidentally register it twice
        LOGGER.info("Creating jar file system for URI {}", uri);
      }
    }

    return Paths.get(uri);
  }

  /**
   * Returns the content of a private key file as an array of characters.
   *
   * @param keyName the path to the key file
   * @throws IOException if the key file is not found
   */
  public void getPrivateKey(String keyName) throws IOException {

    /*
     * this is still not a very good approach to handle this because
     * this way the SSH key file always needs to be bundled with the
     * test case, which does not make too much sense
     */
    final Path path = this.getPathForResource(KEYS_DIR + keyName);

    Files.readString(path);
  }
}
