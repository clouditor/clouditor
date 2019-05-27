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
