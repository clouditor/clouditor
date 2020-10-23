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

package io.clouditor.discovery;

import static com.mongodb.client.model.Filters.eq;

import io.clouditor.events.DiscoveryResultSubscriber;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Modifier;
import java.util.Collection;
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
import java.util.stream.StreamSupport;

import io.clouditor.util.PersistenceManager;
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
  private final Map<String, ScheduledFuture<?>> futures = new HashMap<>();
  private final Map<String, Scanner> scanners = new HashMap<>();

  private final ScheduledThreadPoolExecutor scheduler =
      (ScheduledThreadPoolExecutor) Executors.newScheduledThreadPool(1);

  private final SubmissionPublisher<DiscoveryResult> assetPublisher = new SubmissionPublisher<>();
  private final Map<String, HashSet<EventOutput>> subscriptions = new HashMap<>();

  public DiscoveryService() {
    LOGGER.info("Initializing {}", this.getClass().getSimpleName());
  }

  public void init() {
    var scans = PersistenceManager.getInstance().find(Scan.class);

    // first, init list of scanner classes from Java via reflections
    var classes =
        REFLECTIONS_SUBTYPE_SCANNER.getSubTypesOf(Scanner.class).stream()
            .filter(
                clazz -> !Modifier.isAbstract(clazz.getModifiers()) && !clazz.isAnonymousClass())
            .collect(Collectors.toList());

    // loop through existing scans to make sure that the associated scanner still exists
    for (var scan : scans) {
      if (!classes.contains(scan.getScannerClass())) {
        LOGGER.info(
            "Scan {} contains old or invalid scanner class {}. Removing entry from database.",
            scan.getId(),
            scan.getScannerClass());

        PersistenceManager.getInstance().delete(Scan.class, scan.getId());
      }
    }

    // loop through all Scanner classes, to make sure a Scan object exists for each class
    for (var clazz : classes) {
      var scan =
          PersistenceManager.getInstance()
              .find(Scan.class, eq(Scan.FIELD_SCANNER_CLASS, clazz.getName()))
              .limit(1)
              .first();

      if (scan == null) {
        // create new scanner object
        scan = Scan.fromScanner(clazz);

        // update database
        PersistenceManager.getInstance().persist(scan);
      }
    }
  }

  public Map<String, Scan> getScans() {
    return StreamSupport.stream(
            PersistenceManager.getInstance().find(Scan.class).spliterator(), false)
        .collect(Collectors.toMap(Scan::getId, scan -> scan));
  }

  public void start() {
    var scans = PersistenceManager.getInstance().find(Scan.class);

    // loop through all enabled scans and start them
    for (var scan : scans) {
      if (scan.isEnabled()) {
        this.startScan(scan);
      }
    }
  }

  private void startScan(Scan scan) {
    LOGGER.info("Starting scan {}", scan.getId());

    try {
      // create the associated scanner object, that handles the actual scanning
      var scanner = scan.instantiateScanner();

      // check, if it is somehow already running, and cancel it
      var future = this.futures.get(scan.getId());
      if (future != null) {
        LOGGER.info("It seems this scan is already running, cancelling previous execution.");
        future.cancel(true);
      }

      // increase thread pool size
      int size =
          Math.min(this.scheduler.getCorePoolSize() + 1, this.scheduler.getMaximumPoolSize());

      this.scheduler.setCorePoolSize(size);
      LOGGER.info("Adjusting thread pool size to {}", this.scheduler.getCorePoolSize());

      LOGGER.info("Starting scan {}. Now waiting for next execution", scan.getId());

      future =
          scheduler.scheduleAtFixedRate(
              () -> {
                // set discovering flag to enable its indication in the discovery view
                scan.setDiscovering(true);

                // scan
                var result = scanner.scan(scan.getId());

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

                // update database
                PersistenceManager.getInstance().persist(scan);
              },
              0,
              scan.getInterval(),
              TimeUnit.SECONDS);

      // store the future, so we can cancel it later
      this.futures.put(scan.getId(), future);

      // store the scanner, so we can access it later
      this.scanners.put(scan.getId(), scanner);
    } catch (InstantiationException
        | IllegalAccessException
        | InvocationTargetException
        | NoSuchMethodException e) {
      LOGGER.error("Cannot instantiate scanner from {}: {}", scan.getId(), e.getMessage());

      // disable it again
      disableScan(scan);
    }
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

    // decrease thread pool size
    int size = Math.max(this.scheduler.getCorePoolSize() - 1, 0);

    this.scheduler.setCorePoolSize(size);
    LOGGER.info("Adjusting thread pool size to {}", this.scheduler.getCorePoolSize());

    // clean up associated objects
    this.futures.remove(scan.getId());
    this.scanners.remove(scan.getId());
  }

  public void subscribe(DiscoveryResultSubscriber subscriber) {
    this.assetPublisher.subscribe(subscriber);
  }

  public Scan getScan(String id) {
    return PersistenceManager.getInstance().getById(Scan.class, id);
  }

  public void enableScan(Scan scan) {
    scan.setEnabled(true);

    // update database
    PersistenceManager.getInstance().persist(scan);

    this.startScan(scan);
  }

  public void disableScan(Scan scan) {
    scan.setEnabled(false);

    // update database
    PersistenceManager.getInstance().persist(scan);

    this.stopScan(scan);
  }

  public EventOutput subscribeToEvents(String assetType) {
    // create new event output
    var output = new EventOutput();

    // add to the subscription map for given asset type

    // make sure a hash set exists
    var set = this.subscriptions.putIfAbsent(assetType, new HashSet<>());
    Objects.requireNonNull(set).add(output);

    LOGGER.info("Subscribed an SSE client to asset type {}", assetType);

    return output;
  }

  /** Returns a list of currently running scanners. */
  public Collection<Scanner> getScanners() {
    return this.scanners.values();
  }
}
