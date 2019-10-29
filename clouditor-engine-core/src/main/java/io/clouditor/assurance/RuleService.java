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

package io.clouditor.assurance;

import static io.clouditor.assurance.RuleService.RuleVisitor.renderText;

import io.clouditor.assurance.ccl.CCLDeserializer;
import io.clouditor.discovery.AssetService;
import io.clouditor.discovery.DiscoveryResult;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.events.DiscoveryResultSubscriber;
import io.clouditor.util.FileSystemManager;
import io.clouditor.util.PersistenceManager;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.stream.Collectors;
import javax.inject.Inject;
import org.commonmark.node.AbstractVisitor;
import org.commonmark.node.BulletList;
import org.commonmark.node.FencedCodeBlock;
import org.commonmark.node.Heading;
import org.commonmark.node.Node;
import org.commonmark.parser.Parser;
import org.commonmark.renderer.html.HtmlRenderer;
import org.commonmark.renderer.text.TextContentRenderer;
import org.jvnet.hk2.annotations.Service;

/** The {@link RuleService} keeps a map of all rules and their associated asset type. */
@Service
public class RuleService extends DiscoveryResultSubscriber {

  @Inject private AssetService assetService;
  @Inject private CertificationService certificationService;
  @Inject private DiscoveryService discoveryService;

  private Map<String, Set<Rule>> rules = new HashMap<>();

  private RuleService() {
    LOGGER.info("Initializing {}...", this.getClass().getSimpleName());
  }

  public void load(Path path) {
    try {
      if (Files.isDirectory(path)) {
        LOGGER.info("Looping through path {}...", path);

        try (var stream = Files.list(path)) {
          for (var p : stream.collect(Collectors.toList())) {
            if (Files.isDirectory(path) || p.endsWith(".md")) {
              this.load(p);
            }
          }
        }
      } else {
        var rule = loadRule(path);

        if (rule.getAssetType() != null) {
          // make sure, the set exists
          this.rules.putIfAbsent(rule.getAssetType(), new HashSet<>());

          // getScanners the set to update it
          var set = this.rules.get(rule.getAssetType());

          set.add(rule);

          LOGGER.info("Added rule {} for asset type {}", path.getFileName(), rule.getAssetType());
        }
      }
    } catch (IOException ex) {
      LOGGER.error("Could not load file or path {}: {}", path, ex.getMessage());
    }
  }

  public RuleEvaluation getStatus(Rule rule) {
    var status = new RuleEvaluation(rule);

    var assets = this.assetService.getAssetsWithType(rule.getAssetType());

    for (var asset : assets) {
      var o =
          asset.getEvaluationResults().stream()
              .filter(result -> Objects.equals(result.getRule().getId(), rule.getId()))
              .findAny();

      if (o.isPresent() && o.get().isOk()) {
        status.addCompliant(asset);
      } else {
        status.addNonCompliant(asset);
      }
    }

    return status;
  }

  public Rule getWithId(String id) {
    return rules.entrySet().stream()
        .flatMap(e -> e.getValue().stream().filter(rule -> Objects.equals(rule.getId(), id)))
        .findFirst()
        .orElse(null);
  }

  public static class ControlsVisitor extends AbstractVisitor {

    private Rule rule;

    ControlsVisitor(Rule rule) {
      this.rule = rule;
    }

    @Override
    public void visit(BulletList bulletList) {
      var node = bulletList.getFirstChild();

      this.rule.getControls().add(renderText(node.getFirstChild()));

      while (node.getNext() != null) {
        node = node.getNext();

        this.rule.getControls().add(renderText(node.getFirstChild()));
      }
    }
  }

  public static class RuleVisitor extends AbstractVisitor {

    private Rule rule;

    RuleVisitor(Rule rule) {
      this.rule = rule;
    }

    @Override
    public void visit(Heading heading) {
      var title = renderText(heading).trim();

      switch (heading.getLevel()) {
        case 1:
          // update the name
          this.rule.setName(title);

          // update the description
          var next = heading.getNext();

          // if the next one is a Heading, then we have no description
          if (!(next instanceof Heading)) {
            this.rule.setDescription(renderHTML(next));
          }

          // additionally, visit all children of the heading to parse code blocks
          this.visitChildren(heading);

          break;
        case 2:
          if (title.equals("Controls")) {
            var node = heading.getNext();

            node.accept(new ControlsVisitor(this.rule));
          }
          break;
        default:
          super.visit(heading);
      }
    }

