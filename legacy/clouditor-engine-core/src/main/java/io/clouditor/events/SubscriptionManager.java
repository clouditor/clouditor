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

package io.clouditor.events;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Modifier;
import java.util.HashSet;
import java.util.Set;
import java.util.stream.Collectors;
import org.reflections.Reflections;
import org.reflections.scanners.SubTypesScanner;
import org.reflections.util.ClasspathHelper;
import org.reflections.util.ConfigurationBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SubscriptionManager {

  private static final Reflections REFLECTIONS_SUBTYPE_SCANNER =
      new Reflections(
          new ConfigurationBuilder()
              .addUrls(ClasspathHelper.forPackage(Subscriber.class.getPackage().getName()))
              .setScanners(new SubTypesScanner()));

  private static final Logger LOGGER = LoggerFactory.getLogger(SubscriptionManager.class);

  private static SubscriptionManager instance;

  public static synchronized SubscriptionManager getInstance() {
    if (instance == null) {
      instance = new SubscriptionManager();
    }
    return instance;
  }

  public <T extends Subscriber> Set<T> loadSubscribers(Class<T> subscriberClass) {
    var subscribers = new HashSet<T>();

    // first, load subscription classes from Java via reflections
    var clazzes =
        REFLECTIONS_SUBTYPE_SCANNER.getSubTypesOf(subscriberClass).stream()
            .filter(
                clazz -> !Modifier.isAbstract(clazz.getModifiers()) && !clazz.isAnonymousClass())
            .collect(Collectors.toList());

    for (var clazz : clazzes) {
      try {
        var subscriber = clazz.getDeclaredConstructor().newInstance();

        subscribers.add(subscriber);
      } catch (InstantiationException
          | IllegalAccessException
          | InvocationTargetException
          | NoSuchMethodException e) {
        LOGGER.error("Unable to load subscriber from class {}", clazz.toString());
      }
    }

    return subscribers;
  }
}
