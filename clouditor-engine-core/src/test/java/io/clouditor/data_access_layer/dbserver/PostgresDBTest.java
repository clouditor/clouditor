package io.clouditor.data_access_layer.dbserver;

import com.opentable.db.postgres.embedded.EmbeddedPostgres;
import io.clouditor.data_access_layer.PostgresDB;
import io.clouditor.data_access_layer.utils.BoomingRunnable;
import io.clouditor.data_access_layer.utils.Table;
import io.clouditor.data_access_layer.DBServer;
import io.clouditor.metric_api.pb.*;
import org.junit.jupiter.api.Test;

import java.sql.Connection;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import static org.junit.Assert.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

public class PostgresDBTest {

    private static final DBServer EMPTY_DB;
    private static final DBServer UNINITIATED_EMPTY_DB;

    static {
        try {
            final EmbeddedPostgres pg = EmbeddedPostgres.start();
            final Connection connection = pg.getPostgresDatabase().getConnection();
            EMPTY_DB = new PostgresDB(connection);
            Table.init(connection);

            final EmbeddedPostgres pg2 = EmbeddedPostgres.start();
            final Connection connection2 = pg2.getPostgresDatabase().getConnection();
            UNINITIATED_EMPTY_DB = new PostgresDB(connection2);
            final BoomingRunnable boomingRunnable = () -> {
                EMPTY_DB.close();
                UNINITIATED_EMPTY_DB.close();
            };
            Runtime.getRuntime()
                    .addShutdownHook(
                            new Thread(
                                    boomingRunnable.handleException()
                            )
                );
        } catch (final Exception e) {
            e.printStackTrace();
            throw new AssertionError();
        }
    }

    @Test
    public void testPostgresDBGetEmptyScale() throws Exception {
        // arrange
        final int lowerBound = 0;
        final int upperBound = 1;
        final Optional<Scale> want = Optional.empty();
        // act
        final Optional<Scale> have = EMPTY_DB.loadScale(lowerBound, upperBound);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBStoreScale() throws Exception {
        // arrange
        final int lowerBound = 0;
        final int upperBound = 1;
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("This is a description")
                .build();
        final Optional<Scale> want = Optional.of(scale);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            // act
            sut.storeScale(scale);
            final Optional<Scale> have = sut.loadScale(lowerBound, upperBound);
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBDeleteScale() throws Exception {
        // arrange
        final int lowerBound = 0;
        final int upperBound = 1;
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("This is a description")
                .build();
        final Optional<Scale> want1 = Optional.of(scale);
        final Optional<Scale> want2 = Optional.empty();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            // act
            sut.storeScale(scale);
            final Optional<Scale> have1 = sut.loadScale(lowerBound, upperBound);
            sut.deleteScale(lowerBound, upperBound);
            final Optional<Scale> have2 = sut.loadScale(lowerBound, upperBound);
            // assert
            assertEquals(want1, have1);
            assertEquals(want2, have2);
        }
    }

    @Test
    public void testPostgresDBLoadAllScales() throws Exception {
        // arrange
        final List<Scale> want = new ArrayList<>();
        // act
        final List<Scale> have = EMPTY_DB.loadAllScales();
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadAllScales2() throws Exception {
        // arrange
        final Scale scale1 = Scale.newBuilder()
                .setLowerBound(1)
                .setUpperBound(2)
                .setDescription("first Scale")
                .build();
        final Scale scale2 = Scale.newBuilder()
                .setLowerBound(2)
                .setUpperBound(2)
                .setDescription("second Scale")
                .build();
        final Scale scale3 = Scale.newBuilder()
                .setLowerBound(2)
                .setUpperBound(3)
                .setDescription("third Scale")
                .build();
        final List<Scale> want = new ArrayList<>();
        want.add(scale1);
        want.add(scale2);
        want.add(scale3);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            // act
            sut.storeScale(scale1);
            sut.storeScale(scale2);
            sut.storeScale(scale3);
            final List<Scale> have = sut.loadAllScales();
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBLoadScaleWithoutInitBooms() {
        // arrange
        final Scale scale = Scale.newBuilder()
                .setLowerBound(1)
                .setUpperBound(2)
                .setDescription("first Scale")
                .build();
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.storeScale(scale));
    }

    @Test
    public void testPostgresDBStoreScaleNotMatchingTheRegexBooms() {
        // arrange
        final Scale scale = Scale.newBuilder()
                .setLowerBound(1)
                .setUpperBound(2)
                .setDescription("first Scale;")
                .build();
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.storeScale(scale));
    }

    @Test
    public void testPostgresDBStoreTwiceTheSamePrimaryKeyBooms() throws Exception {
        // arrange
        final Scale scale1 = Scale.newBuilder()
                .setLowerBound(1)
                .setUpperBound(2)
                .setDescription("first Scale")
                .build();
        final Scale scale2 = Scale.newBuilder()
                .setLowerBound(1)
                .setUpperBound(2)
                .setDescription("second Scale")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            // act
            sut.storeScale(scale1);
            assertThrows(SQLException.class, () -> sut.storeScale(scale2));
        }
    }

