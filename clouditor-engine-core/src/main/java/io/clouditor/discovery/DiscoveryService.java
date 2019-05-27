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

package io.clouditor.discovery;

import io.clouditor.events.DiscoveryResultSubscriber;
import io.clouditor.util.PersistenceManager;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Modifier;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.ScheduledThreadPoolExecutor;
import java.util.concurrent.SubmissionPublisher;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.glassfish.jersey.media.sse.EventOutput;
import org.jvnet.hk2.annotations.Service;
import org.reflections.Reflections;
import org.reflections.scanners.SubTypesScanner;
import org.reflections.util.ClasspathHelper;
import org.reflections.util.ConfigurationBuilder;

@Service
public class DiscoveryService {

  private static final Logger LOGGER = LogManager.getLogger();

  private static final Reflections REFLECTIONS_SUBTYPE_SCANNER =
      new Reflections(
          new ConfigurationBuilder()
              .addUrls(ClasspathHelper.forPackage(Scanner.class.getPackage().getName()))
              .setScanners(new SubTypesScanner()));

  private Map<String, Scan> scans = new HashMap<>();
  private Map<String, ScheduledFuture<?>> futures = new HashMap<>();

  private final ScheduledThreadPoolExecutor scheduler =
      (ScheduledThreadPoolExecutor) Executors.newScheduledThreadPool(1);

  private SubmissionPublisher<DiscoveryResult> assetPublisher = new SubmissionPublisher<>();
  private Map<String, HashSet<EventOutput>> subscriptions = new HashMap<>();

  public DiscoveryService() {
    LOGGER.info("Initializing {}", this.getClass().getSimpleName());
  }

  public void load() {
    // TODO: load from database

    // first, load checks from Java via reflections
    var classes =
        REFLECTIONS_SUBTYPE_SCANNER.getSubTypesOf(Scanner.class).stream()
            .filter(
                clazz -> !Modifier.isAbstract(clazz.getModifiers()) && !clazz.isAnonymousClass())
            .collect(Collectors.toList());

    // loop through all Scanner classes, to make sure a Scan object exists for each class
    for (var clazz : classes) {
      try {
        var constructor = clazz.getDeclaredConstructor();
        constructor.setAccessible(true);
        var scan = Scan.fromScanner(constructor.newInstance());

        // update database
        PersistenceManager.getInstance().persist(scan);

        // TODO: do not overwrite custom settings from DB (maybe load db later?)
        this.scans.put(scan.getId(), scan);
      } catch (InstantiationException
          | IllegalAccessException
          | InvocationTargetException
          | NoSuchMethodException e) {
        LOGGER.error("Ignoring instantiate scanner class {}: {}", clazz.getName(), e);
        continue;
      }
    }

    LOGGER.info("Loaded {} scans", this.scans.size());

    // adjust the thread pool
    this.scheduler.setCorePoolSize(this.scans.size());
    LOGGER.info("Adjusting thread pool size to {}", this.scheduler.getCorePoolSize());
  }

  public Map<String, Scan> getScans() {
    return this.scans;
  }

  public void start() {
    // loop through all enabled scans and start them
    for (var scan : this.scans.values()) {
      if (scan.isEnabled()) {
        this.startScan(scan);
      }
    }
  }

  private void startScan(Scan scan) {
    var scanner = scan.getScanner();

    // check, if it is somehow already running, and cancel it
    var future = this.futures.get(scan.getId());
    if (future != null) {
      LOGGER.info("It seems this scan is already running, cancelling previous execution.");
      future.cancel(true);
    }

    LOGGER.info("Starting scan {}. Now waiting for next execution", scan.getId());

    // store the future, so we can cancel it later
    future =
        scheduler.scheduleAtFixedRate(
            () -> {
              // set discovering flag to enable its indication in the discovery view
              scan.setDiscovering(true);

              // scan
              var result = scanner.scan();

              submit(scan, result);

              // TODO: route this through pub/sub
              /*var subscribers = this.subscriptions.get(scanner.getId());
              for (Iterator<EventOutput> iterator = subscribers.iterator(); iterator.hasNext(); ) {
                var subscriber = iterator.next();
                var event =
                    new OutboundEvent.Builder()
                        .mediaType(MediaType.APPLICATION_JSON_TYPE)
                        .name("discovery-complete")
                        .data(DiscoveryResult.class, result)
                        .build();

                try {
                  if (!subscriber.isClosed()) {
                    subscriber.write(event);
                  }
                  if (subscriber.isClosed()) {
                    LOGGER.debug("Removing " + subscriber + " for type " + scanner.getId() + "...");
                    iterator.remove();
                  }

                } catch (IOException e) {
                  LOGGER.error("Could not send event {} to subscriber", event.getName());
                }
              }*/

              scan.setLastResult(result);

              LOGGER.info("Scan {} is now waiting for next execution.", scan.getId());
              scan.setDiscovering(false);
            },
            0,
            scan.getInterval(),
            TimeUnit.SECONDS);

    this.futures.put(scanner.getId(), future);
  }

  public int submit(Scan scan, DiscoveryResult result) {
    var assets = result.getDiscoveredAssets();

    LOGGER.info(
        "Publishing discovery result with {} asset(s) of type {}.",
        assets.size(),
        scan.getAssetType());

    return this.assetPublisher.submit(result);
  }

  private void stopScan(Scan scan) {
    // look for a future
    var future = this.futures.get(scan.getId());
    if (future == null) {
      LOGGER.info("It seems this scan is not running, no need to stop it.");
      return;
    }

    future.cancel(true);

    LOGGER.info("Stopped scan {}", scan.getId());
  }

  public void subscribe(DiscoveryResultSubscriber subscriber) {
    this.assetPublisher.subscribe(subscriber);
  }

  public Scan getScan(String id) {
    return this.scans.get(id);
  }

  public void enableScan(Scan scan) {
    scan.setEnabled(true);

    this.startScan(scan);
  }

  public void disableScan(Scan scan) {
    scan.setEnabled(false);

    this.stopScan(scan);
  }

  public EventOutput subscribeToEvents(String assetType) {
    // create new event output
    var output = new EventOutput();

    // add to the subscription map for given asset type

    // make sure a hashset exists
    var set = this.subscriptions.putIfAbsent(assetType, new HashSet<>());
    Objects.requireNonNull(set).add(output);

    LOGGER.info("Subscribed an SSE client to asset type {}", assetType);

    return output;
  }
}
