package io.clouditor.data_access_layer;

import java.sql.SQLException;
import java.util.List;
import java.util.Optional;

import io.clouditor.metric_api.pb.*;

/**
 * The interface for a database Server.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public interface DBServer extends AutoCloseable {

    int getDefaultPort();

    String getDefaultHost();

    /**
     * Stores a Metric in the Database.
     *
     * The metric should contain a <code>metricName</code> matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * The metric should contain a <code>description</code> matching the regex: "^[\\s0-9a-zA-Z_\\-.]*$".
     * The metric should contain a <code>lowerBound</code>.
     * The metric should contain a <code>upperBound</code>.
     *      The lower and the upper bound should exist as a primary key in the table scale.
     *
     * @param metric: the metric to store.
     * @throws SQLException if one of the conditions of the metric and the related scale was not fulfilled
     *              or the table metric or the table scale does not exist.
     * @throws NullPointerException if the <code>metric</code> is <code>null</code>.
     */
    void storeMetric(final io.clouditor.metric_api.pb.Metric metric) throws SQLException;

    /**
     * Load a metric from the database.
     *
     * @param metricName the unique name of a metric, matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @return An Optional empty optional if no metric with the key <code>metricName</code> does exist in the metric table.
     *         An Optional containing the metric if the metric exist in the table.
     * @throws SQLException if the metric name dos not fulfill the condition or the table metric does not exist.
     * @throws NullPointerException if the <code>metricName</code> is <code>null</code>.
     */
    Optional<Metric> loadMetric(final String metricName) throws SQLException;

    /**
     * Delete the metric from the database.
     *
     * @param metricName the unique name of a metric, matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @throws SQLException if the metric name dos not fulfill the condition or the table metric does not exist.
     * @throws NullPointerException if the <code>metricName</code> is <code>null</code>.
     */
    void deleteMetric(final String metricName) throws SQLException;

    /**
     * Load all metrics from the database.
     *
     * @return an unmodifiable list containing all metrics that are stored in the database.
     * @throws SQLException if the table metric does not exist.
     */
    List<Metric> loadAllMetrics() throws SQLException;

    /**
     * Store an evaluation in the database.
     *
     * The evaluation should contain a <code>metricName</code> matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * The evaluation should contain a <code>timeStamp</code> matching the regex:
     *      "^(Sun|Mon|Tue|Wed|Fri|Sat) "
     *          + "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) "
     *          + "([0-3][0-9]) "
     *          + "([0-2][0-9]):([0-5][0-9]):([0-6][0-9]) "
     *          + "([A-Z]* |)([0-9][0-9][0-9][0-9])$".
     * The evaluation should contain a <code>resultValue</code>.
     *      The metric relating to the <code>metricName</code> should exist.
     *      The asset relating to the <code>assetName</code> should exist.
     *
     * @param evaluation the evaluation to store.
     * @throws SQLException if one of the conditions are not fulfilled or
     *      the evaluation does already exist in the database or
     *      the table does not exist.
     * @throws NullPointerException if the <code>evaluation</code> is <code>null</code>.
     */
     void storeEvaluation(final Evaluation evaluation) throws SQLException;

     /**
      * Load an evaluation from the database.
      *
      * @param metricName the unique metric name matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
      * @param timeStamp the time stamp of te evaluation matching the regex:
      *     "^(Sun|Mon|Tue|Wed|Fri|Sat) "
      *            + "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) "
      *            + "([0-3][0-9]) "
      *            + "([0-2][0-9]):([0-5][0-9]):([0-6][0-9]) "
      *            + "([A-Z]* |)([0-9][0-9][0-9][0-9])$".
      * @return An empty optional if the evaluation, matching the primary key (metricName, assetName,
      *     timeStamp), does not exist in de table evaluation. An optional containing an evaluation if
      *     the primary key does exist.
      * @throws SQLException if the table evaluation does not exist.
      * @throws NullPointerException if the <code>metricName</code>, the <code>assetName</code> or the
      *     <code>timeStamp</code> is <code>null</code>.
      */
     Optional<Evaluation> loadEvaluation(final String metricName, final String timeStamp) throws SQLException;

    /**
     * Delete an evaluation from the database
     *
     * @param metricName the unique metric name matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @param timeStamp the time stamp of te evaluation matching the regex: "^(Sun|Mon|Tue|Wed|Fri|Sat) "
     *                  + "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) "
     *                  + "([0-3][0-9]) "
     *                  + "([0-2][0-9]):([0-5][0-9]):([0-6][0-9]) "
     *                  + "([A-Z]* |)([0-9][0-9][0-9][0-9])$".
     * @throws SQLException if one of the conditions are not fulfilled or
     *      the evaluation does not exist in the database or
     *      the table does not exist.
     * @throws NullPointerException if the <code>metricName</code>,
     *          the <code>assetName</code> or
     *          the <code>timeStamp</code> is <code>null</code>.
     */
    void deleteEvaluation(final String metricName, final String timeStamp) throws SQLException;

    /**
     * Load all evaluations form the database.
     *
     * @return an unmodifiable list containing all metrics that are stored in the database.
     * @throws SQLException if the table evaluation does not exist.
     */
    List<Evaluation> loadAllEvaluations() throws SQLException;



    /**
     * Store an evaluationToAsset in the database.
     *
     * The evaluation should contain a <code>metricName</code> matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * The evaluation should contain a <code>assetName</code> matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * The evaluation should contain a <code>timeStamp</code> matching the regex:
     *      "^(Sun|Mon|Tue|Wed|Fri|Sat) "
     *          + "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) "
     *          + "([0-3][0-9]) "
     *          + "([0-2][0-9]):([0-5][0-9]):([0-6][0-9]) "
     *          + "([A-Z]* |)([0-9][0-9][0-9][0-9])$".
     *
     * @param evaluationToAsset the evaluationToAsset to store.
     * @throws SQLException if one of the conditions are not fulfilled or
     *      the evaluation_to_asset does already exist in the database or
     *      the table does not exist.
     * @throws NullPointerException if the <code>evaluation</code> is <code>null</code>.
     */
    void storeEvaluationToAsset(final EvaluationToAsset evaluationToAsset) throws SQLException;

    /**
     * Load an evaluationToAsset from the database.
     *
     * @param metricName the unique metric name matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @param assetName the unique asset name matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @param timeStamp the time stamp of te evaluation matching the regex:
     *     "^(Sun|Mon|Tue|Wed|Fri|Sat) "
     *            + "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) "
     *            + "([0-3][0-9]) "
     *            + "([0-2][0-9]):([0-5][0-9]):([0-6][0-9]) "
     *            + "([A-Z]* |)([0-9][0-9][0-9][0-9])$".
     * @return An empty optional if the evaluationToAsset, matching the primary key (metricName, assetName,
     *     timeStamp), does not exist in de table evaluation. An optional containing an evaluation if
     *     the primary key does exist.
     * @throws SQLException if the table evaluation does not exist.
     * @throws NullPointerException if the <code>metricName</code>, the <code>assetName</code> or the
     *     <code>timeStamp</code> is <code>null</code>.
     */
    Optional<EvaluationToAsset> loadEvaluationToAsset(
            final String metricName,
            final String assetName,
            final String timeStamp
    ) throws SQLException;

    /**
     * Delete an evaluationToAsset from the database
     *
     * @param metricName the unique metric name matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @param assetName the unique asset name matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @param timeStamp the time stamp of te evaluation matching the regex: "^(Sun|Mon|Tue|Wed|Fri|Sat) "
     *                  + "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) "
     *                  + "([0-3][0-9]) "
     *                  + "([0-2][0-9]):([0-5][0-9]):([0-6][0-9]) "
     *                  + "([A-Z]* |)([0-9][0-9][0-9][0-9])$".
     * @throws SQLException if one of the conditions are not fulfilled or
     *      the evaluation does not exist in the database or
     *      the table does not exist.
     * @throws NullPointerException if the <code>metricName</code>,
     *          the <code>assetName</code> or
     *          the <code>timeStamp</code> is <code>null</code>.
     */
    void deleteEvaluationToAsset(
            final String metricName,
            final String assetName,
            final String timeStamp
    ) throws SQLException;

    /**
     * Load all evaluationsToAssets form the database.
     *
     * @return an unmodifiable list containing all metrics that are stored in the database.
     * @throws SQLException if the table evaluation does not exist.
     */
    List<EvaluationToAsset> loadAllEvaluationsToAssets() throws SQLException;


    /**
     * Store and asset in the database.
     *
     * The asset should contain a <code>assetName</code> matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * The asset should contain a <code>assetType</code> matching the regex: "^[\\s0-9a-zA-Z_\\-.]*$".
     *
     * @param asset the asset to store.
     * @throws SQLException if the asset already exist in the table asset or
     *          the table asset does not exist or
     *          one of the conditions are not fulfilled.
     * @throws NullPointerException if the <code>asset</code> is <code>null</code>.
     */
    void storeAsset(final Asset asset) throws SQLException;

    /**
     * Load the asset from the database.
     *
     * @param assetName the unique name of the asset, matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @return an empty optional if in the table is no asset matching the primary key assetName.
     *          an optional containing the asset if it exist in the table.
     * @throws SQLException if the asset name does not match the regex or
     *          the table does not exist.
     * @throws NullPointerException if the <code>assetName</code> is null.
     */
    Optional<Asset> loadAsset(final String assetName) throws SQLException;

    /**
     * Delete an asset from the database.
     *
     * @param assetName the unique name of the asset, matching the regex: "^[a-zA-Z][0-9a-zA-Z_\\-]*$".
     * @throws SQLException if the asset name does not match the regex or
     *          the table does not exist.
     * @throws NullPointerException if the <code>assetName</code> is null.
     */
    void deleteAsset(final String assetName) throws SQLException;

    /**
     * Load all assets form the database.
     *
     * @return an unmodifiable list containing all assets that are stored in the database.
     * @throws SQLException if the table asset does not exist.
     */
    List<Asset> loadAllAssets() throws SQLException;


    /**
     * Store a scale in the database
     *
     * The scale should contain a <code>lowerBound</code>.
     * The scale should contain a <code>upperBound</code>.
     * The scale should contain a <code>description</code> matching the regex: "^[\\s0-9a-zA-Z_\\-.]*$".
     *
     * @param scale the scale to store.
     * @throws SQLException if the scale does already exist or
     *          the table does not exist.
     * @throws NullPointerException if the <code>scale</code> is <code>null</code>.
     */
    void storeScale(final Scale scale) throws SQLException;

    /**
     * Load a scale from the database.
     *
     * @param lowerBound the lower bound of the scale.
     * @param upperBound the upper bound of the scale.
     * @return an empty optional if in the table scale is no dataset matching the primary key (lowerBound, upperBound)
     *          an optional containing the scale if it exists.
     * @throws SQLException if the table does not exists.
     */
    Optional<Scale> loadScale(final int lowerBound, final int upperBound) throws SQLException;

    /**
     * Delete the scale from the database.
     *
     * @param lowerBound the lower bound of the scale.
     * @param upperBound the upper bound of the scale.
     * @throws SQLException if the scale or the table does not exists.
     */
    void deleteScale(final int lowerBound, final int upperBound) throws SQLException;

    /**
     * Load all scales from the database.
     *
     * @return an unmodifiable list containing all scales.
     * @throws SQLException if the table does not exist.
     */
    List<Scale> loadAllScales() throws SQLException;
}