    @Override
    public void visit(FencedCodeBlock fencedCodeBlock) {
      if (fencedCodeBlock.getInfo().equals("ccl")) {
        var code = fencedCodeBlock.getLiteral();
        // parse CCL
        // for now, one line is one condition, later me might have more complex statements
        var lines = code.split("\\n");

        for (var line : lines) {
          // TODO: this should actually be handled by the grammar
          if (!line.startsWith("#")) {
            var condition = new CCLDeserializer().parse(line);

            this.rule.addCondition(condition);
          }
        }
      } else {
        super.visit(fencedCodeBlock);
      }
    }

    private static String renderHTML(Node node) {
      var renderer = HtmlRenderer.builder().build();

      return renderer.render(node);
    }

    static String renderText(Node node) {
      var renderer = TextContentRenderer.builder().build();

      return renderer.render(node);
    }
  }

  /**
   * Loads a rule from a Markdown-style document specified by the path.
   *
   * @param path The path to the Markdown file
   * @return a parsed {@link Rule}
   * @throws IOException if a parsing error occurred
   */
  public Rule loadRule(Path path) throws IOException {
    var rule = new Rule();

    LOGGER.info("Trying to load rule from {}...", path.getFileName());

    var parser = Parser.builder().build();

    var doc = parser.parseReader(new InputStreamReader(Files.newInputStream(path)));

    doc.accept(new RuleVisitor(rule));

    rule.setId(
        path.getParent().getParent().getFileName()
            + "-"
            + path.getParent().getFileName()
            + "-"
            + path.getFileName().toString().split("\\.")[0]);

    rule.setActive(true);

    return rule;
  }

  public Set<Rule> get(String assetType) {
    return this.rules.getOrDefault(assetType, Collections.emptySet()).stream()
        .filter(Rule::isActive)
        .collect(Collectors.toSet());
  }

  @Override
  public void handle(DiscoveryResult result) {
    LOGGER.info("Handling scan result from {}", result.getTimestamp());

    for (var asset : result.getDiscoveredAssets().values()) {
      // find rules for the asset
      var assetType = asset.getType();

      var rulesForAsset = this.get(assetType);

      // update it regardless of rules, so even an asset with an empty rule set gets recognized
      // TODO: pub/sub?
      assetService.update(asset);

      LOGGER.debug("Evaluating {} rules for asset {}", rulesForAsset.size(), asset.getId());

      // evaluate all rules
      rulesForAsset.forEach(
          rule -> {
            EvaluationResult eval;
            if (!rule.evaluateApplicability(asset)) {
              // simply add an empty EvaluationResult
              eval = new EvaluationResult(rule, asset.getProperties());
            } else {
              eval = rule.evaluate(asset);
            }
            // TODO: can we really update the asset?
            asset.addEvaluationResult(eval);
            if (!eval.isOk()) {
              LOGGER.info(
                  "The following rule failed for asset {} with {} failed conditions: {}",
                  asset.getId(),
                  eval.getFailedConditions().size(),
                  rule.getName());
            }

            // update asset
            assetService.update(asset);
          });
    }

    // now all assets should be evaluated, now we can update the certification
    // TODO: would be nice to only update relevant controls
    this.certificationService.updateCertification();

    // update the scanner with latest result
    var scan = this.discoveryService.getScan(result.getScanId());

    if (scan != null) {
      scan.setLastResult(result);

      PersistenceManager.getInstance().persist(scan);
    }
  }

  public Map<String, Set<Rule>> getRules() {
    return this.rules;
  }

  public List<Rule> getRulesForControl(String controlId) {
    return this.rules.values().stream()
        .flatMap(Collection::stream)
        .filter(rule -> rule.getControls() != null && rule.getControls().contains(controlId))
        .collect(Collectors.toList());
  }

  public void loadAll() {
    try {
      load(FileSystemManager.getInstance().getPathForResource("rules/aws"));
      load(FileSystemManager.getInstance().getPathForResource("rules/azure"));
    } catch (IOException e) {
      LOGGER.error("Could not load rules", e);
    }

    LOGGER.info("Loaded {} rules", this.rules.size());
  }
}
