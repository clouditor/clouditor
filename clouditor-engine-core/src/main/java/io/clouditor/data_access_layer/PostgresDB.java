package io.clouditor.data_access_layer;

import com.google.protobuf.GeneratedMessageV3;
import io.clouditor.data_access_layer.utils.BoomingFunction;
import io.clouditor.data_access_layer.utils.DBUtils;
import io.clouditor.data_access_layer.utils.SQLMethod;
import io.clouditor.data_access_layer.utils.Table;
import org.postgresql.util.PGobject;
import io.clouditor.metric_api.pb.*;

import java.sql.*;
import java.util.*;

/**
 * A Implementation of a java Postgre db client.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public class PostgresDB implements DBServer {
    private static final String URL_SEPARATOR = "/";
    private static final String PROTOCOL = "jdbc:postgresql://";
    private static final String DEFAULT_HOST = "127.0.0.1";
    private static final int DEFAULT_PORT = 5432;

    private static final BoomingFunction<ResultSet, Scale> EXTRACT_SCALE = resultSet -> Scale
            .newBuilder()
            .setLowerBound(
                    resultSet.getInt("lower_bound")
            ).setUpperBound(
                    resultSet.getInt("upper_bound")
            ).setDescription(
                    resultSet.getString("description")
            ).build();

    private static final BoomingFunction<ResultSet, Asset> EXTRACT_ASSET = resultSet -> Asset
            .newBuilder()
            .setAssetName(
                    resultSet.getString("asset_name")
            )
            .setAssetType(
                    resultSet.getString("asset_type")
            )
            .build();

    private static final BoomingFunction<ResultSet, Metric> EXTRACT_METRIC = resultSet -> Metric
            .newBuilder()
            .setMetricName(
                    resultSet.getString("metric_name")
            )
            .setDescription(
                    resultSet.getString("description")
            )
            .setLowerBound(
                    resultSet.getInt("lower_bound")
            )
            .setUpperBound(
                    resultSet.getInt("upper_bound")
            )
            .build();

    private static final BoomingFunction<ResultSet, Evaluation> EXTRACT_EVALUATION = resultSet -> Evaluation
            .newBuilder()
            .setMetricName(
                    resultSet.getString("metric_name")
            )
            .setTimeStamp(
                    resultSet.getString("time_stamp")
            )
            .setResultValue(
                    resultSet.getString("result_value")
            )
            .build();

    private static final BoomingFunction<ResultSet, EvaluationToAsset> EXTRACT_EVALUATION_TO_ASSET = resultSet -> EvaluationToAsset
            .newBuilder()
            .setMetricName(
                    resultSet.getString("metric_name")
            )
            .setTimeStamp(
                    resultSet.getString("time_stamp")
            )
            .setAssetName(
                    resultSet.getString("asset_name")
            )
            .build();

    private final Connection dbConnection;

    public PostgresDB(final Connection dbConnection) {
        this.dbConnection = dbConnection;
    }

    public PostgresDB(
            final String dbName,
            final String userName,
            final String password
    ) throws SQLException {
        this(dbName, userName, password, DEFAULT_HOST);
    }

    public PostgresDB(
            final String dbName,
            final String userName,
            final String password,
            final String host
    ) throws SQLException {
        this(dbName, userName, password, host, DEFAULT_PORT);
    }

    public PostgresDB(
            final String dbName,
            final String userName,
            final String password,
            final String host,
            final int port
    ) throws SQLException {
        Objects.requireNonNull(dbName);
        Objects.requireNonNull(userName);
        Objects.requireNonNull(password);
        Objects.requireNonNull(host);
        if (port <= 1024 || port > 49151)
            throw new IllegalArgumentException(
                    "The port: " + port + " is not allowed."
            );
        final String url = PROTOCOL +
                host + ":" + port +
                URL_SEPARATOR + dbName;
        try {
            this.dbConnection = DriverManager.getConnection(url, userName, password);
        } catch (SQLException sqlException) {
            System.err.println("Unable to connect to: " + url + ".");
            throw sqlException;
        }
    }

    @Override
    public int getDefaultPort() {
        return PostgresDB.DEFAULT_PORT;
    }

    @Override
    public String getDefaultHost() {
        return PostgresDB.DEFAULT_HOST;
    }


    @Override
    public void close() throws SQLException {
        getDBConnection().close();
    }

    @Override
    public void storeMetric(final Metric metric) throws SQLException {
        final String metricName = metric.getMetricName();
        final String description = metric.getDescription();
        final int lowerBound = metric.getLowerBound();
        final int upperBound = metric.getUpperBound();
        DBUtils.testName(metricName);
        DBUtils.testDescription(description);
        final String query = DBUtils.getQuery(SQLMethod.ADD, Table.METRIC);
        final boolean isNotPresent = loadMetric(metricName).isEmpty();
        final PreparedStatement preparedStatement = getPrepared(query, metricName, description, lowerBound, upperBound);
        execUpdate(isNotPresent, preparedStatement);
    }

    @Override
    public Optional<Metric> loadMetric(final String metricName) throws SQLException {
        DBUtils.testName(metricName);
        final String query = DBUtils.getQuery(SQLMethod.GET, Table.METRIC);
        final PreparedStatement preparedStatement = getPrepared(query, metricName);
        return get(preparedStatement, EXTRACT_METRIC);
    }

    @Override
    public void deleteMetric(final String metricName) throws SQLException {
        DBUtils.testName(metricName);
        final String query = DBUtils.getQuery(SQLMethod.DELETE, Table.METRIC);
        final PreparedStatement preparedStatement = getPrepared(query, metricName);
        final boolean isPresent = loadMetric(metricName).isPresent();
        execUpdate(isPresent, preparedStatement);
    }

    @Override
    public List<Metric> loadAllMetrics() throws SQLException {
        return Collections.unmodifiableList(
                getAll(Table.METRIC.getTableName(), EXTRACT_METRIC)
        );
    }

    @Override
    public void storeAsset(final Asset asset) throws SQLException {
        final String assetName = asset.getAssetName();
        final String assetType = asset.getAssetType();
        DBUtils.testName(assetName);
        DBUtils.testDescription(assetType);
        final boolean isNotPresent = loadAsset(assetName).isEmpty();
        final String query = DBUtils.getQuery(SQLMethod.ADD, Table.ASSET);
        final PreparedStatement preparedAdd = getPrepared(query, assetName, assetType);
        execUpdate(isNotPresent, preparedAdd);
    }

    @Override
    public Optional<Asset> loadAsset(final String assetName) throws SQLException {
        DBUtils.testName(assetName);
        final String query = DBUtils.getQuery(SQLMethod.GET, Table.ASSET);
        final PreparedStatement preparedStatement = getPrepared(query, assetName);
        return get(preparedStatement, EXTRACT_ASSET);
    }

    @Override
    public void deleteAsset(final String assetName) throws SQLException {
        DBUtils.testName(assetName);
        final boolean isPresent = loadAsset(assetName).isPresent();
        final String query = DBUtils.getQuery(SQLMethod.DELETE, Table.ASSET);
        final PreparedStatement preparedAdd = getPrepared(query, assetName);
        execUpdate(isPresent, preparedAdd);
    }

    @Override
    public List<Asset> loadAllAssets() throws SQLException {
        return Collections.unmodifiableList(
                getAll(Table.ASSET.getTableName(), EXTRACT_ASSET)
        );
    }

    @Override
    public void storeScale(final Scale scale) throws SQLException {
        final int lowerBound = scale.getLowerBound();
        final int upperBound = scale.getUpperBound();
        final String description = scale.getDescription();
        DBUtils.testDescription(description);
        final boolean isNotPresent = loadScale(lowerBound, upperBound).isEmpty();
        final String query = DBUtils.getQuery(SQLMethod.ADD, Table.SCALE);
        final PreparedStatement preparedStatement = getPrepared(query, lowerBound, upperBound, description);
        execUpdate(isNotPresent, preparedStatement);
    }

    @Override
    public Optional<Scale> loadScale(final int lowerBound, final int upperBound) throws SQLException {
        final String query = DBUtils.getQuery(SQLMethod.GET, Table.SCALE);
        final PreparedStatement preparedStatement = getPrepared(query, lowerBound, upperBound);
        return get(preparedStatement, EXTRACT_SCALE);
    }

    @Override
    public void deleteScale(final int lowerBound, final int upperBound) throws SQLException {
        final String query = DBUtils.getQuery(SQLMethod.DELETE, Table.SCALE);
        final boolean isPresent = loadScale(lowerBound, upperBound).isPresent();
        final PreparedStatement preparedStatement = getPrepared(query, lowerBound, upperBound);
        execUpdate(isPresent, preparedStatement);
    }

    @Override
    public List<Scale> loadAllScales() throws SQLException {
        return Collections.unmodifiableList(
                getAll(Table.SCALE.getTableName(), EXTRACT_SCALE)
        );
    }

    @Override
    public void storeEvaluation(final Evaluation evaluation) throws SQLException {
        final String metricName = evaluation.getMetricName();
        final String timeStamp = evaluation.getTimeStamp();
        final String resultValue = evaluation.getResultValue();
        DBUtils.testName(metricName);
        DBUtils.testTimeStamp(timeStamp);
        DBUtils.getJSON(resultValue);
        final PGobject json = new PGobject();
        json.setType("jsonb");
        json.setValue(resultValue);
        final boolean isNotPresent = loadEvaluation(metricName, timeStamp).isEmpty();
        final String query = DBUtils.getQuery(SQLMethod.ADD, Table.EVALUATION);
        final PreparedStatement preparedStatement = getPrepared(query, metricName, timeStamp, json);
        execUpdate(isNotPresent, preparedStatement);
    }

    @Override
    public Optional<Evaluation> loadEvaluation(final String metricName, final String timeStamp) throws SQLException {
        DBUtils.testName(metricName);
        DBUtils.testTimeStamp(timeStamp);
        final String query = DBUtils.getQuery(SQLMethod.GET, Table.EVALUATION);
        final PreparedStatement preparedStatement = getPrepared(query, metricName, timeStamp);
        return get(preparedStatement, EXTRACT_EVALUATION);
    }

    @Override
    public void deleteEvaluation(final String metricName, final String timeStamp) throws SQLException {
        DBUtils.testName(metricName);
        DBUtils.testTimeStamp(timeStamp);
        final boolean isPresent = loadEvaluation(metricName, timeStamp).isPresent();
        final String query = DBUtils.getQuery(SQLMethod.DELETE, Table.EVALUATION);
        final PreparedStatement preparedStatement = getPrepared(query, metricName, timeStamp);
        execUpdate(isPresent, preparedStatement);
    }

    @Override
    public List<Evaluation> loadAllEvaluations() throws SQLException {
        return Collections.unmodifiableList(
                getAll(Table.EVALUATION.getTableName(), EXTRACT_EVALUATION)
        );
    }

    @Override
    public void storeEvaluationToAsset(final EvaluationToAsset evaluationToAsset) throws SQLException {
        final String metricName = evaluationToAsset.getMetricName();
        final String assetName = evaluationToAsset.getAssetName();
        final String timeStamp = evaluationToAsset.getTimeStamp();
        DBUtils.testName(metricName);
        DBUtils.testName(assetName);
        DBUtils.testTimeStamp(timeStamp);
        final boolean isNotPresent = loadEvaluationToAsset(metricName, assetName, timeStamp).isEmpty();
        final String query = DBUtils.getQuery(SQLMethod.ADD, Table.EVALUATION_TO_ASSET);
        final PreparedStatement preparedStatement = getPrepared(query, metricName, assetName, timeStamp);
        execUpdate(isNotPresent, preparedStatement);
    }

    @Override
    public Optional<EvaluationToAsset> loadEvaluationToAsset(
            final String metricName,
            final String assetName,
            final String timeStamp
    ) throws SQLException {
        DBUtils.testName(metricName);
        DBUtils.testTimeStamp(timeStamp);
        DBUtils.testName(assetName);
        final String query = DBUtils.getQuery(SQLMethod.GET, Table.EVALUATION_TO_ASSET);
        final PreparedStatement preparedStatement = getPrepared(query, metricName, assetName, timeStamp);
        return get(preparedStatement, EXTRACT_EVALUATION_TO_ASSET);
    }

    @Override
    public void deleteEvaluationToAsset(
            final String metricName,
            final String assetName,
            final String timeStamp
    ) throws SQLException {
        DBUtils.testName(metricName);
        DBUtils.testName(assetName);
        DBUtils.testTimeStamp(timeStamp);
        final boolean isPresent = loadEvaluationToAsset(metricName, assetName, timeStamp).isPresent();
        final String query = DBUtils.getQuery(SQLMethod.DELETE, Table.EVALUATION_TO_ASSET);
        final PreparedStatement preparedStatement = getPrepared(query, metricName, assetName, timeStamp);
        execUpdate(isPresent, preparedStatement);
    }

    @Override
    public List<EvaluationToAsset> loadAllEvaluationsToAssets() throws SQLException {
        return Collections.unmodifiableList(
                getAll(Table.EVALUATION_TO_ASSET.getTableName(), EXTRACT_EVALUATION_TO_ASSET)
        );
    }

    private Connection getDBConnection() {
        return this.dbConnection;
    }

    private void execUpdate(
            final boolean condition,
            final PreparedStatement prepared
            ) throws SQLException {
        if (condition) prepared.executeUpdate();
        else throw new SQLException(
                "Unable to modify the dataset relating to the query: "
                        + prepared
        );
    }

    private <T extends GeneratedMessageV3> Optional<T> get(
            final PreparedStatement preparedGet,
            final BoomingFunction<ResultSet, T> extractObject
    ) throws SQLException {
        Optional<T> result = Optional.empty();
        final ResultSet resultSet = preparedGet.executeQuery();
        if (resultSet.next()) {
            final T obj = extractObject.handleSQLException().apply(resultSet);
            result = Optional.of(obj);
            final boolean hasNext = resultSet.next();
            assert !hasNext;
        }
        return result;
    }

    private <T extends GeneratedMessageV3> List<T> getAll(
            final String tableName,
            final BoomingFunction<ResultSet, T> extractObject
    ) throws SQLException {
        Objects.requireNonNull(tableName);
        DBUtils.testName(tableName);
        final List<T> result = new ArrayList<>();
        final PreparedStatement preparedStatement = getDBConnection()
                .prepareStatement(SQLMethod.GET + tableName + ";");
        final ResultSet resultSet = preparedStatement.executeQuery();
        while (resultSet.next()) {
            final T obj = extractObject.handleSQLException().apply(resultSet);
            result.add(obj);
        }
        return result;
    }

    private PreparedStatement getPrepared(
            final String query,
            final Object... values
    ) throws SQLException {
        final PreparedStatement preparedStatement = getDBConnection()
                .prepareStatement(query);
        final List<Object> valueList = Arrays.asList(values);
        for (int index = 0; index < valueList.size(); index++)
            setAtIndex(preparedStatement, index + 1, valueList.get(index));
        return preparedStatement;
    }

    private void setAtIndex(
            final PreparedStatement preparedStatement,
            final int index,
            final Object value
    ) throws SQLException {
        if (value instanceof Integer) {
            preparedStatement.setInt(index, (Integer) value);
        } else if (value instanceof String) {
            final String string = (String) value;
            preparedStatement.setString(index, string);
        } else if (value instanceof Double) {
            preparedStatement.setDouble(index, (Double) value);
        } else if (value instanceof PGobject) {
            preparedStatement.setObject(index, value);
        } else {
            throw new IllegalArgumentException("Got a unexpected value: " + value);
        }
    }
}