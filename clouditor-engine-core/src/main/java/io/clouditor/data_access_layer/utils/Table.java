package io.clouditor.data_access_layer.utils;

import java.sql.Connection;
import java.sql.SQLException;
import java.util.Arrays;
import java.util.List;
import java.util.Objects;

/**
 * A enumeration of all tables used in the PostgresDB service.
 * The arguments follow the order:
 *      ( <NameOfTheTable>, <CountOfAttributes>, <NameOfThePrimaryKeyAttribute>... )
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public enum Table {
    SCALE("scale",  3, "lower_bound", "upper_bound"),
    ASSET("asset", 2, "asset_name"  ),
    METRIC("metric", 4, "metric_name"),
    EVALUATION("evaluation",3,"metric_name", "time_stamp"),
    EVALUATION_TO_ASSET("evaluation_to_asset", 3,"metric_name", "asset_name", "time_stamp");

    private static final String PREPARED_EQUAL = "=?";
    private static final String AND = " AND ";


    private final String tableName;
    private final int attributes;
    private final List<String> keyList;


    Table(
            final String tableName,
            final int attributes,
            final String... keys
    ) {
        Objects.requireNonNull(keys);
        if (keys.length < 1)
            throw new IllegalArgumentException(
                    "Every table should have at least one primary key."
            );
        this.tableName = tableName;
        this.attributes = attributes;
        this.keyList = Arrays.asList(keys);
    }

    public String getTableName() {
        return this.tableName;
    }

    public String getKeyQuery() {
        final StringBuilder stringBuilder = new StringBuilder();
        for (int index = 0; index < keyList.size() - 1; index++) {
            stringBuilder.append(keyList.get(index));
            stringBuilder.append(PREPARED_EQUAL);
            stringBuilder.append(AND);
        }
        final String lastElement = keyList.get(keyList.size()-1);
        stringBuilder.append(lastElement);
        stringBuilder.append(PREPARED_EQUAL);
        return stringBuilder.toString();
    }

    public String getAttributes() {
        return "values("
                + "?,".repeat(Math.max(0, this.attributes - 1))
                + "?)";
    }

    public static Connection init(final Connection connection) throws SQLException {
        final String lineSeparator = System.lineSeparator();
        connection
                .prepareStatement(
                        "CREATE TABLE IF NOT EXISTS scale (" + lineSeparator
                                + "\tlower_bound INT NOT NULL," + lineSeparator
                                + "\tupper_bound INT NOT NULL," + lineSeparator
                                + "\tdescription VARCHAR ( 200 )," + lineSeparator
                                + "\tPRIMARY KEY (lower_bound, upper_bound)" + lineSeparator
                                + ");")
                .executeUpdate();
        connection
                .prepareStatement(
                        "CREATE TABLE IF NOT EXISTS metric (" + lineSeparator
                                + "\tmetric_name VARCHAR ( 20 ) NOT NULL," + lineSeparator
                                + "\tdescription VARCHAR ( 200 ) NOT NULL," + lineSeparator
                                + "\tlower_bound INT NOT NULL," + lineSeparator
                                + "\tupper_bound INT NOT NULL," + lineSeparator
                                + "\tFOREIGN KEY (lower_bound, upper_bound)" + lineSeparator
                                + "\t\tREFERENCES scale (lower_bound, upper_bound)," + lineSeparator
                                + "\tPRIMARY KEY (metric_name)" + lineSeparator
                                + ");")
                .executeUpdate();
        connection
                .prepareStatement(
                        "CREATE TABLE IF NOT EXISTS asset (" + lineSeparator
                                + "\tasset_name VARCHAR ( 20 ) NOT NULL," + lineSeparator
                                + "\tasset_type VARCHAR ( 20 ) NOT NULL," + lineSeparator
                                + "\tPRIMARY KEY (asset_name)" + lineSeparator
                                + ");")
                .executeUpdate();
        connection
                .prepareStatement(
                        "CREATE TABLE IF NOT EXISTS evaluation (" + lineSeparator
                                + "\tmetric_name VARCHAR ( 20 ) NOT NULL," + lineSeparator
                                + "\ttime_stamp VARCHAR ( 29 ) NOT NULL," + lineSeparator
                                + "\tresult_value JSONB NOT NULL," + lineSeparator
                                + "\tFOREIGN KEY (metric_name)" + lineSeparator
                                + "\t\tREFERENCES metric (metric_name)," + lineSeparator
                                + "\tPRIMARY KEY (metric_name, time_stamp)" + lineSeparator
                                + ");")
                .executeUpdate();
        connection
                .prepareStatement(
                        "CREATE TABLE IF NOT EXISTS evaluation_to_asset(" + lineSeparator
                                + "    metric_name VARCHAR ( 30 ) NOT NULL," + lineSeparator
                                + "    asset_name VARCHAR ( 30 ) NOT NULL," + lineSeparator
                                + "    time_stamp VARCHAR ( 29 ) NOT NULL," + lineSeparator
                                + "    FOREIGN KEY (asset_name)" + lineSeparator
                                + "    \t\tREFERENCES asset (asset_name)," + lineSeparator
                                + "    FOREIGN KEY (metric_name, time_stamp)" + lineSeparator
                                + "    \t\tREFERENCES evaluation (metric_name, time_stamp)," + lineSeparator
                                + "    PRIMARY KEY (metric_name, time_stamp)" + lineSeparator
                                + ");")
                .executeUpdate();
        return connection;
    }
}