    @Test
    public void testPostgresDBStoreNullBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.storeScale(null));
    }

    @Test
    public void testPostgresDBDeleteScaleWithoutInitBooms() {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.deleteScale(lowerBound, upperBound));
    }

    @Test
    public void testPostgresDBDeleteNotExistingScaleBooms() {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteScale(lowerBound, upperBound));
    }

    @Test
    public void testPostgresDBLoadAllScalesWithoutInitBooms() {
        // arrange // act
        assertThrows(SQLException.class, UNINITIATED_EMPTY_DB::loadAllScales);
    }

    @Test
    public void testPostgresDBLoadEmptyAsset() throws Exception {
        // arrange
        final String assetName = "asset_name";
        final Optional<Asset> want = Optional.empty();
        // act
        final Optional<Asset> have = EMPTY_DB.loadAsset(assetName);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBStoreAsset() throws Exception {
        // arrange
        final String assetName = "asset_name";
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("This is a asset_type")
                .build();
        final Optional<Asset> want = Optional.of(asset);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeAsset(asset);
            // act
            final Optional<Asset> have = sut.loadAsset(assetName);
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBDeleteAsset() throws Exception {
        // arrange
        final String assetName = "asset_name";
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("This is a asset_type")
                .build();
        final Optional<Asset> want1 = Optional.of(asset);
        final Optional<Asset> want2 = Optional.empty();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeAsset(asset);
            // act
            final Optional<Asset> have1 = sut.loadAsset(assetName);
            sut.deleteAsset(assetName);
            final Optional<Asset> have2 = sut.loadAsset(assetName);
            // assert
            assertEquals(want1, have1);
            assertEquals(want2, have2);
        }
    }
    @Test
    public void testPostgresDBLoadAllAssetsEmpty() throws Exception {
        // arrange
        final List<Asset> want = new ArrayList<>();
        // act
        final List<Asset> have = EMPTY_DB.loadAllAssets();
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadAllAssets() throws Exception {
        // arrange
        final Asset asset1 = Asset.newBuilder()
                .setAssetName("asset1")
                .setAssetType("This is a asset_type")
                .build();
        final Asset asset2 = Asset.newBuilder()
                .setAssetName("asset2")
                .setAssetType("This is a asset_type")
                .build();
        final Asset asset3 = Asset.newBuilder()
                .setAssetName("asset3")
                .setAssetType("This is a asset_type")
                .build();
        final List<Asset> want = new ArrayList<>();
        want.add(asset1);
        want.add(asset2);
        want.add(asset3);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeAsset(asset1);
            sut.storeAsset(asset2);
            sut.storeAsset(asset3);
            // act
            final List<Asset> have = sut.loadAllAssets();
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBLoadAssetNullBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.loadAsset(null));
    }

    @Test
    public void testPostgresDBStoreAssetNullBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.storeAsset(null));
    }

    @Test
    public void testPostgresDBDeleteAssetNullBooms() {
        // arrange
        assertThrows(NullPointerException.class, () -> EMPTY_DB.deleteAsset(null));
    }

    @Test
    public void testPostgresDBLoadAssetWithoutInitBooms() {
        // arrange
        final String assetName = "assetName";
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.loadAsset(assetName));
    }

    @Test
    public void testPostgresDBStoreAssetWithoutInitBooms() {
        // arrange
        final Asset asset = Asset.newBuilder()
                .setAssetName("assetName")
                .setAssetType("assetType")
                .build();
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.storeAsset(asset));
    }

    @Test
    public void testPostgresDBDeleteAssetWithoutInitBooms() {
        // arrange
        final String assetName = "assetName";
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.deleteAsset(assetName));
    }

    @Test
    public void testPostgresDBStoreTwoAssetsWithSamePrimaryKeyBooms() throws Exception {
        // arrange
        final Asset asset1 = Asset.newBuilder()
                .setAssetName("Asset")
                .setAssetType("Type1")
                .build();
        final Asset asset2 = Asset.newBuilder()
                .setAssetName("Asset")
                .setAssetType("Type2")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            // act
            sut.storeAsset(asset1);
            assertThrows(SQLException.class, () -> sut.storeAsset(asset2));
        }
    }

    @Test
    public void testPostgresDBDeleteNotExistingPrimaryKeyBooms() {
        // arrange
        final String assetName = "Asset";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteAsset(assetName));
    }

    @Test
    public void testPostgresDBStoreAssetNameNotMatchingTheRegexBooms() {
        // arrange
        final Asset asset = Asset.newBuilder()
                .setAssetName("Asset;")
                .setAssetType("Type")
                .build();
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.storeAsset(asset));
    }

    @Test
    public void testPostgresDBStoreAssetTypeNotMatchingTheRegexBooms() {
        // arrange
        final Asset asset = Asset.newBuilder()
                .setAssetName("Asset")
                .setAssetType("Type;")
                .build();
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.storeAsset(asset));
    }

    @Test
    public void testPostgresDBLoadAssetNameNotMatchingTheRegexBooms() {
        // arrange
        final String assetName = "Asset;";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.loadAsset(assetName));
    }

    @Test
    public void testPostgresDBDeleteAssetNameNotMatchingTheRegexBooms() {
        // arrange
        final String assetName = "Asset;";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteAsset(assetName));
    }

    @Test
    public void testPostgresDBStoreMetric() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName("MetricName")
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            // act
            sut.storeMetric(metric);
        }
    }

    @Test
    public void testPostgresDBLoadEmptyMetric() throws Exception {
        // arrange
        final String metricName = "MetricName";
        final Optional<Metric> want = Optional.empty();
        // act
        final Optional<Metric> have = EMPTY_DB.loadMetric(metricName);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadMetric() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final String metricName = "MetricName";
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Optional<Metric> want = Optional.of(metric);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            // act
            final Optional<Metric> have = sut.loadMetric(metricName);
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBDeleteMetric() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final String metricName = "MetricName";
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Optional<Metric> want1 = Optional.of(metric);
        final Optional<Metric> want2 = Optional.empty();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            // act
            final Optional<Metric> have1 = sut.loadMetric(metricName);
            sut.deleteMetric(metricName);
            final Optional<Metric> have2 = sut.loadMetric(metricName);
            // assert
            assertEquals(want1, have1);
            assertEquals(want2, have2);
        }
    }

    @Test
    public void testPostgresDBLoadAllMetrics() throws Exception {
        // arrange
        final List<Metric> want = new ArrayList<>();
        // act
        final List<Metric> have = EMPTY_DB.loadAllMetrics();
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadAllMetrics2() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric1 = Metric.newBuilder()
                .setMetricName("MetricName1")
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Metric metric2 = Metric.newBuilder()
                .setMetricName("MetricName2")
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Metric metric3 = Metric.newBuilder()
                .setMetricName("MetricName3")
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final List<Metric> want = new ArrayList<>();
        want.add(metric1);
        want.add(metric2);
        want.add(metric3);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric1);
            sut.storeMetric(metric2);
            sut.storeMetric(metric3);
            // act
            final List<Metric> have = sut.loadAllMetrics();
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBStoreMetricNullBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.storeMetric(null));
    }

    @Test
    public void testPostgresDBLoadMetricNullBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.loadMetric(null));
    }

    @Test
    public void testPostgresDBDeleteMetricNullBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.deleteMetric(null));
    }

    @Test
    public void testPostgresDBStoreMetricWithoutInitBooms() {
        // arrange
        final Metric metric = Metric.newBuilder()
                .setMetricName("MetricName")
                .setDescription("Description")
                .setLowerBound(1)
                .setUpperBound(2)
                .build();
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.storeMetric(metric));
    }

    @Test
    public void testPostgresDBLoadMetricWithoutInitBooms() {
        // arrange
        final String metricName = "MetricName";
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.loadMetric(metricName));
    }

    @Test
    public void testPostgresDBDeleteMetricWithoutInitBooms() {
        // arrange
        final String metricName = "MetricName";
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.deleteMetric(metricName));
    }

    @Test
    public void testPostgresDBStoreTwoMetricsWithSamePrimaryKeyBooms() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Metric metric1 = Metric.newBuilder()
                .setMetricName("MetricName")
                .setDescription("Description1")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Metric metric2 = Metric.newBuilder()
                .setMetricName("MetricName")
                .setDescription("Description2")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            // act
            sut.storeMetric(metric1);
            assertThrows(SQLException.class, () -> sut.storeMetric(metric2));
        }
    }

    @Test
    public void testPostgresDBDeleteMetricWithNotExistingPrimaryKeyBooms() {
        // arrange
        final String metricName ="MetricName";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteMetric(metricName));
    }

    @Test
    public void testPostgresDBStoreMetricWithNotExistingScaleBooms() {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final Metric metric = Metric.newBuilder()
                .setMetricName("MetricName")
                .setDescription("Description1")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.storeMetric(metric));
    }

    @Test
    public void testPostgresDBStoreEvaluation() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            // act
            sut.storeEvaluation(evaluation);
        }
    }

    @Test
    public void testPostgresDBLoadNotExistingEvaluation() throws Exception {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Optional<Evaluation> want = Optional.empty();
        // act
        final Optional<Evaluation> have = EMPTY_DB.loadEvaluation(metricName, timeStamp);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadEvaluation() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Optional<Evaluation> want = Optional.of(evaluation);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation);
            // act
            final Optional<Evaluation> have = sut.loadEvaluation(metricName, timeStamp);
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBDeleteEvaluation() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Optional<Evaluation> want1 = Optional.of(evaluation);
        final Optional<Evaluation> want2 = Optional.empty();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation);
            final Optional<Evaluation> have1 = sut.loadEvaluation(metricName, timeStamp);
            // act
            sut.deleteEvaluation(metricName, timeStamp);
            final Optional<Evaluation> have2 = sut.loadEvaluation(metricName, timeStamp);
            // assert
            assertEquals(want1, have1);
            assertEquals(want2, have2);
        }
    }

    @Test
    public void testPostgresDBLoadAllEvaluations() throws Exception {
        // arrange
        final List<Evaluation> want = new ArrayList<>();
        // act
        final List<Evaluation> have = EMPTY_DB.loadAllEvaluations();
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadAllEvaluations2() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String timeStamp1 = "Thu Sep 22 12:05:31 CEST 2020";
        final String timeStamp2 = "Thu Sep 22 12:05:32 CEST 2020";
        final String timeStamp3 = "Thu Sep 22 12:05:33 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation1 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp1)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Evaluation evaluation2 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp2)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Evaluation evaluation3 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp3)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final List<Evaluation> want = new ArrayList<>();
        want.add(evaluation1);
        want.add(evaluation2);
        want.add(evaluation3);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation1);
            sut.storeEvaluation(evaluation2);
            sut.storeEvaluation(evaluation3);
            // act
            final List<Evaluation> have = sut.loadAllEvaluations();
            // assert
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBStoreNullEvaluationBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.storeEvaluation(null));
    }

    @Test
    public void testPostgresDBLoadEvaluationNullMetricNameBooms() {
        // arrange
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.loadEvaluation(null, timeStamp));
    }

    @Test
    public void testPostgresDBLoadEvaluationNullTimeStampBooms() {
        // arrange
        final String metricName = "metricName";
        // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.loadEvaluation(metricName, null));
    }

    @Test
    public void testPostgresDBDeleteEvaluationNullMetricNameBooms() {
        // arrange
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.deleteEvaluation(null, timeStamp));
    }

    @Test
    public void testPostgresDBDeleteEvaluationNullTimeStampBooms() {
        // arrange
        final String metricName = "metricName";
        // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.deleteEvaluation(metricName, null));
    }

    @Test
    public void testPostgresDBStoreEvaluationWithoutInitBooms() {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.storeEvaluation(evaluation));
    }

    @Test
    public void testPostgresDBLoadEvaluationWithoutInitBooms() {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.loadEvaluation(metricName, timeStamp));
    }

    @Test
    public void testPostgresDBDeleteEvaluationWithoutInitBooms() {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.deleteEvaluation(metricName, timeStamp));
    }

    @Test
    public void testPostgresDBLoadAllEvaluationsWithoutInitBooms() {
        // arrange // act
        assertThrows(SQLException.class, UNINITIATED_EMPTY_DB::loadAllEvaluations);
    }

    @Test
    public void testPostgresDBStoreEvaluationWithoutRelatingMetricBooms() throws Exception {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            // act
            assertThrows(SQLException.class, () -> sut.storeEvaluation(evaluation));
        }
    }

    @Test
    public void testPostgresDBStoreEvaluationSamePrimaryKeyBooms() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation1 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Evaluation evaluation2 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            // act
            sut.storeEvaluation(evaluation1);
            assertThrows(SQLException.class, () -> sut.storeEvaluation(evaluation2));
        }
    }

    @Test
    public void testPostgresDBStoreEvaluationFalseMetricNameBooms() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName;";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            // act
            assertThrows(SQLException.class, () -> sut.storeEvaluation(evaluation));
        }
    }

    @Test
    public void testPostgresDBStoreEvaluationFalseTimeStampBooms() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020;";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            // act
            assertThrows(SQLException.class, () -> sut.storeEvaluation(evaluation));
        }
    }

    @Test
    public void testPostgresDBLoadEvaluationFalseMetricNameBooms() {
        // arrange
        final String metricName = "MetricName;";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.loadEvaluation(metricName, timeStamp));
    }

    @Test
    public void testPostgresDBLoadEvaluationFalseTimeStampBooms() {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020;";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.loadEvaluation(metricName, timeStamp));
    }

    @Test
    public void testPostgresDBDeleteEvaluationFalseMetricNameBooms() {
        // arrange
        final String metricName = "MetricName;";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteEvaluation(metricName, timeStamp));
    }

    @Test
    public void testPostgresDBDeleteEvaluationFalseAssetNameBooms() {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteEvaluation(metricName, timeStamp));
    }

    @Test
    public void testPostgresDBDeleteEvaluationFalseTimeStampBooms() {
        // arrange
        final String metricName = "MetricName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020;";
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteEvaluation(metricName, timeStamp));
    }

    @Test
    public void testPostgresDBStoreEvaluationToAsset() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("AssetType")
                .build();
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation);
            sut.storeAsset(asset);
            // act
            sut.storeEvaluationToAsset(evaluationToAsset);
        }
    }

    @Test
    public void testPostgresDBLoadEmptyEvaluationToAsset() throws Exception {
        // arrange
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Optional<EvaluationToAsset> want = Optional.empty();
        // act
        final Optional<EvaluationToAsset> have = EMPTY_DB.loadEvaluationToAsset(
                metricName,
                assetName,
                timeStamp
        );
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadStoredEvaluationToAsset() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("AssetType")
                .build();
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .build();
        final Optional<EvaluationToAsset> want = Optional.of(evaluationToAsset);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation);
            sut.storeAsset(asset);
            sut.storeEvaluationToAsset(evaluationToAsset);
            // act
            final Optional<EvaluationToAsset> have = sut.loadEvaluationToAsset(
                    metricName,
                    assetName,
                    timeStamp
            );
            // asset
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBDeleteEvaluationToAsset() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("AssetType")
                .build();
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .build();
        final Optional<EvaluationToAsset> want1 = Optional.of(evaluationToAsset);
        final Optional<EvaluationToAsset> want2 = Optional.empty();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation);
            sut.storeAsset(asset);
            sut.storeEvaluationToAsset(evaluationToAsset);
            final Optional<EvaluationToAsset> have1 = sut.loadEvaluationToAsset(
                    metricName,
                    assetName,
                    timeStamp
            );
            // act
            sut.deleteEvaluationToAsset(metricName, assetName, timeStamp);
            final Optional<EvaluationToAsset> have2 = sut.loadEvaluationToAsset(
                    metricName,
                    assetName,
                    timeStamp
            );
            // asset
            assertEquals(want1, have1);
            assertEquals(want2, have2);
        }
    }

    @Test
    public void testPostgresDBLoadAllEmptyEvaluationToAsset() throws Exception {
        // arrange
        final List<EvaluationToAsset> want = new ArrayList<>();
        // act
        final List<EvaluationToAsset> have = EMPTY_DB.loadAllEvaluationsToAssets();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testPostgresDBLoadAllEmptyEvaluationToAsset2() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp1 = "Thu Sep 22 12:05:31 CEST 2020";
        final String timeStamp2 = "Thu Sep 22 12:05:32 CEST 2020";
        final String timeStamp3 = "Thu Sep 22 12:05:33 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("AssetType")
                .build();
        final Evaluation evaluation1 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp1)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Evaluation evaluation2 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp2)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final Evaluation evaluation3 = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp3)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final EvaluationToAsset evaluationToAsset1 = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp1)
                .build();
        final EvaluationToAsset evaluationToAsset2 = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp2)
                .build();
        final EvaluationToAsset evaluationToAsset3 = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp3)
                .build();
        final List<EvaluationToAsset> want = new ArrayList<>();
        want.add(evaluationToAsset1);
        want.add(evaluationToAsset2);
        want.add(evaluationToAsset3);
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation1);
            sut.storeEvaluation(evaluation2);
            sut.storeEvaluation(evaluation3);
            sut.storeAsset(asset);
            sut.storeEvaluationToAsset(evaluationToAsset1);
            sut.storeEvaluationToAsset(evaluationToAsset2);
            sut.storeEvaluationToAsset(evaluationToAsset3);

            // act
            final List<EvaluationToAsset> have = sut.loadAllEvaluationsToAssets();
            // asset
            assertEquals(want, have);
        }
    }

    @Test
    public void testPostgresDBStoreNullEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.storeEvaluationToAsset(null));
    }

    @Test
    public void testPostgresDBStoreUninitializedDbEvaluationToAssetBooms() {
        // arrange
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setMetricName("MetricName")
                .setAssetName("AssetName")
                .setTimeStamp("Thu Sep 22 12:05:31 CEST 2020")
                .build();
        // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.storeEvaluationToAsset(evaluationToAsset));
    }

    @Test
    public void testPostgresDBStoreEvaluationToAssetWithoutAssetBooms() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation);
            // act
            assertThrows(SQLException.class, () -> sut.storeEvaluationToAsset(evaluationToAsset));
        }
    }

    @Test
    public void testPostgresDBStoreEvaluationToAssetWithoutEvaluationBooms() throws Exception {
        // arrange
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .build();
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("assetType")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeAsset(asset);
            // act
            assertThrows(SQLException.class, () -> sut.storeEvaluationToAsset(evaluationToAsset));
        }
    }

    @Test
    public void testPostgresDBStoreEvaluationTwiceToAssetWithoutAssetBooms() throws Exception {
        // arrange
        final int lowerBound = 1;
        final int upperBound = 2;
        final String metricName = "MetricName";
        final String assetName = "AssetName";
        final String timeStamp = "Thu Sep 22 12:05:31 CEST 2020";
        final Scale scale = Scale.newBuilder()
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .setDescription("Description")
                .build();
        final Metric metric = Metric.newBuilder()
                .setMetricName(metricName)
                .setDescription("Description")
                .setLowerBound(lowerBound)
                .setUpperBound(upperBound)
                .build();
        final Evaluation evaluation = Evaluation.newBuilder()
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .setResultValue("{\"result\": 1.5}")
                .build();
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setAssetName(assetName)
                .setMetricName(metricName)
                .setTimeStamp(timeStamp)
                .build();
        final Asset asset = Asset.newBuilder()
                .setAssetName(assetName)
                .setAssetType("assetType")
                .build();
        try (
                final EmbeddedPostgres pg = EmbeddedPostgres.start();
                final Connection connection = pg.getPostgresDatabase().getConnection();
                final DBServer sut = new PostgresDB(connection)
        ) {
            Table.init(connection);
            sut.storeScale(scale);
            sut.storeMetric(metric);
            sut.storeEvaluation(evaluation);
            sut.storeAsset(asset);
            // act
            sut.storeEvaluationToAsset(evaluationToAsset);
            assertThrows(SQLException.class, () -> sut.storeEvaluationToAsset(evaluationToAsset));
        }
    }

    @Test
    public void testPostgresDBStoreMetricNameNotMatchingBooms() {
        // arrange
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setMetricName("1")
                .setAssetName("AssetName")
                .setTimeStamp("Thu Sep 22 12:05:31 CEST 2020")
                .build();
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.storeEvaluationToAsset(evaluationToAsset));
    }

    @Test
    public void testPostgresDBStoreAssetNameNotMatchingBooms() {
        // arrange
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setMetricName("MetricName")
                .setAssetName("1")
                .setTimeStamp("Thu Sep 22 12:05:31 CEST 2020")
                .build();
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.storeEvaluationToAsset(evaluationToAsset));
    }

    @Test
    public void testPostgresDBStoreTimeStampNameNotMatchingBooms() {
        // arrange
        final EvaluationToAsset evaluationToAsset = EvaluationToAsset.newBuilder()
                .setMetricName("MetricName")
                .setAssetName("AssetName")
                .setTimeStamp("111 Sep 22 12:05:31 CEST 2020")
                .build();
        // act
        assertThrows(SQLException.class, () -> EMPTY_DB.storeEvaluationToAsset(evaluationToAsset));
    }

    @Test
    public void testPostgresDBLoadMetricNameNullEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.loadEvaluationToAsset(
                null,
                "AssetName",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBLoadAssetNameNullEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.loadEvaluationToAsset(
                "MetricName",
                null,
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBLoadTimeStampNameNullEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.loadEvaluationToAsset(
                "MetricName",
                "AssetName",
                null
        ));
    }

    @Test
    public void testPostgresDBLoadMetricNameNotOKEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> EMPTY_DB.loadEvaluationToAsset(
                "1",
                "AssetName",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBLoadAssetNameNotOKEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> EMPTY_DB.loadEvaluationToAsset(
                "MetricName",
                "1",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBTimeStampNotOKEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> EMPTY_DB.loadEvaluationToAsset(
                "MetricName",
                "AssetName",
                "1"
        ));
    }

    @Test
    public void testPostgresDBLoadEvaluationToAssetUninitiatedBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.loadEvaluationToAsset(
                "MetricName",
                "AssetName",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBDeleteMetricNameNullEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.deleteEvaluationToAsset(
                null,
                "AssetName",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBDeleteAssetNameNullEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.deleteEvaluationToAsset(
                "MetricName",
                null,
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBDeleteTimeStampNameNullEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(NullPointerException.class, () -> EMPTY_DB.deleteEvaluationToAsset(
                "MetricName",
                "AssetName",
                null
        ));
    }

    @Test
    public void testPostgresDBDeleteMetricNameNotOKEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteEvaluationToAsset(
                "1",
                "AssetName",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBDeleteAssetNameNotOKEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteEvaluationToAsset(
                "MetricName",
                "1",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBDeleteTimeStampNotOKEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteEvaluationToAsset(
                "MetricName",
                "AssetName",
                "1"
        ));
    }

    @Test
    public void testPostgresDBDeleteEvaluationToAssetUninitiatedBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> UNINITIATED_EMPTY_DB.deleteEvaluationToAsset(
                "MetricName",
                "AssetName",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBDeleteEmptyTableEvaluationToAssetBooms() {
        // arrange // act
        assertThrows(SQLException.class, () -> EMPTY_DB.deleteEvaluationToAsset(
                "MetricName",
                "AssetName",
                "Thu Sep 22 12:05:31 CEST 2020"
        ));
    }

    @Test
    public void testPostgresDBLoadAllEvaluationsToAssetsUninitiatedBooms() {
        // arrange // act
        assertThrows(SQLException.class, UNINITIATED_EMPTY_DB::loadAllEvaluationsToAssets);
    }
}
