package io.clouditor.data_access_layer;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import io.clouditor.AbstractEngineUnitTest;
import io.clouditor.assurance.*;
import io.clouditor.assurance.ccl.AssetType;
import io.clouditor.assurance.ccl.Condition;
import io.clouditor.assurance.ccl.FilteredAssetType;
import io.clouditor.discovery.*;
import java.time.Instant;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

class PersistenceManagerTest extends AbstractEngineUnitTest {

  @Test
  void testComplexRelations() {
    // arrange
    this.engine.setDBName("PersistenceManagerTestComplexDB");
    this.engine.initDB();

    final PersistenceManager sut = new HibernatePersistence();

    final Domain domain = new Domain("domain_name");
    domain.setDescription("domain description");

    final String controlID = "control_id";
    final Control control = new Control();
    control.setControlId(controlID);
    control.setDescription("Control Description");
    control.setFulfilled(Control.Fulfillment.GOOD);
    control.setDomain(domain);
    control.setName("control name");
    control.setAutomated(true);
    control.setActive(true);

    final String certificationID = "certification_id";
    final Certification certification = new Certification();
    certification.setId(certificationID);
    certification.setDescription("Certification Description");
    certification.setPublisher("Certification Publisher");
    certification.setWebsite("website");
    certification.setControls(List.of(control));

    final String assetTypeID = "asset_type_id";
    final AssetType assetType = new AssetType();
    assetType.setValue(assetTypeID);

    final String filteredAssetTypeID = "filtered_asset_type_ID";
    final FilteredAssetType filteredAssetType = new FilteredAssetType();
    filteredAssetType.setValue(filteredAssetTypeID);

    final Condition condition = new Condition();
    condition.setAssetType(assetType);
    condition.setSource("source");

    final String ruleID = "rule_id";
    final Rule rule = new Rule();
    rule.setId(ruleID);
    rule.setActive(true);
    rule.setName("rule name");
    rule.setDescription("rule description");
    rule.addControls(control.getControlId());
    rule.setConditions(List.of(condition));

    final AssetProperties assetProperties = new AssetProperties();
    assetProperties.put("key", "value");

    final EvaluationResult evaluationResult = new EvaluationResult(rule, assetProperties);
    evaluationResult.setFailedConditions(List.of(condition));
    final String evaluationResultID = evaluationResult.getId();

    final String assetID = "asset_id";
    final Asset asset = new Asset("asset_type", assetID, "asset_name", assetProperties);
    asset.setEvaluationResults(List.of(evaluationResult));

    final Scan scan = Scan.fromScanner(FakeScanner.class);
    scan.setAssetType(assetType.getValue());
    scan.setDiscovering(true);
    scan.setEnabled(true);

    final DiscoveryResult discoveryResult = new DiscoveryResult(scan.getId());
    discoveryResult.setFailed(true);
    discoveryResult.setError("error");
    final Map<String, Asset> discoveredAssets = new HashMap<>();
    discoveredAssets.put(assetID, asset);
    discoveryResult.setDiscoveredAssets(discoveredAssets);
    final Instant discoverResultID = discoveryResult.getTimestamp();

    // act
    sut.saveOrUpdate(control);
    sut.saveOrUpdate(certification);
    sut.saveOrUpdate(assetType);
    sut.saveOrUpdate(filteredAssetType);
    sut.saveOrUpdate(rule);
    sut.saveOrUpdate(evaluationResult);
    sut.saveOrUpdate(asset);
    sut.saveOrUpdate(scan);
    sut.saveOrUpdate(discoveryResult);

    final Control haveControl = sut.get(Control.class, controlID).orElseThrow();
    final Certification haveCertification =
        sut.get(Certification.class, certificationID).orElseThrow();
    final AssetType haveAssetType = sut.get(AssetType.class, assetTypeID).orElseThrow();
    final FilteredAssetType haveFilteredAssetType =
        sut.get(FilteredAssetType.class, filteredAssetTypeID).orElseThrow();
    final Rule haveRule = sut.get(Rule.class, ruleID).orElseThrow();
    final EvaluationResult haveEvaluationResult =
        sut.get(EvaluationResult.class, evaluationResultID).orElseThrow();
    final Asset haveAsset = sut.get(Asset.class, assetID).orElseThrow();
    final Scan haveScan = sut.get(Scan.class, scan.getId()).orElseThrow();
    final DiscoveryResult haveDiscoveryResult =
        sut.get(DiscoveryResult.class, discoverResultID).orElseThrow();

    sut.delete(discoveryResult);
    sut.delete(scan);
    sut.delete(asset);
    sut.delete(evaluationResult);
    sut.delete(rule);
    sut.delete(filteredAssetType);
    sut.delete(assetType);
    sut.delete(certification);
    sut.delete(control);

    // assert
    assertEquals(haveControl, control);
    assertEquals(haveCertification, certification);
    assertEquals(haveAssetType, assetType);
    assertEquals(haveFilteredAssetType, filteredAssetType);
    assertEquals(haveRule, rule);
    assertEquals(haveEvaluationResult, evaluationResult);
    assertEquals(haveAsset, asset);
    assertEquals(haveScan, scan);
    assertEquals(haveDiscoveryResult, discoveryResult);

    assertNotNull(haveRule.getConditions());
    assertEquals(1, haveRule.getConditions().size());
    assertEquals(condition, haveRule.getConditions().get(0));

    assertNotNull(haveEvaluationResult.getFailedConditions());
    assertEquals(1, haveEvaluationResult.getFailedConditions().size());
    assertEquals(condition, haveEvaluationResult.getFailedConditions().get(0));

    assertEquals(domain, haveControl.getDomain());

    this.engine.shutdown();
  }

  @Test
  void initWithFalsePortBooms() {
    // arrange // act
    Assertions.assertThrows(
        IllegalArgumentException.class,
        () -> HibernateUtils.init("host", -1, "dbName", "userName", "password"));
  }

  @Test
  void initWithFalsePortBooms2() {
    // arrange // act
    Assertions.assertThrows(
        IllegalArgumentException.class,
        () -> HibernateUtils.init("host", 999999, "dbName", "userName", "password"));
  }
}
