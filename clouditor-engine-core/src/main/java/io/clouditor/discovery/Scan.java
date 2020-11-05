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

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.assurance.ccl.AssetType;
import io.clouditor.data_access_layer.PersistentObject;
import java.lang.reflect.InvocationTargetException;
import javax.persistence.*;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

/**
 * A {@link Scan} holds information and configuration about a scan that is regularly executed. The
 * actual "scanning" is done by an implementing {@link Scanner} class.
 */
@Entity(name = "scan")
@Table(name = "scan")
public class Scan implements PersistentObject<String> {

  static final String FIELD_SCANNER_CLASS = "scannerClass";

  private static final long DEFAULT_INTERVAL = 5 * 60L;
  private static final long serialVersionUID = 4612570095809897261L;

  /** The associated {@link Scanner} class. */
  @JsonProperty(FIELD_SCANNER_CLASS)
  @Column(name = "scanner_class")
  @Convert(converter = ScannerConverter.class)
  private Class<? extends Scanner> scannerClass;

  /**
   * The asset type, this scan is targeting. This is automatically parsed from the {@link
   * ScannerInfo}.
   */
  @JsonProperty
  @ManyToOne
  @JoinColumn(name = "type_value", nullable = false)
  @MapKey
  @Id
  private AssetType assetType = new AssetType();

  /**
   * The group, or cloud provider this scan is belonging to.This is automatically parsed from the
   * {@link ScannerInfo}.
   */
  @JsonProperty
  @Column(name = "scan_group")
  private String group;

  /**
   * The asset icon of the asset, this scan is targeting. This is automatically parsed from the
   * {@link ScannerInfo}.
   */
  @JsonProperty
  @Column(name = "asset_icon")
  private String assetIcon;

  /** The description of the scan. This is automatically parsed from the {@link ScannerInfo}. */
  @JsonProperty
  @Column(name = "scan_description")
  private String description;

  /** The discovery state of the scan. */
  @JsonProperty
  @Column(name = "is_discovering")
  private boolean isDiscovering;

  /**
   * The service this scan is belonging to. This is automatically parsed from the {@link
   * ScannerInfo}.
   */
  @JsonProperty
  @Column(name = "service")
  private String service;

  @ManyToOne
  @JoinColumn(name = "last_result")
  private DiscoveryResult lastResult;

  @Column(name = "enabled")
  private boolean enabled;

  @Column(name = "scan_interval")
  private final long interval = DEFAULT_INTERVAL;

  public Scan() {}

  public static Scan fromScanner(Class<? extends Scanner> clazz) {
    var scan = new Scan();

    var info = clazz.getAnnotation(ScannerInfo.class);

    if (info != null) {
      scan.assetType = new AssetType();
      scan.assetType.setValue(info.assetType());
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

  @Override
  public String getId() {
    // TODO: short asset types are not really unique, in the long run we might need to add group and
    // service as well or create a dedicated asset type class
    return this.assetType.getValue();
  }

  public void setEnabled(boolean enabled) {
    this.enabled = enabled;
  }

  public void setDiscovering(boolean discovering) {
    this.isDiscovering = discovering;
  }

  public AssetType getAssetType() {
    return this.assetType;
  }

  public void setAssetType(AssetType assetType) {
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

  @Override
  public String toString() {
    return "Scan{"
        + "scannerClass="
        + scannerClass
        + ", assetType="
        + assetType
        + ", assetIcon='"
        + assetIcon
        + '\''
        + ", group='"
        + group
        + '\''
        + ", description='"
        + description
        + '\''
        + ", isDiscovering="
        + isDiscovering
        + ", service='"
        + service
        + '\''
        + ", lastResult="
        + lastResult
        + ", enabled="
        + enabled
        + ", interval="
        + interval
        + '}';
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    Scan scan = (Scan) o;

    return new EqualsBuilder()
        .append(isDiscovering(), scan.isDiscovering())
        .append(isEnabled(), scan.isEnabled())
        .append(getInterval(), scan.getInterval())
        .append(getScannerClass(), scan.getScannerClass())
        .append(getAssetType(), scan.getAssetType())
        .append(getGroup(), scan.getGroup())
        .append(assetIcon, scan.assetIcon)
        .append(description, scan.description)
        .append(getService(), scan.getService())
        .append(getLastResult(), scan.getLastResult())
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(getScannerClass())
        .append(getAssetType())
        .append(getGroup())
        .append(assetIcon)
        .append(description)
        .append(isDiscovering())
        .append(getService())
        .append(getLastResult())
        .append(isEnabled())
        .append(getInterval())
        .toHashCode();
  }

  @Converter
  private static class ScannerConverter
      implements AttributeConverter<Class<? extends Scanner>, String> {

    @Override
    public String convertToDatabaseColumn(final Class<? extends Scanner> attribute) {
      if (attribute == null) return "";
      else return attribute.getCanonicalName();
    }

    @Override
    public Class<? extends Scanner> convertToEntityAttribute(final String dbData) {
      Class<? extends Scanner> resultValue;
      try {
        if (dbData.equals("")) resultValue = null;
        else // noinspection unchecked
        resultValue = (Class<? extends Scanner>) Class.forName(dbData);
      } catch (ClassNotFoundException e) {
        e.printStackTrace();
        throw new IllegalStateException();
      }
      return resultValue;
    }
  }
}
