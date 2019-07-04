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

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.util.PersistentObject;
import java.lang.reflect.InvocationTargetException;

/**
 * A {@link Scan} holds information and configuration about a scan that is regularly executed. The
 * actual "scanning" is done by an implementing {@link Scanner} class.
 */
public class Scan implements PersistentObject<String> {

  static final String FIELD_SCANNER_CLASS = "scannerClass";

  private static final long DEFAULT_INTERVAL = 5 * 60L;

  /** The associated {@link Scanner} class. */
  @JsonProperty(FIELD_SCANNER_CLASS)
  private Class<? extends Scanner> scannerClass;

  /**
   * The asset type, this scan is targeting. This is automatically parsed from the {@link
   * ScannerInfo}.
   */
  @JsonProperty private String assetType;

  /**
   * The asset icon of the asset, this scan is targeting. This is automatically parsed from the
   * {@link ScannerInfo}.
   */
  @JsonProperty private String assetIcon;

  /**
   * The group, or cloud provider this scan is belonging to.This is automatically parsed from the
   * {@link ScannerInfo}.
   */
  @JsonProperty private String group;

  /** The description of the scan. This is automatically parsed from the {@link ScannerInfo}. */
  @JsonProperty private String description;

  /** The discovery state of the scan. */
  @JsonProperty private boolean isDiscovering;

  /**
   * The service this scan is belonging to. This is automatically parsed from the {@link
   * ScannerInfo}.
   */
  @JsonProperty private String service;

  private DiscoveryResult lastResult;

  private boolean enabled;

  private long interval = DEFAULT_INTERVAL;

  public Scan() {}

  public static Scan fromScanner(Class<? extends Scanner> clazz) {
    var scan = new Scan();

    var info = clazz.getAnnotation(ScannerInfo.class);

    if (info != null) {
      scan.assetType = info.assetType();
      scan.assetIcon = info.assetIcon();
      scan.group = info.group();
      scan.service = info.service();
      scan.description = info.description();
    }

    scan.scannerClass = clazz;

    return scan;
  }

  public DiscoveryResult getLastResult() {
    return this.lastResult;
  }

  public void setLastResult(DiscoveryResult lastResult) {
    this.lastResult = lastResult;
  }

  public long getInterval() {
    return interval;
  }

  public boolean isEnabled() {
    return this.enabled;
  }

  public boolean isDiscovering() {
    return this.isDiscovering;
  }

  public String getId() {
    // TODO: short asset types are not really unique, in the long run we might need to add group and
    // service as well or create a dedicated asset type class
    return this.assetType;
  }

  public void setEnabled(boolean enabled) {
    this.enabled = enabled;
  }

  public void setDiscovering(boolean discovering) {
    this.isDiscovering = discovering;
  }

  public String getAssetType() {
    return this.assetType;
  }

  public void setAssetType(String assetType) {
    this.assetType = assetType;
  }

  public String getGroup() {
    return group;
  }

  public String getService() {
    return service;
  }

  public Class<? extends Scanner> getScannerClass() {
    return scannerClass;
  }

  public Scanner instantiateScanner()
      throws NoSuchMethodException, IllegalAccessException, InvocationTargetException,
          InstantiationException {
    var constructor = scannerClass.getConstructor();
    constructor.setAccessible(true);

    return constructor.newInstance();
  }
}
