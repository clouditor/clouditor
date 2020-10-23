package io.clouditor.data_access_layer.dbserver.utils;

import io.clouditor.data_access_layer.utils.Table;
import org.junit.jupiter.api.Test;

import static org.junit.Assert.assertEquals;

public class TableTest {

    @Test
    public void testScaleName() {
        // arrange
        final Table sut = Table.SCALE;
        final String want = "scale";
        // act
        final String have = sut.getTableName();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testScaleQuery() {
        // arrange
        final Table sut = Table.SCALE;
        final String want = "lower_bound=? AND upper_bound=?";
        // act
        final String have = sut.getKeyQuery();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testScaleAttributes() {
        // arrange
        final Table sut = Table.SCALE;
        final String want = "values(?,?,?)";
        // act
        final String have = sut.getAttributes();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testAssetName() {
        // arrange
        final Table sut = Table.ASSET;
        final String want = "asset";
        // act
        final String have = sut.getTableName();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testAssetQuery() {
        // arrange
        final Table sut = Table.ASSET;
        final String want = "asset_name=?";
        // act
        final String have = sut.getKeyQuery();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testAssetAttributes() {
        // arrange
        final Table sut = Table.ASSET;
        final String want = "values(?,?)";
        // act
        final String have = sut.getAttributes();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testMetricName() {
        // arrange
        final Table sut = Table.METRIC;
        final String want = "metric";
        // act
        final String have = sut.getTableName();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testMetricQuery() {
        // arrange
        final Table sut = Table.METRIC;
        final String want = "metric_name=?";
        // act
        final String have = sut.getKeyQuery();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testMetricAttributes() {
        // arrange
        final Table sut = Table.METRIC;
        final String want = "values(?,?,?,?)";
        // act
        final String have = sut.getAttributes();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testEvaluationName() {
        // arrange
        final Table sut = Table.EVALUATION;
        final String want = "evaluation";
        // act
        final String have = sut.getTableName();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testEvaluationQuery() {
        // arrange
        final Table sut = Table.EVALUATION;
        final String want = "metric_name=? AND time_stamp=?";
        // act
        final String have = sut.getKeyQuery();
        // asset
        assertEquals(want, have);
    }

    @Test
    public void testEvaluationAttributes() {
        // arrange
        final Table sut = Table.EVALUATION;
        final String want = "values(?,?,?)";
        // act
        final String have = sut.getAttributes();
        // asset
        assertEquals(want, have);
    }

}
