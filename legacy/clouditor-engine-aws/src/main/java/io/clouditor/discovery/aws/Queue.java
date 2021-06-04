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

package io.clouditor.discovery.aws;

import com.fasterxml.jackson.annotation.JsonProperty;
import software.amazon.awssdk.utils.builder.CopyableBuilder;
import software.amazon.awssdk.utils.builder.ToCopyableBuilder;

public class Queue implements ToCopyableBuilder<Queue.Builder, Queue> {

  private final String url;

  private Queue(BuilderImpl builder) {
    this.url = builder.url;
  }

  public static Builder builder() {
    return new BuilderImpl();
  }

  @Override
  public Queue.Builder toBuilder() {
    return new BuilderImpl().url(url);
  }

  public String url() {
    return this.url;
  }

  public interface Builder extends CopyableBuilder<Queue.Builder, Queue> {

    Builder url(String url);
  }

  static final class BuilderImpl implements Builder {
    private String url;

    @JsonProperty
    String getUrl() {
      return this.url;
    }

    BuilderImpl() {}

    public Builder url(String url) {
      this.url = url;

      return this;
    }

    @Override
    public Queue build() {
      return new Queue(this);
    }
  }
}
